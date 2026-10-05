# Implementation Plan: Kata Pod Networking

**Branch**: `003-kata-pod-networking` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

## Summary

Restrict selected runtime classes to pod-network Deployments. Omit the runtime class from host-network
Deployments and reject host-network-only runs before resource creation.

## Technical Context

**Language/Version**: Go 1.25
**Dependencies**: Cobra and Kubernetes client-go
**Testing**: Focused Go unit tests and `make verify-ci`
**Target Platform**: Kubernetes/OpenShift
**Project Type**: CLI
**Constraints**: Preserve all behavior when no runtime class is set.

## Constitution Check

All principles pass: the change preserves measurement semantics, uses existing workload builders, adds
focused tests, reports invalid input clearly, and adds no new abstraction.

## Project Structure

```text
cmd/k8s-netperf/k8s-netperf.go   # Host-network-only validation
pkg/k8s/kubernetes.go            # Runtime class only on pod-network deployments
pkg/k8s/kubernetes_test.go        # Builder tests
cmd/k8s-netperf/k8s-netperf_test.go # CLI validation tests
docs/advanced-usage.md           # Restriction documentation
```
