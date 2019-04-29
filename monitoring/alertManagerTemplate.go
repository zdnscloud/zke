package monitoring

const AlertManagerTemplate = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus-alertmanager
  namespace: kube-monitoring
---
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app: prometheus
    component: alertmanager
    release: prometheus
  name: prometheus-alertmanager
  namespace: kube-monitoring
data:
  alertmanager.yml: |
    global: null
    receivers:
    - name: default-receiver
    route:
      group_interval: 5m
      group_wait: 10s
      receiver: default-receiver
      repeat_interval: 3h
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: prometheus-alertmanager
  namespace: kube-monitoring
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
      component: alertmanager
      release: prometheus
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: prometheus
        component: alertmanager
        release: prometheus
    spec:
      containers:
      - args:
        - --config.file=/etc/config/alertmanager.yml
        - --storage.path=/data
        - --web.external-url=/
        image: {{ .PrometheusAlertManagerImage }}
        imagePullPolicy: IfNotPresent
        name: prometheus-alertmanager
        ports:
        - containerPort: 9093
          protocol: TCP
        readinessProbe:
          failureThreshold: 10
          httpGet:
            path: /#/status
            port: 9093
            scheme: HTTP
          initialDelaySeconds: 30
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 30
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
        - mountPath: /data
          name: storage-volume
      - args:
        - --volume-dir=/etc/config
        - --webhook-url=http://localhost:9093/-/reload
        image: {{ .PrometheusConfigMapReloaderImage }}
        imagePullPolicy: IfNotPresent
        name: prometheus-alertmanager-configmap-reload
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
          readOnly: true
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: prometheus-alertmanager
      serviceAccountName: prometheus-alertmanager
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          name: prometheus-alertmanager
        name: config-volume
      - emptyDir: {}
        name: storage-volume
status: {}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: prometheus
    component: alertmanager
    release: prometheus
  name: prometheus-alertmanager
  namespace: kube-monitoring
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 9093
  selector:
    app: prometheus
    component: alertmanager
    release: prometheus
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: prometheus
    component: alertmanager
    release: prometheus
  name: prometheus-alertmanager
  namespace: kube-monitoring
spec:
  rules:
  - host: {{ .PrometheusAlertManagerIngressEndpoint }}
    http:
      paths:
      - backend:
          serviceName: prometheus-alertmanager
          servicePort: 80
status:
  loadBalancer: {}
`
