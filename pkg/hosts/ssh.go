package hosts

import (
	"bytes"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func (h *Host) GetSSHClient() (*ssh.Client, error) {
	dialer, err := NewDialer(h, "network")
	if err != nil {
		return nil, err
	}
	return dialer.getSSHTunnelConnection()
}

func (h *Host) GetSSHCmdOutput(cli *ssh.Client, cmd string) (string, string, error) {
	var cmdout, cmderr string
	session, err := cli.NewSession()
	if err != nil {
		return cmdout, "error", err
	}
	defer session.Close()
	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	session.Run(cmd)
	cmdout = strings.TrimSpace(stdOut.String())
	cmderr = strings.TrimSpace(stdErr.String())
	return cmdout, cmderr, nil
}

func (h *Host) GetSftpClient(cli *ssh.Client) (*sftp.Client, error) {
	sc, err := sftp.NewClient(cli)
	if err != nil {
		return nil, err
	}
	return sc, nil
}
