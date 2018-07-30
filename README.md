[![CircleCI](https://circleci.com/gh/giantswarm/kubernetes-node-exporter.svg?style=svg&circle-token=0a5aafcebabaed6f39a57293a96427f907674276)](https://circleci.com/gh/giantswarm/kubernetes-node-exporter)

# node-exporter Helm Chart
Helm Chart for Node exporter in Guest Clusters.

* Installs the [node-exporter service](https://github.com/prometheus/node_exporter).

## Installing the Chart

To install the chart locally:

```bash
$ git clone https://github.com/giantswarm/kubernetes-node-exporter.git
$ cd kubernetes-node-exporter
$ helm install helm/kubernetes-node-exporter-chart
```

Provide a custom `values.yaml`:

```bash
$ helm install kubernetes-node-exporter-chart -f values.yaml
```

Deployment to Guest Clusters is handled by [chart-operator](https://github.com/giantswarm/chart-operator).
