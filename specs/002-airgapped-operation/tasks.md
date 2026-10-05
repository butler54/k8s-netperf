# Tasks: Internet-Disconnected Workflows

**Input**: Design documents from `specs/002-airgapped-operation/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/cli-options.md`, and `quickstart.md`

**Tests**: Focused Go unit tests are required by the project constitution for changed behavior. Complete
`make verify-ci` before review.

**Organization**: Tasks are grouped by user story so that image overrides and prebuilt VM disks can be
implemented and validated as separate increments.

## Phase 1: Setup

**Purpose**: Establish the implementation and test surfaces for the revised specification.

- [X] T001 Review the CLI contract and existing flag-registration patterns in `cmd/k8s-netperf/k8s-netperf.go` and `specs/002-airgapped-operation/contracts/cli-options.md`
- [X] T002 Review workload construction and cleanup paths in `pkg/k8s/kubernetes.go` and `pkg/k8s/kubevirt.go`
- [X] T003 [P] Review the existing focused test conventions in `cmd/k8s-netperf/k8s-netperf_test.go`, `pkg/k8s/kubernetes_test.go`, and `pkg/k8s/kubevirt_test.go`

---

## Phase 2: Foundational Prerequisites

**Purpose**: Add shared scenario state and ownership conventions used by both implementation increments.

- [X] T004 Add scenario fields for operator-provided pod images, VM images, source DataVolumes, and run ownership in `pkg/config/config.go`
- [X] T005 Add shared ownership marker constants and clone-name generation in `pkg/k8s/kubernetes.go`
- [X] T006 Add unit coverage for clone ownership identity and names in `pkg/k8s/kubernetes_test.go`

**Checkpoint**: Scenario state and clone ownership conventions are available for both user stories.

---

## Phase 3: User Story 1 - Use Operator-Provided Workload Images (Priority: P1) 🎯 MVP

**Goal**: Let operators provide complete private image references for pod workloads and online VM
container-disk workloads without changing defaults when no override is supplied.

**Independent Test**: Start pod and online VM benchmark setup with supplied image references and verify
created workload definitions use those exact references; repeat with no overrides and verify defaults.

### Tests for User Story 1

- [X] T007 [P] [US1] Add CLI validation tests for empty and valid `--image` and `--vm-image` values in `cmd/k8s-netperf/k8s-netperf_test.go`
- [X] T008 [P] [US1] Add pod workload image-selection tests for defaults and overrides in `pkg/k8s/kubernetes_test.go`
- [X] T009 [P] [US1] Add VMI container-disk image-selection tests in `pkg/k8s/kubevirt_test.go`

### Implementation for User Story 1

- [X] T010 [US1] Add `--image` validation and scenario wiring in `cmd/k8s-netperf/k8s-netperf.go`
- [X] T011 [US1] Use the selected pod image for every benchmark pod workload in `pkg/k8s/kubernetes.go`
- [X] T012 [US1] Preserve and apply the selected online VM image to VMI construction in `pkg/k8s/kubevirt.go`
- [X] T013 [US1] Document pod and online VM image overrides, exact-reference behavior, and no-fallback failures in `docs/setup.md`

**Checkpoint**: Private image references work for both pod and online VM workloads, while invocations
without overrides retain their current images.

---

## Phase 4: User Story 2 - Run VMs from a Prebuilt Disk (Priority: P2)

**Goal**: Run disconnected VM benchmarks from a ready source DataVolume, creating per-VMI owned clones
and preserving the source in any namespace during cleanup.

**Independent Test**: Select a ready source DataVolume, create client and server VMIs, verify each uses
its own owned clone and no startup downloads occur, then verify cleanup retains the source and removes
only owned clones.

### Tests for User Story 2

- [X] T014 [P] [US2] Add CLI tests for `--offline-data-volume`, VM-only validation, and conflict with `--vm-image` in `cmd/k8s-netperf/k8s-netperf_test.go`
- [X] T015 [P] [US2] Add source readiness, clone creation, reuse rejection, and selective-cleanup tests in `pkg/k8s/kubernetes_test.go`
- [X] T016 [P] [US2] Add VMI prebuilt-disk attachment and download-free startup tests in `pkg/k8s/kubevirt_test.go`
- [X] T017 [P] [US2] Add missing local netperf prerequisite reporting tests in `pkg/k8s/kubevirt_test.go`

