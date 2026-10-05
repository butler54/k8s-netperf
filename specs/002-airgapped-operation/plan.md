# Implementation Plan: Internet-Disconnected Workflows

**Branch**: `002-airgapped-operation` | **Date**: 2026-10-05 | **Spec**: [spec.md](spec.md)

## Summary

Enable disconnected benchmark runs by accepting complete operator-provided image references for pod and
VM workloads, and by booting VM benchmarks from a ready, operator-managed source disk. Preserve all
current defaults when these options are absent. The implementation validates source readiness, makes one
benchmark-owned disk copy per VM, verifies preinstalled benchmark prerequisites, and avoids all startup
downloads in prebuilt-disk mode.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Cobra, Kubernetes client-go v0.35.2, KubeVirt API v1.8.4, CDI DataVolume API

**Storage**: An existing operator-managed source DataVolume and temporary benchmark-owned clone
DataVolumes

**Testing**: Focused Go unit tests, `make verify-ci`, and the disconnected-workflow quickstart validation

**Target Platform**: Kubernetes/OpenShift clusters with KubeVirt and CDI for VM workflows

**Project Type**: Command-line Kubernetes network benchmark tool

**Performance Goals**: Preserve current image selection and VM startup behavior when disconnected-workflow
options are absent.

**Constraints**: Image overrides are complete references used verbatim. Prebuilt-disk mode performs no
package, source, or build-artifact downloads. Source DataVolumes may be in any namespace and must be
retained; only explicitly benchmark-owned per-VM clones are temporary. Disk construction and upload are
explicitly out of scope.

**Scale/Scope**: One clone per benchmark VMI. No registry provisioning, registry authentication, image
rewriting, VM disk construction, or VM disk upload.

## Constitution Check

| Principle | Compliance | Status |
|-----------|------------|--------|
| Benchmark Integrity | The feature changes workload source selection only; it does not change benchmark definitions, measurements, or result evaluation. Validation will prove default behavior remains unchanged. | PASS |
| Kubernetes-Native Compatibility | The design uses supported workload image fields and existing CDI/KubeVirt resources, preserving default behavior when options are omitted. | PASS |
| Verified Go Changes | Focused CLI, workload-builder, and source-disk lifecycle tests will be added. `gofmt -s`, `golangci-lint`, and relevant tests are required before review. | PASS |
| Observable and Actionable Results | Invalid image references, unavailable source disks, and missing prerequisites are surfaced before a benchmark is reported as successful. User-facing options are documented. | PASS |
| Minimal, Repository-Consistent Changes | The design extends existing CLI, scenario, pod, and VMI construction paths; per-VM disk copies and explicit ownership identification are necessary to avoid shared-disk attachment conflicts and protect source disks. | PASS |

**Post-design re-check**: PASS. The Phase 1 contract introduces no new configuration format, service, or
external integration. The source disk lifecycle and explicit clone ownership identification are the minimum
required to preserve a reusable operator asset in any namespace while allowing multiple benchmark VMIs.

## Project Structure

### Documentation (this feature)

```text
specs/002-airgapped-operation/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/
    └── cli-options.md
```

### Source Code (repository root)

```text
cmd/k8s-netperf/k8s-netperf.go   # CLI options and validation
pkg/config/config.go             # Benchmark scenario settings
pkg/k8s/kubernetes.go            # Pod image selection and source-disk lifecycle
pkg/k8s/kubevirt.go              # VMI disk selection and prebuilt-disk startup
pkg/k8s/*_test.go                # Focused unit coverage
docs/setup.md                    # Image override documentation
docs/advanced-usage.md           # Prebuilt-disk requirements and workflow
```

**Structure Decision**: Extend the existing CLI-to-scenario-to-workload-builder flow. No new project
module or configuration file is needed.

## Complexity Tracking

No constitution violations or complexity exceptions require justification.
