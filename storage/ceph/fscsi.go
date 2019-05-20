package ceph

const FscsiTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: rook-ceph
---
{{- if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cephfs-csi-attacher
  namespace: rook-ceph
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-external-attacher-runner
  namespace: rook-ceph
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["volumeattachments"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["csi.storage.k8s.io"]
    resources: ["csinodeinfos"]
    verbs: ["get", "list", "watch"]

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-csi-attacher-role
  namespace: rook-ceph
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-attacher
    namespace: rook-ceph
roleRef:
  kind: ClusterRole
  name: cephfs-external-attacher-runner
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cephfs-csi-provisioner
  namespace: rook-ceph
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-external-provisioner-runner
  namespace: rook-ceph
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["list", "watch", "create", "update", "patch"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["csi.storage.k8s.io"]
    resources: ["csinodeinfos"]
    verbs: ["get", "list", "watch"]

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-csi-provisioner-role
  namespace: rook-ceph
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-provisioner
    namespace: rook-ceph
roleRef:
  kind: ClusterRole
  name: cephfs-external-provisioner-runner
  apiGroup: rbac.authorization.k8s.io

---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: rook-ceph
  name: cephfs-external-provisioner-cfg
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "watch", "list", "delete", "update", "create"]
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list", "create", "delete"]

---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-csi-provisioner-role-cfg
  namespace: rook-ceph
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-provisioner
    namespace: rook-ceph
roleRef:
  kind: Role
  name: cephfs-external-provisioner-cfg
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cephfs-csi-nodeplugin
  namespace: rook-ceph
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-csi-nodeplugin
  namespace: rook-ceph
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "update"]
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["volumeattachments"]
    verbs: ["get", "list", "watch", "update"]

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-csi-nodeplugin
  namespace: rook-ceph
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-nodeplugin
    namespace: rook-ceph
roleRef:
  kind: ClusterRole
  name: cephfs-csi-nodeplugin
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
kind: Service
apiVersion: v1
metadata:
  name: csi-cephfsplugin-attacher
  namespace: rook-ceph
  labels:
    app: csi-cephfsplugin-attacher
spec:
  selector:
    app: csi-cephfsplugin-attacher
  ports:
    - name: dummy
      port: 12345

---
kind: StatefulSet
apiVersion: apps/v1beta1
metadata:
  name: csi-cephfsplugin-attacher
  namespace: rook-ceph
spec:
  serviceName: "csi-cephfsplugin-attacher"
  replicas: 1
  template:
    metadata:
      labels:
        app: csi-cephfsplugin-attacher
    spec:
      serviceAccount: cephfs-csi-attacher
      containers:
        - name: csi-cephfsplugin-attacher
          image: {{.StorageCephAttacherImage}}
          args:
            - "--v=5"
            - "--csi-address=$(ADDRESS)"
          env:
            - name: ADDRESS
              value: /var/lib/kubelet/plugins/cephfs.csi.ceph.com/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /var/lib/kubelet/plugins/cephfs.csi.ceph.com
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/cephfs.csi.ceph.com
            type: DirectoryOrCreate
---
kind: Service
apiVersion: v1
metadata:
  name: csi-cephfsplugin-provisioner
  namespace: rook-ceph
  labels:
    app: csi-cephfsplugin-provisioner
spec:
  selector:
    app: csi-cephfsplugin-provisioner
  ports:
    - name: dummy
      port: 12345

---
kind: StatefulSet
apiVersion: apps/v1beta1
metadata:
  name: csi-cephfsplugin-provisioner
  namespace: rook-ceph
spec:
  serviceName: "csi-cephfsplugin-provisioner"
  replicas: 1
  template:
    metadata:
      labels:
        app: csi-cephfsplugin-provisioner
    spec:
      serviceAccount: cephfs-csi-provisioner
      containers:
        - name: csi-provisioner
          image: {{.StorageCephProvisionerImage}}
          args:
            - "--csi-address=$(ADDRESS)"
            - "--v=5"
          env:
            - name: ADDRESS
              value: unix:///csi/csi-provisioner.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /csi
        - name: csi-cephfsplugin
          securityContext:
            privileged: true
            capabilities:
              add: ["SYS_ADMIN"]
          image: {{.StorageCephFsCSIImage}}
          args:
            - "--nodeid=$(NODE_ID)"
            - "--endpoint=$(CSI_ENDPOINT)"
            - "--v=5"
            - "--drivername=cephfs.csi.ceph.com"
            - "--metadatastorage=k8s_configmap"
          env:
            - name: NODE_ID
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: CSI_ENDPOINT
              value: unix:///csi/csi-provisioner.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /csi
            - name: host-sys
              mountPath: /sys
            - name: lib-modules
              mountPath: /lib/modules
              readOnly: true
            - name: host-dev
              mountPath: /dev
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/cephfs.csi.ceph.com
            type: DirectoryOrCreate
        - name: host-sys
          hostPath:
            path: /sys
        - name: lib-modules
          hostPath:
            path: /lib/modules
        - name: host-dev
          hostPath:
            path: /dev
---
kind: DaemonSet
apiVersion: apps/v1beta2
metadata:
  name: csi-cephfsplugin
  namespace: rook-ceph
spec:
  selector:
    matchLabels:
      app: csi-cephfsplugin
  template:
    metadata:
      labels:
        app: csi-cephfsplugin
    spec:
      serviceAccount: cephfs-csi-nodeplugin
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      containers:
        - name: driver-registrar
          image: {{.StorageCephDriverRegistrarImage}}
          args:
            - "--v=5"
            - "--csi-address=/csi/csi.sock"
            - "--kubelet-registration-path=/var/lib/kubelet/plugins/cephfs.csi.ceph.com/csi.sock"
          lifecycle:
            preStop:
              exec:
                command: [
                  "/bin/sh", "-c",
                  "rm -rf /registration/csi-cephfsplugin \
                  /registration/csi-cephfsplugin-reg.sock"
                ]
          env:
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          volumeMounts:
            - name: plugin-dir
              mountPath: /csi
            - name: registration-dir
              mountPath: /registration
        - name: csi-cephfsplugin
          securityContext:
            privileged: true
            capabilities:
              add: ["SYS_ADMIN"]
            allowPrivilegeEscalation: true
          image: {{.StorageCephFsCSIImage}}
          args:
            - "--nodeid=$(NODE_ID)"
            - "--endpoint=$(CSI_ENDPOINT)"
            - "--v=5"
            - "--drivername=cephfs.csi.ceph.com"
            - "--metadatastorage=k8s_configmap"
            - "--mountcachedir=/mount-cache-dir"
          env:
            - name: NODE_ID
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: CSI_ENDPOINT
              value: unix:///csi/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: mount-cache-dir
              mountPath: /mount-cache-dir
            - name: plugin-dir
              mountPath: /csi
            - name: csi-plugins-dir
              mountPath: /var/lib/kubelet/plugins/kubernetes.io/csi
              mountPropagation: "Bidirectional"
            - name: pods-mount-dir
              mountPath: /var/lib/kubelet/pods
              mountPropagation: "Bidirectional"
            - name: host-sys
              mountPath: /sys
            - name: lib-modules
              mountPath: /lib/modules
              readOnly: true
            - name: host-dev
              mountPath: /dev
      volumes:
        - name: mount-cache-dir
          emptyDir: {}
        - name: plugin-dir
          hostPath:
            path: /var/lib/kubelet/plugins/cephfs.csi.ceph.com/
            type: DirectoryOrCreate
        - name: csi-plugins-dir
          hostPath:
            path: /var/lib/kubelet/plugins/kubernetes.io/csi
            type: DirectoryOrCreate
        - name: registration-dir
          hostPath:
            path: /var/lib/kubelet/plugins_registry/
            type: Directory
        - name: pods-mount-dir
          hostPath:
            path: /var/lib/kubelet/pods
            type: Directory
        - name: host-sys
          hostPath:
            path: /sys
        - name: lib-modules
          hostPath:
            path: /lib/modules
        - name: host-dev
          hostPath:
            path: /dev
---
apiVersion: v1
kind: Secret
metadata:
  name: csi-cephfs-secret
  namespace: default
data:
  userID: {{.CephAdminUserEncode}}
  userKey: {{.CephAdminKeyEncode}}

  adminID: {{.CephAdminUserEncode}}
  adminKey: {{.CephAdminKeyEncode}}
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: cephfs
provisioner: cephfs.csi.ceph.com
parameters:
  monitors: {{.CephClusterMonitors}}
  provisionVolume: "true"
  pool: {{.CephFilesystem}}-data0
  csi.storage.k8s.io/provisioner-secret-name: csi-cephfs-secret
  csi.storage.k8s.io/provisioner-secret-namespace: default
  csi.storage.k8s.io/node-stage-secret-name: csi-cephfs-secret
  csi.storage.k8s.io/node-stage-secret-namespace: default
reclaimPolicy: Delete`
