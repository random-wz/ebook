## 一、Deployment

>Deployment 是 Kubernetes 在 1.2 版本中引入的新概念，用于更好的解决 Pod 的编排问题。为此， Deployment 内部使用了 Replica Set 来实现目的，无论从 Deployment 的作用与目的、YAML文件的定义，还是从他具体命令行的操作来看，我们都可以把他看作 RC 的一次升级，两者的相似度超过90%，这里向大家介绍如何使用 Deployment 控制器进行应用的部署。

#### 1. 通过命令创建 Deployment

` kubectl run nginx-deployment --image=nginx:1.7.9 --replicas=3`

上面的命令将部署包含三个副本的 Deployment nginx-deployment，我们通过 `kubectl describe`命令可以看到 nginx-deployment 的详细信息：

```bash
kubectl describe deployment nginx-deployment
...
```

从上面的信息中我们可以看到 Deployment 通过 ReplicaSet 来管理 Pod，我们查看 replicaset 状态：

```bash
kubectl get replicaset
```

副本已经启动完成，我们查看 replicaset 的详细信息：

```bash
kubectl describe deployment nginx-deployment
```

可以看到 Controlled By 指明 ReplicatSet 是由 Deployment nginx-deployment 创建的，接着我们可以看到三个 Pod 已经处于 Running 状态：

```bash
kubectl get pod
...
```

通过 kubectl describe 查看 Pod 的详细信息：

```bash
kubectl describe pod ...
```

可以看到 Controlled By 指明 Pod 是由 ReplicaSet nginx-deployment-... 创建的。

总结一下 Deployment 创建 Pod 的整个过程如下图:

```
kubectl->deployment->replicaset->po
```

#### 2. 通过配置文件创建 Deployment

#### 3. Deployment 配置文件简介

#### 4.  扩容与缩容

#### 5. Failover（故障转移）

#### 6. 通过 Label 控制 Pod 的位置

#### 7. Deployment 应用场景

## 二、 DaemonSet



## 三、Job

## 四、Service

