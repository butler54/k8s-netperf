# Implementation Plan: Workload Labels

**Branch**: `004-workload-labels` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

## Summary

Add repeated custom-label CLI input, validate and protect selector labels, and merge valid labels into all
benchmark pod templates and VMIs using existing resource builders.

## Technical Context

**Language/Version**: Go 1.25
**Dependencies**: Cobra and Kubernetes client-go
**Testing**: Focused Go unit tests and `make verify-ci`
**Target Platform**: Kubernetes and OpenShift
**Constraints**: Tool selector labels always take precedence; defaults remain unchanged.

## Constitution Check

All principles pass: existing typed resource paths are extended minimally, invalid input is actionable,
and focused tests and documentation are required.

## Project Structure

```text
cmd/k8s-netperf/k8s-netperf.go      # Repeated label flag and validation
pkg/config/config.go                # Scenario label map
pkg/k8s/kubernetes.go               # Pod-template label merge
pkg/k8s/kubevirt.go                 # VMI label merge
pkg/k8s/*_test.go                   # Tests
docs/advanced-usage.md              # Documentation
```
