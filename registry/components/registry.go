package components

const RegistryTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  config.yml: "version: 0.1\nlog:\n  level: debug\n  fields:\n    service: registry\nstorage:\n
    \ filesystem:\n    rootdirectory: /storage\n  cache:\n    layerinfo: redis\n  maintenance:\n
    \   uploadpurging:\n      enabled: false\n  delete:\n    enabled: true\nredis:\n
    \ addr: \"harbor-redis:6379\"\n  password: \n  db: 2\nhttp:\n  addr: :5000\n
    \ # set via environment variable\n  # secret: placeholder\n  debug:\n    addr:
    localhost:5001\nauth:\n  token:\n    issuer: harbor-token-issuer\n    realm: \"https://{{ .RegistryIngressURL}}/service/token\"\n
    \   rootcertbundle: /etc/registry/root.crt\n    service: harbor-registry\nnotifications:\n
    \ endpoints:\n    - name: harbor\n      disabled: false\n      url: http://harbor-core/service/notifications\n
    \     timeout: 3000ms\n      threshold: 5\n      backoff: 1s\n"
  ctl-config.yml: |
    protocol: "http"
    port: 8080
    log_level: debug
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: registry
  name: harbor-registry
  namespace: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  REGISTRY_HTTP_SECRET: SGxlZG02TE9TOVp4RVNKVw==
kind: Secret
metadata:
  labels:
    app: harbor
    component: registry
  name: harbor-registry
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
    component: registry
  name: harbor-registry
  namespace: {{ .DeployNamespace }}
spec:
  accessModes:
  - ReadWriteOnce
  dataSource: null
  resources:
    requests:
      storage: {{ .RegistryDiskCapacity}}
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
  namespace: {{ .DeployNamespace }}
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
        checksum/configmap: 03ff3a19cf52c6e28bc1a8e1363a69bced38a092c2de016e459679c6e86c1bb8
        checksum/secret: df96f191f774eca0e36e64468a61980b1db01467deb5f80a0b43c0626fb4f1c8
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
        image: {{ .RegistryImage}}
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
        image: {{ .RegistryctlImage}}
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
  namespace: {{ .DeployNamespace }}
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
