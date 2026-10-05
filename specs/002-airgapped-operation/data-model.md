# Data Model: Internet-Disconnected Workflows

| Entity | Fields | Validation rules | Relationships |
|--------|--------|------------------|---------------|
| Workload image override | Workload type, complete image reference | Reference is non-empty when selected; used verbatim; no fallback substitution. | Applies to all created workloads of the selected type. |
| Prebuilt source disk | Namespace, name, readiness state, benchmark prerequisites | Must exist, be ready, and contain all selected benchmark prerequisites. May be in the benchmark namespace. | Produces one benchmark-owned disk copy for each VMI. |
| Benchmark-owned disk copy | Name, benchmark namespace, source reference, owning run, ownership marker | Must be associated with the selected source and current run. | Attached to exactly one benchmark VMI and removed by cleanup. |
| Prebuilt-disk mode | Source disk selection, selected VM image override | Applies only to VM benchmarks; cannot be combined with an explicit VM image override. | Controls source-disk copy creation and no-download startup. |

## Lifecycle

1. Validate image overrides and prebuilt-disk selection before creating benchmark resources.
2. When prebuilt-disk mode is selected, validate the source disk is ready and record clone ownership
   independently of namespace.
3. Create one benchmark-owned disk copy for every benchmark VMI.
4. Start each VMI from its assigned copy and validate required local benchmark tools.
5. Remove benchmark-owned copies during cleanup; preserve the source disk.
