package nodeagent

const NodeAgentTemplate = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-agent
  namespace: zcloud
  labels:
    app: node-agent
spec:
  selector:
    matchLabels:
      app: node-agent
  template:
    metadata:
      labels:
        app: node-agent
    spec:
      tolerations:
      - operator: Exists
        effect: NoSchedule
      - operator: Exists
        effect: NoExecute
      hostNetwork: true
      containers:
      - name: node-agent
        image: {{.Image}}
        command: ["/bin/sh", "-c","/node-agent -listen :8899"]
        securityContext:
          privileged: true
        volumeMounts:
          - mountPath: /var/lib/kubelet
            name: kubelet
      volumes:
        - name: kubelet
          hostPath:
            path: /var/lib/kubelet`
