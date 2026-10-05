# Tasks: Workload Labels

## Phase 1: Setup

- [ ] T001 Run the baseline quality gate in `Makefile` with `make verify-ci`

## Phase 2: User Story 1 - Apply Custom Labels (P1)

**Goal**: Apply validated custom labels without changing workload selectors.

**Independent Test**: Inspect constructed pod and VMI metadata for custom and protected labels.

- [ ] T002 [US1] Add custom-label parsing and managed-label conflict tests in `cmd/k8s-netperf/k8s-netperf_test.go`
- [ ] T003 [US1] Add pod label merge and selector-preservation tests in `pkg/k8s/kubernetes_test.go`
- [ ] T004 [US1] Add VMI label merge tests in `pkg/k8s/kubevirt_test.go`
- [ ] T005 [US1] Parse and validate repeated `--label=KEY=VALUE` values in `cmd/k8s-netperf/k8s-netperf.go`
- [ ] T006 [US1] Add custom labels to `PerfScenarios` in `pkg/config/config.go`
- [ ] T007 [US1] Merge custom labels into pod templates while preserving selector labels in `pkg/k8s/kubernetes.go`
- [ ] T008 [US1] Merge custom labels into VMI metadata while preserving selector labels in `pkg/k8s/kubevirt.go`

## Phase 3: Polish

- [ ] T009 Document repeated custom labels and protected selector labels in `docs/advanced-usage.md`
- [ ] T010 Run formatting, linting, and tests from `Makefile` with `make verify-ci`

## Dependencies

- T002-T004 precede implementation.
- T005-T008 are sequential shared-path changes.
- T009-T010 follow implementation.
