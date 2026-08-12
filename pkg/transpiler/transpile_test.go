package transpiler

import (
	"strings"
	"testing"
)

func TestTranspile_Basic(t *testing.T) {
	input := `
users:
  - name: core
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2E...
write_files:
  - path: /etc/test.txt
    content: "hello world"
`
	out, err := Transpile([]byte(input))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}
	outStr := string(out)
	if !strings.Contains(outStr, "variant: flatcar") {
		t.Errorf("expected variant: flatcar, got:\n%s", outStr)
	}
}

// TestTranspile_CaCerts_NilGuard verifies that a cloud-config without a
// ca_certs block does not cause a nil pointer dereference panic.
func TestTranspile_CaCerts_NilGuard(t *testing.T) {
	input := `
users:
  - name: admin
`
	_, err := Transpile([]byte(input))
	if err != nil {
		t.Fatalf("Transpile failed unexpectedly: %v", err)
	}
}

// TestTranspile_WriteFiles_OctalPermissions verifies that write_files permissions
// given as octal strings (e.g. "0644") are correctly parsed as octal, not decimal.
func TestTranspile_WriteFiles_OctalPermissions(t *testing.T) {
	input := `
write_files:
  - path: /etc/config.json
    content: "{}"
    permissions: "0644"
`
	out, err := Transpile([]byte(input))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}
	outStr := string(out)
	// octal 0644 = decimal 420; strconv.Atoi would produce 644 (wrong)
	if !strings.Contains(outStr, "mode: 420") {
		t.Errorf("expected mode: 420 (octal 0644), got:\n%s", outStr)
	}
}
