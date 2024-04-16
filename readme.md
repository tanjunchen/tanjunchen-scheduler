# tanjunchen-scheduler

This repo is a example for Kubernetes scheduler framework. The `sample` plugin implements `filter` extension points.

And the custom scheduler name is `tanjunchen-scheduler` which defines in `KubeSchedulerConfiguration` object.

## test enviroment

```bash
Kubernetes version: v1.26.9
```

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
