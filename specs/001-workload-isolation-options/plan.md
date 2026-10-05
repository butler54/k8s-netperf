# Implementation Plan: Workload Isolation Options

**Branch**: `001-workload-isolation-options` | **Date**: 2026-09-14 | **Spec**:
[spec.md](spec.md)

**Input**: Feature specification from `/specs/001-workload-isolation-options/spec.md`

## Summary

Add command-line workload options for repeated custom annotations, a pod-only runtime class, and VM
launch security (`snp` or `tdx`). Thread these values through the existing scenario structures to the
centralized Deployment and VMI builders. Upgrade the KubeVirt API dependency to the stable v1.8 API
release so the typed VMI model can express TDX, without enabling or configuring attestation.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Cobra CLI, Kubernetes API/client-go v0.35.2, KubeVirt API v1.8.4

**Storage**: N/A; settings exist only for the active benchmark run.

**Testing**: Go unit tests with `go test -v ./...`; `gofmt -s` and `golangci-lint` through
`make verify-ci`

**Target Platform**: Kubernetes and OpenShift clusters; KubeVirt 1.8.4 clusters for VM benchmarks

**Project Type**: Command-line Kubernetes network benchmark tool

**Performance Goals**: Preserve the existing benchmark execution and result-collection behavior when
all new options are omitted.

**Constraints**: Settings are command-line options only. Runtime class applies to pod workloads only.
VM launch security accepts only `snp` and `tdx` and MUST NOT configure TDX attestation. Tool-managed
annotations remain protected from custom overrides.

**Scale/Scope**: One set of settings applies uniformly to every client and server workload created by a
single invocation; no persisted configuration, new services, or cluster capability provisioning.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Plan Compliance | Status |
|-----------|-----------------|--------|
| Benchmark Integrity | Settings modify only created workload metadata and scheduling. Existing measurements and result evaluation remain unchanged; unsupported settings fail before workload creation. | PASS |
| Kubernetes-Native Compatibility | Uses typed Kubernetes and KubeVirt workload fields and preserves defaults when options are absent. KubeVirt API is aligned to the target 1.8 release. | PASS |
| Verified Go Changes | The implementation includes focused construction and validation tests and ends with `make verify-ci`. | PASS |
| Observable and Actionable Results | Invalid input and unsupported cluster settings produce explicit errors; new flags are documented. | PASS |
| Minimal, Repository-Consistent Changes | Extends existing CLI-to-`PerfScenarios`-to-builder flow; no new configuration system, service, or unstructured VMI patching. | PASS |

**Post-design re-check**: PASS. The contract confines the change to existing CLI, scenario, Deployment,
no attestation mechanism or new abstraction is introduced.

## Project Structure

### Documentation (this feature)

```text
specs/001-workload-isolation-options/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-options.md
└── tasks.md                 # Created by /speckit.tasks
```

### Source Code (repository root)

```text
cmd/k8s-netperf/
├── k8s-netperf.go           # CLI flags, validation, and PerfScenarios construction
└── k8s-netperf_test.go      # CLI validation tests

pkg/config/
└── config.go                # PerfScenarios settings

pkg/k8s/
├── kubernetes.go            # Deployment parameters and pod template construction
└── kubevirt.go              # VMI parameter propagation and typed VMI construction

docs/
├── setup.md                 # Command-line option reference
└── advanced-usage.md        # Isolation and confidential VM prerequisites
```

**Structure Decision**: Extend the existing command, scenario, and centralized workload builders.
Keep parsing and validation beside existing CLI options, and keep resource-specific construction in the
existing Kubernetes and KubeVirt packages.

## Complexity Tracking

No constitution violations require justification.
