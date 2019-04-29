package monitoring

const StateMetricsTemplate = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus-kube-state-metrics
  namespace: kube-monitoring
---
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app: prometheus
    component: kube-state-metrics
    release: prometheus
  name: prometheus-kube-state-metrics
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  - nodes
  - persistentvolumeclaims
  - pods
  - services
  - resourcequotas
  - replicationcontrollers
  - limitranges
  - persistentvolumeclaims
  - persistentvolumes
  - endpoints
  - secrets
  - configmaps
  verbs:
  - list
  - watch
- apiGroups:
  - extensions
  resources:
  - daemonsets
  - deployments
  - replicasets
  verbs:
  - list
  - watch
- apiGroups:
  - apps
  resources:
  - statefulsets
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - batch
  resources:
  - cronjobs
  - jobs
  verbs:
  - list
  - watch
- apiGroups:
  - autoscaling
  resources:
  - horizontalpodautoscalers
  verbs:
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app: prometheus
    component: kube-state-metrics
    release: prometheus
  name: prometheus-kube-state-metrics
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus-kube-state-metrics
subjects:
- kind: ServiceAccount
  name: prometheus-kube-state-metrics
  namespace: kube-monitoring
{{- end}}
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: prometheus
    component: kube-state-metrics
    release: prometheus
  name: prometheus-kube-state-metrics
  namespace: kube-monitoring
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
      component: kube-state-metrics
      release: prometheus
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: prometheus
        component: kube-state-metrics
        release: prometheus
    spec:
      containers:
      - image: {{ .KubeStateMetricsImage }}
        imagePullPolicy: IfNotPresent
        name: prometheus-kube-state-metrics
        ports:
        - containerPort: 8080
          name: metrics
          protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: prometheus-kube-state-metrics
      serviceAccountName: prometheus-kube-state-metrics
      terminationGracePeriodSeconds: 30
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/scrape: "true"
  labels:
    app: prometheus
    component: kube-state-metrics
    release: prometheus
  name: prometheus-kube-state-metrics
  namespace: kube-monitoring
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: prometheus
    component: kube-state-metrics
    release: prometheus
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
`
