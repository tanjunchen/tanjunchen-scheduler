# tanjunchen-scheduler

This repo is a example for Kubernetes scheduler framework. 

And the custom scheduler name is `tanjunchen-scheduler` which defines in `KubeSchedulerConfiguration` object.

## test enviroment

```bash
Kubernetes version: v1.26.1
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
