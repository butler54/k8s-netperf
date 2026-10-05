<!--
Sync Impact Report
- Version change: unversioned scaffold -> 1.0.0
- Modified principles: none; the scaffold placeholders were replaced with five initial principles.
- Added sections: Benchmark Integrity and Development Workflow.
- Removed sections: none.
- Follow-up TODOs: TODO(RATIFICATION_DATE): The original adoption date is unknown.
-->
# k8s-netperf Constitution

## Core Principles

### I. Benchmark Integrity
Benchmark definitions, execution paths, collected measurements, and pass/fail evaluation MUST preserve
their documented meaning. Changes affecting measurement methodology, defaults, units, or aggregation
MUST update the relevant documentation and include validation that demonstrates the intended behavior.
This protects the comparability and trustworthiness of performance results.

### II. Kubernetes-Native Compatibility
The tool MUST remain compatible with its supported Kubernetes networking scenarios, benchmark tools,
and documented configuration formats. Interface or configuration changes MUST preserve existing behavior
unless an intentional breaking change is documented in the README and configuration documentation.
This keeps established automation and benchmark environments usable.

### III. Verified Go Changes
Production Go changes MUST be formatted with `gofmt -s`, pass `golangci-lint`, and pass relevant Go
tests before review. New or changed behavior MUST have focused automated coverage when it can be tested
without a live cluster; Kubernetes integration behavior MUST include a documented validation procedure.
These checks prevent regressions in a tool used to assess infrastructure performance.

### IV. Observable and Actionable Results
Benchmark output, errors, archived results, and exported metrics MUST make the executed scenario and
outcome clear enough to diagnose failures or regressions. New user-visible output, environment
variables, ports, file locations, or container parameters MUST be documented. This enables reliable
comparison and efficient troubleshooting.

### V. Minimal, Repository-Consistent Changes
Changes MUST use established project patterns and add only the complexity required for the stated
functionality. New abstractions, dependencies, configuration, or deployment surface MUST be justified
by a concrete requirement. This limits maintenance cost and reduces risk to benchmark behavior.

## Benchmark Integrity

Supported benchmark tools include netperf, iperf3, uperf, and ib_write_bw. Changes to supported tools,
network scenarios, result formats, or regression criteria MUST retain the documented behavior or
explicitly document the change. Test configuration and result fixtures MUST be representative and MUST
NOT misstate measured performance.

## Development Workflow

Contributors MUST discuss proposed changes with repository owners before implementation. Reviewers MUST
verify constitution compliance, required validation, and documentation changes. The standard quality
gates are `make verify-fast` for static checks and `make verify-ci` for static checks plus unit tests.
Pull requests require sign-off from at least one other developer before merge.

## Governance

This constitution defines project governance for this workspace and complements the repository's
contribution guidance. Spec Kit artifacts, including this constitution, are private workflow material
and MUST NOT be committed or contributed upstream.

Amendments MUST document the affected principles, rationale, and compatibility impact in the Sync Impact
Report. The constitution version uses semantic versioning: MAJOR for incompatible governance changes,
MINOR for new or materially expanded principles or sections, and PATCH for clarifications and
non-semantic refinements. Every code review MUST assess compliance with these principles; deviations
MUST be explicitly justified and approved by maintainers.

**Version**: 1.0.0 | **Ratified**: TODO(RATIFICATION_DATE): The original adoption date is unknown. | **Last Amended**: 2026-09-14
