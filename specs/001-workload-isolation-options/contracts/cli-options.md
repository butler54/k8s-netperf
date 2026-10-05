# CLI Contract: Workload Isolation Options

## Options

| Option | Value | Applicability | Behavior |
|--------|-------|---------------|----------|
| `--annotation` | Repeated `KEY=VALUE` pair | Pods and VMIs | Adds the pair to each applicable benchmark resource's metadata. The first `=` separates key from value, allowing `=` in values. |
| `--runtime-class` | Runtime class name | Pods only | Sets the selected runtime class on every created benchmark pod. A VM-only run rejects it. |
| `--launch-security` | `snp` or `tdx` | VMIs only | Sets `launchSecurity.snp: {}` or `launchSecurity.tdx: {}` on every created VMI. It does not configure TDX attestation. A non-VM run rejects it. |

## Validation Rules

- Each `--annotation` value MUST have a non-empty key and value.
- A duplicate annotation key MUST fail validation.
- A user annotation key owned by the tool, including a generated network or default workload annotation,
  MUST fail validation rather than replace the managed value.
- `--runtime-class` MUST fail if pod execution is disabled.
- `--launch-security` MUST fail if VM execution is disabled or its value is not `snp` or `tdx`.
- Cluster-side admission and scheduling failures MUST be returned to the operator without claiming valid
  benchmark results.

## Compatibility

Omitting all options preserves current command behavior and generated resources. The options are
independent except that each is validated against the enabled workload types.
