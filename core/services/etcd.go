package services

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/core/pki/cert"
	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/util"
	"github.com/zdnscloud/zke/types"

	etcdclient "github.com/coreos/etcd/client"
	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	"github.com/zdnscloud/cement/errgroup"
)

const (
	EtcdSnapshotPath     = "/opt/zke/etcd-snapshots/"
	EtcdRestorePath      = "/opt/zke/etcd-snapshots-restore/"
	EtcdDataDir          = "/var/lib/zcloud/etcd/"
	EtcdInitWaitTime     = 10
	EtcdSnapshotWaitTime = 5
)

func RunEtcdPlane(
	ctx context.Context,
	etcdHosts []*hosts.Host,
	etcdNodePlanMap map[string]types.ZKENodePlan,
	prsMap map[string]types.PrivateRegistry,
	alpineImage string,
	es types.ETCDService,
	certMap map[string]pki.CertificatePKI) error {
	log.Infof(ctx, "[%s] Building up etcd plane..", ETCDRole)
	for _, host := range etcdHosts {
		etcdProcess := etcdNodePlanMap[host.Address].Processes[EtcdContainerName]
		imageCfg, hostCfg, _ := GetProcessConfig(etcdProcess)
		if err := docker.DoRunContainer(ctx, host.DClient, imageCfg, hostCfg, EtcdContainerName, host.Address, ETCDRole, prsMap); err != nil {
			return err
		}
		if *es.Snapshot == true {
			if err := RunEtcdSnapshotSave(ctx, host, prsMap, alpineImage, EtcdSnapshotContainerName, false, es); err != nil {
				return err
			}
			if err := pki.SaveBackupBundleOnHost(ctx, host, alpineImage, EtcdSnapshotPath, prsMap); err != nil {
				return err
			}
		} else {
			if err := docker.DoRemoveContainer(ctx, host.DClient, EtcdSnapshotContainerName, host.Address); err != nil {
				return err
			}
		}
		if err := createLogLink(ctx, host, EtcdContainerName, ETCDRole, alpineImage, prsMap); err != nil {
			return err
		}
	}
	log.Infof(ctx, "[%s] Successfully started etcd plane.. Checking etcd cluster health", ETCDRole)
	clientCert := cert.EncodeCertPEM(certMap[pki.KubeNodeCertName].Certificate)
	clientkey := cert.EncodePrivateKeyPEM(certMap[pki.KubeNodeCertName].Key)
	var healthy bool
	var checkTimes = 0

	for {
		select {
		case <-ctx.Done():
			return util.CancelErr
		default:
			for _, host := range etcdHosts {
				_, _, healthCheckURL := GetProcessConfig(etcdNodePlanMap[host.Address].Processes[EtcdContainerName])
				if healthy = isEtcdHealthy(ctx, host, clientCert, clientkey, healthCheckURL); healthy {
					break
				}
			}
			if !healthy {
				checkTimes = checkTimes + 1
				log.Warnf(ctx, "[Etcd] Etcd Cluster is not healthy, has checked [%s] times!", strconv.Itoa(checkTimes))
			} else {
				return nil
			}
		}
	}
	return nil
}

func RestartEtcdPlane(ctx context.Context, etcdHosts []*hosts.Host) error {
	log.Infof(ctx, "[%s] Restarting up etcd plane..", ETCDRole)

	_, err := errgroup.Batch(etcdHosts, func(h interface{}) (interface{}, error) {
		runHost := h.(*hosts.Host)
		return nil, docker.DoRestartContainer(ctx, runHost.DClient, EtcdContainerName, runHost.Address)
	})
	if err != nil {
		return err
	}

	log.Infof(ctx, "[%s] Successfully restarted etcd plane..", ETCDRole)
	return nil
}

