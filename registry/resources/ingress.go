package resources

const IngressTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-registry
---
apiVersion: v1
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM5RENDQWR5Z0F3SUJBZ0lRWnFBQThJTS9wSnBEUEh2ZVhUdkEwekFOQmdrcWhraUc5dzBCQVFzRkFEQVUKTVJJd0VBWURWUVFERXdsb1lYSmliM0l0WTJFd0hoY05NVGt3TlRFd01ESTFOVEE1V2hjTk1qQXdOVEE1TURJMQpOVEE1V2pBVU1SSXdFQVlEVlFRREV3bG9ZWEppYjNJdFkyRXdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCCkR3QXdnZ0VLQW9JQkFRQ3AwU0VhZDdHRnYwdm05UjNCL0RSN1hIZmVZKzQxOE9uVTExUUo3ek5UVW01SlFvTzgKYkF1T3hZR3pwT3FMM3dLSHljWXZYZUtFRDIwRDNRTHRydFpvTGNYU3NBeXE5RE11bFZ2YzI0aUo3dXZZK3hpVQpsNmozUWVvZmZPaVpvTWdnckN6bWtkdlY4MVNQZ2xLMGdENTM3YUZUaytsMGdMaEV1aVloZFpVTktDaWkyczVvCiswcHJEZXp6dWE2ZjVWTlpnZ1NNRnRkdmhETWI1Tnd2SGVIdXdNM1Awb2NWa01MbGtOOEF1UmRwZ0hZNVltWEwKUWpLaXE0N0g4aEkwSUpRSkM3TTcxNnQ3elZiaGdOQ1dZbVRoTW9rNStuTldZZHNtRlRZMDN3WmFyc3hZdnVweApWdFF4b0diVEdNNWdENUh5UkNZQnlrM3dlK3dJRHVDN0M4MUhBZ01CQUFHalFqQkFNQTRHQTFVZER3RUIvd1FFCkF3SUNwREFkQmdOVkhTVUVGakFVQmdnckJnRUZCUWNEQVFZSUt3WUJCUVVIQXdJd0R3WURWUjBUQVFIL0JBVXcKQXdFQi96QU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFlYzJ0T1d3bkdjT1dGLzM2WEhTT1JqR1NQdGoySVE1YgpwMzNBVDY5V0FwdjFVdkZtWnhseTIyQ09MRitaci9YOHJlY1U4RUVFSVVua1pkblRNUmthVXhUSmVHQXVsSFJMCmEwTVdTMGtTUThxaHlBZ0lXUlBQOThSSE5rdWF3bzByWmIvYVdsa3BtZy9HNGJIWU9uTU1vSk9RWjBpU3ZTNnQKWEJXa09ucGNoanRqRjI5MTlxdjNaNWRPRER1dk1TTjd1aytwUDk3QUFUbTdqOVVNREFwUGppWjU0dnREOHAxZApmYm92bnRIZVVMME5QekFINkVHdmoraWtWbnJITUpqcC9jcVIrZEYvTjdsVWd4dlh6dmptN05CQXB6T3ZDUTZmCmJsbUVPVWM3RGM3VXBFcWJUNmFqWjV5UUN4bEQ5YjRtMHFIQ1dxMFRHZjhzQUcrYVpIbUFkdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURPVENDQWlHZ0F3SUJBZ0lSQVBuQkNnRXoyeG83SFVBTDVZYW5PVzB3RFFZSktvWklodmNOQVFFTEJRQXcKRkRFU01CQUdBMVVFQXhNSmFHRnlZbTl5TFdOaE1CNFhEVEU1TURVeE1EQXlOVFV3T1ZvWERUSXdNRFV3T1RBeQpOVFV3T1Zvd0lERWVNQndHQTFVRUF4TVZZMjl5WlM1b1lYSmliM0l1WTJ4MWMzUmxjaTUzTUlJQklqQU5CZ2txCmhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBemtVL3RMekVpeGkrTmk0d1RnUFZsV0VrZmFjRmY0Q2MKMkpIbmZuMTU1NWR1akJvZGg2a2tZa2VTNmg2OVI3TEh6Z1hteHJLSkE1dlB1OS8xSGdCVmdSSktKM2JTM2R2QwpsTzV1ZDQ0TzA4S096THdWb0VYMVhkWmFKQWo5YzRqUGlIMVdVb1RVVVpiU0o4WVdpS3FaelJUS1lqRDlydFZBCks2ODNGRnUwa2w4Y3dIajd3aHRsaHBoL09wUHFjdDBwdy9zNEwzMHJVSmtaYm9Qd0tIbXYwMU5UeEtPTHk1UUMKaFR3b0xOdE15cjdIY3hEOW8yakpvSkNZbHFvdGJaMUlXcjdYK3BqRG1mMVB6SUVKbkI2dDVIMkh3ME9DMjJrcgpGc0JBRmNwSXR2eTRNRFI2blZPVFUxSzNUclEzTy9WOHAxNzB0RWpDQzNxK3krUFc1NzJQRHdJREFRQUJvM293CmVEQU9CZ05WSFE4QkFmOEVCQU1DQmFBd0hRWURWUjBsQkJZd0ZBWUlLd1lCQlFVSEF3RUdDQ3NHQVFVRkJ3TUMKTUF3R0ExVWRFd0VCL3dRQ01BQXdPUVlEVlIwUkJESXdNSUlWWTI5eVpTNW9ZWEppYjNJdVkyeDFjM1JsY2k1MwpnaGR1YjNSaGNua3VhR0Z5WW05eUxtTnNkWE4wWlhJdWR6QU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFFSE5LCm0yMHVzdS9VaFVhUGZmMkdaZ21JQk4wY0dhVnBvM0YwdnQwamhnOWU0Qk5Pc0REeEI4YjFkZUJHNS9xTlZNbzgKRFV0UC9NMmcxdlRJc09nYTJseFBhWjNVT2lNWGdCWFZRU2dsa1JGNEx2WEh2cVBwbmJkc28vWVJINEhBU3diYwo5YVI4YWhOVVpGY2w1ak9jZUd4YklVRnduZmNIVTIwZHZFOGphWUZ5aXNHaG5oWjBmN0xIaVpzRGVHRmVHREZICmJlS0FLT0RtQWY4WERob3hDTlRvMUhzTjRhNitrbVpmMGt3WFI0bmk1TzU2TTE1TTlrbXZPaGxmQUw0S2NNaFYKYVY5ODl1SUV3RGN5S3BBeUkrTk03dzh6S2FreUxEclBQQ1RLdHhjTHk3MmhiM0lac0xDbHV5N2RwcnoxSFlGdwpvUVJ5NFhCUzhhVCt4QTkzREE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBemtVL3RMekVpeGkrTmk0d1RnUFZsV0VrZmFjRmY0Q2MySkhuZm4xNTU1ZHVqQm9kCmg2a2tZa2VTNmg2OVI3TEh6Z1hteHJLSkE1dlB1OS8xSGdCVmdSSktKM2JTM2R2Q2xPNXVkNDRPMDhLT3pMd1YKb0VYMVhkWmFKQWo5YzRqUGlIMVdVb1RVVVpiU0o4WVdpS3FaelJUS1lqRDlydFZBSzY4M0ZGdTBrbDhjd0hqNwp3aHRsaHBoL09wUHFjdDBwdy9zNEwzMHJVSmtaYm9Qd0tIbXYwMU5UeEtPTHk1UUNoVHdvTE50TXlyN0hjeEQ5Cm8yakpvSkNZbHFvdGJaMUlXcjdYK3BqRG1mMVB6SUVKbkI2dDVIMkh3ME9DMjJrckZzQkFGY3BJdHZ5NE1EUjYKblZPVFUxSzNUclEzTy9WOHAxNzB0RWpDQzNxK3krUFc1NzJQRHdJREFRQUJBb0lCQVFESzVMTFZSUmpPK1hlZwpNNjZ3RG5WNGlpVXFzNjlreTAxOGVZZ0xrOERsWEw4UWNGKzdvVlI0bDQ2Ylc4RXpWVmZULzFvUStHeHRjRVhWCnQyV1VMUi90NWQrckVlWTQ5SUZobldacmt3QmlxMjFyVlZhd1lDQUtQVjVOTThxYWFtZzVDWkJ2ZXRpZHFJenYKTXBuWHRIZTBrazdBWnhBaGVRRzE5cE5uSXcxckt2NkY5VWpoRG41TyttVVB6OXBGbFBxRi9kb0l2bXhuWTlVdQprSDNJUUE3cWt1cDNvU2g2U3Z3a2N3L0NZZ21SNmFHbFFQS0xlTTZVRnhxYlVKRy9CSXVjUDNyeDdXVWdIcWlnCkRIVFpNNENjcFdmb0ZyQk05OXlrWTM5cXFST0tYQzUyaFV1cytPQnNGVGNOMHJTNWlKOUlnL1pxMEtBWUlud1oKNndvY0VlS0JBb0dCQVBNRUlwN2NVbWxFaytxTFBPWFhGRm5sNFpsaExiZ0hBSGtySjk2TkVjeUhuVy9qdDJNWQpVdE9RRHdqenBsUnlqV3Y4SzQ4UDFRWlFseWVKbTNQeUMzdktlNWUySUFMcFYwU3JxdVVWMTVVS2Q5Tk1LdUsrCk0xMTVzYzhXYkthRmlkTmM4QjJ6ZmlpeERVVjV4azFMUGFKb3lUdkl6NDB2U3dJRml4QXRCNnJsQW9HQkFObEsKaGZWZDZLRm9HMkFZN0R0Ym5OZWpySll6cjRCYWZ5a054Z0RjaXM1WlhUdVZwVS9zQzRmWi9BMkp0RUZtS0tObApPTlAvM2QrWmg4WmY3cFhXdVRiOUlBV05NT3pjR0RCZUVYZ09hTXdvdVQvQTN2WkNTT1B1UWxhQ1NWOGxRQmRECng4RHFVdXJGNVBxL09VTHBDbHJvT05VQ1VkVVNUODdWa1YrMExvN2pBb0dBZlQ1MFdVdFRiYzFhTGxiMFc4QXQKVE9lZERWOTRJSS83UG5kdlpOTnZpT21ieWo2aUZRQkVMNlRmR3MzM3V5NE9sTWl4NmxsT2dLS29SRWczUmRwSQo4Tk82UHNZdVdWSEpHQ3NoT0UrNWU2YVpldHlXZmFWbzg1UDBmN1llNlBPSnhOVHhLMTJHZDVKSU5MWTk3VGdKCjI5b0ZYRHB5UFdGU1Z3aGtVMEhoNGNFQ2dZRUExdEI3K29UcWRKZE5hUVY2bFh3T1pJamxHR0RrZ2duMWFJcTcKVXlLMC95Y05xdGhZWTlqQjFYNUZWc1RxTlRWZnU1bFlRdzNUTlRpOUovcXpEZU1IbkR0R0t2YlloWEZaWUllKwowV2U5Wndoamk0bUxZdjFJdmoxUHkrSmwwdkFxbWxWaDUzQkFwT1ViYUdFZnBPeHFWbEQ1em12S3B0REJvWU0xCmd0M0lvVzBDZ1lBNGttdTMrRSs2Um5CdkhKZFA2eDBQczFya0xUZjZ3a1RlNzQzWEpzeTV6c2JTOG5qelZjR0kKUDZseEdYSFE0cy9odi9UTEF6dVJEVTdmK2dUN3doRy9yQ3o2NWFMaFgrbkN3RHIzZ1MzUllFUVpMNTdaQnhEOQpzS1JDQXRreFlsaGswTmloN1FlejY4TTNDaTlOK3RBMXp4MDRrMGFvd01ORDgzVTZyRythWUE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
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