### Implementation for User Story 2

- [X] T018 [US2] Add `--offline-data-volume` parsing, DataVolume readiness validation, and VM-image conflict handling in `cmd/k8s-netperf/k8s-netperf.go`
- [X] T019 [US2] Implement source DataVolume validation and one owned clone per VMI in `pkg/k8s/kubernetes.go`
- [X] T020 [US2] Update cleanup to preserve the selected source DataVolume in any namespace and delete only benchmark-owned clones in `pkg/k8s/kubernetes.go`
- [X] T021 [US2] Attach the owned clone as the VMI disk and generate prebuilt-disk cloud-init without package, source, or build downloads in `pkg/k8s/kubevirt.go`
- [X] T022 [US2] Validate required locally installed client and server benchmark tools before running VM workloads in `pkg/k8s/kubevirt.go` and `pkg/drivers/netperf.go`
- [X] T023 [US2] Document source DataVolume requirements, no-download behavior, namespace-independent preservation, and out-of-scope disk preparation in `docs/advanced-usage.md`

**Checkpoint**: Disconnected VM benchmarks use a ready source DataVolume safely, fail clearly for missing
prerequisites, and preserve source data during cleanup.

---

## Phase 5: Polish and Cross-Cutting Validation

**Purpose**: Verify compatibility, documentation, and quality gates across both stories.

- [X] T024 [P] Reconcile flag reference output and option descriptions with the CLI contract in `docs/setup.md` and `cmd/k8s-netperf/k8s-netperf.go`
- [X] T025 [P] Review error messages for supplied image references, source DataVolumes, and missing prerequisites in `cmd/k8s-netperf/k8s-netperf.go` and `pkg/k8s/kubernetes.go`
- [ ] T026 Run the quickstart scenarios and record results against `specs/002-airgapped-operation/quickstart.md`
- [ ] T027 Run formatting, static analysis, and unit tests with `Makefile` target `verify-ci`

---

## Dependencies and Execution Order

```text
Phase 1 (T001–T003)
        ↓
Phase 2 (T004–T006)
        ├── User Story 1 / MVP (T007–T013)
        └── User Story 2 (T014–T023)
                    ↓
        Phase 5 (T024–T027)
```

- **User Story 1** depends on T004 and can be delivered as the MVP.
- **User Story 2** depends on T004–T006; it does not depend on User Story 1 behavior.
- **Polish** begins after the desired story increments are complete.

## Parallel Opportunities

- T003 can run in parallel with T001–T002.
- T007–T009 can run in parallel once the shared scenario fields are in place.
- T014–T017 can run in parallel once clone ownership conventions are in place.
- T024–T025 can run in parallel after implementation is complete.

## Parallel Example: User Story 1

```text
T007: CLI validation tests in cmd/k8s-netperf/k8s-netperf_test.go
T008: Pod image-selection tests in pkg/k8s/kubernetes_test.go
T009: VMI image-selection tests in pkg/k8s/kubevirt_test.go
```

## Parallel Example: User Story 2

```text
T014: CLI offline-mode tests in cmd/k8s-netperf/k8s-netperf_test.go
T015: DataVolume lifecycle tests in pkg/k8s/kubernetes_test.go
T016: VMI disk and cloud-init tests in pkg/k8s/kubevirt_test.go
T017: Driver prerequisite tests in pkg/drivers/netperf_test.go
```

## Implementation Strategy

### MVP First

1. Complete setup and foundational ownership tasks.
2. Complete User Story 1 (T007–T013).
3. Verify supplied pod and online VM image references and default behavior independently.

### Incremental Delivery

1. Deliver image overrides without changing existing defaults.
2. Add prebuilt DataVolume support with selective cleanup and no-download startup.
3. Complete quickstart validation and `verify-ci` before review.

## Format Validation

Every implementation task uses the required checklist format: checkbox, sequential task ID, optional
parallel marker, required user-story label within story phases, and an exact file path.