func RemoveEtcdPlane(ctx context.Context, etcdHosts []*hosts.Host, force bool) error {
	log.Infof(ctx, "[%s] Tearing down etcd plane..", ETCDRole)

	_, err := errgroup.Batch(etcdHosts, func(h interface{}) (interface{}, error) {
		runHost := h.(*hosts.Host)
		if err := docker.DoRemoveContainer(ctx, runHost.DClient, EtcdContainerName, runHost.Address); err != nil {
			return nil, err
		}
		if !runHost.IsWorker || !runHost.IsControl || force {
			// remove unschedulable kubelet on etcd host
			if err := removeKubelet(ctx, runHost); err != nil {
				return nil, err
			}
			if err := removeKubeproxy(ctx, runHost); err != nil {
				return nil, err
			}
			if err := removeNginxProxy(ctx, runHost); err != nil {
				return nil, err
			}
			if err := removeSidekick(ctx, runHost); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	log.Infof(ctx, "[%s] Successfully tore down etcd plane..", ETCDRole)
	return nil
}

func AddEtcdMember(ctx context.Context, toAddEtcdHost *hosts.Host, etcdHosts []*hosts.Host, cert, key []byte) error {
	log.Infof(ctx, "[add/%s] Adding member [etcd-%s] to etcd cluster", ETCDRole, toAddEtcdHost.NodeName)
	peerURL := fmt.Sprintf("https://%s:2380", toAddEtcdHost.InternalAddress)
	added := false
	for _, host := range etcdHosts {
		if host.Address == toAddEtcdHost.Address {
			continue
		}
		etcdClient, err := getEtcdClient(ctx, host, cert, key)
		if err != nil {
			log.Debugf(ctx, "Failed to create etcd client for host [%s]: %v", host.Address, err)
			continue
		}
		memAPI := etcdclient.NewMembersAPI(etcdClient)
		if _, err := memAPI.Add(ctx, peerURL); err != nil {
			log.Debugf(ctx, "Failed to Add etcd member [%s] from host: %v", host.Address, err)
			continue
		}
		added = true
		break
	}
	if !added {
		return fmt.Errorf("Failed to add etcd member [etcd-%s] to etcd cluster", toAddEtcdHost.NodeName)
	}
	log.Infof(ctx, "[add/%s] Successfully Added member [etcd-%s] to etcd cluster", ETCDRole, toAddEtcdHost.NodeName)
	return nil
}

func RemoveEtcdMember(ctx context.Context, etcdHost *hosts.Host, etcdHosts []*hosts.Host, cert, key []byte) error {
	log.Infof(ctx, "[remove/%s] Removing member [etcd-%s] from etcd cluster", ETCDRole, etcdHost.NodeName)
	var mID string
	removed := false
	for _, host := range etcdHosts {
		etcdClient, err := getEtcdClient(ctx, host, cert, key)
		if err != nil {
			log.Debugf(ctx, "Failed to create etcd client for host [%s]: %v", host.Address, err)
			continue
		}
		memAPI := etcdclient.NewMembersAPI(etcdClient)
		members, err := memAPI.List(ctx)
		if err != nil {
			log.Debugf(ctx, "Failed to list etcd members from host [%s]: %v", host.Address, err)
			continue
		}
		for _, member := range members {
			if member.Name == fmt.Sprintf("etcd-%s", etcdHost.NodeName) {
				mID = member.ID
				break
			}
		}
		if err := memAPI.Remove(ctx, mID); err != nil {
			log.Debugf(ctx, "Failed to list etcd members from host [%s]: %v", host.Address, err)
			continue
		}
		removed = true
		break
	}
	if !removed {
		return fmt.Errorf("Failed to delete etcd member [etcd-%s] from etcd cluster", etcdHost.NodeName)
	}
	log.Infof(ctx, "[remove/%s] Successfully removed member [etcd-%s] from etcd cluster", ETCDRole, etcdHost.NodeName)
	return nil
}

func ReloadEtcdCluster(ctx context.Context, readyEtcdHosts []*hosts.Host, newHost *hosts.Host, cert, key []byte, prsMap map[string]types.PrivateRegistry, etcdNodePlanMap map[string]types.ZKENodePlan, alpineImage string) error {
	imageCfg, hostCfg, _ := GetProcessConfig(etcdNodePlanMap[newHost.Address].Processes[EtcdContainerName])
	if err := docker.DoRunContainer(ctx, newHost.DClient, imageCfg, hostCfg, EtcdContainerName, newHost.Address, ETCDRole, prsMap); err != nil {
		return err
	}
	if err := createLogLink(ctx, newHost, EtcdContainerName, ETCDRole, alpineImage, prsMap); err != nil {
		return err
	}
	time.Sleep(EtcdInitWaitTime * time.Second)
	var healthy bool
	for _, host := range readyEtcdHosts {
		_, _, healthCheckURL := GetProcessConfig(etcdNodePlanMap[host.Address].Processes[EtcdContainerName])
		if healthy = isEtcdHealthy(ctx, host, cert, key, healthCheckURL); healthy {
			break
		}
	}
	if !healthy {
		return fmt.Errorf("[etcd] Etcd Cluster is not healthy")
	}
	return nil
}

func IsEtcdMember(ctx context.Context, etcdHost *hosts.Host, etcdHosts []*hosts.Host, cert, key []byte) (bool, error) {
	var listErr error
	peerURL := fmt.Sprintf("https://%s:2380", etcdHost.InternalAddress)
	for _, host := range etcdHosts {
		if host.Address == etcdHost.Address {
			continue
		}
		etcdClient, err := getEtcdClient(ctx, host, cert, key)
		if err != nil {
			listErr = errors.Wrapf(err, "Failed to create etcd client for host [%s]", host.Address)
			log.Debugf(ctx, "Failed to create etcd client for host [%s]: %v", host.Address, err)
			continue
		}
		memAPI := etcdclient.NewMembersAPI(etcdClient)
		members, err := memAPI.List(ctx)
		if err != nil {
			listErr = errors.Wrapf(err, "Failed to create etcd client for host [%s]", host.Address)
			log.Debugf(ctx, "Failed to list etcd cluster members [%s]: %v", etcdHost.Address, err)
			continue
		}
		for _, member := range members {
			if strings.Contains(member.PeerURLs[0], peerURL) {
				log.Infof(ctx, "[etcd] member [%s] is already part of the etcd cluster", etcdHost.Address)
				return true, nil
			}
		}
		// reset the list of errors to handle new hosts
		listErr = nil
		break
	}
	if listErr != nil {
		return false, listErr
	}
	return false, nil
}

func RunEtcdSnapshotSave(ctx context.Context, etcdHost *hosts.Host, prsMap map[string]types.PrivateRegistry, etcdSnapshotImage string, name string, once bool, es types.ETCDService) error {
	log.Infof(ctx, "[etcd] Saving snapshot [%s] on host [%s]", name, etcdHost.Address)
	imageCfg := &container.Config{
		Cmd: []string{
			"/opt/zke-tools/etcd-backup",
			"etcd-backup",
			"save",
			"--cacert", pki.GetCertPath(pki.CACertName),
			"--cert", pki.GetCertPath(pki.KubeNodeCertName),
			"--key", pki.GetKeyPath(pki.KubeNodeCertName),
			"--name", name,
			"--endpoints=" + etcdHost.InternalAddress + ":2379",
		},
		Image: etcdSnapshotImage,
	}
	if once {
		imageCfg.Cmd = append(imageCfg.Cmd, "--once")
	} else if es.BackupConfig == nil {
		imageCfg.Cmd = append(imageCfg.Cmd, "--retention="+es.Retention)
		imageCfg.Cmd = append(imageCfg.Cmd, "--creation="+es.Creation)
	}

	if es.BackupConfig != nil {
		imageCfg = configBackupImgCmd(ctx, imageCfg, es.BackupConfig)
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/backup", EtcdSnapshotPath),
			fmt.Sprintf("%s:/etc/kubernetes:z", path.Join(etcdHost.PrefixPath, "/etc/kubernetes"))},
		NetworkMode:   container.NetworkMode("host"),
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}

	if once {
		if err := docker.DoRunContainer(ctx, etcdHost.DClient, imageCfg, hostCfg, EtcdSnapshotOnceContainerName, etcdHost.Address, ETCDRole, prsMap); err != nil {
			return err
		}
		status, _, stderr, err := docker.GetContainerOutput(ctx, etcdHost.DClient, EtcdSnapshotOnceContainerName, etcdHost.Address)
		if status != 0 || err != nil {
			if removeErr := docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdSnapshotOnceContainerName); removeErr != nil {
				log.Warnf(ctx, "Failed to remove container [%s]: %v", removeErr)
			}
			if err != nil {
				return err
			}
			return fmt.Errorf("Failed to take one-time snapshot, exit code [%d]: %v", status, stderr)
		}

		return docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdSnapshotOnceContainerName)
	}

	if err := docker.DoRunContainer(ctx, etcdHost.DClient, imageCfg, hostCfg, EtcdSnapshotContainerName, etcdHost.Address, ETCDRole, prsMap); err != nil {
		return err
	}
	// check if the container exited with error
	snapshotCont, err := docker.InspectContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdSnapshotContainerName)
	if err != nil {
		return err
	}
	time.Sleep(EtcdSnapshotWaitTime * time.Second)
	if snapshotCont.State.Status == "exited" || snapshotCont.State.Restarting {
		log.Warnf(ctx, "Etcd rolling snapshot container failed to start correctly")
		return docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdSnapshotContainerName)

	}
	return nil
}

