package common

import (
	"bytes"
	"context"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/types"
	"golang.org/x/crypto/ssh"
	storagev1 "k8s.io/api/storage/v1"
	"strings"
)

func CheckStorageClassExist(cli client.Client, name string) bool {
	var exist bool
	scs := storagev1.StorageClassList{}
	err := cli.List(context.TODO(), nil, &scs)
	if err != nil {
		return exist
	}
	for _, s := range scs.Items {
		if s.Name == name {
			exist = true
			break
		}
	}
	return exist
}

func MakeSSHClient(node types.ZKEConfigNode) (*ssh.Client, error) {
	var sshKeyString, sshCertString string
	if !node.SSHAgentAuth {
		var err error
		sshKeyString, err = hosts.PrivateKeyPath(node.SSHKeyPath)
		if err != nil {
			return nil, err
		}

		if len(node.SSHCertPath) > 0 {
			sshCertString, err = hosts.CertificatePath(node.SSHCertPath)
			if err != nil {
				return nil, err
			}
		}
	}
	cfg, err := hosts.GetSSHConfig(node.User, sshKeyString, sshCertString, node.SSHAgentAuth)
	if err != nil {
		return nil, err
	}
	addr := node.Address + ":22"
	return ssh.Dial("tcp", addr, cfg)
}

func GetSSHCmdOut(client *ssh.Client, cmd string) (string, string, error) {
	var cmdout, cmderr string
	session, err := client.NewSession()
	if err != nil {
		return cmdout, "error", err
	}
	defer session.Close()
	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	session.Run(cmd)
	cmdout = strings.Replace(stdOut.String(), "\n", "", -1)
	cmderr = strings.Replace(stdErr.String(), "\n", "", -1)
	return cmdout, cmderr, nil
}
