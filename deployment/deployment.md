## **Kubernetes**
### **安装minikube**
```shell
# 安装 minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 验证安装
minikube status
# 启动集群
minikube start
# 关闭集群
minikube stop

# 构建镜像
# in 2501YTC/
docker build -f app/order/Dockerfile -t order-service:latest .
# 本地运行(test)
docker run -d -p 8085:8085 -p 4317:4317 -e GO_ENV=online --name order-service order-service

# 切换到 minikube docker 环境
eval $(minikube docker-env)
# 验证镜像是否存在
docker images | grep order-service
# 将本地镜像加载到 minikube
minikube image load order-service:latest

# 部署服务
kubectl apply -f deployment/XX.yaml
# 删除服务
kubectl delete -f XXXX.yaml

# 验证基础服务状态
kubectl get pods
kubectl get services

# 如果 Pod 仍然失败，查看详细信息
kubectl describe pod $(kubectl get pods -l app=gateway -o jsonpath='{.items[0].metadata.name}')
# 查看容器日志
kubectl logs $(kubectl get pods -l app=gateway -o jsonpath='{.items[0].metadata.name}')
# 进入容器内
kubectl exec -it $(kubectl get pods -l app=gateway -o jsonpath='{.items[0].metadata.name}') -- /bin/sh

# 转发服务端口到本地
kubectl port-forward svc/gateway-service 8080:8080
kubectl port-forward svc/grafana-service 3000:3000
kubectl port-forward svc/jaeger-service 16686:16686
```