# Tasks: Kata Pod Networking

## Phase 1: Setup

- [ ] T001 [P] Run the baseline quality gate documented in `Makefile` with `make verify-ci`

## Phase 2: User Story 1 - Kata Pod Networking (P1)

**Goal**: Apply runtime classes only to pod-network workloads.

**Independent Test**: Construct mixed pod- and host-network deployments and inspect runtime class fields.

- [X] T002 [US1] Add mixed-scenario and host-network runtime-class tests in `pkg/k8s/kubernetes_test.go`
- [X] T003 [US1] Add host-network-only runtime-class validation tests in `cmd/k8s-netperf/k8s-netperf_test.go`
- [X] T004 [US1] Reject `--runtime-class` with `--hostNet` in `cmd/k8s-netperf/k8s-netperf.go`
- [X] T005 [US1] Clear runtime class for host-network Deployment parameters in `pkg/k8s/kubernetes.go`

## Phase 3: Polish

- [X] T006 [P] Document runtime-class host-network restriction in `docs/advanced-usage.md`
- [X] T007 Run all checks in `Makefile` with `make verify-ci`

## Dependencies

- T002-T005 implement the single P1 story in order; tests precede code.
- T006-T007 follow implementation.
