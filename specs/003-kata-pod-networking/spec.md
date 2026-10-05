# Feature Specification: Kata Pod Networking

**Feature Branch**: `003-kata-pod-networking`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "When testing with Kata containers, do not use Kata containers for host
networking; use them only for pod networking."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run Kata Pod Networking Tests (Priority: P1)

As a benchmark operator, I can select a Kata runtime class and run pod-network benchmarks without
creating host-network workloads using that runtime class.

**Why this priority**: Host networking does not represent the intended Kata isolation path and would make
results misleading.

**Independent Test**: Run a benchmark with a Kata runtime class and all scenarios enabled; inspect
created workloads and confirm every runtime-class workload uses pod networking.

**Acceptance Scenarios**:

1. **Given** a Kata runtime class and pod-network testing, **When** the operator starts a benchmark,
   **Then** each applicable pod-network workload uses the selected runtime class.
2. **Given** a Kata runtime class and all scenarios enabled, **When** the operator starts a benchmark,
   **Then** no host-network workload uses the selected runtime class.
3. **Given** a Kata runtime class and host-network-only testing, **When** the operator starts a benchmark,
   **Then** the run fails before creating workloads with an actionable incompatibility error.
4. **Given** no runtime class, **When** the operator starts any supported benchmark scenario, **Then**
   existing host-network and pod-network behavior is unchanged.

### Edge Cases

- A configuration combines a runtime class with a scenario that creates both host-network and pod-network
  workloads; only pod-network workloads receive the runtime class.
- A runtime class is supplied with host-network-only selection; the run fails rather than silently omitting
  the requested runtime class.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: When a runtime class is selected, the system MUST apply it only to pod-network benchmark
  workloads.
- **FR-002**: When a runtime class is selected, the system MUST NOT apply it to any host-network benchmark
  workload.
- **FR-003**: The system MUST reject runtime-class use with host-network-only benchmark selection before
  creating workloads.
- **FR-004**: The system MUST preserve existing behavior for all scenarios when no runtime class is
  selected.
- **FR-005**: The project MUST document the runtime-class restriction for host-network benchmarks.

### Key Entities *(include if feature involves data)*

- **Runtime-class workload**: A benchmark pod scheduled with the operator-selected runtime class.
- **Pod-network workload**: A benchmark pod using the cluster pod network.
- **Host-network workload**: A benchmark pod using node networking directly.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In acceptance testing, 100% of runtime-class workloads created for a mixed-scenario run use
  pod networking and 0% use host networking.
- **SC-002**: 100% of host-network-only runs with a runtime class fail before creating workloads.
- **SC-003**: Existing runs without a runtime class retain their current workload-network behavior.

## Assumptions

- Kata is selected through the existing runtime-class option.
- This feature applies the restriction to any selected runtime class, not only a runtime class named
  `kata`, because host networking is incompatible with the requested isolated workload behavior.
