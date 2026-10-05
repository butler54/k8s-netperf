# Feature Specification: Internet-Disconnected Workflows

**Feature Branch**: `002-airgapped-operation`

**Created**: 2026-09-14

**Updated**: 2026-10-05

**Status**: Draft

**Input**: Issue #286, "Support internet disconnected workflows."

## Clarifications

### Session 2026-10-05

- Q: Should image overrides cover both benchmark pods and online VM container-disk images? → A: Override pod and online VM images.
- Q: Should the prebuilt VM source be limited to a ready CDI DataVolume? → A: CDI DataVolume only.
- Q: Must the source DataVolume reside outside the benchmark namespace? → A: Do not impose a namespace restriction.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Use Operator-Provided Workload Images (Priority: P1)

As a benchmark operator in a disconnected environment, I can provide the image references for the
containerized workloads used by my benchmark so that the run does not depend on public image
registries or platform-specific image mirroring services.

**Why this priority**: A benchmark cannot begin unless every required workload image is available to
the cluster.

**Independent Test**: An operator provides images hosted in an accessible private registry, starts
pod and VM benchmarks with public registry access blocked, and verifies that each created workload
uses the supplied image.

**Acceptance Scenarios**:

1. **Given** an accessible private image reference for pod workloads, **When** the operator starts a
   pod benchmark, **Then** every benchmark pod uses that reference.
2. **Given** an accessible private image reference for VM workloads, **When** the operator starts a
   VM benchmark, **Then** every created VM uses that reference when a prebuilt VM disk is not selected.
3. **Given** no image override, **When** an operator starts a benchmark, **Then** existing image
   selection behavior is unchanged.
4. **Given** an invalid or unavailable supplied image, **When** the operator starts a benchmark,
   **Then** the run reports the supplied reference and does not substitute a public image.

---

### User Story 2 - Run VMs from a Prebuilt Disk (Priority: P2)

As a VM benchmark operator in a disconnected environment, I can select a prebuilt VM disk already
available in the cluster so that benchmark VMs start without downloading or building software.

**Why this priority**: The current VM startup workflow requires internet access to install benchmark
software, which prevents VM benchmarks in disconnected environments.

**Independent Test**: An operator supplies a ready prebuilt VM disk, blocks outbound internet access,
starts a VM benchmark, and confirms that the benchmark reaches its ready state without download
attempts.

**Acceptance Scenarios**:

1. **Given** a ready prebuilt VM disk in an operator-managed location, **When** the operator starts a
   VM benchmark, **Then** the system creates an independent benchmark disk for each created VM and
   starts each VM from it.
2. **Given** a requested prebuilt VM disk is missing or not ready, **When** the operator starts a VM
   benchmark, **Then** the run fails before creating a benchmark VM and identifies the requested disk.
3. **Given** a prebuilt VM disk lacks a required benchmark tool, **When** the operator starts a VM
   benchmark, **Then** the run identifies the missing prerequisite and does not attempt to download it.
4. **Given** a benchmark run using a prebuilt VM disk completes, **When** cleanup runs, **Then** the
   operator-provided source remains available and only benchmark-created disks are removed.

### Edge Cases

- A supplied image reference is empty, malformed, inaccessible, or does not contain the required
  workload image; the run identifies the supplied reference and does not fall back to a public image.
- A prebuilt VM disk cannot be shared directly by multiple benchmark VMs; each benchmark VM receives
  an independent benchmark-owned copy.
- A prebuilt VM disk is requested together with a VM image override; the run rejects the incompatible
  choices before creating benchmark resources.
- The source disk is in the benchmark namespace; the run preserves it while removing only
  benchmark-owned disk copies.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow operators to provide complete image references for benchmark pod
  workloads and online VM container-disk workloads.
- **FR-002**: When an image reference is provided, the system MUST use that exact reference for every
  applicable workload created by the benchmark.
- **FR-003**: The system MUST preserve existing image selection behavior when an operator does not
  provide an image reference.
- **FR-004**: The system MUST report an image retrieval failure with the supplied image reference and
  MUST NOT silently substitute a public image.
- **FR-005**: The system MUST allow VM benchmark operators to select a named, ready CDI DataVolume
  that already exists in the target cluster as the prebuilt source disk.
- **FR-006**: For a prebuilt source disk, the system MUST verify readiness before creating any
  benchmark VM or benchmark-owned disk copy.
- **FR-007**: The system MUST create a separate benchmark-owned disk copy for every VM that uses the
  selected source disk.
- **FR-008**: The system MUST start prebuilt-disk VMs without attempting to download packages, source
  code, or build artifacts.
- **FR-009**: The system MUST report missing required VM benchmark prerequisites without attempting
  to obtain them from the internet.
- **FR-010**: Cleanup MUST preserve the operator-provided source disk and remove only benchmark-owned
  disk copies, regardless of the source disk's namespace.
- **FR-011**: The system MUST reject a prebuilt-disk selection that conflicts with an explicitly
  selected VM image before creating benchmark resources.
- **FR-012**: The project MUST document the required properties of a prebuilt VM disk and the image
  override options, without providing instructions to build or upload VM disks.

### Key Entities *(include if feature involves data)*

- **Operator-provided image reference**: A complete pod or online VM container-disk image location
  selected by an operator instead of the default image location.
- **Prebuilt source disk**: An operator-managed, ready CDI DataVolume in the target cluster containing
  the operating system and benchmark prerequisites needed for disconnected operation.
- **Benchmark-owned disk copy**: A temporary independent copy of the prebuilt source disk used by one
  benchmark VM and eligible for benchmark cleanup.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In acceptance testing, 100% of workloads created with an operator-provided image
  reference use the specified reference.
- **SC-002**: A VM benchmark using a ready prebuilt source disk reaches its ready state with outbound
  internet access blocked and makes zero package, source-code, or build-artifact download attempts
  during startup.
- **SC-003**: In cleanup acceptance testing, 100% of operator-provided source disks remain available
  and 100% of benchmark-owned disk copies are removed.
- **SC-004**: In validation testing, 100% of missing, unready, incompatible, or incomplete prebuilt
  disk selections fail before a benchmark result is reported as valid.
- **SC-005**: Existing benchmark invocations without disconnected-workflow options retain their
  current image selection and VM startup behavior.

## Assumptions

- Operators manage image registries, image access, VM disk construction, and VM disk delivery to the
  cluster; these activities are outside this feature's scope.
- Operators provide complete image references for every selected workload type and have permission to
  retrieve those images.
- A selected prebuilt source disk is reusable, may reside in any namespace, and contains all software
  required by the selected benchmark drivers.
- The feature selects existing images and source disks but does not create, synchronize, authenticate
  to, build, or upload them.
