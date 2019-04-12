package templates

const LVMStorageTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: storage
{{range .LVMList}}
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: csi-lvmd-{{.Host}}
  namespace: storage
spec:
  selector:
    matchLabels:
      app: csi-lvmd-{{.Host}}
  template:
    metadata:
      labels:
        app: csi-lvmd-{{.Host}}
    spec:
      nodeName: "{{.Host}}"
      hostNetwork: true
      containers:
      - name: lvmd
        image: {{$.StorageLvmdImage}}
        command: ["/lvmd.sh"]
        env:
          - name: MOUNT_PATH
            value: "/host/dev"
          - name: VG_NAME
            value: "k8s"
          - name: DEVICE
            value: "{{.Devs}}"
        securityContext:
          privileged: true
          capabilities:
            add: ["SYS_ADMIN"]
          allowPrivilegeEscalation: true
        volumeMounts:
          - mountPath: /host/dev
            name: host-dev
      volumes:
        - name: host-dev
          hostPath:
            path: /dev
{{end}}
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: csi-attacher
  namespace: storage
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: external-attacher-runner
  namespace: storage
rules:
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["volumeattachments"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["extensions"]
    resourceNames:
    - privileged 
    resources: ["podsecuritypolicies"]
    verbs:
    - use
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: csi-attacher-role
  namespace: storage
subjects:
  - kind: ServiceAccount
    name: csi-attacher
    namespace: storage
roleRef:
  kind: ClusterRole
  name: external-attacher-runner
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
kind: Service
apiVersion: v1
metadata:
  namespace: storage
  name: csi-attacher
  labels:
    app: csi-attacher
spec:
  selector:
    app: csi-attacher
  ports:
    - name: dummy
      port: 12345
---
kind: StatefulSet
apiVersion: apps/v1beta1
metadata:
  namespace: storage
  name: csi-attacher
spec:
  serviceName: "csi-attacher"
  replicas: 1
  template:
    metadata:
      labels:
        app: csi-attacher
    spec:
      nodeSelector: 
        {{.NodeSelector}}: "true"
      serviceAccount: csi-attacher
      containers:
        - name: csi-attacher
          image: {{.StorageCSIAttacherImage}}
          args:
            - "--v=5"
            - "--csi-address=$(ADDRESS)"
          env:
            - name: ADDRESS
              value: /var/lib/kubelet/plugins/csi-lvmplugin/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /var/lib/kubelet/plugins/csi-lvmplugin
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvmplugin
            type: DirectoryOrCreate
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: csi-provisioner
  namespace: storage
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: storage
  name: external-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["list", "watch", "create", "update", "patch"]
  - apiGroups: ["extensions"]
    resourceNames:
    - privileged 
    resources: ["podsecuritypolicies"]
    verbs:
    - use
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: storage
  name: csi-provisioner-role
subjects:
  - kind: ServiceAccount
    name: csi-provisioner
    namespace: storage
roleRef:
  kind: ClusterRole
  name: external-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
kind: Service
apiVersion: v1
metadata:
  namespace: storage
  name: csi-provisioner
  labels:
    app: csi-provisioner
spec:
  selector:
    app: csi-provisioner
  ports:
    - name: dummy
      port: 12345

---
kind: StatefulSet
apiVersion: apps/v1beta1
metadata:
  namespace: storage
  name: csi-provisioner
spec:
  serviceName: "csi-provisioner"
  replicas: 1
  template:
    metadata:
      labels:
        app: csi-provisioner
    spec:
      nodeSelector: 
        {{.NodeSelector}}: "true"
      serviceAccount: csi-provisioner
      tolerations:
      - key: "node-role.kubernetes.io/master"
        operator: "Equal"
        value: ""
        effect: "NoSchedule"
      containers:
        - name: csi-provisioner
          image: {{.StorageCSIProvisionerImage}}
          args:
            - "--provisioner=csi-lvmplugin"
            - "--csi-address=$(ADDRESS)"
            - "--v=50"
            - "--logtostderr"
          env:
            - name: ADDRESS
              value: /var/lib/kubelet/plugins/csi-lvmplugin/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /var/lib/kubelet/plugins/csi-lvmplugin
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvmplugin
            type: DirectoryOrCreate
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: csi-lvmplugin
  namespace: storage
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: storage
  name: csi-lvmplugin
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "update"]
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["volumeattachments"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["extensions"]
    resourceNames:
    - privileged 
    resources: ["podsecuritypolicies"]
    verbs:
    - use
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: storage
  name: csi-lvmplugin
subjects:
  - kind: ServiceAccount
    name: csi-lvmplugin
    namespace: storage
roleRef:
  kind: ClusterRole
  name: csi-lvmplugin
  apiGroup: rbac.authorization.k8s.io          
{{- end}}
---
kind: DaemonSet
apiVersion: apps/v1beta2
metadata:
  namespace: storage
  name: csi-lvmplugin
spec:
  selector:
    matchLabels:
      app: csi-lvmplugin
  template:
    metadata:
      labels:
        app: csi-lvmplugin
    spec:
      nodeSelector: 
        {{.NodeSelector}}: "true"
      serviceAccount: csi-lvmplugin
      hostNetwork: true
      containers:
        - name: driver-registrar
          image: {{.StorageDriverRegistrarImage}}
          args:
            - "--v=5"
            - "--csi-address=$(ADDRESS)"
            - "--kubelet-registration-path=/var/lib/kubelet/plugins/csi-lvmplugin/csi.sock"
          env:
            - name: ADDRESS
              value: /var/lib/kubelet/plugins/csi-lvmplugin/csi.sock
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          volumeMounts:
            - name: socket-dir
              mountPath: /var/lib/kubelet/plugins/csi-lvmplugin
            - name: registration-dir
              mountPath: /registration/
        - name: csi-lvmplugin 
          securityContext:
            privileged: true
            capabilities:
              add: ["SYS_ADMIN"]
            allowPrivilegeEscalation: true
          image: {{.StorageCSILvmpluginImage}}
          args :
            - "--nodeid=$(NODE_ID)"
            - "--endpoint=$(CSI_ENDPOINT)"
            - "--v=5"
            - "--vgname=$(VG_NAME)"
            - "--drivername=csi-lvmplugin"
          env:
            - name: VG_NAME
              value: "k8s"
            - name: NODE_ID
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: CSI_ENDPOINT
              value: unix://var/lib/kubelet/plugins/csi-lvmplugin/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: plugin-dir
              mountPath: /var/lib/kubelet/plugins/csi-lvmplugin
            - name: pods-mount-dir
              mountPath: /var/lib/kubelet/pods
              mountPropagation: "Bidirectional"
            - mountPath: /dev
              name: host-dev
            - mountPath: /sys
              name: host-sys
            - mountPath: /lib/modules
              name: lib-modules
              readOnly: true
      volumes:
        - name: registration-dir
          hostPath:
            path: /var/lib/kubelet/plugins/
            type: DirectoryOrCreate
        - name: plugin-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvmplugin
            type: DirectoryOrCreate
        - name: pods-mount-dir
          hostPath:
            path: /var/lib/kubelet/pods
            type: Directory
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvmplugin
            type: DirectoryOrCreate
        - name: host-dev
          hostPath:
            path: /dev
        - name: host-sys
          hostPath:
            path: /sys
        - name: lib-modules
          hostPath:
            path: /lib/modules
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  name: lvm
provisioner: csi-lvmplugin
reclaimPolicy: Delete`
