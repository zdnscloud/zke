package ceph

const FilesystemTemplate = `
---
apiVersion: ceph.rook.io/v1
kind: CephFilesystem
metadata:
  name: {{.CephFilesystem}}
  namespace: rook-ceph
spec:
  metadataPool:
    replicated:
      size: {{.Replicas}}
  dataPools:
    - failureDomain: host
      replicated:
        size: {{.Replicas}}
  metadataServer:
    activeCount: 1
    activeStandby: true
    placement:
    annotations:
    resources:`
