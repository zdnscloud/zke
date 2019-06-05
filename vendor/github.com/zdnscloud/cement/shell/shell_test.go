package shell

import (
	ut "github.com/zdnscloud/cement/unittest"
	"testing"
)

func TestShell(t *testing.T) {
	fs, err := Shell("ls")
	ut.Assert(t, err == nil, "ls shouldn't return any error,but %v", err)
	ut.Equal(t, fs, "os_cmd.go\nshell.go\nshell_test.go\n")

	errMsg := ShellReturnErrMsg("ls", "goodboy")
	ut.Equal(t, errMsg, "ls: cannot access 'goodboy': No such file or directory\n")
}
