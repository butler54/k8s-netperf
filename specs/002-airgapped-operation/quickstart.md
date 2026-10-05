# Quickstart: Validate Internet-Disconnected Workflows

## Prerequisites

- Access to a registry that contains the required benchmark workload images.
- For VM validation, a ready prebuilt source disk. The disk may be in the benchmark namespace and must include
  the operating system and all tools required by the benchmark drivers selected for the run.
- A Kubernetes/OpenShift environment with KubeVirt and CDI for VM validation.

Preparing or uploading the prebuilt source disk is outside this project’s scope.

## Validate Workload Image Overrides

1. Block or remove public-registry access for the benchmark environment.
2. Start a pod benchmark with `--image` set to the complete private image reference.
3. For an online VM benchmark, start a VM benchmark with `--vm-image` set to the complete private image
   reference.
4. Inspect the created workloads and verify every applicable workload uses the supplied image reference.
5. Repeat without image overrides and verify the existing defaults are retained.

## Validate a Prebuilt VM Source Disk

1. Confirm the selected source disk is ready; it may be in the benchmark namespace.
2. Block outbound internet access.
3. Start a VM benchmark with `--offline-data-volume=NAMESPACE/NAME`.
4. Verify that each benchmark VMI starts from its own benchmark-owned disk copy and no startup download
   occurs.
5. Verify that a missing local benchmark prerequisite causes an actionable failure without a download
   attempt.
6. Complete cleanup and verify that the source disk remains while only explicitly benchmark-owned copies
   are removed.

## Repository Validation

Run `make verify-ci` after implementation. See [the CLI contract](contracts/cli-options.md) and
[data model](data-model.md) for expected option and lifecycle behavior.
