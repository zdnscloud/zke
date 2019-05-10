package resources

const NotarySignerTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  generation: 1
  labels:
    app: harbor
    component: notary-signer
  name: harbor-notary-signer
  namespace: kube-registry
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: notary-signer
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: d67064d07e39447f3fba0902d40d76df20dc9dbbb3d55accfc0e669b751b356b
      labels:
        app: harbor
        component: notary-signer
    spec:
      containers:
      - env:
        - name: MIGRATIONS_PATH
          value: migrations/signer/postgresql
        - name: DB_URL
          value: postgres://postgres:changeit@harbor-database:5432/notarysigner?sslmode=disable
        - name: NOTARY_SIGNER_DEFAULTALIAS
          value: defaultalias
        image: goharbor/notary-signer-photon:v0.6.1-v1.7.5
        imagePullPolicy: IfNotPresent
        name: notary-signer
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/notary
          name: notary-config
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          name: harbor-notary-server
        name: notary-config
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: notary-signer
  name: harbor-notary-signer
  namespace: kube-registry
spec:
  ports:
  - port: 7899
    protocol: TCP
    targetPort: 7899
  selector:
    app: harbor
    component: notary-signer
  sessionAffinity: None
  type: ClusterIP
`
