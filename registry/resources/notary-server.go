package resources

const NotaryServerTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  notary-signer-ca.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDAzCCAeugAwIBAgIRAOAhivwD9W/TqfqLqdzfFT0wDQYJKoZIhvcNAQELBQAw
    GzEZMBcGA1UEAxMQaGFyYm9yLW5vdGFyeS1jYTAeFw0xOTA1MDcwOTEyNDVaFw0y
    MDA1MDYwOTEyNDVaMBsxGTAXBgNVBAMTEGhhcmJvci1ub3RhcnktY2EwggEiMA0G
    CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDrY2CFDnv55fDKP5zcUnbpTwv9o4+K
    U7a39aOr8TFDkYd1zs4w8WWPDjrCOTFvJzG63ieCrGhUizwDxFyh7BSY4WD/rjBn
    RaXHWvMJz4nj8obIAE+WVPWx3Hy8U8NIpb1d+5b7acYZlSA8a+OJt48e2XmFAC7f
    y57yQOIi/oww2OKJut5DEtRzkCuukQwc8NepCNHoW6nSn+es0WZRwPuUipvy1r0G
    22r4tijBT4TgMh8DDSY/z7Xeqi0DNRK5KSbci6Zdtop0JNztleVNwpt0sEXdqmaq
    LJryHDMNMrnGazYWkL+mpKxRYgvRiHJVkfk+4ni8HYD1lM5dnCDCkPKVAgMBAAGj
    QjBAMA4GA1UdDwEB/wQEAwICpDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
    AwIwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEA2phTqYvNt0VO
    o01g2Fp5AqnxUNeLOjJhf8i6lgvtQsgj46zOxjVvqHrtXXpPBeccX6ST2BzRAJ1a
    EnDZzQ2o/OP8J/lNJCznzjxl32EzuRPGeZcMS5gOrj8Gyquv6ocbBVawf9fSwOGE
    bVC9HWylCd/b90am4H/rsmHArWjOQoABZwHBT8ilZf6XFjwI4mo3UwuCw7U5PHr4
    f+UCRmZAF9XsOIb/qDrXGrCCA11zmvTVCn5+MNm7Sd6VWTsZrQMJVD9mpNjmgrrj
    ARjRdqZthtk+VrDPECMPdRGwSScESfWIryel9at2Jhz5ECWcYH35THL0XbyaVi7G
    qOfUbxFZcQ==
    -----END CERTIFICATE-----
  notary-signer.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDCjCCAfKgAwIBAgIQDBc7164TZy3BkSasdNVKfDANBgkqhkiG9w0BAQsFADAb
    MRkwFwYDVQQDExBoYXJib3Itbm90YXJ5LWNhMB4XDTE5MDUwNzA5MTI0NloXDTIw
    MDUwNjA5MTI0NlowJjEkMCIGA1UEAxMbaGFyYm9yLWhhcmJvci1ub3Rhcnktc2ln
    bmVyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1XyWiarcH0aWgk45
    7iZ16NBCnGiYzhotA2nLShFpgvkgQMH4AHcvvYrwhVcDLplM8bIk6VEmCPZc3AWU
    08Qfmm7e+2cGAZwyalW2R7a6vBqWn3Rxo2SeDwPNUYfm61ZTlJ91bbXAXBu9njt6
    mT2fKMdEnfzArlSRmCAwynaccybFkaRrpwU/IqNYDOblaTtmMRzfhxvtjvcbyGp+
    StzbqCgLJfhTtq/HHWPmBgVwyiuAlpDiNw6TEjJChk9tNplngG5COEdtwXdmTrAo
    P3IF0OYOChEsIiHxH9L5hpqA558U4qHSvK/M7JIwn8QJlZkcpILxaU9ex7sp2FBC
    HXjB3wIDAQABoz8wPTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUH
    AwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQELBQADggEBAEqi
    jtq27YgaWI4Tqy62BHuV4mlFARBDj6nU2hjA1bFyVicQR2eFd9o36ErwPLhuWvGR
    78K5s/C6YyCHbj5CfiRLMFjp0s3k4t+b1caaS9uA9fCh2Q9B5lJyJQUTMQl3saFi
    dHt/9azvX7AQ2NiGg6Iq87pP3ufyctfdiYcNJszObqfs2ENu9iNH/wtx7ykFSHl8
    vxEw3Z9ccqYa7mKARktW1TfdmRU4SlZXIVd1pcH/noA6iU88TYw1inguo6VfsHni
    exYV3kzG6a81/jv8IUgIgaUKlHlj3oZSSfhP7EF03NlWeZ72InPS2dgLzZUbkpCm
    hByLWAhIp9CDW7oXTA8=
    -----END CERTIFICATE-----
  notary-signer.key: |
    -----BEGIN RSA PRIVATE KEY-----
    MIIEpAIBAAKCAQEA1XyWiarcH0aWgk457iZ16NBCnGiYzhotA2nLShFpgvkgQMH4
    AHcvvYrwhVcDLplM8bIk6VEmCPZc3AWU08Qfmm7e+2cGAZwyalW2R7a6vBqWn3Rx
    o2SeDwPNUYfm61ZTlJ91bbXAXBu9njt6mT2fKMdEnfzArlSRmCAwynaccybFkaRr
    pwU/IqNYDOblaTtmMRzfhxvtjvcbyGp+StzbqCgLJfhTtq/HHWPmBgVwyiuAlpDi
    Nw6TEjJChk9tNplngG5COEdtwXdmTrAoP3IF0OYOChEsIiHxH9L5hpqA558U4qHS
    vK/M7JIwn8QJlZkcpILxaU9ex7sp2FBCHXjB3wIDAQABAoIBAQCjFIQub92s6pAo
    xDcOjES/7u8jaedocah3Fgbb8scl7MbNkR6wxFsssIkhYqGkpCiZ7RqzPHEQoZm3
    3W+eARCfORiO9VkqO7ZrckRHLfghnzH2Zs40IbV4BNB/+o/UsGIyg0kB4Lgr5GkK
    CaeSjfcaAHaTNTO/OAzsJ5L95nOGpe8Vf5EA7+vKBelFEhny8P0HeRnJrrelxdAH
    KJJQ5r+Nxsb+v4ta4uvqpC7BNRn1NVNl3cOidzMY1yjmylrMRPhVN06a94nUi7u0
    qwkzRwsMuzceqUiVgnCGnvJkA7n1OErlv6j/Sspkh9zQxxSHSZn045rd84dT49vV
    wK3HGZXBAoGBAPKvP/sa23Hq+tkcr7F8SDv1NFw3oKhukBE3O/rcIiEOX/bpESed
    paxN69q2yiKD37vHdaPNMTIgRzr1hJOJplsNsCxQtpEB0a+gGxUnRyuf03sD11yQ
    1HmQC3ZC8+vUuxNr2oiYtY0dXYXMOq65evx2tMrw0/PPhCn9jJUzcVzhAoGBAOEz
    OWA1BgdRZCGXJ/GiGmGI5wS4agu+rhmsNsjjN0W6SWg9oFqk8ena/95S08VOb3AX
    e53PlaCjpg8o080Wgn8EOZ1bLXNfhRE/AX4uUB5FBjiBOFt4b20Nj2lvlfahaMij
    BfKDzoJuqmPLtU3rsBNaA6OrCa5BufQz5riW7Ta/AoGAPhADVLwxkph9PjjP1Zvq
    /SpgEZVISMq9nSl69VSGhd2fPQ2tjWwLil0DDBPi7aC7/tGrjBBVnHQUw0c2eGSj
    XnXJsAuUJNFKRpezVV2OHeHpu3PoB4wiSlREGiJVLuJgVT8ny/cBtuzjlev8teJJ
    SXcyFRQxoBBZxENLSHy3aQECgYEAzDcNYqbyrpQqPyO5fy9GyQfCps8sqzXg3zsB
    +y3Ao6SIiNTJoylMjoqf2NY3YAb+myFQYg0qXJ/KKJkXaDVvZQtJy94w2xzVqIwA
    KJKK6MgjGf5kQt51/Oh9Elm0HhDE2pyq+f54uGLudMz3vo9p/kJ0ZmjlwHWt0Tt3
    kBCdUDUCgYBhLnVrbL9jSBb3Yw9ri3AP17XgwOXhUSSVAIgZeueclzbXs7W64x2W
    lts0mcGwv+QJci4Vs6WWSoeCZGQwIVpXL7X2zSs3GiBlrJIOlI2xWrgvTFd/OFQS
    R9TmbOKtjCcOVvYAAjyT6q2DecjpgzOycQB/hu2vCENgoOav9M3lfg==
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
              "realm": "https://harbor.cluster.w/service/token",
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
    component: notary-server
  name: harbor-notary-server
  namespace: kube-registry
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
        checksum/configmap: 9dd0ef0ed8079083135c784c1016b56326d45a506bfc4194e0e33f2ffa0f7be2
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
        image: goharbor/notary-server-photon:v0.6.1-v1.7.5
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
  namespace: kube-registry
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
