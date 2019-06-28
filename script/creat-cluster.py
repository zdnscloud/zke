#!/usr/bin/env python
import requests
import json

url = 'http://127.0.0.1:8080/apis/zkemanager.zcloud.cn/v1/zkeconfigs'
ssh_key_path = '/Users/wangyanwei/.ssh/id_rsa'

zke_config = {
    "name": "wang",
    "option": {
        "sshUser": "zcloud",
        "sshKey": "",
        "sshPort": "22",
        "dockerSocket": "/var/run/docker.sock",
        "ignoreDockerVersion": True,
        "clusterCidr": "10.42.0.0/16",
        "serviceClusterIpRange": "10.43.0.0/16",
        "clusterDomain": "cluster.local",
        "clusterDNSServiceIP": "10.43.0.10",
        "upstreamnameservers": [
            "114.114.114.114",
            "223.5.5.5"
        ],
    },
    "nodes": [
        {
            "name": "master",
            "address": "202.173.9.61",
            "roles": [
                "controlplane",
                "etcd"
            ],
        },
        {
            "name": "worker1",
            "address": "202.173.9.62",
            "roles": [
                "worker"
            ],
        },
        {
            "name": "worker2",
            "address": "202.173.9.63",
            "roles": [
                "worker"
            ],
        },
    ]
}

with open(ssh_key_path, 'r') as f:
    zke_config['option']['sshKey'] = f.read()


def creatCluster(url, params):
    headers = {'Content-type': 'application/json'}
    r = requests.post(url, data=json.dumps(params), headers=headers)
    print(json.dumps(params))
    return r


print(creatCluster(url, zke_config))
