package resources

const RedisTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  generation: 1
  labels:
    app: harbor
    component: redis
  name: harbor-redis
  namespace: kube-registry
spec:
  podManagementPolicy: OrderedReady
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: redis
  serviceName: harbor-redis
  template:
    metadata:
      labels:
        app: harbor
        component: redis
    spec:
      containers:
      - image: goharbor/redis-photon:v1.7.5
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          tcpSocket:
            port: 6379
          timeoutSeconds: 1
        name: redis
        readinessProbe:
          failureThreshold: 3
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          tcpSocket:
            port: 6379
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/lib/redis
          name: data
      dnsPolicy: ClusterFirst
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
        component: redis
      name: data
    spec:
      accessModes:
      - ReadWriteOnce
      dataSource: null
      resources:
        requests:
          storage: 1Gi
      volumeMode: Filesystem
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: redis
  name: harbor-redis
  namespace: kube-registry
spec:
  ports:
  - port: 6379
    protocol: TCP
    targetPort: 6379
  selector:
    app: harbor
    component: redis
  sessionAffinity: None
  type: ClusterIP
`
