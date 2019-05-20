package lvmd

const LVMDTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-storage
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: storage-agent-lvmd
  namespace: kube-storage
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-lvmd-runner
  namespace: kube-storage
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-lvmd-role
  namespace: kube-storage
subjects:
  - kind: ServiceAccount
    name: storage-agent-lvmd
    namespace: kube-storage
roleRef:
  kind: ClusterRole
  name: storage-agent-lvmd-runner
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: storage-agent-lvmd
  namespace: kube-storage
  labels:
    app: storage-agent-lvmd
spec:
  selector:
    matchLabels:
      app: storage-agent-lvmd
  template:
    metadata:
      labels:
        app: storage-agent-lvmd
    spec:
      serviceAccount: storage-agent-lvmd
      hostNetwork: true
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      containers:
      - name: storage-agent-lvmd
        image: {{.StorageLvmdImage}}
        command: ["/lvmd.sh"]
        env:
          - name: MOUNT_PATH
            value: "/host/dev"
          - name: VG_NAME
            value: "k8s"
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
