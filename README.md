# k8s-image-credential-helper

监听 ns 的创建，然后在 ns 级别配置能够拉取私有仓库的私有镜像权限。
所以作用范围是私有仓库的私有镜像。参考
[pull-image-private-registry](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/pull-image-private-registry/)
[add-imagepullsecrets-to-a-service-account](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account)

目前已支持的私有镜像如下:

- Harbor

```
环境变量 loader
	envHost     = "IMAGE_HOST"
	envUser     = "IMAGE_USER"
	envPassword = "IMAGE_PASSWORD"

	serviceAccounts = "SERVICE_ACCOUNTS"
	watchNamespaces = "WATCH_NAMESPACES"
```

# 功能

监听 ns 创建，在新 ns 上面创建 secret docker-registry，绑定 secret docker-registry 给指定服务账号。
pod 使用上指定的服务账号

# 流程图

# Feature

1. watch ns create event, and auto add image credential
2.

# 权限

# 部署步骤

## yaml 部署

## helm 部署
