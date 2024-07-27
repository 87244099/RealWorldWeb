# k8s

## what


## why

## how

### 安装

它有几种方式：
- 学习环境（推荐）
- 生产环境
- xxx

#### 学习环境

- 安装kubectrl，是一个命令行工具，可以对集群进行管理，监控等，相当于shell
- 安装minikube
  - 是一个mini的学习环境，它可以让你本地运行k8s
  - minikube start启动服务
    - 然后它会下载一堆服务，不过因为墙的原因，你需要搞一个代理：使用docker proxy解决
- 使用minikube
  - 先看下有哪些上下文：kubectl config get-contexts
  - 使设置上下文: kubectl config set-context minikube
  - 查看pod：kubectl get pods -A
  - 查看命名空间：kubectl get namespace
- 部署应用
  - 创建一个一个实例
    - kubectl create deployment hello-minikube --image=kicbase/echo-server:1.0
  - 暴露实例的端口
    - kubectl expose deployment hello-minikube --type=NodePort --port=8080
  - 查看服务是否存在
    - minikube service hello-minikube
  - 做端口转发
    - kubectl port-forward service/hello-minikube 7080:8080
- 拉取echo-server镜像：
  - `docker pull kicbase/echo-server:1.0`
  -  docker exec -it minikube bash
    - 进入后查看下镜像信息：docker inspect kicbase/echo-server:1.0
- 上面会提到一种叫NodePort的东西，他是服务的一种类型，用于映射到某个端口，怎么证明这一点？
  - 我们在echo-server看到的NodePort是30462

## deployment的增删查改
- 如何把服务从1个节点扩容为3个节点
  - 先查看你有多少个节点
  
      $ kc get deployments
      NAME             READY   UP-TO-DATE   AVAILABLE   AGE
      hello-minikube   1/1     1            1           53m

  -  kc edit deployment hello-minikube: 会用默认编辑器打开配置文件进行修改，你也可以设置默认编辑器为某个IDE

## FQA

### Q Unable to resolve the current Docker CLI context "default": context "default"
运行：`docker context use default`

### Q：如果设置打开k8s配置文件的默认编辑器
<https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-DC2BB6E0-A327-4DB8-9A87-5F3376E70033.html?hWord=N4IghgNiBcICYFMBmYCuEAuACDCAe2CcAlhgPYBOIAvkA>