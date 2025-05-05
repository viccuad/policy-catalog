# priority-class-policy

Kubernetes allows users to define `PriorityClasses` for their pods. This helps
the control plane properly schedule workloads based on their priorities and
determine which pods should be preempted to ensure critical workloads are
always scheduled. However, users can misuse the available `PriorityClasses` by
setting improper classes.

To avoid this kind of situation, this policy ensures that the
`priorityClassName` defined in a pod is only one of the expected ones.
Therefore, cluster administrators can use this policy to enforce a set of
allowed `PriorityClass` names, rejecting deployments if users try to create a
Pod with an unexpected `priorityClassName`.

This policy can inspect `Pod` resources, but can also operate against
"higher-order" Kubernetes resources like `Deployment`, `ReplicaSet`,
`DaemonSet`, `ReplicationController`, `Job`, and `CronJob`.

## Settings

This policy setting is simply a list of the `PriorityClass` names that users
can use in the Pod's `priorityClassName` field. The user **must** provide a
list of priority classes names when deploying the policy

Here's an example of the settings:

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  annotations:
    io.kubewarden.policy.category: Resource validation
    io.kubewarden.policy.severity: medium
  name: priority-class-policy
spec:
  module: ghcr.io/kubewarden/policies/priority-class-policy:latest
  rules:
    - apiGroups:
        - ""
      apiVersions:
        - v1
      resources:
        - pods
      operations:
        - CREATE
        - UPDATE
  mutating: false
  settings:
    allowed_priority_classes:
      - low-priority
      - med-priority
      - high-priority
```

With the above configuration, any `Pod` that defines a `priorityClassName` that
is not present in the list will be rejected. If the `priorityClassName` is not
defined, the resource will be accepted.

Furthermore, the cluster operator can use this policy to enforce specific
classes based on namespace or object selectors. For example, to apply the
policy to specific namespaces, the `namespaceSelector` can be used:

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  annotations:
    io.kubewarden.policy.category: Resource validation
    io.kubewarden.policy.severity: medium
  name: priority-class-policy
spec:
  module: ghcr.io/kubewarden/policies/priority-class-policy:latest
  rules:
    - apiGroups:
        - ""
      apiVersions:
        - v1
      resources:
        - pods
      operations:
        - CREATE
        - UPDATE
  mutating: false
  namespaceSelector:
    matchExpressions:
      - key: "kubernetes.io/metadata.name"
        operator: In
        values: [team1, team2]
  settings:
    allowed_priority_classes:
      - low-priority
      - med-priority
      - high-priority
```

In the previous example, the same `priorityClassName` configuration validation
will be applied to `Pods` only deployed in the `team1` and `team2` namespaces.
