package ceph

const commonTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cephclusters.ceph.rook.io
spec:
  group: ceph.rook.io
  names:
    kind: CephCluster
    listKind: CephClusterList
    plural: cephclusters
    singular: cephcluster
  scope: Namespaced
  version: v1
  validation:
    openAPIV3Schema:
      properties:
        spec:
          properties:
            cephVersion:
              properties:
                allowUnsupported:
                  type: boolean
                image:
                  type: string
                name:
                  pattern: ^(luminous|mimic|nautilus)$
                  type: string
            dashboard:
              properties:
                enabled:
                  type: boolean
                urlPrefix:
                  type: string
                port:
                  type: integer
            dataDirHostPath:
              pattern: ^/(\S+)
              type: string
            mon:
              properties:
                allowMultiplePerNode:
                  type: boolean
                count:
                  maximum: 9
                  minimum: 1
                  type: integer
                preferredCount:
                  maximum: 9
                  minimum: 0
                  type: integer
              required:
              - count
            network:
              properties:
                hostNetwork:
                  type: boolean
            storage:
              properties:
                nodes:
                  items: {}
                  type: array
                useAllDevices: {}
                useAllNodes:
                  type: boolean
          required:
          - mon
  additionalPrinterColumns:
    - name: DataDirHostPath
      type: string
      description: Directory used on the K8s nodes
      JSONPath: .spec.dataDirHostPath
    - name: MonCount
      type: string
      description: Number of MONs
      JSONPath: .spec.mon.count
    - name: Age
      type: date
      JSONPath: .metadata.creationTimestamp
    - name: State
      type: string
      description: Current State
      JSONPath: .status.state
    - name: Health
      type: string
      description: Ceph Health
      JSONPath: .status.ceph.health
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cephfilesystems.ceph.rook.io
spec:
  group: ceph.rook.io
  names:
    kind: CephFilesystem
    listKind: CephFilesystemList
    plural: cephfilesystems
    singular: cephfilesystem
  scope: Namespaced
  version: v1
  additionalPrinterColumns:
    - name: MdsCount
      type: string
      description: Number of MDSs
      JSONPath: .spec.metadataServer.activeCount
    - name: Age
      type: date
      JSONPath: .metadata.creationTimestamp
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cephnfses.ceph.rook.io
spec:
  group: ceph.rook.io
  names:
    kind: CephNFS
    listKind: CephNFSList
    plural: cephnfses
    singular: cephnfs
    shortNames:
    - nfs
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cephobjectstores.ceph.rook.io
spec:
  group: ceph.rook.io
  names:
    kind: CephObjectStore
    listKind: CephObjectStoreList
    plural: cephobjectstores
    singular: cephobjectstore
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cephobjectstoreusers.ceph.rook.io
spec:
  group: ceph.rook.io
  names:
    kind: CephObjectStoreUser
    listKind: CephObjectStoreUserList
    plural: cephobjectstoreusers
    singular: cephobjectstoreuser
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cephblockpools.ceph.rook.io
spec:
  group: ceph.rook.io
  names:
    kind: CephBlockPool
    listKind: CephBlockPoolList
    plural: cephblockpools
    singular: cephblockpool
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: volumes.rook.io
spec:
  group: rook.io
  names:
    kind: Volume
    listKind: VolumeList
    plural: volumes
    singular: volume
    shortNames:
    - rv
  scope: Namespaced
  version: v1alpha2
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: rook-ceph-cluster-mgmt
  labels:
    operator: rook
    storage-backend: ceph
aggregationRule:
  clusterRoleSelectors:
  - matchLabels:
      rbac.ceph.rook.io/aggregate-to-rook-ceph-cluster-mgmt: "true"
rules: []
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: rook-ceph-cluster-mgmt-rules
  labels:
    operator: rook
    storage-backend: ceph
    rbac.ceph.rook.io/aggregate-to-rook-ceph-cluster-mgmt: "true"
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  - pods
  - pods/log
  - services
  - configmaps
  verbs:
  - get
  - list
  - watch
  - patch
  - create
  - update
  - delete
- apiGroups:
  - apps
  resources:
  - deployments
  - daemonsets
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: rook-ceph-system
  namespace: {{.StorageNamespace}}
  labels:
    operator: rook
    storage-backend: ceph
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - configmaps
  - services
  verbs:
  - get
  - list
  - watch
  - patch
  - create
  - update
  - delete
- apiGroups:
  - apps
  resources:
  - daemonsets
  - statefulsets
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: rook-ceph-global
  labels:
    operator: rook
    storage-backend: ceph
aggregationRule:
  clusterRoleSelectors:
  - matchLabels:
      rbac.ceph.rook.io/aggregate-to-rook-ceph-global: "true"
