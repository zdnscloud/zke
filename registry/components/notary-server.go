package components

const NotaryServerTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}
---
apiVersion: v1
data:
  notary-signer-ca.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDAzCCAeugAwIBAgIRAO5ZGsMfcfIrCCx93tA0QdwwDQYJKoZIhvcNAQELBQAw
    GzEZMBcGA1UEAxMQaGFyYm9yLW5vdGFyeS1jYTAeFw0xOTA1MTAwMjU1MDlaFw0y
    MDA1MDkwMjU1MDlaMBsxGTAXBgNVBAMTEGhhcmJvci1ub3RhcnktY2EwggEiMA0G
    CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDMYlgco4yYlYTCn6uEy6ynahPj470P
    F69dL22JIeLabuR95NbtMovQXrAnXeSvX9z0Zo4FS2CN2AMrX09iZ3d/xoVDMLX6
    cESpQzSOHI4Pc/B7Uu8AxWws5Uh8bp/Xu6zjH+/UVXoJyXcoWq38qzR9B5tAH/Hn
    uVss3JyN9BcWqRcREOIYc11VovIXrVfUv2LyOJ5/vhvYn7uJmXU90mzKgzor1V+y
    Gtg7uD319mMv6kjxurAJ50jH+I8WtGxqdkUrRihbVEK5gATiDv742ztqQHfWedyP
    RBP1aUQsHJVphoZM2tx03cp4xZWjWo5vc3ev4FefiEwR+mQBFX8+O5T7AgMBAAGj
    QjBAMA4GA1UdDwEB/wQEAwICpDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
    AwIwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAXUwsCfrRRnWn
    sla22F5NuYI/BPe09o9iWQqxRRWV6BQ2I17n17VkooV8pTh1mgbcrK45G9RP/Flz
    aYj2IJ5uu4b++SHjq4jlIfIchloZdbxbM7BRM+5IZ+q715W8789X7WyPobL63n7b
    3n2MIVitbxVPwhMRPlC7r5Nb14MCTLBJ52f9agwNXAuZUqWWqs0NNVe3Z5Ba/WF0
    Q6iGRrkIQ6qRLCHqgM8fNn50vhijej8QO3C3JfbD2BZqH87695r9tmWDIFF6Dv4z
    qUuXgqyMsoFV9Si1tJrUgcnz7YGlcCBgXE3ohyLydfrJLFq+lEASEHODCn5oa9/6
    y84AWFgiWw==
    -----END CERTIFICATE-----
  notary-signer.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDCzCCAfOgAwIBAgIRAPh8MvY9CjzK/SSu7zR5+ZIwDQYJKoZIhvcNAQELBQAw
    GzEZMBcGA1UEAxMQaGFyYm9yLW5vdGFyeS1jYTAeFw0xOTA1MTAwMjU1MDlaFw0y
    MDA1MDkwMjU1MDlaMCYxJDAiBgNVBAMTG2hhcmJvci1oYXJib3Itbm90YXJ5LXNp
    Z25lcjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMULYH/g7xBt8Vr/
    LZ8p/PJ8y7sTZnqWMQsVu4FMu3ZeuvaW/o6eZdXu3KhUTSa80xfkIcpNlUnD0M7W
    bem6MUMowRgXdf6sFdO5RoznaMmeF7JL0WjUvEebh3Z8TGn1oz7LscZXRjM8Qjb6
    CQ0bwMeyoAttBpboDSkTr3LMyYZcJ35yR3IFPcHzGrIj+CPpu1cNa3Btb7JabD+C
    IAbNTYzNF+4khwurDkIA7BVKbfPocnFN+a1U77UoKjvdl3pTUVmRo6P9YmnQiWw1
    d2RGacs9eKswCTU4pRjTNR++hFq8dWVKipLrLgCDC2JNFwk0iXpE0uVRtZA3pqdV
    KZtCCc8CAwEAAaM/MD0wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUF
    BwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQB2
    oYLSWCUTCelMEZGhG50PkY1+6KT5uTOdJ4GqzIl1VbgA16cPkkcfNHaQPb7gHZcI
    QQ1nrR1p63kGySd9GRCw09QV8n23iVki8xLMJ8EIjC9qIHhqAhzxDF3jlWWxlVjA
    XQw3hiGJIBKfJMDNyTFVjNg3Z6OiaHmRh1Vbvh71iXWcCoBf8+/MCusp8x7o0rEH
    wMhtOU3GI/y6a9FlhNrNSL6Sf5DYyZ6ekp1dWMCbuFop+/d3TmrMN1S7n0XE0Sgd
    lB9njJlWjQdbrNkIKmPW8fJP7ZpOR9kPecZHb+1LZQcl6+Uf74N1ZkP6GUBobvn4
    TZ2vmBTYm4dUp+1qAJOG
    -----END CERTIFICATE-----
  notary-signer.key: |
    -----BEGIN RSA PRIVATE KEY-----
    MIIEowIBAAKCAQEAxQtgf+DvEG3xWv8tnyn88nzLuxNmepYxCxW7gUy7dl669pb+
    jp5l1e7cqFRNJrzTF+Qhyk2VScPQztZt6boxQyjBGBd1/qwV07lGjOdoyZ4XskvR
    aNS8R5uHdnxMafWjPsuxxldGMzxCNvoJDRvAx7KgC20GlugNKROvcszJhlwnfnJH
    cgU9wfMasiP4I+m7Vw1rcG1vslpsP4IgBs1NjM0X7iSHC6sOQgDsFUpt8+hycU35
    rVTvtSgqO92XelNRWZGjo/1iadCJbDV3ZEZpyz14qzAJNTilGNM1H76EWrx1ZUqK
    kusuAIMLYk0XCTSJekTS5VG1kDemp1Upm0IJzwIDAQABAoIBAQCP6uZZoHWb32FI
    YFb3CJjql4HKKrpP0QETIpVoNB47r6cI0nIswr6IicT64U/UelgH/CU9+HqQfQg2
    +mEfNFIkxlB6gzA4iYILuSgxZBeiIukV3dCeq7q05oEDZnf9cF8CT46R2k64v0tx
    kiAoRdwBP5MrpT8J24U5OlqME80hZ0dYpjrpT8YKnsq8sZexOpnZpYEM/ZE3OPie
    uKqeRhciyuqYS9jCOdt+SrErl8M6MpEX6aQ7n2lIaQyC+cehVrSpz9i/28KnpFgz
    zzAYJlga3EQIN4qIj+vbGrF/rg2gbLsD7+MPR8M7mfzkzFWI1UGNI07CLGEHpLSh
    uHScyc5xAoGBANil1Ugnu1H6uvGmYUGycsnQJR8YPnhe46xaQP7KvJTlmwR9Z4TG
    gLMoX30d+q0C/kvvq5zV4XfovIFTGEPqbpeO1Nes82LpBt2Gya7T+didSYKfWpRU
    41XhWSaTMzdcv8IOzjfPld+BpHVax/8iRRsLeyzDK8sF/jftsUdzkNnHAoGBAOjV
    +/d6nycSMM9gUCiVndcUPFyyS+jCP/VnymNpuVtnQhdx7Fm07CiJszwoNyR3DYFo
    B1h/W7IolABXOYGlf4noPDPobNhftEDBfXRKgNwndrPe9xDP04+f1nLCu+iRPGJM
    j1t9hPdz1GfnF8PA3OdAMgIMrK7kUhuAEns+ng+5AoGAdyzKJXYNuiv1uEZxC6Wx
    NUj1kqRQgQCZt06yoDZABJxFaPPfBQ/47hTQalAafB/AaV8/BPg2njJ9t5pRJ9MG
    4QImHTo7bHaJW0TxHuXmc30aWet09VG4+J4M34ZrzxGLPqHMWLEtXZTANfopODTO
    1PC84kO+jGEQlg1/zrFIxjMCgYB6YK40n0Czen4pIUhAbJMvjrVDS3tWdXLEe68G
    nXUNM7KrO/esFsnhbK7GOaTyB5kToSfrPdVmSKmxnCbfm6rzQxsRdWJwP60wNALK
    crZUAHIFjHVzYqih3rMKUowNavi/+dmHjuuqXDkR+4akHuR8r2MZbKv+qIb3aVNN
    b9YIEQKBgBywbi9mTWesKIWROXUwzTNwkDoXZPTb6BHOZEao6e6KFfxi7NnqRRNB
    U+yrEsUW2bGlq96TyzSc/Hov6lKK69oC40glPUyV29id6rIjIrFfGX15zJIwpCy6
    v081YWYml3zVkLHMMxJZ70KxVkjwm9aa9SKtmnYF1iSMOAdZDex8
    -----END RSA PRIVATE KEY-----
  server-config.postgres.json: |
    {
      "server": {
        "http_addr": ":4443"
      },
      "trust_service": {
        "type": "remote",
        "hostname": "harbor-notary-signer",
        "port": "7899",
        "tls_ca_file": "./notary-signer-ca.crt",
        "key_algorithm": "ecdsa"
      },
      "logging": {
        "level": "debug"
      },
      "storage": {
        "backend": "postgres",
        "db_url": "postgres://postgres:changeit@harbor-database:5432/notaryserver?sslmode=disable"
      },
      "auth": {
          "type": "token",
          "options": {
              "realm": "https://{{ .RegistryIngressURL}}/service/token",
              "service": "harbor-notary",
              "issuer": "harbor-token-issuer",
              "rootcertbundle": "/root.crt"
          }
      }
    }
  signer-config.postgres.json: |
    {
      "server": {
        "grpc_addr": ":7899",
        "tls_cert_file": "./notary-signer.crt",
        "tls_key_file": "./notary-signer.key"
      },
      "logging": {
        "level": "debug"
      },
      "storage": {
        "backend": "postgres",
        "db_url": "postgres://postgres:changeit@harbor-database:5432/notarysigner?sslmode=disable",
        "default_alias": "defaultalias"
      }
    }
