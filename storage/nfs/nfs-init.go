package nfs

const NFSInitTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-storage
---
  {{- if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: storage-agent-nfs-init
  namespace: kube-storage
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-nfs-init-runner
  namespace: kube-storage
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-nfs-init-role
  namespace: kube-storage
subjects:
  - kind: ServiceAccount
    name: storage-agent-nfs-init
    namespace: kube-storage
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
  namespace: kube-storage
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
