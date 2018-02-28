[![CircleCI](https://circleci.com/gh/giantswarm/kubernetes-node-exporter.svg?style=svg&circle-token=0a5aafcebabaed6f39a57293a96427f907674276)](https://circleci.com/gh/giantswarm/kubernetes-node-exporter)

# node-exporter Helm Chart
Helm Chart for Node exporter in Guest Clusters

* Installs the [node-exporter service](https://github.com/prometheus/node_exporter).

## Installing the Chart

To install the chart locally:

```bash
$ git clone https://github.com/giantswarm/kubernetes-node-exporter.git
$ cd kubernetes-node-exporter
$ helm install kubernetes-node-exporter/helm/kubernetes-node-exporter-chart
```

Deployment to Guest Clusters will be handled by [chart-operator](https://github.com/giantswarm/chart-operator)

## Configuration

| Parameter          | Description                                | Default                          |
|--------------------|--------------------------------------------|----------------------------------|
| `name`             | The name of the service                    | node-exporter                    |
| `namespace`        | The namespaces the services runs in        | kube-system                      |
| `image.repository` | The image repository to pull from          | quay.io/giantswarm/node-exporter |
| `image.tag`        | The image tag to pull from                 | v0.15.1                          |
| `port`             | The port of the container                  | 10300                            |
| `resources`        | node-exporter resource requests and limits | cpu:55m  - memory:75Mi           |
