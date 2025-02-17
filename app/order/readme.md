# *** Project

## introduce

- Use the [Kitex](https://github.com/cloudwego/kitex/) framework
- Generating the base code for unit tests.
- Provides basic config functions
- Provides the most basic MVC code hierarchy.

## Directory structure

|  catalog   | introduce  |
|  ----  | ----  |
| conf  | Configuration files |
| main.go  | Startup file |
| handler.go  | Used for request processing return of response. |
| kitex_gen  | kitex generated code |
| biz/service  | The actual business logic. |
| biz/dal  | Logic for operating the storage layer |

## How to run

```shell
sh build.sh
sh output/bootstrap.sh

```

## Note!
```shell
# 查看所有注册的服务
curl http://localhost:8500/v1/catalog/services
# 查看订单服务的详细信息
curl http://localhost:8500/v1/catalog/service/order
```
```shell
# 在rabbitMQ容器中
# 查看所有队列
rabbitmqctl list_queues
# 删除队列
rabbitmqadmin delete queue name="order.delay.queue"
rabbitmqadmin delete queue name="order.dlx.queue"
```