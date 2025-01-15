# k8sgpt-single-replica-analyzer

# K8sGPT Single Replica Analyzer

A custom analyzer for K8sGPT that identifies Kubernetes deployments running with single replicas, which might indicate potential availability risks.

## Description

This analyzer checks all deployments in a Kubernetes cluster and raises warnings for those configured with only one replica, suggesting improvements for high availability.

## Installation

```bash
go get github.com/ksatchit/k8sgpt-single-replica-analyzer
```
