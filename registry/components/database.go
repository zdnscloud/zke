package components

const DatabaseTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  POSTGRES_PASSWORD: Y2hhbmdlaXQ=
kind: Secret
metadata:
  labels:
    app: harbor
    component: database
  name: harbor-database
  namespace: {{ .DeployNamespace }}
type: Opaque
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  generation: 1
  labels:
    app: harbor
    component: database
  name: harbor-database
  namespace: {{ .DeployNamespace }}
spec:
  podManagementPolicy: OrderedReady
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: database
  serviceName: harbor-database
  template:
    metadata:
      annotations:
        checksum/secret: 8accf53dfbf601db7e1b0fb79823d057c19b43014181ecc32c20f049f948b187
      labels:
        app: harbor
        component: database
    spec:
      containers:
      - envFrom:
        - secretRef:
            name: harbor-database
        image: {{ .DatabaseImage}}
        imagePullPolicy: IfNotPresent
        livenessProbe:
          exec:
            command:
            - /docker-healthcheck.sh
          failureThreshold: 3
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: database
        readinessProbe:
          exec:
            command:
            - /docker-healthcheck.sh
          failureThreshold: 3
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/lib/postgresql/data
          name: database-data
      dnsPolicy: ClusterFirst
      initContainers:
      - command:
        - rm
        - -Rf
        - /var/lib/postgresql/data/lost+found
        image: {{ .DatabaseImage}}
        imagePullPolicy: IfNotPresent
        name: remove-lost-found
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/lib/postgresql/data
          name: database-data
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
  updateStrategy:
    rollingUpdate:
      partition: 0
    type: RollingUpdate
  volumeClaimTemplates:
  - metadata:
      labels:
        app: harbor
        chart: database
      name: database-data
    spec:
      accessModes:
      - ReadWriteOnce
      dataSource: null
      resources:
        requests:
          storage: {{ .DatabaseDiskCapacity}}
      volumeMode: Filesystem
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: database
  name: harbor-database
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - port: 5432
    protocol: TCP
    targetPort: 5432
  selector:
    app: harbor
    component: database
  sessionAffinity: None
  type: ClusterIP
`
