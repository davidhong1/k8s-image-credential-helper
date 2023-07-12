# k8s-image-credential-helper
在 Namespace 维度配置 pod 拉取私有镜像的权限

参考:

[pull-image-private-registry](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/pull-image-private-registry/)

[add-imagepullsecrets-to-a-service-account](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account)


# Feature

1. Watch ns create event, and auto add image credential
2. Support registry
   - Harbor
   - Docker(TODO)
3. 支持配置注入方式
   - 环境变量
3. 支持指定多个 Namespace 添加 Image credential，默认所有 Namespace
4. 支持指定多个 ServiceAccount，默认 default ServiceAccount

# 如何注入配置
## 环境变量

|环境变量名|是否必填|默认值|
|----|----|----|
|INIT_CONFIG|Yes|environment 表示从环境变量读取配置|
|HTTP_HEALTH_CHECK_PORT|No|8080|
|IMAGE_PROVIDER|Yes|harbor|
|IMAGE_HOST|Yes||
|IMAGE_USER|Yes||
|IMAGE_PASSWORD|Yes||
|SERVICE_ACCOUNTS|No|default|
|WATCH_NAMESPACES|No|*|
|FORCE_UPDATE_SECRET|No| Just '1\|yes\|y' mean need force update secret


Harbor + 环境变量配置例子见 deploy/all.yaml

# Pod 如何使用 Image credential

Pod 需要和授权了 Image credential 的 ServiceAccounts 进行绑定；
或者对应 Deployment & StatefulSet & DaemonSets 和授权了 Image credential 的 ServiceAccounts 进行绑定。

# 具体实现

监听 ns 创建，在新 ns 上面创建 secret docker-registry，绑定 secret docker-registry 给指定服务账号。

pod 使用上指定的服务账号

## 流程图


## k8s RBAC 说明

# 快速部署

## yaml

## helm
