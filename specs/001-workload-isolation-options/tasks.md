---

description: "Task list for workload isolation options"
---

# Tasks: Workload Isolation Options

**Input**: Design documents from `/specs/001-workload-isolation-options/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`,
`contracts/cli-options.md`, and `quickstart.md`

**Tests**: Focused tests are required by the project constitution for new behavior. Write each test
before its corresponding implementation task.

**Organization**: Tasks are grouped by user story so each user-facing capability can be completed and
validated independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel with other marked tasks because it changes a different file and has no
  incomplete-task dependency.
- **[Story]**: Maps the task to its user story in `spec.md`.

## Phase 1: Setup

**Purpose**: Confirm the existing repository baseline before feature changes.

- [ ] T001 [P] Run and record the baseline quality result for the feature in `specs/001-workload-isolation-options/quickstart.md` using `make verify-ci`

---

## Phase 2: Foundational

**Purpose**: Add shared run-scoped option storage without changing workload behavior.

**⚠️ CRITICAL**: Complete this phase before wiring any option into workload creation.

- [X] T002 Add fields for custom annotations, pod runtime class, and VM launch security to `pkg/config/config.go`

**Checkpoint**: The scenario can carry all feature settings while existing runs remain unchanged.

---

## Phase 3: User Story 1 - Configure Isolated Workloads (Priority: P1) 🎯 MVP

**Goal**: Let operators select a runtime class for every benchmark pod and reject the setting for
VM-only runs.

**Independent Test**: Run a pod benchmark with `--runtime-class`, inspect client and server pod
templates, and confirm the selected class is present; run without the option and confirm the field is
unset.

### Tests for User Story 1

- [X] T003 [US1] Add runtime-class construction and default-preservation tests in `pkg/k8s/kubernetes_test.go`
- [X] T004 [US1] Add command-line validation tests for runtime class and pod-disabled runs in `cmd/k8s-netperf/k8s-netperf_test.go`

### Implementation for User Story 1

- [X] T005 [US1] Add runtime-class CLI parsing, pod-mode validation, and `PerfScenarios` wiring in `cmd/k8s-netperf/k8s-netperf.go`
- [X] T006 [US1] Propagate the optional runtime class through deployment parameters and pod template specs in `pkg/k8s/kubernetes.go`

**Checkpoint**: Pod benchmarks honor an installed runtime class, VM-only runs reject the option, and
default pod behavior is unchanged.

---

## Phase 4: User Story 2 - Apply Custom Workload Annotations (Priority: P2)

**Goal**: Let operators apply validated custom metadata to every created benchmark pod and VMI without
overwriting tool-managed annotations.

**Independent Test**: Run pod and VM benchmarks with repeated `--annotation=KEY=VALUE`, inspect all
created client and server resources, and confirm the annotations are present; verify malformed,
duplicate, and reserved keys fail before creation.

### Tests for User Story 2

- [X] T007 [P] [US2] Add pod annotation merge, reserved-key collision, and default-preservation tests in `pkg/k8s/kubernetes_test.go`
- [X] T008 [P] [US2] Add VMI metadata annotation tests in `pkg/k8s/kubevirt_test.go`
- [X] T009 [US2] Add repeated annotation parsing and invalid-input tests in `cmd/k8s-netperf/k8s-netperf_test.go`

### Implementation for User Story 2

- [X] T010 [US2] Parse repeated `--annotation=KEY=VALUE` values, validate empty or duplicate keys, reject tool-managed keys, and wire annotations into `PerfScenarios` in `cmd/k8s-netperf/k8s-netperf.go`
- [X] T011 [US2] Merge validated custom annotations into deployment pod-template metadata while preserving default and network annotations in `pkg/k8s/kubernetes.go`
- [X] T012 [US2] Propagate validated custom annotations through VM creation and set VMI object metadata in `pkg/k8s/kubevirt.go`

**Checkpoint**: Pod and VMI benchmarks receive custom annotations uniformly, while invalid or reserved
annotations cannot create resources or alter required network behavior.

---

## Phase 5: User Story 3 - Select VM Launch Security (Priority: P3)

**Goal**: Let VM benchmark operators select SNP or TDX launch security, with no TDX attestation
configuration.

