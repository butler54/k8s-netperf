# Research: Internet-Disconnected Workflows

## Operator-Provided Image References

**Decision**: Accept complete operator-provided image references for benchmark pod and VM workloads,
without rewriting registry prefixes or repository paths.

**Rationale**: Private registries commonly use repository layouts that differ from public registries.
Using the complete supplied reference lets operators use their existing image distribution approach and
prevents an unavailable public image from being selected implicitly.

**Alternatives considered**: Registry-prefix replacement was rejected because it assumes that repository
paths and tags remain unchanged. Automatic registry discovery was rejected because registry configuration
and credentials are operator-managed and outside the feature scope.

## Prebuilt VM Source Disks

**Decision**: Select an existing, ready source DataVolume, validate it before VMI creation, and create one
benchmark-owned clone DataVolume per VMI in the benchmark namespace.

**Rationale**: Separate clones let client and server VMIs use a reusable source disk without depending on
shared-access storage. Explicitly identifying benchmark-owned clones lets cleanup preserve the source even
when it shares the benchmark namespace.

**Alternatives considered**: Directly attaching one shared source disk was rejected because common
ReadWriteOnce disks cannot be attached to multiple VMIs. Rejecting sources in the benchmark namespace was
rejected because operators may store a reusable source there. Creating or uploading disks during a
benchmark was rejected because the updated issue makes those operator responsibilities.

## Prebuilt-Disk Startup

**Decision**: In prebuilt-disk mode, startup only configures required guest access and networking,
verifies locally available benchmark prerequisites, and starts needed local services. It must not download
packages, source code, or build artifacts.

**Rationale**: Disconnected operation is only reliable when all selected benchmark dependencies are
already present on the source disk. Explicit validation produces an actionable failure when the source is
incomplete.

**Alternatives considered**: Downloading a missing dependency or using an offline cache was rejected
because both add image-preparation responsibility and undermine the no-internet guarantee.