func RestoreEtcdSnapshot(ctx context.Context, etcdHost *hosts.Host, prsMap map[string]types.PrivateRegistry, etcdRestoreImage, snapshotName, initCluster string) error {
	log.Infof(ctx, "[etcd] Restoring [%s] snapshot on etcd host [%s]", snapshotName, etcdHost.Address)
	nodeName := pki.GetEtcdCrtName(etcdHost.InternalAddress)
	snapshotPath := fmt.Sprintf("%s%s", EtcdSnapshotPath, snapshotName)

	// make sure that restore path is empty otherwise etcd restore will fail
	imageCfg := &container.Config{
		Cmd: []string{
			"sh", "-c", strings.Join([]string{
				"rm -rf", EtcdRestorePath,
				"&& /usr/local/bin/etcdctl",
				fmt.Sprintf("--endpoints=[%s:2379]", etcdHost.InternalAddress),
				"--cacert", pki.GetCertPath(pki.CACertName),
				"--cert", pki.GetCertPath(nodeName),
				"--key", pki.GetKeyPath(nodeName),
				"snapshot", "restore", snapshotPath,
				"--data-dir=" + EtcdRestorePath,
				"--name=etcd-" + etcdHost.NodeName,
				"--initial-cluster=" + initCluster,
				"--initial-cluster-token=etcd-cluster-1",
				"--initial-advertise-peer-urls=https://" + etcdHost.InternalAddress + ":2380",
				"&& mv", EtcdRestorePath + "*", EtcdDataDir,
				"&& rm -rf", EtcdRestorePath,
			}, " "),
		},
		Env:   []string{"ETCDCTL_API=3"},
		Image: etcdRestoreImage,
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			"/opt/zke/:/opt/zke/:z",
			fmt.Sprintf("%s:/var/lib/zcloud/etcd:z", path.Join(etcdHost.PrefixPath, "/var/lib/etcd")),
			fmt.Sprintf("%s:/etc/kubernetes:z", path.Join(etcdHost.PrefixPath, "/etc/kubernetes"))},
		NetworkMode: container.NetworkMode("host"),
	}
	if err := docker.DoRunContainer(ctx, etcdHost.DClient, imageCfg, hostCfg, EtcdRestoreContainerName, etcdHost.Address, ETCDRole, prsMap); err != nil {
		return err
	}
	status, err := docker.WaitForContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdRestoreContainerName)
	if err != nil {
		return err
	}
	if status != 0 {
		containerLog, _, err := docker.GetContainerLogsStdoutStderr(ctx, etcdHost.DClient, EtcdRestoreContainerName, "5", false)
		if err != nil {
			return err
		}
		if err := docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdRestoreContainerName); err != nil {
			return err
		}
		// printing the restore container's logs
		return fmt.Errorf("Failed to run etcd restore container, exit status is: %d, container logs: %s", status, containerLog)
	}
	return docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdRestoreContainerName)
}

