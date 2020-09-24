# kubectl-resource-view kubectl

A `kubectl` plugin to  view node resources 

## Quick Start

```
kubectl krew install kubectl-resource-view
kubectl resource-view
```
## Usage

Get node  resources (request,limit)

```
kubectl resource-view -node  host103052172
```
Result

```
Node:           CPU     Memory          CPURequests     CPULimits       MemoryRequests          MemoryLimits
host103052172   32      125.000000G     5.85core (18%)  5.10core(15%)   20670Mi (16%)           20870Mi (16%)

```

Use label filter

```
kubectl resource-view -l  host103052172
```
Result

```
Node:           CPU     Memory          CPURequests     CPULimits       MemoryRequests          MemoryLimits
host103052172   32      125.000000G     5.85core (18%)  5.10core(15%)   20670Mi (16%)           20870Mi (16%)


```
