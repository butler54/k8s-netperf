# CLI Contract

`--runtime-class` applies only to pod-network workloads. Combining it with `--hostNet` fails before
creating workloads. In mixed runs, host-network workloads omit the runtime class.
