# Research: Workload Isolation Options

## CLI Input and Propagation

**Decision**: Add command-line-only options and propagate their parsed values through
`config.PerfScenarios` to the existing Deployment and VMI builders.

**Rationale**: Existing VM and pod behavior is selected through Cobra command-line options and passed
into `PerfScenarios` in `cmd/k8s-netperf/k8s-netperf.go`. This avoids creating a second configuration
model for a run-scoped feature.

**Alternatives considered**:

- Add fields to the benchmark YAML: rejected because the existing YAML describes test profiles, while
  workload behavior is controlled by CLI options.
- Add a separate workload configuration file: rejected because it adds indirection without a required
  reuse or persistence need.

## Annotation Placement and Precedence

**Decision**: Use one repeated annotation option for all benchmark resources. Apply it to pod template
metadata for pod workloads and VMI object metadata for VM workloads. Reject keys owned by the tool.

**Rationale**: `CreateDeployment` centralizes pod template annotations, and `CreateVMI` centralizes VMI
metadata. A single option applies uniformly to benchmark client and server resources while retaining the
existing resource-specific creation paths. Rejecting collisions protects the sidecar and network
annotations the tool requires.

**Alternatives considered**:

- Separate pod and VM annotation options: rejected because one shared user intent can be applied at the
  appropriate resource layer without duplicate interface surface.
- Let custom values overwrite managed annotations: rejected because it can silently remove required
  network configuration and invalidate benchmark results.

## Runtime Class Scope

**Decision**: Support an optional runtime class for pod benchmarks only. Reject it when pods are
disabled.

**Rationale**: The specification clarifies that runtime class is pod-only. Kubernetes provides this on
the pod specification, and the pinned KubeVirt VMI model does not expose an equivalent VMI runtime class
field.

**Alternatives considered**:

- Apply runtime class to VM benchmarks: rejected by the clarified scope and unavailable in the current
  typed VMI model.
- Silently ignore it for VM-only runs: rejected because no-op settings make benchmark conditions unclear.

## VM Launch Security

**Decision**: Upgrade `kubevirt.io/api` from v1.2.2 to v1.8.4 and set exactly one typed launch-security
member: `SNP` for `snp` or `TDX` for `tdx`. Do not set an attestation field or add TDX attestation
configuration.

**Rationale**: The target KubeVirt 1.8.4 release supports TDX. The stable v1.8.4 Go API exposes
`LaunchSecurity.TDX` and `LaunchSecurity.SNP`. The v1.8.4
module requires Go 1.24, which is compatible with this repository's Go 1.25 baseline.

**Alternatives considered**:

- Keep KubeVirt API v1.2.2: rejected because it cannot build a typed TDX launch security request.
- Use an unstructured VMI patch: rejected because it bypasses the repository's typed VMI construction
  and could diverge from the installed CRD schema.
- Configure TDX attestation: rejected by explicit requirement; `tdx` is emitted without attestation
  configuration.

## Validation and Failure Reporting

**Decision**: Validate annotation syntax, reserved annotation collisions, workload-mode compatibility,
and the `snp`/`tdx` launch-security enum before resource creation. Let cluster admission and scheduling
errors remain visible to the operator.

**Rationale**: Existing CLI validation occurs before workload creation. Local validation gives immediate,
actionable failures; the cluster remains authoritative for installed runtime classes, KubeVirt feature
gates, and node hardware availability.

**Alternatives considered**:

- Discover all cluster capabilities before each run: rejected because it expands scope and cannot fully
  replace admission and scheduling validation.
- Treat unsupported cluster settings as successful empty runs: rejected because it would misrepresent
  benchmark conditions.
