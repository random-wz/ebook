> 在使用k8s进行服务的部署过程中我们会使用到Deployment、Service、Pod等资源，在 yaml 文件中我们需要指定对应的 API 版本，我们可以通过访问相应的接口来管理相应的资源信息，在 k8s 中为了提高 API 的可扩展性，采用了 API Groups 进行标识这些接口，在 client-go 源码中就是通过指定的 API Groups 来访问 k8s 集群的，这里向大家介绍 API Groups 都有哪些，希望对你有帮助。

##### 当前 k8s 支持两类 API Groups：

##### 1. Core Groups（核心组）

该分组也可以称之为 Legacy Groups，作为 k8s 最核心的 API ，其特点是没有组的概念，例如 “v1”，在资源对象的定义中表示为 "apiVersion: v1"，属于核心组的资源主要有下面几种：

- Container
- Pod
- ReplicationController
- Endpoint
- Service
- ConfigMap
- Secret
- Volume
- PersistentVolumeClaim
- Event
- LimitRange
- PodTemplate
- Binding
- ComponentStatus
- Namespace
- Node

##### 2. 具有分组信息的 API 

这种 API 接口以`/apis/$GROUP_NAME/$VERSION` URL 路径进行标识，在资源对象的定义中表示为 "apiVersion: $GROUP_NAME/$VERSION"， 例如 “apiVersion: batch/v1”，常见的 Group 及资源主要有下面几种：

- apps/v1
  - [ ] DaemonSet
  - [ ] Deployment
  - [ ] StatefulSet
  - [ ] ReplicaSet
- batch/v1
  - [ ] Job 

- batch/v1beta
  - [ ] CronJob

更多 API 接口信息请参考官网：[k8s1.17 API 接口文档](https://v1-17.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.17)

