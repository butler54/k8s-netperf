# Research

**Decision**: Clear `RuntimeClass` for `HostNetwork` deployment parameters and reject `--runtime-class`
with `--hostNet`.

**Rationale**: Existing deployment construction receives both values; this is the smallest way to prevent
Kata from being used for host networking while retaining mixed-run pod-network coverage.

**Alternative**: Reject all mixed runs was rejected because pod-network Kata measurements remain valid.
