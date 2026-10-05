# CLI Contract: Internet-Disconnected Workflows

| Option | Value | Behavior |
|--------|-------|----------|
| `--image` | Complete pod workload image reference | Replaces the default benchmark pod image for the run. |
| `--vm-image` | Complete VM workload image reference | Replaces the default VM image for an online VM run. |
| `--offline-data-volume` | `NAMESPACE/NAME` | Selects a ready, prebuilt source disk for a VM run and creates one benchmark-owned copy per VMI. |

## Validation

- Image-reference values must be non-empty complete references when selected.
- `--offline-data-volume` applies only to VM benchmarks and requires a non-empty namespace and name.
- An explicit `--vm-image` and `--offline-data-volume` are mutually exclusive.
- The prebuilt source disk must be ready before any benchmark-owned copy or VMI is created.
- The source disk may be in any namespace and is never removed by benchmark cleanup; cleanup removes only
  explicitly benchmark-owned disk copies.
- An unavailable image or incomplete prebuilt source disk fails the run with an actionable error; no
  public-image fallback or startup download is attempted.
