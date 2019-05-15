package cluster

const clusterTemplate = `
---
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: rook-ceph
spec:
  cephVersion:
    image: ceph/ceph:v14.2.1-20190430
    allowUnsupported: false
  dataDirHostPath: /var/lib/rook
  mon:
    count: 3
    allowMultiplePerNode: false
  dashboard:
    enabled: false
  network:
    hostNetwork: false
  rbdMirroring:
    workers: 0
  placement:
    all:
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
          - matchExpressions:
            - key: node-role.kubernetes.io/storage
              operator: In
              values:
              - "true"
  annotations:
  resources:
  storage:
    useAllNodes: false
    useAllDevices: false
    deviceFilter:
    location:
    config:
    nodes:
    - name: "k8s-04"
      devices:
      - name: "vdb"
    - name: "k8s-05"
      devices:
      - name: "vdb"
    - name: "k8s-06"
      devices:
      - name: "vdb"`
