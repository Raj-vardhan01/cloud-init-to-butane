package transpiler

// CloudConfig represents the cloud-init / cloud-config YAML structure.
type CloudConfig struct {
	Users             []CloudUser     `yaml:"users,omitempty"`
	Groups            []CloudGroup    `yaml:"groups,omitempty"`
	WriteFiles        []WriteFile     `yaml:"write_files,omitempty"`
	RunCmd            []any           `yaml:"runcmd,omitempty"`
	BootCmd           []any           `yaml:"bootcmd,omitempty"`
	CaCerts           *CaCertsConfig  `yaml:"ca_certs,omitempty"`
	SSHAuthorizedKeys []string        `yaml:"ssh_authorized_keys,omitempty"`
}

type CloudUser struct {
	Name              string   `yaml:"name,omitempty"`
	Gecos             string   `yaml:"gecos,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	HomeDir           string   `yaml:"homedir,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
	PrimaryGroup      string   `yaml:"primary_group,omitempty"`
	Sudo              string   `yaml:"sudo,omitempty"`
	Passwd            string   `yaml:"passwd,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	System            bool     `yaml:"system,omitempty"`
	NoCreateHome      bool     `yaml:"no_create_home,omitempty"`
}

type CloudGroup struct {
	Name    string   `yaml:"name,omitempty"`
	Members []string `yaml:"members,omitempty"`
	System  bool     `yaml:"system,omitempty"`
}

type WriteFile struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Owner       string `yaml:"owner"`
	Permissions string `yaml:"permissions"`
	Encoding    string `yaml:"encoding"`
	Append      bool   `yaml:"append"`
}

type CaCertsConfig struct {
	RemoveDefaults bool     `yaml:"remove_defaults,omitempty"`
	Trusted        []string `yaml:"trusted,omitempty"`
}

// ButaneConfig represents the target Butane v1.1.0 (Flatcar variant) YAML structure.
type ButaneConfig struct {
	Variant string         `yaml:"variant"`
	Version string         `yaml:"version"`
	Passwd  ButanePasswd   `yaml:"passwd,omitempty"`
	Storage ButaneStorage  `yaml:"storage,omitempty"`
	Systemd ButaneSystemd  `yaml:"systemd,omitempty"`
}

type ButanePasswd struct {
	Users  []ButaneUser  `yaml:"users,omitempty"`
	Groups []ButaneGroup `yaml:"groups,omitempty"`
}

type ButaneUser struct {
	Name              string   `yaml:"name"`
	PasswdHash        string   `yaml:"password_hash,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	HomeDir           string   `yaml:"home_dir,omitempty"`
	PrimaryGroup      string   `yaml:"primary_group,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	System            *bool    `yaml:"system,omitempty"`
	NoCreateHome      *bool    `yaml:"no_create_home,omitempty"`
}

type ButaneGroup struct {
	Name   string `yaml:"name"`
	System *bool  `yaml:"system,omitempty"`
}

type ButaneStorage struct {
	Files []ButaneFile `yaml:"files,omitempty"`
}

type ButaneFile struct {
	Path      string             `yaml:"path"`
	Mode      *int               `yaml:"mode,omitempty"`
	Overwrite *bool              `yaml:"overwrite,omitempty"`
	Contents  ButaneFileContents `yaml:"contents"`
	User      *ButaneFileUser    `yaml:"user,omitempty"`
	Group     *ButaneFileGroup   `yaml:"group,omitempty"`
}

type ButaneFileContents struct {
	Inline string `yaml:"inline,omitempty"`
}

type ButaneFileUser struct {
	Name string `yaml:"name,omitempty"`
}

type ButaneFileGroup struct {
	Name string `yaml:"group,omitempty"`
}

type ButaneSystemd struct {
	Units []ButaneUnit `yaml:"units,omitempty"`
}

type ButaneUnit struct {
	Name     string `yaml:"name"`
	Enabled  *bool  `yaml:"enabled,omitempty"`
	Contents string `yaml:"contents,omitempty"`
}
