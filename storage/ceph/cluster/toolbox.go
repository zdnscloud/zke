package cluster

const clustertoolboxTemplate = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rook-ceph-tools
  namespace: rook-ceph
  labels:
    app: rook-ceph-tools
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rook-ceph-tools
  template:
    metadata:
      labels:
        app: rook-ceph-tools
    spec:
      dnsPolicy: ClusterFirstWithHostNet
      containers:
      - name: rook-ceph-tools
        image: rook/ceph:master
        command: ["/tini"]
        args: ["-g", "--", "/usr/local/bin/toolbox.sh"]
        imagePullPolicy: IfNotPresent
        env:
          - name: ROOK_ADMIN_SECRET
            valueFrom:
              secretKeyRef:
                name: rook-ceph-mon
                key: admin-secret
        securityContext:
          privileged: true
        volumeMounts:
          - mountPath: /dev
            name: dev
          - mountPath: /sys/bus
            name: sysbus
          - mountPath: /lib/modules
            name: libmodules
          - name: mon-endpoint-volume
            mountPath: /etc/rook
      hostNetwork: true
      volumes:
        - name: dev
          hostPath:
            path: /dev
        - name: sysbus
          hostPath:
            path: /sys/bus
        - name: libmodules
          hostPath:
            path: /lib/modules
        - name: mon-endpoint-volume
          configMap:
            name: rook-ceph-mon-endpoints
            items:
            - key: data
              path: mon-endpoints`
