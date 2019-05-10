package resources

const AdminServerTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  CLAIR_DB_PASSWORD: Y2hhbmdlaXQ=
  HARBOR_ADMIN_PASSWORD: SGFyYm9yMTIzNDU=
  POSTGRESQL_PASSWORD: Y2hhbmdlaXQ=
  secretKey: bm90LWEtc2VjdXJlLWtleQ==
kind: Secret
metadata:
  labels:
    app: harbor
    component: adminserver
  name: harbor-adminserver
  namespace: kube-registry
type: Opaque
---
apiVersion: v1
data:
  ADMIRAL_URL: NA
  AUTH_MODE: db_auth
  CFG_EXPIRATION: "5"
  CHART_REPOSITORY_URL: http://harbor-chartmuseum
  CLAIR_DB: postgres
  CLAIR_DB_HOST: harbor-database
  CLAIR_DB_PORT: "5432"
  CLAIR_DB_SSLMODE: disable
  CLAIR_DB_USERNAME: postgres
  CLAIR_URL: http://harbor-clair:6060
  CORE_URL: http://harbor-core
  DATABASE_TYPE: postgresql
  EMAIL_FROM: admin <sample_admin@mydomain.com>
  EMAIL_HOST: smtp.mydomain.com
  EMAIL_PORT: "25"
  EXT_ENDPOINT: https://harbor.cluster.w
  IMAGE_STORE_PATH: /
  JOBSERVICE_URL: http://harbor-jobservice
  LOG_LEVEL: debug
  NOTARY_URL: http://harbor-notary-server:4443
  POSTGRESQL_DATABASE: registry
  POSTGRESQL_HOST: harbor-database
  POSTGRESQL_PORT: "5432"
  POSTGRESQL_SSLMODE: disable
  POSTGRESQL_USERNAME: postgres
  PROJECT_CREATION_RESTRICTION: everyone
  REGISTRY_STORAGE_PROVIDER_NAME: filesystem
  REGISTRY_URL: http://harbor-registry:5000
  RESET: "false"
  SELF_REGISTRATION: "on"
  TOKEN_EXPIRATION: "30"
  TOKEN_SERVICE_URL: http://harbor-core/service/token
  UAA_CLIENTID: ""
  UAA_CLIENTSECRET: ""
  UAA_ENDPOINT: ""
  UAA_VERIFY_CERT: "True"
  WITH_CHARTMUSEUM: "true"
  WITH_CLAIR: "true"
  WITH_NOTARY: "true"
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: adminserver
  name: harbor-adminserver
  namespace: kube-registry
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  generation: 1
  labels:
    app: harbor
    component: adminserver
  name: harbor-adminserver
  namespace: kube-registry
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: adminserver
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: 049f3f89a0f736e612cc0dc30d93be8883ec30ab030e3b4d9fa0c03eaf052477
        checksum/secret: 1ec1bce946a4884b14334d24464b9fbba2652bfabf3da1dcca8901d1a86313ba
        checksum/secret-core: 4894968bd85028eeb6a94bf8e4b0acd1b01e5300a2955a4c70b8f2c903e0f2ca
        checksum/secret-jobservice: 39d9f61a0c4e05110951162a49f703b55925977e92a7a088e07735975ed20619
      labels:
        app: harbor
        component: adminserver
    spec:
      containers:
      - env:
        - name: PORT
          value: "8080"
        - name: JSON_CFG_STORE_PATH
          value: /etc/adminserver/config/config.json
        - name: KEY_PATH
          value: /etc/adminserver/key
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
        envFrom:
        - configMapRef:
            name: harbor-adminserver
        - secretRef:
            name: harbor-adminserver
        image: goharbor/harbor-adminserver:v1.7.5
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/ping
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: adminserver
        ports:
        - containerPort: 8080
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/ping
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/adminserver/key
          name: adminserver-key
          subPath: key
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - name: adminserver-key
        secret:
          defaultMode: 420
          items:
          - key: secretKey
            path: key
          secretName: harbor-adminserver
---
apiVersion: v1
kind: Service
metadata:
  name: harbor-adminserver
  namespace: kube-registry
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: harbor
    component: adminserver
  sessionAffinity: None
  type: ClusterIP
`
