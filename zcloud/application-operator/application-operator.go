package applicationoperator

const ApplicationOperatorTemplate = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: applications.app.zcloud.cn
spec:
  group: app.zcloud.cn
  names:
    kind: Application
    listKind: ApplicationList
    plural: applications
    singular: application
    shortNames:
    - app
    - apps
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Application is the Schema for the applications API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ApplicationSpec defines the desired state of Application
          properties:
            crdManifests:
              items:
                properties:
                  content:
                    type: string
                  duplicate:
                    type: boolean
                  file:
                    type: string
                type: object
              type: array
            createdByAdmin:
              type: boolean
            injectServiceMesh:
              type: boolean
            manifests:
              items:
                properties:
                  content:
                    type: string
                  duplicate:
                    type: boolean
                  file:
                    type: string
                type: object
              type: array
            ownerChart:
              description: 'INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
                Important: Run "operator-sdk generate k8s" to regenerate code after
                modifying this file Add custom validation using kubebuilder tags:
                https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html'
              properties:
                icon:
                  type: string
                name:
                  type: string
                systemChart:
                  type: boolean
                version:
                  type: string
              required:
              - icon
              - name
              - systemChart
              - version
              type: object
          required:
          - ownerChart
          type: object
        status:
          description: ApplicationStatus defines the observed state of Application
          properties:
            appResources:
              items:
                properties:
                  exists:
                    type: boolean
                  link:
                    type: string
                  name:
                    type: string
                  namespace:
                    type: string
                  readyReplicas:
                    type: integer
                  replicas:
                    type: integer
                  type:
                    type: string
                required:
                - exists
                - name
                - namespace
                - type
                type: object
              type: array
            readyWorkloadCount:
              type: integer
            state:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: string
            workloadCount:
              type: integer
          required:
          - state
          type: object
      type: object
  version: v1beta1
  versions:
  - name: v1beta1
    served: true
    storage: true
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: application-operator
  namespace: zcloud
spec:
  replicas: 1
  selector:
    matchLabels:
      app: application-operator
  template:
    metadata:
      name: application-operator
      labels:
        app: application-operator
    spec:
      serviceAccount: zcloud-cluster-admin
      containers:
      - name: application-operator
        image: {{.Image}}
`
