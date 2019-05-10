package resources

const ChartMuseumTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  CACHE_REDIS_PASSWORD: ""
kind: Secret
metadata:
  labels:
    app: harbor
    chart: harbor
  name: harbor-chartmuseum
  namespace: kube-registry
type: Opaque
---
apiVersion: v1
data:
  ALLOW_OVERWRITE: "true"
  AUTH_ANONYMOUS_GET: "false"
  BASIC_AUTH_USER: chart_controller
  CACHE: redis
  CACHE_REDIS_ADDR: harbor-redis:6379
  CACHE_REDIS_DB: "3"
  CHART_POST_FORM_FIELD_NAME: chart
  CHART_URL: ""
  CONTEXT_PATH: ""
  DEBUG: "true"
  DEPTH: "1"
  DISABLE_API: "false"
  DISABLE_METRICS: "false"
  DISABLE_STATEFILES: "false"
  INDEX_LIMIT: "0"
  LOG_JSON: "true"
  MAX_STORAGE_OBJECTS: "0"
  MAX_UPLOAD_SIZE: "20971520"
  PORT: "9999"
  PROV_POST_FORM_FIELD_NAME: prov
  STORAGE: local
  STORAGE_LOCAL_ROOTDIR: /chart_storage
  TLS_CERT: ""
  TLS_KEY: ""
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: chartmuseum
  name: harbor-chartmuseum
  namespace: kube-registry
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    volume.beta.kubernetes.io/storage-provisioner: csi-lvmplugin
  finalizers:
  - kubernetes.io/pvc-protection
  labels:
    app: harbor
    component: chartmuseum
  name: harbor-chartmuseum
  namespace: kube-registry
spec:
  accessModes:
  - ReadWriteOnce
  dataSource: null
  resources:
    requests:
      storage: 5Gi
  storageClassName: lvm
  volumeMode: Filesystem
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  generation: 1
  labels:
    app: harbor
    component: chartmuseum
  name: harbor-chartmuseum
  namespace: kube-registry
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: chartmuseum
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: 750928f282f9bef8056e91bfc8e42b3a687396aeadf193032a219632dadad001
        checksum/secret: 4e65b2fa8b4ab04c4627cf1b1132453622716fd8e5bcb808adc8a2eadaabfd1a
      labels:
        app: harbor
        component: chartmuseum
    spec:
      containers:
      - env:
        - name: BASIC_AUTH_PASS
          valueFrom:
            secretKeyRef:
              key: secret
              name: harbor-core
        envFrom:
        - configMapRef:
            name: harbor-chartmuseum
        - secretRef:
            name: harbor-chartmuseum
        image: goharbor/chartmuseum-photon:v0.8.1-v1.7.5
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: 9999
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: chartmuseum
        ports:
        - containerPort: 9999
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: 9999
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /chart_storage
          name: chartmuseum-data
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - name: chartmuseum-data
        persistentVolumeClaim:
          claimName: harbor-chartmuseum
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: chartmuseum
  name: harbor-chartmuseum
  namespace: kube-registry
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 9999
  selector:
    app: harbor
    component: chartmuseum
  sessionAffinity: None
  type: ClusterIP
`
