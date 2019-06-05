package shell

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"syscall"
)

func Shell(cmd string, args ...string) (string, error) {
	result := exec.Command(cmd, args...)
	out, err := result.Output()
	return string(out), err
}

func ShellReturnErrMsg(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err.Error()
	}
	if err := cmd.Start(); err != nil {
		return err.Error()
	}
	ret, _ := ioutil.ReadAll(stderr)
	if err := cmd.Wait(); err != nil {
		return string(ret)
	} else {
		return ""
	}
}

func FindByPid(pidFile string) (int, error) {
	dataTmp, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return -1, err
	}

	data := []byte{}
	for _, d := range dataTmp {
		if d != '\n' && d != '\r' {
			data = append(data, d)
		}
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return -1, err
	}
	err = syscall.Kill(pid, 0)
	if err != nil {
		return -1, err
	}
	return pid, nil
}

func SetULimit(limit uint64) error {
	var rlimit syscall.Rlimit
	rlimit.Max = limit
	rlimit.Cur = limit
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
}
