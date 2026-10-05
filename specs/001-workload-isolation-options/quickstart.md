# Quickstart: Validate Workload Isolation Options

## Prerequisites

- Access to a Kubernetes or OpenShift cluster and a namespace/service account usable by k8s-netperf.
- For VM scenarios, KubeVirt 1.8.4 installed with the selected launch security capability enabled and
  nodes capable of the selected mode.
- For pod runtime-class scenarios, an installed runtime class usable by the benchmark service account.
- A built binary: run `make build` from the repository root.

Refer to the [CLI contract](contracts/cli-options.md) for complete option and validation rules and the
[data model](data-model.md) for workload applicability.

## Pod Annotation and Runtime-Class Scenario

1. Run a pod benchmark with one or more `--annotation` values and `--runtime-class`.
2. Inspect the created client and server pods.
3. Confirm every pod template-derived workload has each annotation and the selected runtime class.
4. Run the same benchmark without these options and confirm the workload metadata and scheduling behavior
   match the previous default behavior.

Expected result: the configured pod settings are present only in the first run; benchmark result
collection continues normally.

## VM Annotation and SNP Scenario

1. Run a VM benchmark with `--annotation` and `--launch-security=snp`.
2. Inspect each created VirtualMachineInstance.
3. Confirm each VMI has the custom metadata and exactly `launchSecurity.snp` selected.
4. Confirm no TDX attestation configuration is present.

Expected result: VMIs request SNP and progress normally when cluster feature gates, firmware, and nodes
support it. Cluster admission errors are returned to the operator.

## VM TDX Scenario

1. Run a VM benchmark with `--launch-security=tdx`.
2. Inspect each created VirtualMachineInstance.
3. Confirm each VMI contains `launchSecurity.tdx: {}` and no attestation configuration.
4. Confirm no `launchSecurity.snp` member is set.

Expected result: VMIs request TDX without enabling attestation. Cluster admission errors are returned to
the operator; confirm KubeVirt, feature-gate, and node capability prerequisites before the run.

## Negative Validation Scenarios

1. Supply an empty, malformed, duplicate, or managed `--annotation` key and confirm resource creation is
   prevented with an actionable error.
2. Run with `--runtime-class` while pods are disabled and confirm validation fails before creation.
3. Run with `--launch-security` while VMs are disabled, or with a value other than `snp` or `tdx`, and
   confirm validation fails before creation.

## Repository Verification

Run `make verify-ci` after implementation. It must complete successfully, including formatting, linting,
