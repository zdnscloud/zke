package fscsi

const fscsisecretTemplate = `
---
apiVersion: v1
kind: Secret
metadata:
  name: csi-cephfs-secret
  namespace: default
data:
  userID: YWRtaW4=
  userKey: QVFEbWt0cGNsdmZ0SkJBQXBtY2tRYk03UjIxVGdUaEhUQTZOeXc9PQ==

  adminID: YWRtaW4=
  adminKey: QVFEbWt0cGNsdmZ0SkJBQXBtY2tRYk03UjIxVGdUaEhUQTZOeXc9PQ==
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-cephfs
provisioner: cephfs.csi.ceph.com
parameters:
  monitors: 10.43.214.49:6789,10.43.204.53:6789,10.43.153.203:6789
  provisionVolume: "true"
  pool: myfs-data0
  csi.storage.k8s.io/provisioner-secret-name: csi-cephfs-secret
  csi.storage.k8s.io/provisioner-secret-namespace: default
  csi.storage.k8s.io/node-stage-secret-name: csi-cephfs-secret
  csi.storage.k8s.io/node-stage-secret-namespace: default
reclaimPolicy: Delete`
