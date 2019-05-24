package lvm

const LVMDTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: storage-agent-lvmd
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-lvmd-runner
  namespace: {{.StorageNamespace}}
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: storage-agent-lvmd-role
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: storage-agent-lvmd
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: storage-agent-lvmd-runner
  apiGroup: rbac.authorization.k8s.io
{{- end}}
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: storage-agent-lvmd
  namespace: {{.StorageNamespace}}
  labels:
    app: storage-agent-lvmd
spec:
  selector:
    matchLabels:
      app: storage-agent-lvmd
  template:
    metadata:
      labels:
        app: storage-agent-lvmd
    spec:
      serviceAccount: storage-agent-lvmd
      hostNetwork: true
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      containers:
      - name: storage-agent-lvmd
        image: {{.StorageLvmdImage}}
        command: ["/lvmd.sh"]
        env:
          - name: MOUNT_PATH
            value: "/host/dev"
          - name: VG_NAME
            value: "k8s"
          - name: NodeName
            valueFrom:
              fieldRef: 
                fieldPath: spec.nodeName
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
            path: /dev`

const LVMStorageTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.StorageNamespace}}
{{- if eq .RBACConfig "rbac"}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: csi-lvmplugin-attacher
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: external-attacher-runner
  namespace: {{.StorageNamespace}}
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
  name: csi-lvmplugin-attacher-role
  namespace: {{.StorageNamespace}}
subjects:
  - kind: ServiceAccount
    name: csi-lvmplugin-attacher
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: external-attacher-runner
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: csi-lvmplugin-provisioner
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: external-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
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
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin-provisioner-role
subjects:
  - kind: ServiceAccount
    name: csi-lvmplugin-provisioner
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: external-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: csi-lvmplugin
  namespace: {{.StorageNamespace}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "update", "watch"]
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["list", "watch"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["list", "watch"]
  - apiGroups: ["apps"]
    resources: ["statefulsets"]
    verbs: ["list", "watch"]
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
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin
subjects:
  - kind: ServiceAccount
    name: csi-lvmplugin
    namespace: {{.StorageNamespace}}
roleRef:
  kind: ClusterRole
  name: csi-lvmplugin
  apiGroup: rbac.authorization.k8s.io  
{{- end}}
---
kind: DaemonSet
apiVersion: apps/v1
metadata:
  namespace: {{.StorageNamespace}}
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
        {{.LabelKey}}: {{.LabelValue}}
      serviceAccount: csi-lvmplugin
      hostNetwork: true
      containers:
        - name: csi-lvmplugin-driver-registrar
          image: {{.StorageLvmDriverRegistrarImage}}
          args:
            - "--v=5"
            - "--csi-address=$(ADDRESS)"
            - "--kubelet-registration-path=/var/lib/kubelet/plugins/csi-lvm/csi.sock"
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "rm -rf /registration/ /csi/"]
          env:
            - name: ADDRESS
              value: /csi/csi.sock
          volumeMounts:
            - name: plugin-dir
              mountPath: /csi
            - name: registration-dir
              mountPath: /registration
        - name: csi-lvmplugin
          securityContext:
            privileged: true
            capabilities:
              add: ["SYS_ADMIN"]
            allowPrivilegeEscalation: true
          image: {{.StorageLvmCSIImage}}
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
              value: unix://csi/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: plugin-dir
              mountPath: /csi
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
            path: /var/lib/kubelet/plugins_registry/
            type: DirectoryOrCreate
        - name: pods-mount-dir
          hostPath:
            path: /var/lib/kubelet/pods
            type: Directory
        - name: plugin-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvm/
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
kind: Service
apiVersion: v1
metadata:
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin-attacher
  labels:
    app: csi-lvmplugin-attacher
spec:
  selector:
    app: csi-lvmplugin-attacher
  ports:
    - name: dummy
      port: 12345
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin-attacher
spec:
  serviceName: csi-lvmplugin-attacher
  replicas: 1
  selector:
    matchLabels:
      app: csi-lvmplugin-attacher
  template:
    metadata:
      labels:
        app: csi-lvmplugin-attacher
    spec:
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      serviceAccount: csi-lvmplugin-attacher
      hostNetwork: true
      containers:
        - name: csi-lvmplugin-attacher
          image: {{.StorageLvmAttacherImage}}
          args:
            - "--v=5"
            - "--csi-address=$(ADDRESS)"
          env:
            - name: ADDRESS
              value: /csi/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /csi
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvm
            type: DirectoryOrCreate
---
kind: Service
apiVersion: v1
metadata:
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin-provisioner
  labels:
    app: csi-lvmplugin-provisioner
spec:
  selector:
    app: csi-lvmplugin-provisioner
  ports:
    - name: dummy
      port: 12345

---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  namespace: {{.StorageNamespace}}
  name: csi-lvmplugin-provisioner
spec:
  serviceName: csi-lvmplugin-provisioner
  replicas: 1
  selector:
    matchLabels:
      app: csi-lvmplugin-provisioner
  template:
    metadata:
      labels:
        app: csi-lvmplugin-provisioner
    spec:
      nodeSelector: 
        {{.LabelKey}}: {{.LabelValue}}
      serviceAccount: csi-lvmplugin-provisioner
      hostNetwork: true
      containers:
        - name: csi-lvmplugin-provisioner
          image: {{.StorageLvmProvisionerImage}}
          args:
            - "--provisioner=csi-lvmplugin"
            - "--csi-address=$(ADDRESS)"
            - "--v=50"
            - "--logtostderr"
            - "--feature-gates=Topology=true"
          env:
            - name: ADDRESS
              value: /csi/csi.sock
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: socket-dir
              mountPath: /csi
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-lvm
            type: DirectoryOrCreate
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  name: {{.StorageClassName}}
provisioner: csi-lvmplugin
reclaimPolicy: Delete`
