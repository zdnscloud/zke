# 使用指导
## 创建集群
1. 生成zke配置
```shell
./zke generateconfig
```
执行命令后会生成cluster.yml
```yaml
cluster_name: ""
option:
  ssh_user: ""
  ssh_key: ""
  ssh_key_path: ""
  ssh_port: ""
  ignore_docker_version: false
  cluster_cidr: ""
  service_cidr: ""
  cluster_domain: ""
  up_stream_name_servers: []
  disable_port_check: false
nodes:
- name: ""
  address: ""
  roles: []
config_version: v1.0.10
```
2. 修改配置文件
```yaml
cluster_name: local
option:
  ssh_user: zcloud
  ssh_key: ""
  ssh_key_path: ~/.ssh/id_rsa
  ssh_port: "22"
  ignore_docker_version: false
  cluster_cidr: 10.42.0.0/16
  service_cidr: 10.43.0.0/16
  cluster_domain: cluster.w
  up_stream_name_servers:
  - 223.5.5.5
  - 114.114.114.114
  disable_port_check: false
network:
  plugin: flannel
  iface: ""
nodes:
- name: master
  address: 202.173.9.61
  roles:
  - controlplane
  - etcd
- name: worker1
  address: 202.173.9.62
  roles:
  - etcd
  - worker
  - edge
- name: worker2
  address: 202.173.9.63
  roles:
  - etcd
  - worker
  - edge
config_version: v1.0.10
```
> 其中node.roles有controlplane、etcd、worker、edge可选，edge需依赖于contronplane或worker，不可单独存在，若需要手动指定flannel网卡名称，可自行如上在配置文件中添加network部分配置
3. 执行zke up命令创建集群
```shell
./zke up
```

