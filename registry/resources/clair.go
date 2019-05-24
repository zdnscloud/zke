package resources

const ClairTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  config.yaml: |
    clair:
      database:
        type: pgsql
        options:
          source: "postgres://postgres:changeit@harbor-database:5432/postgres?sslmode=disable"
          # Number of elements kept in the cache
          # Values unlikely to change (e.g. namespaces) are cached in order to save prevent needless roundtrips to the database.
          cachesize: 16384

      api:
        # API server port
        port: 6060
        healthport: 6061

        # Deadline before an API request will respond with a 503
        timeout: 300s
      updater:
        interval: 12h

      notifier:
        attempts: 3
        renotifyinterval: 2h
        http:
          endpoint: "http://harbor-core/service/notifications/clair"
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: clair
  name: harbor-clair
  namespace: {{ .DeployNamespace }}
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  generation: 1
  labels:
    app: harbor
    component: clair
  name: harbor-clair
  namespace: {{ .DeployNamespace }}
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: clair
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: 6276aba6568a256b0ae8faf4722952f4528fb91d5a5287ff165c07d3caf8297c
      creationTimestamp: null
      labels:
        app: harbor
        component: clair
    spec:
      containers:
      - args:
        - -log-level
        - debug
        env:
        - name: NO_PROXY
          value: harbor-registry,harbor-core
        image: {{ .ClairImage}}
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: 6061
            scheme: HTTP
          initialDelaySeconds: 30
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: clair
        ports:
        - containerPort: 6060
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: 6061
            scheme: HTTP
          initialDelaySeconds: 30
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/clair/config.yaml
          name: clair-config
          subPath: config.yaml
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          items:
          - key: config.yaml
            path: config.yaml
          name: harbor-clair
        name: clair-config
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    chart: clair
  name: harbor-clair
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - port: 6060
    protocol: TCP
    targetPort: 6060
  selector:
    app: harbor
    component: clair
  sessionAffinity: None
  type: ClusterIP
`