func GetEtcdSnapshotChecksum(ctx context.Context, etcdHost *hosts.Host, prsMap map[string]types.PrivateRegistry, alpineImage, snapshotName string) (string, error) {
	var checksum string
	var err error

	snapshotPath := fmt.Sprintf("%s%s", EtcdSnapshotPath, snapshotName)
	imageCfg := &container.Config{
		Cmd: []string{
			"sh", "-c", strings.Join([]string{
				"md5sum", snapshotPath,
				"|", "cut", "-f1", "-d' '", "|", "tr", "-d", "'\n'"}, " "),
		},
		Image: alpineImage,
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			"/opt/zke/:/opt/zke/:z",
		}}

	if err := docker.DoRunContainer(ctx, etcdHost.DClient, imageCfg, hostCfg, EtcdChecksumContainerName, etcdHost.Address, ETCDRole, prsMap); err != nil {
		return checksum, err
	}
	if _, err := docker.WaitForContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdChecksumContainerName); err != nil {
		return checksum, err
	}
	_, checksum, err = docker.GetContainerLogsStdoutStderr(ctx, etcdHost.DClient, EtcdChecksumContainerName, "1", false)
	if err != nil {
		return checksum, err
	}
	if err := docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdChecksumContainerName); err != nil {
		return checksum, err
	}
	return checksum, nil
}

