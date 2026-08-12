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
