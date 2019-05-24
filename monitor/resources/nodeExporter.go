package resources

const NodeExporterTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-monitor
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus-node-exporter
  namespace: kube-monitor
---
apiVersion: apps/v1beta2
kind: DaemonSet
metadata:
  labels:
    app: prometheus
    component: node-exporter
    release: prometheus
  name: prometheus-node-exporter
  namespace: kube-monitor
spec:
  selector:
    matchLabels:
      app: prometheus
      component: node-exporter
      release: prometheus
  template:
    metadata:
      labels:
        app: prometheus
        component: node-exporter
        release: prometheus
    spec:
      containers:
      - args:
        - --path.procfs=/host/proc
        - --path.sysfs=/host/sys
        image: {{ .PrometheusNodeExporterImage }}
        imagePullPolicy: IfNotPresent
        name: prometheus-node-exporter
        ports:
        - containerPort: 9100
          hostPort: 9100
          name: metrics
          protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /host/proc
          name: proc
          readOnly: true
        - mountPath: /host/sys
          name: sys
          readOnly: true
      dnsPolicy: ClusterFirst
      hostNetwork: true
      hostPID: true
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: prometheus-node-exporter
      serviceAccountName: prometheus-node-exporter
      terminationGracePeriodSeconds: 30
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/controlplane
        operator: Exists
      - effect: NoExecute
        key: node-role.kubernetes.io/etcd
        operator: Exists
      - effect: NoExecute
        key: node.kubernetes.io/not-ready
        operator: Exists
      - effect: NoExecute
        key: node.kubernetes.io/unreachable
        operator: Exists
      volumes:
      - hostPath:
          path: /proc
          type: ""
        name: proc
      - hostPath:
          path: /sys
          type: ""
        name: sys
  updateStrategy:
    type: OnDelete
status:
  currentNumberScheduled: 0
  desiredNumberScheduled: 0
  numberMisscheduled: 0
  numberReady: 0
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/scrape: "true"
  labels:
    app: prometheus
    component: node-exporter
    release: prometheus
  name: prometheus-node-exporter
  namespace: kube-monitor
spec:
  ports:
  - name: metrics
    port: 9100
    protocol: TCP
    targetPort: 9100
  selector:
    app: prometheus
    component: node-exporter
    release: prometheus
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
`
