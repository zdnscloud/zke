package resources

const IngressTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  ca.crt: {{ .IngresscaCertBase64 }}
  tls.crt: {{ .IngresstlsCertBase64 }}
  tls.key: {{ .IngresstlsKeyBase64 }}
kind: Secret
metadata:
  labels:
    app: harbor
    component: ingress
  name: harbor-ingress
  namespace: kube-registry
type: kubernetes.io/tls
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    ingress.kubernetes.io/proxy-body-size: "0"
    ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "0"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
  generation: 1
  labels:
    app: harbor
    component: ingress
  name: harbor-ingress
  namespace: kube-registry
spec:
  rules:
  - host: {{ .RegistryIngressURL}}
    http:
      paths:
      - backend:
          serviceName: harbor-portal
          servicePort: 80
        path: /
      - backend:
          serviceName: harbor-core
          servicePort: 80
        path: /api/
      - backend:
          serviceName: harbor-core
          servicePort: 80
        path: /service/
      - backend:
          serviceName: harbor-core
          servicePort: 80
        path: /v2/
      - backend:
          serviceName: harbor-core
          servicePort: 80
        path: /chartrepo/
      - backend:
          serviceName: harbor-core
          servicePort: 80
        path: /c/
  - host: {{ .NotaryIngressURL}}
    http:
      paths:
      - backend:
          serviceName: harbor-notary-server
          servicePort: 4443
        path: /
  tls:
  - hosts:
    - {{ .RegistryIngressURL}}
    secretName: harbor-ingress
  - hosts:
    - {{ .NotaryIngressURL}}
    secretName: harbor-ingress
`
