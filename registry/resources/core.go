package resources

const CoreTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  app.conf: |
    appname = Harbor
    runmode = prod
    enablegzip = true

    [prod]
    httpport = 8080
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: core
  name: harbor-core
  namespace: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  secret: bnpOUWZEd2c1V09zdzNpZA==
  secretKey: bm90LWEtc2VjdXJlLWtleQ==
  tokenServicePrivateKey: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBK21UbWRiUFEwbi9GOGZqQVVKMEU0UE9wd3NXb09Qc3pKVy9TeC9yRjUrRkV4emZNClVPUVJUUTlBWEdnWmJCb3liN1oreDdSVFZjOUFNMWZ3eSt2d2NPMElMa3ZTckVuMGRadEY3V2JzOVltL3U3NVYKMUJoRkFhbnYyR2hhVGw1dmFBTTE2dHk4MTdrMHA1UWVGOEozY3RLVkVBSnhCT0g5N3kvajhvVUh4eWE1TW4wTAp2RnFqVXMyK1NWS21VQWNwNTlCaVFGYU9jQkErTnhWU2phTVB0QUFIYVU2RThYSVNsSjRPbU4yaFpuTWdLNVhCCmxYQ0hKdWtIVjUwcERWMi9QR05hSzhpSkFJcEVlMi8zSk1XVUpxT3pRSHpHMlVVQS9FSnBHVTFTdTRkT25EYVkKaEFRTHR4Zkxza3k5VU96NnNJMnM1Qjg0aFlkRi9ZYnA5VVJncHdJREFRQUJBb0lCQUVWMW9BaWVyUnIzbnUyVQoxNlNGS2tsTXpwYmRSZUVvcmZPQXBiUWIrTEp1Wnluc2JKMHo1eWk5UmxsYjkxRnBvdndpWCtEK1FPL1k0akpiCi9zeFMwd3hBZExpRENCb0xHZWxWL1h2eHhXSUhxRXhvYkY2OXJNYmRZVTlqQTBvaUlEMTJSb3EwV1I1dW5oQ3UKb2ZwdFU3MWlkVDlQcmhKd3JvU2ZnRmhTVnVvRFkzd1A2NUxTbGdnamt6KzMwZEpIQ0p3VDd6dGFDWkFmbTlVQwpGUjF5TTBNR09FQVNseUptV1RRcG41UWdFbWlGKzFpTjFXaTM3dnVBQ2M3ZURzZllGc1FqRjBSL2FLOUJRdVRBCjlEWkhPeld2R0ZxZHNWUU5mYXJXRVBiNlAwU0dQc1JjWFY3MVVla0VPaWxmbDV0MWk0SjUrOG5FSmxQMWZuczcKOHBzSmxLRUNnWUVBL3RNenl5VDFuWWxvMkFTaUFqRXlOYlorL1BOcU1PNS9wR3crV1UrZjNPbjdFckFPRS9ZNQp1cUtQa3kvM3ZacVZyM1hyK0V3ZXJZbjdjNG16NG9yUDZPYy9ueEthRkJQV0pWWUtsL1RIblFFdjdiMjJtVUlmCjdIWXd3aXBQUnp4RDZTUlB0K2VWNlR2b1NobzczeHFWTzQ3Zk4rUDVQRndmRk5iMEJya2tyNXNDZ1lFQSs0eDMKdWVEUUZlbnRmMjE3bmZPUHpUcFNpQzBVMmJyalNaTngxdFY3Vkx4NjkrYmtUY09NcXlqU1ZIYlF1MmNKWUVoOAp6NDhXbVlXOFZKRHBseFhPSDlXWG1kODl0RTlOYU9kL0N1OWs5c3VqQUFKU3A2NEhtUDg2cVBKVG9OVnJrUWNZCkJuaXVqeG0yVGhiL3laeEdvMXpiMHRlMUcydkMxNGErWnNUK0VlVUNnWUFOeTZsV0toNFI5VXB6eDJ4dDZmUHAKN0lOYmRtSWRYQXdVL3JjeFpwb2svNVhVSVN2aDhNYVhVQTJ0emo3L1NNc3B3SnlSeUswd2YvUFpBVzkzcUVReApPN082RE92Q2dvQnBiUXNOeHZhM2pVVG0vZ3BRcWIvSXNXMWFWYWdORnpvbCtRMUh5NFhXSnFRZ3Z1TFc1VDJICkIra1Z3WVhRdXJ3RUNNOFZQaGk2V1FLQmdRRHV4R3RValZjV3BkL2dKNHpCNWRHbWJPaXdCNUtXQlBYKy9heDkKems5dHBDWllydG9nRWpDd3VtUEM5ckMwWVY0ZC9WQXpOOCtzMDZ0cTNjNUxzYy9nbWI1M3VOWDZFNWdYcmowQgpwVEJCcmhNL1MvVW00bUtsMEFYZkhYMVIyYUpybTc4clRWdnJ5dTBuZkY4NUFGUndkaERXTmhmYk9sTk1mc0syCkMrQmFIUUtCZ0dlUWJ3WTh6RXplcCt1VnBPa1I1anM2b3FMbnprVFdiN0w1Y3BsTnBwdFZnaysxYk5FZEZVNU8KMGVMN2hOTEF1Z3NnT3RtQ2h4c0VnTGVBNlhBa3JrbytGSGhrbGFMQmZyQ3RZeWovQ2NpWE15MitLOUU3SWxtWApyTjNxaUJGcEFZNGRIZHpKcWJPS3ZxVXJzYkxTVkpsdTlvTE5mOFNYdUpabCtJSHBtSmlsCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
  tokenServiceRootCertBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM3RENDQWRTZ0F3SUJBZ0lSQVBhUm5QR3czOHBWNXcyYjd6SzQ1QTh3RFFZSktvWklodmNOQVFFTEJRQXcKRVRFUE1BMEdBMVVFQXhNR2FHRnlZbTl5TUI0WERURTVNRFV4TURBeU5UVXhNRm9YRFRJd01EVXdPVEF5TlRVeApNRm93RVRFUE1BMEdBMVVFQXhNR2FHRnlZbTl5TUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCCkNnS0NBUUVBK21UbWRiUFEwbi9GOGZqQVVKMEU0UE9wd3NXb09Qc3pKVy9TeC9yRjUrRkV4emZNVU9RUlRROUEKWEdnWmJCb3liN1oreDdSVFZjOUFNMWZ3eSt2d2NPMElMa3ZTckVuMGRadEY3V2JzOVltL3U3NVYxQmhGQWFudgoyR2hhVGw1dmFBTTE2dHk4MTdrMHA1UWVGOEozY3RLVkVBSnhCT0g5N3kvajhvVUh4eWE1TW4wTHZGcWpVczIrClNWS21VQWNwNTlCaVFGYU9jQkErTnhWU2phTVB0QUFIYVU2RThYSVNsSjRPbU4yaFpuTWdLNVhCbFhDSEp1a0gKVjUwcERWMi9QR05hSzhpSkFJcEVlMi8zSk1XVUpxT3pRSHpHMlVVQS9FSnBHVTFTdTRkT25EYVloQVFMdHhmTApza3k5VU96NnNJMnM1Qjg0aFlkRi9ZYnA5VVJncHdJREFRQUJvejh3UFRBT0JnTlZIUThCQWY4RUJBTUNCYUF3CkhRWURWUjBsQkJZd0ZBWUlLd1lCQlFVSEF3RUdDQ3NHQVFVRkJ3TUNNQXdHQTFVZEV3RUIvd1FDTUFBd0RRWUoKS29aSWh2Y05BUUVMQlFBRGdnRUJBTVlwdXgrUlNyWFAxeGdNYnFFUjZNV204dSs1VG9ldVFoR3RCRU5STWxDYwp5T1RncEZKNUs1UjdVandzSklsTm5RUW5nYmxIRVpqUGdvYzQvUWdXc0dJZWxMUDFBSFVIaUhvclZuQnUwY3d6Ckl5ZVJ2ZGhzRzUzR1RvYWNEbjV2ZXl1OURoV0k5ZHZtWDdRd3hyODc2aGpPb1JqdHZHNFJuYXpETUFPVHpoNnEKODIyUHBJNERaYVl4L284WGJzNTA0a2NHQTBPNFdTZUprKzBhOTVHUTEvYVc2bVIrM2FjRWdzYlBtNUVWV284Vwpla0hTa1g0MktjVHh3TitEakJFemdqcEg5V0xkUS8yNTJHbW0reW83ZmFrZFgvQjNoanB6ZFpWN1lObjgvM2hzCkpybVVFOHBSWC9peC9GalBPbzVYWlVQU3hSTHJLRngzbFJJYTBXenIrMXM9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