func configBackupImgCmd(ctx context.Context, imageCfg *container.Config, bc *types.BackupConfig) *container.Config {
	cmd := []string{
		"--creation=" + fmt.Sprintf("%dh", bc.IntervalHours),
		"--retention=" + fmt.Sprintf("%dh", bc.Retention*bc.IntervalHours),
	}
	imageCfg.Cmd = append(imageCfg.Cmd, cmd...)
	return imageCfg
}

func StartBackupServer(ctx context.Context, etcdHost *hosts.Host, prsMap map[string]types.PrivateRegistry, etcdSnapshotImage string, name string) error {
	log.Infof(ctx, "[etcd] starting backup server on host [%s]", etcdHost.Address)
	imageCfg := &container.Config{
		Cmd: []string{
			"/opt/zke-tools/etcd-backup",
			"etcd-backup",
			"serve",
			"--name", name,
			"--cacert", pki.GetCertPath(pki.CACertName),
			"--cert", pki.GetCertPath(pki.KubeNodeCertName),
			"--key", pki.GetKeyPath(pki.KubeNodeCertName),
		},
		Image: etcdSnapshotImage,
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/backup", EtcdSnapshotPath),
			fmt.Sprintf("%s:/etc/kubernetes:z", path.Join(etcdHost.PrefixPath, "/etc/kubernetes"))},
		NetworkMode:   container.NetworkMode("host"),
		RestartPolicy: container.RestartPolicy{Name: "on-failure"},
	}
	return docker.DoRunContainer(ctx, etcdHost.DClient, imageCfg, hostCfg, EtcdServeBackupContainerName, etcdHost.Address, ETCDRole, prsMap)
}

func DownloadEtcdSnapshotFromBackupServer(ctx context.Context, etcdHost *hosts.Host, prsMap map[string]types.PrivateRegistry, etcdSnapshotImage, name string, backupServer *hosts.Host) error {
	log.Infof(ctx, "[etcd] Get snapshot [%s] on host [%s]", name, etcdHost.Address)
	imageCfg := &container.Config{
		Cmd: []string{
			"/opt/zke-tools/etcd-backup",
			"etcd-backup",
			"download",
			"--name", name,
			"--local-endpoint", backupServer.Address,
			"--cacert", pki.GetCertPath(pki.CACertName),
			"--cert", pki.GetCertPath(pki.KubeNodeCertName),
			"--key", pki.GetKeyPath(pki.KubeNodeCertName),
		},
		Image: etcdSnapshotImage,
	}

	hostCfg := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/backup", EtcdSnapshotPath),
			fmt.Sprintf("%s:/etc/kubernetes:z", path.Join(etcdHost.PrefixPath, "/etc/kubernetes"))},
		NetworkMode:   container.NetworkMode("host"),
		RestartPolicy: container.RestartPolicy{Name: "on-failure"},
	}

	if err := docker.DoRunContainer(ctx, etcdHost.DClient, imageCfg, hostCfg, EtcdDownloadBackupContainerName, etcdHost.Address, ETCDRole, prsMap); err != nil {
		return err
	}

	status, _, stderr, err := docker.GetContainerOutput(ctx, etcdHost.DClient, EtcdDownloadBackupContainerName, etcdHost.Address)
	if status != 0 || err != nil {
		if removeErr := docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdDownloadBackupContainerName); removeErr != nil {
			log.Warnf(ctx, "Failed to remove container [%s]: %v", removeErr)
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to download etcd snapshot from backup server [%s], exit code [%d]: %v", backupServer.Address, status, stderr)
	}
	return docker.RemoveContainer(ctx, etcdHost.DClient, etcdHost.Address, EtcdDownloadBackupContainerName)
}
