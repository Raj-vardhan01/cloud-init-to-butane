package transpiler

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Transpile converts cloud-config YAML input bytes into Butane YAML output bytes.
func Transpile(cloudInitYAML []byte) ([]byte, error) {
	var cc CloudConfig
	if err := yaml.Unmarshal(cloudInitYAML, &cc); err != nil {
		return nil, fmt.Errorf("failed to parse cloud-config YAML: %w", err)
	}

	butane := ButaneConfig{
		Variant: "flatcar",
		Version: "1.1.0",
	}

	// 1. Convert Users
	for _, u := range cc.Users {
		if u.Name == "" {
			continue
		}
		bu := ButaneUser{
			Name:              u.Name,
			PasswdHash:        u.Passwd,
			Shell:             u.Shell,
			HomeDir:           u.HomeDir,
			PrimaryGroup:      u.PrimaryGroup,
			SSHAuthorizedKeys: u.SSHAuthorizedKeys,
		}

		if u.Groups != "" {
			groups := strings.Split(u.Groups, ",")
			bu.Groups = groups
		}

		if u.System {
			sys := true
			bu.System = &sys
		}

		butane.Passwd.Users = append(butane.Passwd.Users, bu)
	}

	// 2. Convert WriteFiles
	for _, f := range cc.WriteFiles {
		bf := ButaneFile{
			Path: f.Path,
		}

		content := f.Content
		if f.Encoding == "base64" || f.Encoding == "b64" {
			decoded, err := base64.StdEncoding.DecodeString(f.Content)
			if err == nil {
				content = string(decoded)
			}
		}
		bf.Contents.Inline = content

		if f.Permissions != "" {
			mode, err := strconv.Atoi(f.Permissions)
			if err == nil {
				bf.Mode = &mode
			}
		}

		if f.Owner != "" {
			parts := strings.Split(f.Owner, ":")
			bf.User = &ButaneFileUser{Name: parts[0]}
			bf.Group = &ButaneFileGroup{Name: parts[1]}
		}

		butane.Storage.Files = append(butane.Storage.Files, bf)
	}

	// 3. Convert RunCmd
	if len(cc.RunCmd) > 0 {
		var cmdLines []string
		for _, rawCmd := range cc.RunCmd {
			cmdList := rawCmd.([]any)
			var strParts []string
			for _, p := range cmdList {
				strParts = append(strParts, fmt.Sprintf("%v", p))
			}
			cmdLines = append(cmdLines, strings.Join(strParts, " "))
		}

		unitContent := fmt.Sprintf("[Unit]\nDescription=Cloud-init runcmd replacement\nAfter=network.target\n\n[Service]\nType=oneshot\nExecStart=/bin/sh -c '%s'\n\n[Install]\nWantedBy=multi-user.target\n", strings.Join(cmdLines, "; "))
		enabled := true
		butane.Systemd.Units = append(butane.Systemd.Units, ButaneUnit{
			Name:     "cloud-init-runcmd.service",
			Enabled:  &enabled,
			Contents: unitContent,
		})
	}

	// 4. Convert CaCerts
	if len(cc.CaCerts.Trusted) > 0 {
		for i, cert := range cc.CaCerts.Trusted {
			path := fmt.Sprintf("/etc/ssl/certs/ca-custom-%d.pem", i)
			butane.Storage.Files = append(butane.Storage.Files, ButaneFile{
				Path:     path,
				Contents: ButaneFileContents{Inline: cert},
			})
		}
	}

	out, err := yaml.Marshal(butane)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Butane YAML: %w", err)
	}
	return out, nil
}