kind: Secret
metadata:
  labels:
    app: harbor
    component: core
  name: harbor-core
  namespace: {{ .DeployNamespace }}
type: Opaque
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  generation: 1
  labels:
    app: harbor
    component: core
  name: harbor-core
  namespace: {{ .DeployNamespace }}
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: core
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: 3b4121f192f68cd8015d18f223b9af3894de6aac6381660f1bc1b69db3239ba6
        checksum/secret: b6634612d2e4d7bca1f97fbcac91188a32edcdee19d919fadb3ca79ffca47b39
      creationTimestamp: null
      labels:
        app: harbor
        component: core
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
        - name: _REDIS_URL
          value: harbor-redis:6379,100,
        - name: _REDIS_URL_REG
          value: redis://harbor-redis:6379/2
        - name: LOG_LEVEL
          value: debug
        - name: CONFIG_PATH
          value: /etc/core/app.conf
        - name: SYNC_REGISTRY
          value: "false"
        - name: ADMINSERVER_URL
          value: http://harbor-adminserver
        - name: CHART_CACHE_DRIVER
          value: redis
        image: {{ .CoreImage}}
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/ping
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 20
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: core
        ports:
        - containerPort: 8080
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/ping
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
        - mountPath: /etc/core/app.conf
          name: config
          subPath: app.conf
        - mountPath: /etc/core/key
          name: secret-key
          subPath: key
        - mountPath: /etc/core/private_key.pem
          name: token-service-private-key
          subPath: tokenServicePrivateKey
        - mountPath: /etc/core/ca/ca.crt
          name: ca-download
          subPath: ca.crt
        - mountPath: /etc/core/token
          name: psc
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          name: harbor-core
        name: config
      - name: secret-key
        secret:
          defaultMode: 420
          items:
          - key: secretKey
            path: key
          secretName: harbor-core
      - name: token-service-private-key
        secret:
          defaultMode: 420
          secretName: harbor-core
      - name: ca-download
        secret:
          defaultMode: 420
          items:
          - key: ca.crt
            path: ca.crt
          secretName: harbor-ingress
      - emptyDir: {}
        name: psc
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: core
  name: harbor-core
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: harbor
    component: core
  sessionAffinity: None
  type: ClusterIP
`
