package resources

const CoreTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  secret: STdQdW5vVHZoaGtxRnZNYg==
  secretKey: bm90LWEtc2VjdXJlLWtleQ==
  tokenServicePrivateKey: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBdWNBaWY5NUVsQ25EQ29rY1NoMHUwd051SWYyUnU4aEUvQ29mZjR1c0JvRE9tZGZnCnNMd0x3aXZNYXdlWHNFZG5CSGdLakdTTk85ZjlIOVN3ZjROYVJ3d0JTWWFCclRuK0RoOTR4ckJuUmhON3Q5UlcKYi9hM25pdDQ1SFJkRFZLcnRVandoQVp5bE9nMzd0cTc3RkRwNGNJbGE0Q3FPcENqSDV6eXYzY0xTUWE4eVVaegpPMU9aTnpvM3ZodGFtQWJJTklqaHNEdGEwZis5RTM1eWZGcEswSVBJaWxFNzRqZUo1bERxaW5mZ2tTakJxa0lzCnlQMVhxOVcvbys5b0d2NVRhVklKM2xteGpldVhoek81cDNXUWNDZGQyS09pQWN4MnY3VkNoUmVXbGNlM2RNeGoKd2tNNzlBUzdDNGluUHRjQ3ZUMndlZW11ckFReExpcC9Qd3VlQXdJREFRQUJBb0lCQUJMalBwN1I4eGM5eDk5aQpZY1lIZ2FOalRlZVc2U0szRW95eG05SlVXWUd2eDFKTmFVT1BXNEs3QVdaQXQyUHRYc1JOa0lVR0ZIWnBXQXZNCmpNWHFqVnQ4RlFWczcvSzFXdWdFdXpzNDVNZHpZN2tsbGtSOFNkK0RmQjkrMVpkdE9KaU9laVp5b1dRUzJOMG8KME1NTzF6bGxqSGlKOU1WUHY4YVBKdi9oRXZ5eHpqQnE0S1JxR1NBNnZsVWF1OFdZOWtuSXlJbTQwRjRpTFNLRAo4WGE0YWFuejhqL1VsVDB5WXc3Um80RzgvR2k2ajR6RFJhUWsydlgrWHNzRFQ1cjhBQzRnQ0dNZHE0Si9YZmRaCnkxejQ0Vmk3dW5ITlBGOUl4WXJkVXZPaDBhVksrZ2t0dDdsa3JNdGwvTi9uOVZIYi9wcHNIKy8zSDJHcUs1bmkKMTZvYkZJRUNnWUVBeUxISnZ1bTBPcVlBdlRJbHExcEtDNHdLRUN1TEVTTkM0c0hEcmhYQWxyQVV2SUsrMWVqRApTa0lCdHJCeXZIZ0NmK3ZKMjRuVUFoUGc2ZzhOSDVUUU9YaHNVcUljdEY4eW5KU2QrM21rL1dmNEhmMjdLNVVKCi9MUmxrTk40ckxoaXpCU1AyMVgyMkxwNEcvb0JFd2lWSnlpU2s3WEZ3QXpuSldzc0NwZFFiYWtDZ1lFQTdQQWIKMjJCRXpicVFnWjc1Q0F0Mk5mSDdNZXZscm9aaGM4enN0VDFZMS80YS9JU3llUFhIMUtXcVkxdGlWTDQ0cjN2NgoxbGNKYURzWC9oaFdXY3k1ejh0bFQ5cVZqUURTS2JXK2d6aG9kTm95MExMMzVoM2pMTVExcDlLZzNLOGVRczJuCkk2ZkFBcmo3UDBPK3NrSi82MUNsTHIyNUFrS010WTR4dUJOL0Fjc0NnWUVBa3Q1SW5ZVzVkeEgwaUlBaVFQdWEKSkVrZk5DWXBaeWsrMFdLcktNS1NaYlFGK001VmlZVUZKVnFZbG5FYUJnSnRFZUFqb0oyRW9PQ2JQNjQwRkdCNgo0UlBYY2NGZzhENmFjeXZ2VVJEOFJOWEpKV21CaDZ0UjI3VElmdXZDNitNanFlV0NRU2p2dERzQm1yZWlBYVBPClF4SFY4bktiZktmMG45V0dMVm4rYWNFQ2dZRUF1MzZpUEs1b28vaFBwQk05OUF6RjVaaHdkQ2U5WUtjOGROdWsKTVNPenEzQ013R2p0cG1Td1Zta21kV3Q2VzU4UDBtWWtyL3ErR2ZveFdVUy9DRjdHWjFZSC9QSDNTRlp1K015MgpUcmFUaW15a3E1d0VLZGhhempFU1dKU0g4VHF1a3FTVTc5VXVUN2s4TU9zVis1QStFK09FWTRHRTV2SHMwVHNOCnd5SENicE1DZ1lFQW1xVmtEV1d3cHl4R3VjWXFsN2pzZ0lTOTVtdUdxbVRxa25yK0ZJcGYvdTlNZEx2emU5V3cKajBOU2NjY1I1TTBnVWJQeHFJcEl2TFFiSUgxaFpBcWtzMkdRczI5QkQwMnBhQmgzVWhhMjYvcGROd2tDaThhUAp6U2NmOVVKdHdoU3Y1RlVuSnZzK2dYVHJzSlk3TWhkWDl2QjQ1MUVJVFFCeERxZ1pkVU5sVFlvPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
  tokenServiceRootCertBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM2ekNDQWRPZ0F3SUJBZ0lRZEtEcWpNVTB6UDVaK0p0d1o5WHlJakFOQmdrcWhraUc5dzBCQVFzRkFEQVIKTVE4d0RRWURWUVFERXdab1lYSmliM0l3SGhjTk1Ua3dOVEEzTURreE1qUTJXaGNOTWpBd05UQTJNRGt4TWpRMgpXakFSTVE4d0RRWURWUVFERXdab1lYSmliM0l3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLCkFvSUJBUUM1d0NKLzNrU1VLY01LaVJ4S0hTN1RBMjRoL1pHN3lFVDhLaDkvaTZ3R2dNNloxK0N3dkF2Q0s4eHIKQjVld1IyY0VlQXFNWkkwNzEvMGYxTEIvZzFwSERBRkpob0d0T2Y0T0gzakdzR2RHRTN1MzFGWnY5cmVlSzNqawpkRjBOVXF1MVNQQ0VCbktVNkRmdTJydnNVT25od2lWcmdLbzZrS01mblBLL2R3dEpCcnpKUm5NN1U1azNPamUrCkcxcVlCc2cwaU9Hd08xclIvNzBUZm5KOFdrclFnOGlLVVR2aU40bm1VT3FLZCtDUktNR3FRaXpJL1ZlcjFiK2oKNzJnYS9sTnBVZ25lV2JHTjY1ZUhNN21uZFpCd0oxM1lvNklCekhhL3RVS0ZGNWFWeDdkMHpHUENRenYwQkxzTAppS2MrMXdLOVBiQjU2YTZzQkRFdUtuOC9DNTREQWdNQkFBR2pQekE5TUE0R0ExVWREd0VCL3dRRUF3SUZvREFkCkJnTlZIU1VFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSEF3SXdEQVlEVlIwVEFRSC9CQUl3QURBTkJna3EKaGtpRzl3MEJBUXNGQUFPQ0FRRUFCTWZaOWh2S21NZk9oT0RoQkdHbW05UUdOVVhEeE1lcmxOM0ZLWS9wMkNRawpHZ1ZuOWtVc0FIZ016bllweC9FWWNOdVBOZVZlMmp3ZGw1aVExZzdyZUluU1VnNEVmTk92YTdnQzRBSmxrN2poCk5EaGQ5RkZDdlhDa29WeUl3emdTTDVGcnRxemM3M0xWbFFWME5XSmZ2cVRkMldZRW92L3l5MVd1NUNCUHZzVCsKQ1lUam1xeTBRdjNPb3kvTkdtK3hwR2Q0TmNDa084eng4OVpQR2pjSlRUQUliRGs2QVpDTDV2d0RNTWFaak9XWQpKTkx6Qm9FZXVJdWNUVFVLc3pFVFlWSGNuNDE0N2FZWWU0WUlKMUpXY0VaTHRBVWN3d0RxZUhmMkdYN1NKdkVGCnpRUDBna2lkV0hKTG1pakhsbW4zbnMvQWtDZEtnWHlVcXJ0TWdEVjdPdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
kind: Secret
metadata:
  labels:
    app: harbor
    component: core
  name: harbor-core
  namespace: kube-registry
type: Opaque
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
    component: core
  name: harbor-core
  namespace: kube-registry
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
        checksum/secret: 0ded61ec2e077ee766595dcc27be46a69963e84ba11b59966ae55def2bd02d7f
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
        image: goharbor/harbor-core:v1.7.5
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
  namespace: kube-registry
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
