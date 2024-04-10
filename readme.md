# tanjunchen-scheduler-sample

This repo is a sample for Kubernetes scheduler framework. The `sample` plugin implements `filter` and `prebind` extension points.

And the custom scheduler name is `tanjunchen-framework` which defines in `KubeSchedulerConfiguration` object.

## Build

### binary
```shell
$ make local
```

### image
```shell
$ make image
```

## Deploy

```shell
$ kubectl apply -f ./deploy/
```
