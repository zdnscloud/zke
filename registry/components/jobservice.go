package components

const JobserviceTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  config.yml: |
    protocol: "http"
    port: 8080
    worker_pool:
      workers: 10
      backend: "redis"
      redis_pool:
        redis_url: "harbor-redis:6379/1"
        namespace: "harbor_job_service_namespace"
    job_loggers:
      - name: "FILE"
        level: DEBUG
        settings: # Customized settings of logger
          base_dir: "/var/log/jobs"
        sweeper:
          duration: 14 #days
          settings: # Customized settings of sweeper
            work_dir: "/var/log/jobs"
    #Loggers for the job service
    loggers:
      - name: "STD_OUTPUT"
        level: DEBUG
    admin_server: "http://harbor-adminserver"
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: jobservice
  name: harbor-jobservice
  namespace: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  secret: TlFieVNXdDlWb0pzdFVjdQ==
kind: Secret
metadata:
  labels:
    app: harbor
    component: jobservice
  name: harbor-jobservice
  namespace: {{ .DeployNamespace }}
type: Opaque
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
    component: jobservice
  name: harbor-jobservice
  namespace: {{ .DeployNamespace }}
spec:
  accessModes:
  - ReadWriteOnce
  dataSource: null
  resources:
    requests:
      storage: {{ .JobserviceDiskCapacity}}
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
    component: jobservice
  name: harbor-jobservice
  namespace: {{ .DeployNamespace }}
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: jobservice
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: 65a06469e8f69f18602f9d9083d4706b2e9c84b0ee71cec263f6068be30ae03d
        checksum/secret: a5b5f87462f05fd9b75a7c285e7e80b99f6a2fcdfcc96d23b17d883dad74dd6a
      creationTimestamp: null
      labels:
        app: harbor
        component: jobservice
    spec:
      containers:
      - env:
        - name: CORE_SECRET
          valueFrom:
            secretKeyRef:
              key: secret
              name: harbor-core
        - name: JOBSERVICE_SECRET
          valueFrom:
            secretKeyRef:
              key: secret
              name: harbor-jobservice
        - name: ADMINSERVER_URL
          value: http://harbor-adminserver
        - name: REGISTRY_CONTROLLER_URL
          value: http://harbor-registry:8080
        - name: LOG_LEVEL
          value: debug
        image: {{ .JobserviceImage}}
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/v1/stats
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 20
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: jobservice
        ports:
        - containerPort: 8080
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/v1/stats
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 20
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/jobservice/config.yml
          name: jobservice-config
          subPath: config.yml
        - mountPath: /var/log/jobs
          name: job-logs
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          name: harbor-jobservice
        name: jobservice-config
      - name: job-logs
        persistentVolumeClaim:
          claimName: harbor-jobservice
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: jobservice
  name: harbor-jobservice
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: harbor
    component: jobservice
  sessionAffinity: None
  type: ClusterIP
`
