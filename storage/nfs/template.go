package nfs

const NFSInitTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
---
  {{- if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: storage-agent-nfs-init
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-nfs-init-runner
  namespace: {{.StorageNamespace}}
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-nfs-init-role
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: storage-agent-nfs-init
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: storage-agent-nfs-init-runner
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nfs-init
  namespace: {{.StorageNamespace}}
  labels:
    app: nfs-init
spec:
  selector:
    matchLabels:
      app: nfs-init
  template:
    metadata:
      labels:
        app: nfs-init
    spec:
      serviceAccount: storage-agent-nfs-init
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      containers:
      - name: nfs-init
        image: {{.StorageNFSInitImage}}
        command: ["/init.sh"]
        env:
          - name: MOUNT_PATH
            value: "/host/dev"
          - name: VG_NAME
            value: "nfs"
          - name: LVM_NAME
            value: "data"
          - name: NodeName
            valueFrom:
              fieldRef: 
                fieldPath: spec.nodeName
        securityContext:
          privileged: true
          capabilities:
            add: ["SYS_ADMIN"]
          allowPrivilegeEscalation: true
        volumeMounts:
          - mountPath: /host/dev
            name: host-dev
      volumes:
        - name: host-dev
          hostPath:
            path: /dev`

const NFSStorageTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
---
{{- if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-provisioner
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nfs-provisioner-runner
  namespace: {{.StorageNamespace}}
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "update", "patch"]
  - apiGroups: [""]
    resources: ["services", "endpoints"]
    verbs: ["get"]
  - apiGroups: ["extensions"]
    resources: ["podsecuritypolicies"]
    resourceNames: ["nfs-provisioner"]
    verbs: ["use"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-nfs-provisioner
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: nfs-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: leader-locking-nfs-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: leader-locking-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    namespace: {{.StorageNamespace}}
roleRef:
  kind: Role
  name: leader-locking-nfs-provisioner
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
kind: Service
apiVersion: v1
metadata:
  namespace: {{.StorageNamespace}}
  name: nfs-provisioner
  labels:
    app: nfs-provisioner
spec:
  ports:
    - name: nfs
      port: 2049
    - name: mountd
      port: 20048
    - name: rpcbind
      port: 111
    - name: rpcbind-udp
      port: 111
      protocol: UDP
  selector:
    app: nfs-provisioner
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: nfs-provisioner
spec:
  selector:
    matchLabels:
      app: nfs-provisioner
  replicas: 1
  serviceName: nfs-provisioner
  template:
    metadata:
      labels:
        app: nfs-provisioner
    spec:
      serviceAccount: nfs-provisioner
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      containers:
        - name: nfs-provisioner
          image: {{.StorageNFSProvisionerImage}}
          ports:
            - name: nfs
              containerPort: 2049
            - name: mountd
              containerPort: 20048
            - name: rpcbind
              containerPort: 111
            - name: rpcbind-udp
              containerPort: 111
              protocol: UDP
          securityContext:
            capabilities:
              add:
                - DAC_READ_SEARCH
                - SYS_RESOURCE
          args:
            - "-provisioner=example.com/nfs"
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: SERVICE_NAME
              value: nfs-provisioner
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: nfs-data
              mountPath: /export
      volumes:
      - name: nfs-data
        hostPath:
          path: /var/lib/singlecloud/nfs-export
---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: {{.StorageClassName}}
provisioner: example.com/nfs
reclaimPolicy: Retain
mountOptions:
  - vers=4.1`