**Independent Test**: Run a VM benchmark with each `--launch-security` value and inspect every VMI for
exactly `launchSecurity.snp: {}` or `launchSecurity.tdx: {}`; confirm no attestation field is set and
non-VM or invalid selections fail before creation.

### Dependency Update for User Story 3

- [X] T013 [US3] Upgrade the typed KubeVirt API dependency to v1.8.4 in `go.mod` and update resolved checksums in `go.sum`

### Tests for User Story 3

- [X] T014 [US3] Add SNP, TDX-without-attestation, unset default, and invalid launch-security tests in `pkg/k8s/kubevirt_test.go`
- [X] T015 [US3] Add launch-security command-line validation tests for VM-disabled and unsupported values in `cmd/k8s-netperf/k8s-netperf_test.go`

### Implementation for User Story 3

- [X] T016 [US3] Parse `--launch-security`, enforce VM-only `snp` or `tdx` values, and wire the setting into `PerfScenarios` in `cmd/k8s-netperf/k8s-netperf.go`
- [X] T017 [US3] Propagate launch security to centralized VMI construction and set exactly the typed SNP or TDX member without attestation configuration in `pkg/k8s/kubevirt.go`

**Checkpoint**: VM benchmarks request SNP or TDX exactly as selected, do not configure TDX attestation,
and surface cluster admission or scheduling failures without reporting valid benchmark results.

---

## Phase 6: Polish and Cross-Cutting Concerns

**Purpose**: Document the final command interface and run full repository and cluster validation.

- [X] T018 [P] Document the annotation, runtime-class, and launch-security options in `docs/setup.md`
- [X] T019 [P] Document KubeVirt 1.8.4 prerequisites, SNP/TDX behavior, limitations, and no-attestation scope in `docs/advanced-usage.md`
- [ ] T020 Verify all contract scenarios and expected outcomes in `specs/001-workload-isolation-options/quickstart.md`
- [X] T021 Run formatting, linting, and all unit tests through `Makefile` with `make verify-ci`

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on T001 and blocks option wiring.
- **User Story 1 (Phase 3)**: Depends on T002 and is the MVP.
- **User Story 2 (Phase 4)**: Depends on T002; it can be implemented after the MVP without modifying its
  runtime-class behavior.
- **User Story 3 (Phase 5)**: Depends on T002 and T013 before T017; it can be implemented after the MVP
  without modifying annotation or runtime-class behavior.
- **Polish (Phase 6)**: Depends on the desired user stories being complete.

### User Story Dependencies

- **US1 (P1)**: Depends only on the shared scenario fields from T002.
- **US2 (P2)**: Depends only on the shared scenario fields from T002.
- **US3 (P3)**: Depends on T002 and the KubeVirt API upgrade in T013.

### Parallel Opportunities

- T007 and T008 can proceed together because they test distinct resource builders.
- T018 and T019 can proceed together because they update distinct documentation files.
- Once T002 completes, separate contributors can prepare US1, US2, and the T013 dependency update for
  US3, but changes to `cmd/k8s-netperf/k8s-netperf.go` and `pkg/config/config.go` must be serialized.

## Parallel Example: User Story 2

```text
Task: "Add pod annotation merge tests in pkg/k8s/kubernetes_test.go"
Task: "Add VMI metadata annotation tests in pkg/k8s/kubevirt_test.go"
```

## Implementation Strategy

### MVP First

1. Complete T001-T002.
2. Complete T003-T006 for pod runtime-class support.
3. Validate the Phase 3 independent test before proceeding.

### Incremental Delivery

1. Deliver pod runtime-class support with preserved defaults.
2. Add uniform, protected custom annotations for pods and VMIs.
3. Upgrade the typed KubeVirt API and add SNP/TDX launch security without attestation.
4. Complete documentation and full repository verification.

## Phase 7: Convergence

- [X] T022 Add bounded Deployment readiness waiting, return Deployment failure conditions, and cover unavailable runtime-class scheduling in `pkg/k8s/kubernetes.go` and `pkg/k8s/kubernetes_test.go` per FR-006, US1/AC3, and SC-002 (partial)

## Phase 8: Convergence

- [X] T023 Read the named VMI before watching, watch from its resource version, and cover the already-running VMI path in `pkg/k8s/kubevirt.go` and `pkg/k8s/kubevirt_test.go` per FR-008, US3/AC1-2, and SC-001 (partial)
