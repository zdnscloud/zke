package cicd

const TektonTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .DeployNamespace }}

---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: tekton-pipelines
spec:
  privileged: false
  allowPrivilegeEscalation: false
  volumes:
  - 'emptyDir'
  - 'configMap'
  - 'secret'
  hostNetwork: false
  hostIPC: false
  hostPID: false
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'MustRunAs'
    ranges:
    - min: 1
      max: 65535
  fsGroup:
    rule: 'MustRunAs'
    ranges:
    - min: 1
      max: 65535

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tekton-pipelines-admin
rules:
- apiGroups: [""]
  resources: ["pods", "pods/log", "namespaces", "secrets", "events", "serviceaccounts",
    "configmaps", "persistentvolumeclaims"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments/finalizers"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["admissionregistration.k8s.io"]
  resources: ["mutatingwebhookconfigurations"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["tekton.dev"]
  resources: ["tasks", "clustertasks", "taskruns", "pipelines", "pipelineruns", "pipelineresources",
    "conditions"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["tekton.dev"]
  resources: ["taskruns/finalizers", "pipelineruns/finalizers"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["tekton.dev"]
  resources: ["tasks/status", "clustertasks/status", "taskruns/status", "pipelines/status",
    "pipelineruns/status", "pipelineresources/status"]
  verbs: ["get", "list", "create", "update", "delete", "patch", "watch"]
- apiGroups: ["policy"]
  resources: ["podsecuritypolicies"]
  resourceNames: ["tekton-pipelines"]
  verbs: ["use"]

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tekton-pipelines-controller
  namespace: {{ .DeployNamespace }}

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: tekton-pipelines-controller-admin
subjects:
- kind: ServiceAccount
  name: tekton-pipelines-controller
  namespace: {{ .DeployNamespace }}
roleRef:
  kind: ClusterRole
  name: tekton-pipelines-admin
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: clustertasks.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: ClusterTask
    plural: clustertasks
    categories:
    - tekton
    - tekton-pipelines
  scope: Cluster
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: conditions.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: Condition
    plural: conditions
    categories:
    - tekton
    - tekton-pipelines
  scope: Namespaced
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: images.caching.internal.knative.dev
  labels:
    knative.dev/crd-install: "true"
spec:
  group: caching.internal.knative.dev
  version: v1alpha1
  names:
    kind: Image
    plural: images
    singular: image
    categories:
    - knative-internal
    - caching
    shortNames:
    - img
  scope: Namespaced
  subresources:
    status: {}

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: pipelines.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: Pipeline
    plural: pipelines
    categories:
    - tekton
    - tekton-pipelines
  scope: Namespaced
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: pipelineruns.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: PipelineRun
    plural: pipelineruns
    categories:
    - tekton
    - tekton-pipelines
    shortNames:
    - pr
    - prs
  scope: Namespaced
  additionalPrinterColumns:
  - name: Succeeded
    type: string
    JSONPath: ".status.conditions[?(@.type==\"Succeeded\")].status"
  - name: Reason
    type: string
    JSONPath: ".status.conditions[?(@.type==\"Succeeded\")].reason"
  - name: StartTime
    type: date
    JSONPath: .status.startTime
  - name: CompletionTime
    type: date
    JSONPath: .status.completionTime
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: pipelineresources.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: PipelineResource
    plural: pipelineresources
    categories:
    - tekton
    - tekton-pipelines
  scope: Namespaced
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: tasks.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: Task
    plural: tasks
    categories:
    - tekton
    - tekton-pipelines
  scope: Namespaced
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: taskruns.tekton.dev
spec:
  group: tekton.dev
  names:
    kind: TaskRun
    plural: taskruns
    categories:
    - tekton
    - tekton-pipelines
    shortNames:
    - tr
    - trs
  scope: Namespaced
  additionalPrinterColumns:
  - name: Succeeded
    type: string
    JSONPath: ".status.conditions[?(@.type==\"Succeeded\")].status"
  - name: Reason
    type: string
    JSONPath: ".status.conditions[?(@.type==\"Succeeded\")].reason"
  - name: StartTime
    type: date
    JSONPath: .status.startTime
  - name: CompletionTime
    type: date
    JSONPath: .status.completionTime
  # Opt into the status subresource so metadata.generation
  # starts to increment
  subresources:
    status: {}
  version: v1alpha1

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: tekton-pipelines-controller
  name: tekton-pipelines-controller
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - name: http-metrics
    port: 9090
    protocol: TCP
    targetPort: 9090
  selector:
    app: tekton-pipelines-controller

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: tekton-pipelines-webhook
  name: tekton-pipelines-webhook
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - name: https-webhook
    port: 443
    targetPort: 8443
  selector:
    app: tekton-pipelines-webhook

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: tekton-aggregate-edit
  labels:
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
rules:
- apiGroups:
  - tekton.dev
  resources:
  - tasks
  - taskruns
  - pipelines
  - pipelineruns
  - pipelineresources
  - conditions
  verbs:
  - create
  - delete
  - deletecollection
  - get
  - list
  - patch
  - update
  - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: tekton-aggregate-view
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
rules:
- apiGroups:
  - tekton.dev
  resources:
  - tasks
  - taskruns
  - pipelines
  - pipelineruns
  - pipelineresources
  - conditions
  verbs:
  - get
  - list
  - watch

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-artifact-bucket
  namespace: {{ .DeployNamespace }}

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-artifact-pvc
  namespace: {{ .DeployNamespace }}

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-defaults
  namespace: {{ .DeployNamespace }}
data:
  _example: |-
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # default-timeout-minutes contains the default number of
    # minutes to use for TaskRun and PipelineRun, if none is specified.
    default-timeout-minutes: "60"  # 60 minutes

    # default-service-account contains the default service account name
    # to use for TaskRun and PipelineRun, if none is specified.
    default-service-account: "default"

    # default-managed-by-label-value contains the default value given to the
    # "app.kubernetes.io/managed-by" label applied to all Pods created for
    # TaskRuns. If a user's requested TaskRun specifies another value for this
    # label, the user's request supercedes.
    default-managed-by-label-value: "tekton-pipelines"

    # default-pod-template contains the default pod template to use
    # TaskRun and PipelineRun, if none is specified. If a pod template
    # is specified, the default pod template is ignored.
    # default-pod-template:

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-logging
  namespace: {{ .DeployNamespace }}
data:
  # Common configuration for all knative codebase
  zap-logger-config: |
    {
      "level": "info",
      "development": false,
      "sampling": {
        "initial": 100,
        "thereafter": 100
      },
      "outputPaths": ["stdout"],
      "errorOutputPaths": ["stderr"],
      "encoding": "json",
      "encoderConfig": {
        "timeKey": "",
        "levelKey": "level",
        "nameKey": "logger",
        "callerKey": "caller",
        "messageKey": "msg",
        "stacktraceKey": "stacktrace",
        "lineEnding": "",
        "levelEncoder": "",
        "timeEncoder": "",
        "durationEncoder": "",
        "callerEncoder": ""
      }
    }
  # Log level overrides
  loglevel.controller: "info"
  loglevel.webhook: "info"

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-observability
  namespace: {{ .DeployNamespace }}
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################

    # This block is not actually functional configuration,
    # but serves to illustrate the available configuration
    # options and document them in a way that is accessible
    # to users that kubectl edit this config map.
    #
    # These sample configuration options may be copied out of
    # this example block and unindented to be in the data block
    # to actually change the configuration.

    # metrics.backend-destination field specifies the system metrics destination.
    # It supports either prometheus (the default) or stackdriver.
    # Note: Using Stackdriver will incur additional charges.
    metrics.backend-destination: prometheus

    # metrics.stackdriver-project-id field specifies the Stackdriver project ID. This
    # field is optional. When running on GCE, application default credentials will be
    # used and metrics will be sent to the cluster's project if this field is
    # not provided.
    metrics.stackdriver-project-id: "<your stackdriver project id>"

    # metrics.allow-stackdriver-custom-metrics indicates whether it is allowed
    # to send metrics to Stackdriver using "global" resource type and custom
    # metric type. Setting this flag to "true" could cause extra Stackdriver
    # charge.  If metrics.backend-destination is not Stackdriver, this is
    # ignored.
    metrics.allow-stackdriver-custom-metrics: "false"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tekton-pipelines-controller
  namespace: {{ .DeployNamespace }}
  labels:
    app.kubernetes.io/name: tekton-pipelines
    app.kubernetes.io/component: controller
    tekton.dev/release: "v0.10.1"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tekton-pipelines-controller
  template:
    metadata:
      annotations:
        cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
        # tekton.dev/release value replaced with inputs.params.versionTag in pipeline/tekton/publish.yaml
        tekton.dev/release: "v0.10.1"
      labels:
        app: tekton-pipelines-controller
        app.kubernetes.io/name: tekton-pipelines
        app.kubernetes.io/component: controller
    spec:
      serviceAccountName: tekton-pipelines-controller
      containers:
      - name: tekton-pipelines-controller
        image: {{.ControllerImage}}
        args: ["-kubeconfig-writer-image", "{{.KubeConfigWriterImage}}",
          "-creds-image", "{{.CredsIniterImage}}",
          "-git-image", "{{.GitIniterImage}}",
          "-nop-image", "tianon/true", "-shell-image", "busybox", "-gsutil-image",
          "google/cloud-sdk", "-entrypoint-image", "{{.EntrypointerImage}}",
          "-imagedigest-exporter-image", "{{.ImageDigestExporterImage}}",
          "-pr-image", "{{.PullRequestIniterImage}}",
          "-build-gcs-fetcher-image", "{{.GCSFetcherImage}}"]
        volumeMounts:
        - name: config-logging
          mountPath: /etc/config-logging
        env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CONFIG_LOGGING_NAME
          value: config-logging
        - name: CONFIG_OBSERVABILITY_NAME
          value: config-observability
        - name: CONFIG_ARTIFACT_BUCKET_NAME
          value: config-artifact-bucket
        - name: CONFIG_ARTIFACT_PVC_NAME
          value: config-artifact-pvc
        - name: METRICS_DOMAIN
          value: tekton.dev/pipeline
      volumes:
      - name: config-logging
        configMap:
          name: config-logging

---
apiVersion: apps/v1
kind: Deployment
metadata:
  # Note: the Deployment name must be the same as the Service name specified in
  # config/400-webhook-service.yaml. If you change this name, you must also
  # change the value of WEBHOOK_SERVICE_NAME below.
  name: tekton-pipelines-webhook
  namespace: {{ .DeployNamespace }}
  labels:
    app.kubernetes.io/name: tekton-pipelines
    app.kubernetes.io/component: webhook-controller
    tekton.dev/release: "v0.10.1"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tekton-pipelines-webhook
  template:
    metadata:
      annotations:
        cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
        tekton.dev/release: "v0.10.1"
      labels:
        app: tekton-pipelines-webhook
        app.kubernetes.io/name: tekton-pipelines
        app.kubernetes.io/component: webhook-controller
    spec:
      serviceAccountName: tekton-pipelines-controller
      containers:
      - name: webhook
        # This is the Go import path for the binary that is containerized
        # and substituted here.
        image: {{.WebhookImage}}
        volumeMounts:
        - name: config-logging
          mountPath: /etc/config-logging
        env:
        - name: SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: WEBHOOK_SERVICE_NAME
          value: tekton-pipelines-webhook
      volumes:
      - name: config-logging
        configMap:
          name: config-logging`

const TektonDashBoardTemplate = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: extensions.dashboard.tekton.dev
spec:
  group: dashboard.tekton.dev
  names:
    categories:
    - tekton
    - tekton-dashboard
    kind: Extension
    plural: extensions
  scope: Namespaced
  subresources:
    status: {}
  version: v1alpha1
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: tekton-dashboard
  name: tekton-dashboard
  namespace: {{ .DeployNamespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: tekton-dashboard-minimal
  namespace: {{ .DeployNamespace }}
rules:
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - pods
  - services
  verbs:
  - create
  - update
  - delete
  - patch
- apiGroups:
  - ""
  resources:
  - secrets
  - configmaps
  verbs:
  - create
  - update
  - delete
- apiGroups:
  - extensions
  - apps
  resources:
  - deployments
  verbs:
  - create
  - update
  - delete
  - patch
- apiGroups:
  - tekton.dev
  resources:
  - tasks
  - clustertasks
  - taskruns
  - pipelines
  - pipelineruns
  - pipelineresources
  - conditions
  verbs:
  - create
  - update
  - delete
  - patch
- apiGroups:
  - tekton.dev
  resources:
  - taskruns/finalizers
  - pipelineruns/finalizers
  verbs:
  - create
  - update
  - delete
  - patch
- apiGroups:
  - tekton.dev
  resources:
  - tasks/status
  - clustertasks/status
  - taskruns/status
  - pipelines/status
  - pipelineruns/status
  verbs:
  - create
  - update
  - delete
  - patch
- apiGroups:
  - dashboard.tekton.dev
  resources:
  - extensions
  verbs:
  - create
  - update
  - delete
  - patch
- apiGroups:
  - tekton.dev
  resources:
  - eventlisteners
  - triggerbindings
  - triggertemplates
  verbs:
  - create
  - update
  - delete
  - patch
  - add
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  - get
  - list
- apiGroups:
  - security.openshift.io
  resources:
  - securitycontextconstraints
  verbs:
  - use
- apiGroups:
  - route.openshift.io
  resources:
  - routes
  verbs:
  - get
  - list
- apiGroups:
  - extensions
  - apps
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - get
  - list
- apiGroups:
  - ""
  resources:
  - pods
  - services
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - pods/log
  - namespaces
  - events
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - secrets
  - configmaps
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - extensions
  - apps
  resources:
  - deployments
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - tekton.dev
  resources:
  - tasks
  - clustertasks
  - taskruns
  - pipelines
  - pipelineruns
  - pipelineresources
  - conditions
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - tekton.dev
  resources:
  - taskruns/finalizers
  - pipelineruns/finalizers
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - tekton.dev
  resources:
  - tasks/status
  - clustertasks/status
  - taskruns/status
  - pipelines/status
  - pipelineruns/status
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - dashboard.tekton.dev
  resources:
  - extensions
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - tekton.dev
  resources:
  - eventlisteners
  - triggerbindings
  - triggertemplates
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tekton-dashboard-minimal
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: tekton-dashboard-minimal
subjects:
- kind: ServiceAccount
  name: tekton-dashboard
  namespace: {{ .DeployNamespace }}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: tekton-dashboard
  name: tekton-dashboard
  namespace: {{ .DeployNamespace }}
spec:
  ports:
  - name: http
    port: 9097
    protocol: TCP
    targetPort: 9097
  selector:
    app: tekton-dashboard
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: tekton-dashboard
    version: v0.5.1
  name: tekton-dashboard
  namespace: {{ .DeployNamespace }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tekton-dashboard
  template:
    metadata:
      labels:
        app: tekton-dashboard
      name: tekton-dashboard
    spec:
      containers:
      - env:
        - name: PORT
          value: "9097"
        - name: WEB_RESOURCES_DIR
          value: /var/run/ko/web
        - name: PIPELINE_RUN_SERVICE_ACCOUNT
          value: ""
        - name: INSTALLED_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: {{.DashboardImage}}
        livenessProbe:
          httpGet:
            path: /health
            port: 9097
        name: tekton-dashboard
        ports:
        - containerPort: 9097
        readinessProbe:
          httpGet:
            path: /readiness
            port: 9097
      serviceAccountName: tekton-dashboard
---
apiVersion: tekton.dev/v1alpha1
kind: Pipeline
metadata:
  name: pipeline0
  namespace: {{ .DeployNamespace }}
spec:
  params:
  - default: /workspace/git-source
    description: The path to the resource files to apply
    name: pathToResourceFiles
    type: string
  - default: .
    description: The directory from which resources are to be applied
    name: apply-directory
    type: string
  - default: tekton-pipelines
    description: The namespace in which to create the resources being imported
    name: target-namespace
    type: string
  resources:
  - name: git-source
    type: git
  tasks:
  - name: pipeline0-task
    params:
    - name: pathToResourceFiles
      value: $(params.pathToResourceFiles)
    - name: apply-directory
      value: $(params.apply-directory)
    - name: target-namespace
      value: $(params.target-namespace)
    resources:
      inputs:
      - name: git-source
        resource: git-source
    taskRef:
      name: pipeline0-task
---
apiVersion: tekton.dev/v1alpha1
kind: Task
metadata:
  name: pipeline0-task
  namespace: {{ .DeployNamespace }}
spec:
  inputs:
    params:
    - default: /workspace/git-source
      description: The path to the resource files to apply
      name: pathToResourceFiles
      type: string
    - default: .
      description: The directory from which resources are to be applied
      name: apply-directory
      type: string
    - default: tekton-pipelines
      description: The namespace where created resources will go
      name: target-namespace
      type: string
    resources:
    - name: git-source
      type: git
  steps:
  - args:
    - apply
    - -f
    - $(inputs.params.pathToResourceFiles)/$(inputs.params.apply-directory)
    - -n
    - $(inputs.params.target-namespace)
    command:
    - kubectl
    image: lachlanevenson/k8s-kubectl@sha256:6944392a5ab6f762addbe3b3d8fcca7f47abdc44ff6076a6d18dc08510597e30
    name: kubectl-apply
`
