package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"

	etcdclient "github.com/coreos/etcd/client"
)

func getEtcdClient(ctx context.Context, etcdHost *hosts.Host, cert, key []byte) (etcdclient.Client, error) {
	dialer, err := getEtcdDialer(etcdHost)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a dialer for host [%s]: %v", etcdHost.Address, err)
	}
	tlsConfig, err := getEtcdTLSConfig(cert, key)
	if err != nil {
		return nil, err
	}

	var DefaultEtcdTransport etcdclient.CancelableTransport = &http.Transport{
		Dial:                dialer,
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cfg := etcdclient.Config{
		Endpoints: []string{"https://" + etcdHost.InternalAddress + ":2379"},
		Transport: DefaultEtcdTransport,
	}

	return etcdclient.New(cfg)
}

func isEtcdHealthy(ctx context.Context, host *hosts.Host, cert, key []byte, url string) bool {
	log.Debugf(ctx, "[etcd] Check etcd cluster health")
	for i := 0; i < 3; i++ {
		dialer, err := getEtcdDialer(host)
		if err != nil {
			return false
		}
		tlsConfig, err := getEtcdTLSConfig(cert, key)
		if err != nil {
			log.Debugf(ctx, "[etcd] Failed to create etcd tls config for host [%s]: %v", host.Address, err)
			return false
		}

		hc := http.Client{
			Transport: &http.Transport{
				Dial:                dialer,
				TLSClientConfig:     tlsConfig,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}
		healthy, err := getHealthEtcd(hc, host, url)
		if err != nil {
			log.Debugf(ctx, "", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if healthy == "true" {
			log.Debugf(ctx, "[etcd] etcd cluster is healthy")
			return true
		}
	}
	return false
}

func getHealthEtcd(hc http.Client, host *hosts.Host, url string) (string, error) {
	healthy := struct{ Health string }{}
	resp, err := hc.Get(url)
	if err != nil {
		return healthy.Health, fmt.Errorf("Failed to get /health for host [%s]: %v", host.Address, err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return healthy.Health, fmt.Errorf("Failed to read response of /health for host [%s]: %v", host.Address, err)
	}
	resp.Body.Close()
	if err := json.Unmarshal(bytes, &healthy); err != nil {
		return healthy.Health, fmt.Errorf("Failed to unmarshal response of /health for host [%s]: %v", host.Address, err)
	}
	return healthy.Health, nil
}

func GetEtcdInitialCluster(hosts []*hosts.Host) string {
	initialCluster := ""
	for i, host := range hosts {
		initialCluster += fmt.Sprintf("etcd-%s=https://%s:2380", host.NodeName, host.InternalAddress)
		if i < (len(hosts) - 1) {
			initialCluster += ","
		}
	}
	return initialCluster
}

func getEtcdDialer(etcdHost *hosts.Host) (func(network, address string) (net.Conn, error), error) {
	etcdHost.LocalConnPort = 2379
	var etcdFactory hosts.DialerFactory
	etcdFactory = hosts.LocalConnFactory
	return etcdFactory(etcdHost)
}

func GetEtcdConnString(hosts []*hosts.Host) string {
	connString := ""
	for i, host := range hosts {
		connString += "https://" + host.InternalAddress + ":2379"
		if i < (len(hosts) - 1) {
			connString += ","
		}
	}
	return connString
}

func getEtcdTLSConfig(certificate, key []byte) (*tls.Config, error) {
	// get tls config
	x509Pair, err := tls.X509KeyPair([]byte(certificate), []byte(key))
	if err != nil {
		return nil, err

	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{x509Pair},
	}
	if err != nil {
		return nil, err
	}
	return tlsConfig, nil
}
