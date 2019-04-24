package templates

const GrafanaTemplate = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus-grafana
  namespace: zcloud
---
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app: prometheus-grafana
    release: prometheus
  name: psp-prometheus-grafana
rules:
- apiGroups:
  - extensions
  resourceNames:
  - prometheus-grafana
  resources:
  - podsecuritypolicies
  verbs:
  - use
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app: prometheus-grafana
    chart: grafana-0.0.31
    release: prometheus
  name: psp-prometheus-grafana
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: psp-prometheus-grafana
subjects:
- kind: ServiceAccount
  name: prometheus-grafana
  namespace: zcloud
{{- end}}
---
apiVersion: v1
data:
  password: emRuc2Nsb3Vk
  user: YWRtaW4=
kind: Secret
metadata:
  labels:
    app: prometheus-grafana
    release: prometheus
  name: prometheus-grafana
  namespace: zcloud
type: Opaque
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: prometheus-grafana
  namespace: zcloud
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 2
  selector:
    matchLabels:
      app: prometheus-grafana
      release: prometheus
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: prometheus-grafana
        release: prometheus
    spec:
      containers:
      - env:
        - name: GF_AUTH_BASIC_ENABLED
          value: "true"
        - name: GF_AUTH_ANONYMOUS_ENABLED
          value: "true"
        - name: GF_SECURITY_ADMIN_USER
          valueFrom:
            secretKeyRef:
              key: user
              name: prometheus-grafana
        - name: GF_SECURITY_ADMIN_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: prometheus-grafana
        - name: GF_USERS_DEFAULT_THEME
          value: "light"
        image: {{ .GrafanaImage }}
        imagePullPolicy: IfNotPresent
        name: grafana
        ports:
        - containerPort: 3000
          name: web
          protocol: TCP
        readinessProbe:
          failureThreshold: 10
          httpGet:
            path: /api/health
            port: 3000
            scheme: HTTP
          periodSeconds: 1
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          limits:
            cpu: 200m
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 100Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/lib/grafana
          name: grafana-storage
      - args:
        - --watch-dir=/var/grafana-dashboards
        - --watch-dir=/var/grafana-dashboards/k8s
        - --watch-dir=/var/grafana-dashboards/k8s-resources
        - --grafana-url=http://127.0.0.1:3000
        env:
        - name: GRAFANA_USER
          valueFrom:
            secretKeyRef:
              key: user
              name: prometheus-grafana
        - name: GRAFANA_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: prometheus-grafana
        image: {{ .GrafanaWatcherImage }}
        imagePullPolicy: IfNotPresent
        name: grafana-watcher
        resources:
          limits:
            cpu: 100m
            memory: 32Mi
          requests:
            cpu: 50m
            memory: 16Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/grafana-dashboards/k8s
          name: grafana-dashboards-k8s
        - mountPath: /var/grafana-dashboards/k8s-resources
          name: grafana-dashboards-resources
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: prometheus-grafana
      serviceAccountName: prometheus-grafana
      terminationGracePeriodSeconds: 30
      volumes:
      - emptyDir: {}
        name: grafana-storage
      - configMap:
          defaultMode: 420
          name: prometheus-grafana
        name: grafana-dashboards-k8s
      - configMap:
          defaultMode: 420
          name: prometheus-grafana-resources
        name: grafana-dashboards-resources
status: {}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: prometheus-grafana
    release: prometheus
  name: prometheus-grafana
  namespace: zcloud
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 3000
  selector:
    app: prometheus-grafana
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: grafana
    release: prometheus
  name: prometheus-grafana
  namespace: zcloud
spec:
  rules:
  - host: {{ .GrafanaIngressEndpoint }}
    http:
      paths:
      - backend:
          serviceName: prometheus-grafana
          servicePort: 80
status:
  loadBalancer: {}
`
