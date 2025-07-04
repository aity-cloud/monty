# Multi Cluster Observability with AIOps 
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![codecov](https://codecov.io/gh/rancher/monty/branch/main/graph/badge.svg?token=9TYXGQ54FM)](https://codecov.io/gh/rancher/monty)
[![Go Report Card](https://goreportcard.com/badge/ github.com/aity-cloud/monty-monitoring)](https://goreportcard.com/report/ github.com/aity-cloud/monty-monitoring)

Observability data comes in the form of ***logs, metrics and traces***.
The collection and storage of observability data is handled by observability backends and agents.
AIOps helps makes sense of this observability data.
[Monty](https://monty.aity.cloud/) comes with all these nuts and bolts and can be used to self monitor a single cluster or be a centralized observability data sink for multiple clusters.

You can easily create the following with Monty:
* Backends
  - **Monty Logging** - extends [Opensearch](https://opensearch.org) to make it easy to search, visualize and analyze **logs**, **traces** and **Kubernetes events**
  - **Monty Monitoring** - extends [Cortex](https://cortexmetrics.io) to enable multi cluster, long term storage for **Prometheus metrics**

* Monty Agent
  - Collects logs, Kubernetes events, OpenTelemetry traces and Prometheus metrics with the click of a button
* Agentic AIOps
* Alerting and SLOs

Check out the [docs page](https://monty.aity.cloud/) to get started!

----


## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
