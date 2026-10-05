# Feature Specification: Workload Isolation Options

**Feature Branch**: `001-workload-isolation-options`

**Created**: 2026-09-14

**Status**: Draft

**Input**: User description: "There are a number enhancements I want to be able to build in the
neatest way: custom annotations for benchmark pods and VM instances; a runtime class for Kata or
confidential containers; and KubeVirt VM launch security set to TDX or SNP."

## Clarifications

### Session 2026-09-14

- Q: Which benchmark workloads must accept the selected runtime class? → A: Pods only.
- Q: How should operators provide annotations and isolation settings for a benchmark run? → A:
  Command-line options only, following existing behavior.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure Isolated Workloads (Priority: P1)

As a benchmark operator, I can select a runtime class for benchmark workloads so that I can measure
network performance in Kata or confidential-container environments.

**Why this priority**: Runtime isolation is required to benchmark the environments the feature targets.

**Independent Test**: An operator runs a benchmark with a configured runtime class and confirms every
created benchmark workload uses that class; a run without the setting retains current behavior.

**Acceptance Scenarios**:

1. **Given** a valid runtime class and a pod benchmark configuration, **When** the operator starts the
   benchmark, **Then** each created benchmark pod uses the selected runtime class.
2. **Given** no runtime class is configured, **When** the operator starts a benchmark, **Then** workload
   scheduling behavior is unchanged.
3. **Given** an unavailable or unsupported runtime class, **When** the operator starts a pod benchmark,
   **Then** the run fails with a clear message identifying the setting and workload type.

---

### User Story 2 - Apply Custom Workload Annotations (Priority: P2)

As a benchmark operator, I can attach custom annotations to benchmark pods and VM instances so that
cluster integrations can apply the policies or behavior required by my test environment.

**Why this priority**: Annotations let operators use existing cluster capabilities without changing
benchmark definitions for each environment.

**Independent Test**: An operator configures multiple annotations, runs pod and VM benchmarks, and
confirms the requested annotations appear on every corresponding created workload.

**Acceptance Scenarios**:

1. **Given** custom annotations for a pod benchmark, **When** the operator starts the benchmark,
   **Then** every created benchmark pod contains those annotations.
2. **Given** custom annotations for a VM benchmark, **When** the operator starts the benchmark,
   **Then** every created benchmark VM instance contains those annotations.
3. **Given** custom annotations that duplicate an annotation managed by the tool, **When** the operator
   starts the benchmark, **Then** the run fails before creating workloads and identifies the conflict.
4. **Given** no custom annotations are configured, **When** the operator starts a benchmark, **Then**
   existing annotations and behavior are unchanged.

---

### User Story 3 - Select VM Launch Security (Priority: P3)

As a VM benchmark operator, I can select TDX or SNP launch security so that I can measure networking
in the corresponding confidential VM environment.

**Why this priority**: Launch security is specific to VM environments but is necessary for confidential
VM benchmark coverage.

**Independent Test**: An operator starts VM benchmarks with each supported launch security selection and
confirms the created VM instance requests that selection; an incompatible run reports a clear error.

**Acceptance Scenarios**:

1. **Given** a VM benchmark and TDX launch security, **When** the operator starts the benchmark,
   **Then** each created VM instance requests TDX launch security.
2. **Given** a VM benchmark and SNP launch security, **When** the operator starts the benchmark,
   **Then** each created VM instance requests SNP launch security.
3. **Given** launch security is configured for a non-VM benchmark, **When** the operator starts the
   benchmark, **Then** the run fails before creating workloads and explains that the setting is VM-only.
4. **Given** an unsupported launch security value, **When** the operator starts a VM benchmark, **Then**
   the run fails before creating workloads and lists the supported values.

### Edge Cases

- A configured annotation key or value is empty, malformed, or repeated; the run fails before creating
  workloads and identifies the invalid entry.
- A requested runtime class or launch security mode cannot be honored by the target cluster; the run
  fails without reporting benchmark results as valid.
- The benchmark creates both client and server workloads; all applicable settings are applied uniformly
  to both sides.
- Existing network-related and system-managed annotations are preserved and cannot be replaced by custom
  configuration.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow operators to define one or more custom key-value annotations for
  benchmark workloads through command-line options.
- **FR-002**: The system MUST apply configured custom annotations to every benchmark pod and every
  benchmark VM instance created by a run, including client and server workloads.
- **FR-003**: The system MUST preserve annotations managed by the tool and reject custom annotations
  that would replace a managed annotation before creating any workload.
- **FR-004**: The system MUST allow operators to select one runtime class for a benchmark run through
  command-line options.
- **FR-005**: The system MUST apply the selected runtime class to every benchmark pod and retain existing
  scheduling behavior when no runtime class is selected.
- **FR-006**: The system MUST reject a runtime class that is unavailable or unsupported for the selected
  pod benchmark before reporting a successful benchmark result.
- **FR-007**: The system MUST allow VM benchmark operators to select exactly one launch security mode,
  TDX or SNP, through command-line options.
- **FR-008**: The system MUST apply the selected launch security mode to every VM instance created by a
  VM benchmark run.
- **FR-009**: The system MUST reject launch security configuration for non-VM benchmarks and reject
  values other than TDX or SNP before creating workloads.
- **FR-010**: The system MUST report configuration validation failures clearly enough for an operator to
  identify the invalid setting, its value, and the applicable workload type.
- **FR-011**: The system MUST preserve current behavior when none of the new workload options are set.
- **FR-012**: The system MUST document the new workload options, their applicable workload types, valid
  values, and incompatibility behavior.

### Key Entities *(include if feature involves data)*

- **Workload options**: Operator-selected annotations, runtime class, and launch security settings for a
  benchmark run.
- **Managed annotation**: An annotation controlled by the tool to configure networking or required
  workload behavior; it cannot be overridden by a custom annotation.
- **Benchmark workload**: A client or server pod or VM instance created to execute a benchmark.
- **Launch security mode**: The confidential-VM protection requested for a VM benchmark; valid values are
  TDX and SNP.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In acceptance testing, 100% of created applicable client and server workloads carry every
  configured custom annotation, pod runtime class, and VM launch security selection.
- **SC-002**: In acceptance testing, 100% of invalid annotation, runtime class, or launch security
  configurations fail before workload creation and identify the invalid setting in the error.
- **SC-003**: Existing benchmark configurations that omit all new options complete with the same workload
  scheduling and metadata behavior as before the feature.
- **SC-004**: An operator can configure and run each supported isolation scenario using documented
  instructions without requiring source changes.

## Assumptions

- Operators have permission to create the requested workload type and to use the selected runtime class
  and launch security capability in the target cluster.
- Workload options follow existing VM and pod controls and are provided through command-line options.
- Custom annotations apply uniformly to all benchmark client and server workloads of the selected type.
- Tool-managed annotations take precedence over custom annotations to preserve required networking and
  benchmark behavior.
- Runtime classes apply only to pod benchmarks and are accepted only when the target cluster supports
  them.
- This feature does not create runtime classes, enable confidential computing on clusters, or validate
  node hardware capabilities beyond reporting that a requested setting cannot be honored.