kind: ConfigMap
metadata:
  labels:
    app: harbor
    component: notary
  name: harbor-notary-server
  namespace: {{ .DeployNamespace }}
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  generation: 1
  labels:
    app: harbor
    component: notary-server
  name: harbor-notary-server
  namespace: {{ .DeployNamespace }}
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: harbor
      component: notary-server
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configmap: 27e5fa3f9c03881532c72b494d782dc95371e467d7d069f81abac3a3678f7508
      labels:
        app: harbor
        component: notary-server
    spec:
      containers:
      - env:
        - name: MIGRATIONS_PATH
          value: migrations/server/postgresql
        - name: DB_URL
          value: postgres://postgres:changeit@harbor-database:5432/notaryserver?sslmode=disable
        image: {{ .NotaryServerImage}}
        imagePullPolicy: IfNotPresent
        name: notary-server
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/notary
          name: notary-config
        - mountPath: /root.crt
          name: root-certificate
          subPath: tokenServiceRootCertBundle
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
      - name: root-certificate
        secret:
          defaultMode: 420
          secretName: harbor-core
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: harbor
    component: notary-server
  name: harbor-notary-server
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - port: 4443
    protocol: TCP
    targetPort: 4443
  selector:
    app: harbor
    component: notary-server
  sessionAffinity: None
  type: ClusterIP
`
