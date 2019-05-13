package lvmd

const LVMDTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-storage
{{range .LVMList}}
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: storage-lvmd-{{.Host}}
  namespace: kube-storage
spec:
  selector:
    matchLabels:
      app: storage-lvmd-{{.Host}}
  template:
    metadata:
      labels:
        app: storage-lvmd-{{.Host}}
    spec:
      nodeName: "{{.Host}}"
      hostNetwork: true
      containers:
      - name: lvmd
        image: {{$.StorageLvmdImage}}
        command: ["/lvmd.sh"]
        env:
          - name: MOUNT_PATH
            value: "/host/dev"
          - name: VG_NAME
            value: "k8s"
          - name: DEVICE
            value: "{{.Devs}}"
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
            path: /dev
{{end}}
---`
