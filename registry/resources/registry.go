package resources

const RegistryTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  REGISTRY_HTTP_SECRET: MnpXMkwzUlNOSVc3Zm1MNQ==
kind: Secret
metadata:
  labels:
    app: harbor
    component: registry
  name: harbor-registry
  namespace: kube-registry
type: Opaque
---
apiVersion: v1
data:
  config.yml: "version: 0.1\nlog:\n  level: debug\n  fields:\n    service: registry\nstorage:\n
    \ filesystem:\n    rootdirectory: /storage\n  cache:\n    layerinfo: redis\n  maintenance:\n
    \   uploadpurging:\n      enabled: false\n  delete:\n    enabled: true\nredis:\n
    \ addr: \"harbor-redis:6379\"\n  password: \n  db: 2\nhttp:\n  addr: :5000\n
    \ # set via environment variable\n  # secret: placeholder\n  debug:\n    addr:
    localhost:5001\nauth:\n  token:\n    issuer: harbor-token-issuer\n    realm: \"https://harbor.cluster.w/service/token\"\n
    \   rootcertbundle: /etc/registry/root.crt\n    service: harbor-registry\nnotifications:\n
    \ endpoints:\n    - name: harbor\n      disabled: false\n      url: http://harbor-core/service/notifications\n
    \     timeout: 3000ms\n      threshold: 5\n      backoff: 1s\n"
  ctl-config.yml: |
    ---
    protocol: "http"
    port: 8080
    log_level: debug
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: registry
  name: harbor-registry
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
    component: registry
  name: harbor-registry
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
    component: registry
  name: harbor-registry
  namespace: kube-registry
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: registry
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: d31ec1993e974889cecec80317081c34bdd81982bc7630af18aa97a11e79f230
        checksum/secret: 23c448540f2784f6ae4f0330d9d6a41aef34bc515d197a1cfd3e354b992eb8f1
      creationTimestamp: null
      labels:
        app: harbor
        component: registry
    spec:
      containers:
      - args:
        - serve
        - /etc/registry/config.yml
        envFrom:
        - secretRef:
            name: harbor-registry
        image: goharbor/registry-photon:v2.6.2-v1.7.5
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: 5000
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: registry
        ports:
        - containerPort: 5000
          protocol: TCP
        - containerPort: 5001
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: 5000
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /storage
          name: registry-data
        - mountPath: /etc/registry/root.crt
          name: registry-root-certificate
          subPath: tokenServiceRootCertBundle
        - mountPath: /etc/registry/config.yml
          name: registry-config
          subPath: config.yml
      - args:
        - serve
        - /etc/registry/config.yml
        env:
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
        - secretRef:
            name: harbor-registry
        image: goharbor/harbor-registryctl:v1.7.5
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/health
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 1
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: registryctl
        ports:
        - containerPort: 8080
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/health
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
        - mountPath: /storage
          name: registry-data
        - mountPath: /etc/registry/config.yml
          name: registry-config
          subPath: config.yml
        - mountPath: /etc/registryctl/config.yml
          name: registry-config
          subPath: ctl-config.yml
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - name: registry-root-certificate
        secret:
          defaultMode: 420
          secretName: harbor-core
      - configMap:
          defaultMode: 420
          name: harbor-registry
        name: registry-config
      - name: registry-data
        persistentVolumeClaim:
          claimName: harbor-registry
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: registry
  name: harbor-registry
  namespace: kube-registry
spec:
  ports:
  - name: registry
    port: 5000
    protocol: TCP
    targetPort: 5000
  - name: controller
    port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: harbor
    component: registry
  sessionAffinity: None
  type: ClusterIP
`