rules: []
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: rook-ceph-global-rules
  labels:
    operator: rook
    storage-backend: ceph
    rbac.ceph.rook.io/aggregate-to-rook-ceph-global: "true"
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - nodes
  - nodes/proxy
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - events
  - persistentvolumes
  - persistentvolumeclaims
  - endpoints
  verbs:
  - get
  - list
  - watch
  - patch
  - create
  - update
  - delete
- apiGroups:
  - storage.k8s.io
  resources:
  - storageclasses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - batch
  resources:
  - jobs
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
- apiGroups:
  - ceph.rook.io
  resources:
  - "*"
  verbs:
  - "*"
- apiGroups:
  - rook.io
  resources:
  - "*"
  verbs:
  - "*"
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr-cluster
  labels:
    operator: rook
    storage-backend: ceph
aggregationRule:
  clusterRoleSelectors:
  - matchLabels:
      rbac.ceph.rook.io/aggregate-to-rook-ceph-mgr-cluster: "true"
rules: []
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr-cluster-rules
  labels:
    operator: rook
    storage-backend: ceph
    rbac.ceph.rook.io/aggregate-to-rook-ceph-mgr-cluster: "true"
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - nodes
  - nodes/proxy
  verbs:
  - get
  - list
  - watch
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rook-ceph-system
  namespace: {{.StorageNamespace}}
  labels:
    operator: rook
    storage-backend: ceph
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-system
  namespace: {{.StorageNamespace}}
  labels:
    operator: rook
    storage-backend: ceph
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: rook-ceph-system
subjects:
- kind: ServiceAccount
  name: rook-ceph-system
  namespace: {{.StorageNamespace}}
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-global
  namespace: {{.StorageNamespace}}
  labels:
    operator: rook
    storage-backend: ceph
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: rook-ceph-global
subjects:
- kind: ServiceAccount
  name: rook-ceph-system
  namespace: {{.StorageNamespace}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rook-ceph-osd
  namespace: {{.StorageNamespace}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rook-ceph-mgr
  namespace: {{.StorageNamespace}}
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-osd
  namespace: {{.StorageNamespace}}
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: [ "get", "list", "watch", "create", "update", "delete" ]
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr-system
  namespace: {{.StorageNamespace}}
aggregationRule:
  clusterRoleSelectors:
  - matchLabels:
      rbac.ceph.rook.io/aggregate-to-rook-ceph-mgr-system: "true"
rules: []
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr-system-rules
  namespace: {{.StorageNamespace}}
  labels:
      rbac.ceph.rook.io/aggregate-to-rook-ceph-mgr-system: "true"
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr
  namespace: {{.StorageNamespace}}
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - services
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - batch
  resources:
  - jobs
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
- apiGroups:
  - ceph.rook.io
  resources:
  - "*"
  verbs:
  - "*"
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-cluster-mgmt
  namespace: {{.StorageNamespace}}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: rook-ceph-cluster-mgmt
subjects:
- kind: ServiceAccount
  name: rook-ceph-system
  namespace: {{.StorageNamespace}}
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-osd
  namespace: {{.StorageNamespace}}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: rook-ceph-osd
subjects:
- kind: ServiceAccount
  name: rook-ceph-osd
  namespace: {{.StorageNamespace}}
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr
  namespace: {{.StorageNamespace}}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: rook-ceph-mgr
subjects:
- kind: ServiceAccount
  name: rook-ceph-mgr
  namespace: {{.StorageNamespace}}
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr-system
  namespace: {{.StorageNamespace}}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: rook-ceph-mgr-system
subjects:
- kind: ServiceAccount
  name: rook-ceph-mgr
  namespace: {{.StorageNamespace}}
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: rook-ceph-mgr-cluster
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: rook-ceph-mgr-cluster
subjects:
- kind: ServiceAccount
  name: rook-ceph-mgr
  namespace: {{.StorageNamespace}}
{{- end}}`

const clusterTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rook-ceph-operator
  namespace: {{.StorageNamespace}}
  labels:
    operator: rook
    storage-backend: ceph
spec:
  selector:
    matchLabels:
      app: rook-ceph-operator
  replicas: 1
  template:
    metadata:
      labels:
        app: rook-ceph-operator
    spec:
      serviceAccountName: rook-ceph-system
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      containers:
      - name: rook-ceph-operator
        image: {{.StorageCephOperatorImage}}
        args: ["ceph", "operator"]
        volumeMounts:
        - mountPath: /var/lib/rook
          name: rook-config
        - mountPath: /etc/ceph
          name: default-config-dir
        env:
        - name: ROOK_CURRENT_NAMESPACE_ONLY
          value: "true"
        - name: ROOK_ALLOW_MULTIPLE_FILESYSTEMS
          value: "false"
        - name: ROOK_LOG_LEVEL
          value: "INFO"
        - name: ROOK_CEPH_STATUS_CHECK_INTERVAL
          value: "60s"
        - name: ROOK_MON_HEALTHCHECK_INTERVAL
          value: "45s"
        - name: ROOK_MON_OUT_TIMEOUT
          value: "600s"
        - name: ROOK_DISCOVER_DEVICES_INTERVAL
          value: "60m"
        - name: ROOK_HOSTPATH_REQUIRES_PRIVILEGED
          value: "false"
        - name: ROOK_ENABLE_SELINUX_RELABELING
          value: "true"
        - name: ROOK_ENABLE_FSGROUP
          value: "true"
        - name: ROOK_DISABLE_DEVICE_HOTPLUG
          value: "false"
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      volumes:
      - name: rook-config
        emptyDir: {}
      - name: default-config-dir
        emptyDir: {}
---
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: {{.StorageNamespace}}
spec:
  cephVersion:
    image: {{.StorageCephClusterImage}}
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
            - key: {{.LabelKey}}
              operator: In
              values:
              - "{{.LabelValue}}"
  annotations:
  resources:
  storage:
    useAllNodes: false
    useAllDevices: false
    deviceFilter:
    location:
    config:
    nodes:{{range .CephList}}
    - name: {{.Host}}
      devices:{{range .Devs}}
      - name: {{.Dev}}{{end}}{{end}}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rook-ceph-tools
  namespace: {{.StorageNamespace}}
  labels:
    app: rook-ceph-tools
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rook-ceph-tools
  template:
    metadata:
      labels:
        app: rook-ceph-tools
    spec:
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      dnsPolicy: ClusterFirstWithHostNet
      containers:
      - name: rook-ceph-tools
        image: {{.StorageCephToolsImage}}
        command: ["/tini"]
        args: ["-g", "--", "/usr/local/bin/toolbox.sh"]
        imagePullPolicy: IfNotPresent
        env:
          - name: ROOK_ADMIN_SECRET
            valueFrom:
              secretKeyRef:
                name: rook-ceph-mon
                key: admin-secret
        securityContext:
          privileged: true
        volumeMounts:
          - mountPath: /dev
            name: dev
          - mountPath: /sys/bus
            name: sysbus
          - mountPath: /lib/modules
            name: libmodules
          - name: mon-endpoint-volume
            mountPath: /etc/rook
      volumes:
        - name: dev
          hostPath:
            path: /dev
        - name: sysbus
          hostPath:
            path: /sys/bus
        - name: libmodules
          hostPath:
            path: /lib/modules
        - name: mon-endpoint-volume
          configMap:
            name: rook-ceph-mon-endpoints
            items:
            - key: data
              path: mon-endpoints`

const filesystemTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
---
apiVersion: ceph.rook.io/v1
kind: CephFilesystem
metadata:
  name: {{.CephFilesystem}}
  namespace: {{.StorageNamespace}}
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

const fscsiTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
---
{{- if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cephfs-csi-attacher
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-external-attacher-runner
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-attacher
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: cephfs-external-attacher-runner
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cephfs-csi-provisioner
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-external-provisioner-runner
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-provisioner
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: cephfs-external-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-provisioner
    namespace: {{.StorageNamespace}}
roleRef:
  kind: Role
  name: cephfs-external-provisioner-cfg
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cephfs-csi-nodeplugin
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cephfs-csi-nodeplugin
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: cephfs-csi-nodeplugin
    namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
spec:
  serviceName: "csi-cephfsplugin-attacher"
  replicas: 1
  template:
    metadata:
      labels:
        app: csi-cephfsplugin-attacher
    spec:
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
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
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
spec:
  serviceName: "csi-cephfsplugin-provisioner"
  replicas: 1
  template:
    metadata:
      labels:
        app: csi-cephfsplugin-provisioner
    spec:
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
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
  namespace: {{.StorageNamespace}}
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
  namespace: {{.StorageNamespace}}
data:
  userID: {{.CephAdminUserEncode}}
  userKey: {{.CephAdminKeyEncode}}

  adminID: {{.CephAdminUserEncode}}
  adminKey: {{.CephAdminKeyEncode}}
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: {{.StorageClassName}}
provisioner: cephfs.csi.ceph.com
parameters:
  monitors: {{.CephClusterMonitors}}
  provisionVolume: "true"
  pool: {{.CephFilesystem}}-data0
  csi.storage.k8s.io/provisioner-secret-name: csi-cephfs-secret
  csi.storage.k8s.io/provisioner-secret-namespace: {{.StorageNamespace}}
  csi.storage.k8s.io/node-stage-secret-name: csi-cephfs-secret
  csi.storage.k8s.io/node-stage-secret-namespace: {{.StorageNamespace}}
reclaimPolicy: Delete`
