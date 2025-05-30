package ssh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddHostSection(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		execPath   string
		host       string
		user       string
		context    string
		workspace  string
		workdir    string
		command    string
		gpgagent   bool
		devSpaceHome string
		expected   string
	}{
		{
			name:       "Basic host addition",
			config:     "",
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "",
			gpgagent:   false,
			devSpaceHome: "",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace
  User testuser
# DevSpace End testhost`,
		},
		{
			name:       "Basic host addition with DEVSPACE_HOME",
			config:     "",
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "",
			gpgagent:   false,
			devSpaceHome: "C:\\\\White Space\\devspace\\test",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace --devspace-home "C:\\White Space\devspace\test"
  User testuser
# DevSpace End testhost`,
		},
		{
			name:       "Host addition with workdir",
			config:     "",
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "/path/to/workdir",
			command:    "",
			gpgagent:   false,
			devSpaceHome: "",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace --workdir "/path/to/workdir"
  User testuser
# DevSpace End testhost`,
		},
		{
			name:       "Host addition with gpg agent",
			config:     "",
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "",
			gpgagent:   true,
			devSpaceHome: "",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace --gpg-agent-forwarding
  User testuser
# DevSpace End testhost`,
		},
		{
			name:       "Host addition with custom command",
			config:     "",
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "ssh -W %h:%p bastion",
			gpgagent:   false,
			devSpaceHome: "",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "ssh -W %h:%p bastion"
  User testuser
# DevSpace End testhost`,
		},
		{
			name: "Host addition to existing config",
			config: `Host existinghost
  User existinguser`,
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "",
			gpgagent:   false,
			devSpaceHome: "",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace
  User testuser
# DevSpace End testhost
Host existinghost
  User existinguser`,
		},
		{
			name: "Host addition to existing config with DevSpace host",
			config: `# DevSpace Start existingtesthost
Host existingtesthost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace
  User testuser
# DevSpace End testhost

Host existinghost
  User existinguser`,
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "",
			gpgagent:   false,
			devSpaceHome: "",
			expected: `# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace
  User testuser
# DevSpace End testhost
# DevSpace Start existingtesthost
Host existingtesthost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace
  User testuser
# DevSpace End testhost

Host existinghost
  User existinguser`,
		},
		{
			name: "Host addition after top level includes",
			config: `Include ~/config1 

Include ~/config2



Include ~/config3`,
			execPath:   "/path/to/exec",
			host:       "testhost",
			user:       "testuser",
			context:    "testcontext",
			workspace:  "testworkspace",
			workdir:    "",
			command:    "",
			gpgagent:   false,
			devSpaceHome: "",
			expected: `Include ~/config1 

Include ~/config2



Include ~/config3
# DevSpace Start testhost
Host testhost
  ForwardAgent yes
  LogLevel error
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  HostKeyAlgorithms rsa-sha2-256,rsa-sha2-512,ssh-rsa
  ProxyCommand "/path/to/exec" ssh --stdio --context testcontext --user testuser testworkspace
  User testuser
# DevSpace End testhost`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := addHostSection(tt.config, tt.execPath, tt.host, tt.user, tt.context, tt.workspace, tt.workdir, tt.command, tt.gpgagent, tt.devSpaceHome)
			if err != nil {
				t.Errorf("Failed with err: %v", err)
			}

			if result != tt.expected {
				t.Errorf("addHostSection result does not match expected.\nGot:\n%s\nExpected:\n%s", result, tt.expected)
				t.Errorf("addHostSection result does not match expected:\n%s", cmp.Diff(result, tt.expected))
			}

			if !strings.Contains(result, MarkerEndPrefix+tt.host) {
				t.Errorf("Result does not contain the end marker: %s", MarkerEndPrefix+tt.host)
			}

			if !strings.Contains(result, "Host "+tt.host) {
				t.Errorf("Result does not contain the Host line: Host %s", tt.host)
			}

			if !strings.Contains(result, "User "+tt.user) {
				t.Errorf("Result does not contain the User line: User %s", tt.user)
			}

			if tt.command != "" && !strings.Contains(result, fmt.Sprintf("ProxyCommand \"%s\"", tt.command)) {
				t.Errorf("Result does not contain the custom ProxyCommand: %s", tt.command)
			}

			if tt.workdir != "" && !strings.Contains(result, fmt.Sprintf("--workdir \"%s\"", tt.workdir)) {
				t.Errorf("Result does not contain the workdir: %s", tt.workdir)
			}

			if tt.gpgagent && !strings.Contains(result, "--gpg-agent-forwarding") {
				t.Errorf("Result does not contain gpg-agent-forwarding flag")
			}

			if tt.config != "" && !strings.Contains(result, tt.config) {
				t.Errorf("Result does not contain the original config")
			}
		})
	}
}
