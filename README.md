<img width="878" height="536" alt="image" src="https://github.com/user-attachments/assets/3b60994b-68a0-4003-bd6a-71611480c8b4" /># 一、项目介绍

> 该项目基于 Go 语言和微服务架构打造，采用 Kitex 与 Hertz 构建高性能电商平台，涵盖用户认证、产品管理、购物车、订单支付及 AI 大模型接入等核心功能。项目通过 JWT 与 Casbin 实现安全认证，结合 MySQL、Redis、RabbitMQ 和 Consul 等关键技术，并借助 OpenTelemetry、Docker/Kubernetes 部署和 CloudWeGo Eino 驱动的 AI 智能处理，各模块高效协同、无缝整合，全面展现了现代微服务架构在高并发电商场景下的卓越性能与智能化应用能力
>
> 项目服务地址 ：容器化部署
>
> 项目地址 ：https://github.com/1850298154/2501YTC
>
> Apifox接口文档地址 ：https://apifox.com/apidoc/shared-bca22d02-4b8a-48de-8fa3-60f17842bec8

# 二、项目分工

| **团队成员** | **主要贡献**                                                 |
| ------------ | ------------------------------------------------------------ |
| 郭东翌       | 负责处理认证中心服务及测试                                   |
| 任丹妮       | 团队主要负责人，负责处理项目结构设计、用户服务、AI大模型接入等功能及测试 |
| 贾世飞       | 负责处理用户服务及测试                                       |
| 彭海林       | 负责处理商品服务及测试                                       |
| 赵威         | 负责处理购物车服务及测试                                     |
| 郑伊杰       | 负责处理订单服务及测试                                       |
| 赵雨腾       | 负责处理结算服务及测试                                       |
| 张景远       | 负责处理支付服务及测试                                       |

# 三、项目实现

### 3.1 技术选型与相关开发文档

#### 3.1.1 问题和目标

构建一个高性能、稳定且易扩展的电商平台，支持用户浏览、下单、支付、以及商品管理等业务场景。同时，平台需要具备良好的可观察性和容错性，保障系统在高并发场景下依然能提供稳定服务。

#### 3.1.2 前提假设

- **业务需求**  

1. 用户能够在平台上浏览商品、添加购物车、提交订单并完成支付；
2. 后台支持商品管理、订单处理、用户认证与授权等功能；
3. 业务需支持高并发读写操作，比如商品查询、订单创建等情况。

- **存储需求**

1. 随着用户浏览和下单数据的不断增加，MySQL 数据库需要保证数据的持久性；  
2. Redis 用作缓存以提升读写性能，预计需要根据实际流量规划内存；  
3. 商品图片和其他静态资源存储依赖 MinIO 或分布式存储系统，考虑到每份数据大约几 MB，每日上传量较大时建议至少预留数十 GB 空间；
4. 整体数据存储按照开发环境和线上环境分别规划，开发环境可满足基本 5G 存储需求，线上部署则建议预留 20G 甚至更大存储资源以应对高并发读写及海量数据存储。

- **服务器需求**

1. 根据微服务架构及高并发访问需求，至少需要 1 台主服务器作为基础节点，同时根据各模块的流量和业务量水平进行水平扩展；  
2. 每个微服务均能独立部署，支持通过容器编排工具进行动态伸缩，以应对峰值流量；  
3. 网络带宽、日志存储和监控等其他服务器资源也需要一并规划，确保整个平台在高压下依然稳定运行。

#### 3.1.3 开发流程

本项目按照以下流程进行开发：

1. 定义接口(IDL)
2. 生成代码框架
3. 实现业务逻辑
4. 编写单元测试
5. 代码审查
6. 构建部署
7. 部署方案

#### 3.1.4 编码规范

本项目具有统一的编码规范，并使用gitlint进行审查，并将其部署到GitHub Action中。

具体的编码规范如下所示。

暂时无法在飞书文档外展示此内容

同时，本项目具有良好的Git规范，统一分支创建、命名、合并等，统一commit。

具体的Git规范如下所示。

暂时无法在飞书文档外展示此内容

### 3.2 架构设计

暂时无法在飞书文档外展示此内容

#### 3.2.1 核心功能模块

1. 各模块介绍
   1. 用户服务(User Service)：用户注册、登录、信息管理
   2. 商品服务(Product Service)：商品CRUD、搜索、图片上传
   3. 购物车服务(Cart Service)：购物车管理
   4. 订单服务(Order Service)：订单创建、支付、取消
   5. 网关服务(Gateway Service)：统一API入口
   6. 认证服务(Auth Service): 认证授权
   7. 支付服务(Payment Service)：支付处理
   8. 结账服务(Checkout Service)：订单结算
   9. AI大模型(AI Service)：查询订单、自动下单
2. 网关 gateway 接受前端UX的HTTP请求。 网关 gateway 进行数据验证及格式转换，通过集成各个 Kitex-Client 转发给各个微服务 Kitex-Server。

<img width="878" height="536" alt="image" src="https://github.com/user-attachments/assets/292bc136-4f9a-4d51-85b1-680849326fba" />


1. IDL生成Kitex框架代码。

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=N2E1OGQ1ZWNiMTBlODUyY2Q5OTJmYjc4NjNlNjBkNDNfNjZ4a1ZBNnZxbnpnNWRjcjNmbUV2SlkyZ0JVYTFjaG5fVG9rZW46U2wzVWI4UVF0b0FGTEZ4SUVNb2NvM3lKblBkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

1. 业务流程图——以结算服务为例

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NjlmM2E3NzE3ZGUxODIzZTg4ZjNmYjJiMmFkZDEzYmJfYXQwS1Q0Z0xmaW41cjB2QjNQUlBYeHdIU0c4bWlzbklfVG9rZW46Skl3UWJ4eHlKbzUxbkd4bWJwYWN6aWFJbkZjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

#### 3.2.2**技术选型**  

- 使用Kitex框架开发微服务，Hertz框架实现api网关；
- 数据存储采用 MySQL和 Redis 进行缓存支持；  
- 使用 RabbitMQ 实现异步消息处理、实现定时任务（如取消超时订单、定时取消支付）；
- 采用 Consul 进行服务注册与发现；  
- 对图片处理采用 MinIO、Meilisearch 等技术实现对象存储和检索；  
- 使用JWT认证，Casbin进行权限控制；
- 使用Eino框架调用Doubao-pro-32k模型实现查询订单和模拟自动下单。

#### 3.2.3 可观测性

- 链路追踪：OpenTelemetry
- 监控可视化：Grafana
- 链路追踪可视化：Jaeger
- 日志：标准日志+logrus

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=M2IzOWIzNjlmNGQzMGVjZGU5Y2IwNzg0YzczODgyNzJfU3BWVzl3SkZYemxtR292d2JZemtZZ2RQSkF1REZ5OENfVG9rZW46SFRFdGJKNUpYb0ZoQXF4SEExRmM1THNYbkFWXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZDk2MWQxODczMTJjNGMxNDg5YWE2ZmQ2ODAwYjVhOTBfeXVtbzV3UlJlZGpWT3J4TElYZ3pMZXJwMGlUejhlTHdfVG9rZW46Vjhpd2JrcWtrb2tQV0h4dEttUGMweDhHblFNXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=OTNjZjA2ZGI0YWYzOWFkMTRkMTdiNDc5YmE4MzEyYjZfSGxweWRNVmJMNTNMTEg5N2RMdXdDSWI1WXJvMmhwckNfVG9rZW46S0tSZGI0RmRIb21KSjR4UU9sTWN0ZDlTblBlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

#### 3.2.4 **部署与运维** 

- 所有服务采用容器化部署，通过 Docker 和 Docker Compose 快速构建开发和测试环境；

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YjMyMmM2NjhjOTU0NjM2OGQxMzk3MDY3YmY4YTRkN2VfSVRJZTZTeklBbVRpSUJkeVlFbjlrSDFGT2I4dmVSeWtfVG9rZW46TGQ4RWJGQUVub3F5aXJ4dHVXaWNlczlzbmdlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

- 生产环境通过 Kubernetes(Minikube) 实现各微服务的伸缩、监控和自动化运维；
- 集成了限流、熔断、负载均衡策略。

#### 3.2.5 关系数据库设计

##### 3.2.5.1 用户表(user)

| 列名            | 数据类型     | 约束                        | 索引 | 备注                         |
| --------------- | ------------ | --------------------------- | ---- | ---------------------------- |
| id              | INT          | PRIMARY KEY, AUTO_INCREMENT | 是   | 主键，自动递增               |
| created_at      | DATETIME     | NOT NULL                    | 否   | 记录创建时间                 |
| updated_at      | DATETIME     | NOT NULL                    | 否   | 记录更新时间                 |
| deleted_at      | DATETIME     | DEFAULT NULL                | 否   | 记录删除时间，可为空         |
| email           | VARCHAR(255) | NOT NULL                    | 是   | 用户的电子邮件地址           |
| password_hashed | VARCHAR(255) | NOT NULL                    | 否   | 哈希后的密码                 |
| role            | INT          | NOT NULL                    | 否   | 角色，1表示管理员，2表示用户 |
| is_banned       | BOOL         | NOT NULL                    | 否   | 是否被封禁，默认false        |

##### 3.2.5.2 商品表(product)

| 列名        | 数据类型 | 约束       | 索引 | 备注     |
| ----------- | -------- | ---------- | ---- | -------- |
| id          | uint     | 主键，自增 | 无   | 产品 ID  |
| name        | string   | 无         | 无   | 产品名称 |
| description | string   | 无         | 无   | 产品描述 |
| picture     | string   | 无         | 无   | 产品图片 |
| price       | float32  | 无         | 无   | 产品价格 |
| categories  | string   | 无         | 无   | 产品分类 |
| created_at  | datetime | 无         | 无   | 创建时间 |
| updated_at  | datetime | 无         | 无   | 更新时间 |
| deleted_at  | datetime | 无         | 无   | 删除时间 |

##### 3.2.5.3 订单表(order)

| 列名           | 数据类型    | 约束            | 索引        | 备注               |
| -------------- | ----------- | --------------- | ----------- | ------------------ |
| id             | uint        | 主键，自增      | 无          | 订单ID             |
| order_id       | string(256) | 唯一            | uniqueIndex | 订单编号           |
| user_id        | uint32      | 无              | 无          | 用户ID             |
| user_currency  | string      | 无              | 无          | 用户使用的货币类型 |
| email          | string      | 嵌入(Consignee) | 无          | 收货人邮箱         |
| recipient_name | string      | 嵌入(Consignee) | 无          | 收货人姓名         |
| phone_number   | string      | 嵌入(Consignee) | 无          | 收货人电话         |
| street_address | string      | 嵌入(Consignee) | 无          | 街道地址           |
| city           | string      | 嵌入(Consignee) | 无          | 城市               |
| state          | string      | 嵌入(Consignee) | 无          | 州/省              |
| country        | string      | 嵌入(Consignee) | 无          | 国家               |
| zip_code       | int32       | 嵌入(Consignee) | 无          | 邮政编码           |
| order_state    | string      | 无              | 无          | 订单状态           |
| cancel_time    | datetime    | 可为空          | 无          | 取消时间           |
| cancel_type    | string      | 无              | 无          | 取消类型           |
| version        | uint32      | 默认值1         | 无          | 乐观锁版本号       |
| created_at     | datetime    | 无              | 无          | 创建时间           |
| updated_at     | datetime    | 无              | 无          | 更新时间           |
| deleted_at     | datetime    | 无              | 无          | 删除时间           |

##### 3.2.5.4 订单项表(order_item)

| 列名           | 数据类型    | 约束       | 索引  | 备注         |
| -------------- | ----------- | ---------- | ----- | ------------ |
| id             | uint        | 主键，自增 | 无    | 订单项ID     |
| product_id     | uint32      | 无         | 无    | 产品ID       |
| order_id_refer | string(256) | 外键       | index | 关联的订单ID |
| quantity       | int32       | 无         | 无    | 商品数量     |
| cost           | float32     | 无         | 无    | 商品成本     |
| created_at     | datetime    | 无         | 无    | 创建时间     |
| updated_at     | datetime    | 无         | 无    | 更新时间     |
| deleted_at     | datetime    | 无         | 无    | 删除时间     |

### 3.3 项目代码介绍

#### 3.3.1 项目目录介绍

```bash
├─.github                          # GitHub相关配置目录
│  └─workflows                     # GitHub Actions工作流配置目录，用于CI/CD自动化
├─app                              # 核心业务代码目录，包含各业务模块
│  ├─ai                            # AI模块，处理与大模型相关的业务逻辑及服务
│  │  ├─biz                        # 业务逻辑处理目录
│  │  ├─conf                       # 配置目录
│  │  ├─errno                      # 错误码定义目录
│  │  ├─infra                      # 基础设施目录
│  │  │  └─rpc                     # RPC调用目录
│  │  ├─pkg                        # 公共包目录
│  │  │  └─tool                    # 工具函数目录
│  │  └─script                     # 脚本目录，存放自动化、部署、初始化脚本
│  ├─auth                          # 认证/权限模块，负责用户认证、权限校验等
│  │  ├─biz                        # 业务逻辑处理目录
│  │  ├─conf                       # 配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  ├─errno                      # 错误码目录
│  │  └─script                     # 认证相关脚本目录
│  ├─cart                          # 购物车模块，管理用户购物车数据和逻辑
│  │  ├─biz                        # 业务逻辑处理目录
│  │  ├─cmd                        # 微服务启动入口目录
│  │  ├─conf                       # 配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  ├─errno                      # 错误码目录
│  │  ├─infra                      # 基础设施目录
│  │  │  └─rpc                     # RPC调用目录
│  │  ├─script                     # 脚本目录，存放初始化、部署脚本
│  │  └─utils                      #购物车工具目录
│  ├─checkout                      # 结账模块，处理订单支付和结算流程
│  │  ├─biz                        # 业务逻辑处理目录
│  │  ├─conf                       # 配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  ├─infra                      # 基础设施目录
│  │  │  └─rpc                     # RPC调用目录
│  │  ├─testClient                 #客户端测试目录
│  │  └─script                     # 脚本目录
│  ├─gateway                       # API网关，统一入口和请求路由转发
│  │  ├─biz                        # 网关业务逻辑处理目录
│  │  ├─conf                       # 网关配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  ├─infra                      # 基础设施目录
│  │  │  └─rpc                     # RPC调用目录
│  │  ├─hertz_gen                  # Hertz代码生成目录（自动生成API、路由等代码）
│  │  │  ├─api                     # API相关生成代码
│  │  │  ├─cart                    # 购物车相关生成代码
│  │  │  ├─gateway                 # 网关相关生成代码
│  │  │  └─order                   # 订单相关生成代码
│  │  ├─script                     # 脚本目录，存放初始化、部署脚本
│  │  └─utils                      # 网关工具目录，封装常量、公共函数等
│  ├─order                         # 订单模块，管理订单数据和业务流程
│  │  ├─biz                        # 订单业务逻辑处理目录
│  │  ├─conf                       # 订单模块配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  ├─error                      # 订单错误处理目录
│  │  └─script                     # 订单相关脚本目录
│  ├─payment                       # 支付模块，处理支付逻辑和对接第三方接口
│  │  ├─biz                        # 支付业务逻辑处理目录
│  │  ├─conf                       # 支付模块配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  └─script                     # 支付相关脚本目录
│  ├─product                       # 商品模块，管理商品数据、库存及展示逻辑
│  │  ├─biz                        # 商品业务逻辑处理目录
│  │  ├─conf                       # 商品模块配置目录
│  │  │  ├─dev                     # 开发环境配置
│  │  │  ├─online                  # 线上环境配置
│  │  │  └─test                    # 测试环境配置
│  │  ├─script                     # 商品相关脚本目录
│  │  └─utils                      # 商品工具目录，封装常量、公共函数等
│  └─user                          # 用户模块，管理用户信息、注册、登录及相关业务逻辑
│     ├─biz                        # 用户业务逻辑处理目录
│     ├─conf                       # 用户模块配置目录
│     │  ├─dev                     # 开发环境配置
│     │  ├─online                  # 线上环境配置
│     │  └─test                    # 测试环境配置
│     ├─script                     # 用户相关脚本目录
│     └─errno                      # 用户错误码目录
├─common                           # 公共工具库，提供跨模块共享的工具和辅助功能
│  ├─clientsuite                   # 客户端工具库目录，封装对外服务调用
│  ├─healthcheck                   # 系统健康检查模块目录
│  ├─mtl                           # 日志、指标、链路追踪等中间件工具目录
│  ├─serversuite                   # 服务器工具库目录，提供基础服务封装
│  └─utils                         # 通用工具函数集合目录
├─db                               # 数据库相关目录
│  └─sql                           # SQL脚本目录
│     └─ini                        # 数据库初始化脚本和配置文件目录
├─deployment                       # 部署与运维配置目录，存放Docker、Kubernetes等部署配置及文档
├─idl                              # 接口定义文件目录（IDL/Proto），定义模块间通信接口和数据结构
├─rpc_gen                          # RPC代码生成目录，用于存放自动生成的远程调用代码
│  ├─kitex_gen                     # 基于Kitex框架生成的RPC代码目录
│  └─rpc                           # 其他RPC代码生成目录，包含具体RPC实现
├─scripts                          # 辅助脚本目录，存放构建、测试、部署等自动化脚本
├─.gitignore                       # Git忽略配置文件，指定不加入版本控制的文件和目录 
├─.golangci.yml                    # GolangCI配置文件，用于代码静态检查和质量控制 
├─all.sh                           # 一键构建脚本，集成编译、测试、部署等常用任务 
├─docker-compose.yaml              # Docker Compose配置文件，定义多容器应用的编排和运行 
├─dockerps.txt                     # Docker进程状态记录文件，跟踪当前容器运行状态 
├─go.work                          # Go工作区配置文件，用于管理多个模块的依赖关系 
├─go.work.sum                      # Go工作区依赖校验文件，确保依赖版本一致性 
├─LICENSE                          # 开源许可协议文件，规定代码的使用、分发和贡献规则 
├─Makefile                         # Makefile构建文件，定义自动化构建、测试、部署任务 
├─otel-collector-config.yaml       # OpenTelemetry Collector配置文件，用于采集日志、指标和追踪数据 
└─README.md                        # 项目总体说明文档，介绍项目背景、架构设计及使用方法   
```

#### 3.3.2 典型代码介绍

##### 3.3.2.1 微服务启动入口——以user服务为例

```go
func main() {
    // 读取环境变量
    //_ = godotenv.Load()

    // 初始化数据库服务
    dal.Init()

    // 初始化kitex服务
    opts := kitexInit()

    // 解析服务地址
    addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
    if err != nil {
       panic(err)
    }
    opts = append(opts, server.WithServiceAddr(addr))

    // 链路追踪
    p := provider.NewOpenTelemetryProvider(
       provider.WithServiceName(conf.GetConf().Kitex.Service),
       provider.WithExportEndpoint(conf.GetConf().OpenTelemetry.Endpoint),
       provider.WithInsecure(),
    )
    defer func(p provider.OtelProvider, ctx context.Context) {
       err := p.Shutdown(ctx)
       if err != nil {
          klog.Error(err.Error())
       }
    }(p, context.Background())

    svr := userservice.NewServer(new(UserServiceImpl), opts...)

    err = svr.Run()
    if err != nil {
       klog.Error(err.Error())
    }
}

func kitexInit() (opts []server.Option) {
    // service info
    opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
       ServiceName: conf.GetConf().Kitex.Service,
    }))

    // 服务注册与发现
    r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
    if err != nil {
       panic(err)
    }
    opts = append(opts, server.WithRegistry(r))

    // 限流处理
    opts = append(opts, server.WithLimit(&limit.Option{
       MaxConnections: conf.GetConf().Kitex.MaxConnections, // MaxConnections: 最大连接数
       MaxQPS:         conf.GetConf().Kitex.MaxQPS,         // MaxQPS: 每秒最大请求数
    }))

    opts = append(opts, server.WithSuite(tracing.NewServerSuite()))

    // klog
    logger := kitexlogrus.NewLogger()
    klog.SetLogger(logger)
    klog.SetLevel(conf.LogLevel())
    asyncWriter := &zapcore.BufferedWriteSyncer{
       WS: zapcore.AddSync(&lumberjack.Logger{
          Filename:   conf.GetConf().Kitex.LogFileName,
          MaxSize:    conf.GetConf().Kitex.LogMaxSize,
          MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
          MaxAge:     conf.GetConf().Kitex.LogMaxAge,
       }),
       FlushInterval: time.Minute,
    }
    klog.SetOutput(asyncWriter)
    server.RegisterShutdownHook(func() {
       _ = asyncWriter.Sync()
    })
    return
}
```

##### 3.3.2.2 AI模拟自动下单的代码

```Go
type AutoOrderService struct {
    ctx context.Context
} // NewAutoOrderService new AutoOrderService
func NewAutoOrderService(ctx context.Context) *AutoOrderService {
    return &AutoOrderService{ctx: ctx}
}

// Run create note info
func (s *AutoOrderService) Run(req *ai.AutoOrderReq) (resp *ai.AutoOrderResp, err error) {
    // Finish your business logic.
    rpc.InitClient()

    chatModel, err := pkg.CreateARKModel(s.ctx)
    if err != nil {
       err = errno.CreateChatModelErr(err)
       klog.Error(err)
       return
    }
    searchProductTool := autoOrderTool.GetSearchProductTool()
    addToCartTool := autoOrderTool.GetAddToCartTool()
    checkoutTool := autoOrderTool.GetCheckoutTool()
    tools := []einoTool.BaseTool{
       searchProductTool,
       addToCartTool,
       checkoutTool,
    }

    persona := `你是一个帮助用户搜索商品，并且下单的助手，根据用户的需要，查询商品信息，并将查到的商品加入到购物车，等商品都加入到购物车后，进行结算。注意按照用户输入的商品数量进行下单！
请将下单后的订单信息按照json对象的形式进行返回，例如：
       [{
          "order_id": "12345",
          "user_id": 67890,
          "user_currency": "USD",
          "email": "user@example.com",
          "created_at": "2023-10-01T12:34:56Z",
          "order_items": [
             {
             "product_id": 1,
             "product_name": "Product A",
             "quantity": 2,
             "cost": 19.99
             },
             {
             "product_id": 2,
             "product_name": "Product B",
             "quantity": 1,
             "cost": "9.99"
             }
          ],
          "orderState": "placed"
       }]
注意，只返回json形式的数据即可，不要有多余的文字输出，如果没有创建订单，就输出“{}”！
`

    input := fmt.Sprintf("为user_id是%d的用户%s", req.UserId, req.Content)
    ragent, err := react.NewAgent(s.ctx, &react.AgentConfig{
       Model: chatModel,
       ToolsConfig: compose.ToolsNodeConfig{
          Tools: tools,
       },

       MessageModifier: react.NewPersonaModifier(persona),
    })
    if err != nil {
       err = errno.CreateAgentErr(err)
       klog.Error(err)
       return
    }

    sr, err := ragent.Generate(s.ctx, []*schema.Message{
       {
          Role:    schema.User,
          Content: input,
       },
    }, agent.WithComposeOptions(compose.WithCallbacks(&pkg.LoggerCallback{})))
    if err != nil {
       err = errno.StreamErr(err)
       klog.Error(err)
       return
    }
    klog.Infof("===== start streaming =====\n\n")
    order, err := pkg.ConvertToAiOrderView(sr.Content)
    if err != nil {
       err = errno.ConvertToAiOrderViewErr(err)
       klog.Error(err)
       return nil, err
    }
    klog.Infof("===== finished =====\n")

    return &ai.AutoOrderResp{Order: order}, nil
}
```

自动下单相关的tools代码：

```go
func GetSearchProductTool() tool.InvokableTool {
    return &SearchProductsTool{}
}

func GetAddToCartTool() tool.InvokableTool {
    return &AddToCartTool{}
}

func GetCheckoutTool() tool.InvokableTool {
    return &CheckoutTool{}
}

type SearchProductParam struct {
    ProductName string `json:"product_name"`
    Quantity    int32  `json:"quantity"`
    Topn        int64  `json:"topn"`
}

type SearchProductsTool struct{}

func (s *SearchProductsTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    return &schema.ToolInfo{
       Name: "search products",
       Desc: "query the specified product based on the product name",
       ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
          "product_name": {
             Type:     "string",
             Desc:     "The name of one product",
             Required: true,
          },
          "topn": {
             Type: "number",
             Desc: "top n products sorted by prices",
          },
       }),
    }, nil
}

func (s *SearchProductsTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
    // 解析参数
    p := &SearchProductParam{}
    err := json.Unmarshal([]byte(argumentsInJSON), p)
    if err != nil {
       return "", err
    }

    if p.Topn == 0 {
       p.Topn = 1
    }

    // 调用商品服务查找特定名称的商品
    rests, err := rpc.ProductClient.SearchProductsByName(ctx, &rpcproduct.SearchProductsByNameReq{
       Query:    p.ProductName,
       Page:     1,
       PageSize: p.Topn,
       Flag:     false,
    })
    if err != nil {
       klog.Errorf("SearchProductsByName err: %v", err)
       return "", err
    }

    // 序列化结果
    res, err := json.Marshal(rests.Results)
    if err != nil {
       klog.Error(err)
       return "", err
    }

    return string(res), nil
}

type AddToCartParam struct {
    UserID    uint32 `json:"user_id"`
    ProductId uint32 `json:"id"`
    Quantity  int32  `json:"quantity"`
}
type AddToCartTool struct{}

func (a *AddToCartTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    return &schema.ToolInfo{
       Name: "add products to cart",
       Desc: "add the selected items to the shopping cart.",
       ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
          "user_id": {
             Type:     "number",
             Desc:     "The id of user",
             Required: true,
          },
          "id": {
             Type:     "number",
             Desc:     "The id of one product",
             Required: true,
          },
          "quantity": {
             Type:     "number",
             Desc:     "the number of products that the user want to buy",
             Required: true,
          },
       }),
    }, nil
}

func (a *AddToCartTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
    // 解析参数
    p := &AddToCartParam{}
    err := json.Unmarshal([]byte(argumentsInJSON), p)
    if err != nil {
       klog.Error(err)
       return "", err
    }

    // 调用购物车服务将商品添加到购物车
    _, err = rpc.CartClient.AddItem(ctx, &rpccart.AddItemReq{
       Item: &rpccart.CartItem{
          ProductId: p.ProductId,
          Quantity:  p.Quantity,
       },
       UserId: p.UserID,
    })
    if err != nil {
       klog.Errorf("AddItem err: %v", err)
       return "", err
    }

    // 返回购物车信息
    cart, err := rpc.CartClient.GetCart(ctx, &rpccart.GetCartReq{
       UserId: p.UserID,
    })
    if err != nil {
       klog.Errorf("GetCart err: %v", err)
       return "", err
    }
    // 序列化结果
    res, err := json.Marshal(cart.Cart.Items)
    if err != nil {
       klog.Error(err)
       return "", err
    }

    return string(res), nil
}

type CheckoutTool struct{}

func (c *CheckoutTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    return &schema.ToolInfo{
       Name: "checkout",
       Desc: "settle the payment based on the items in the user's shopping cart, create an order, and return the created order information.",
       ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
          "user_id": {
             Type:     "number",
             Desc:     "The id of user",
             Required: true,
          },
       }),
    }, nil
}

func (c *CheckoutTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
    // 解析参数
    p := &rpcorder.PlaceOrderReq{}
    err := json.Unmarshal([]byte(argumentsInJSON), p)
    if err != nil {
       return "", err
    }

    // 调用结算服务进行订单结算
    checkoutResp, err := rpc.CheckoutClient.Checkout(ctx, &rpccheckout.CheckoutReq{
       UserId:    p.UserId,
       Firstname: "user",
       Lastname:  "user",
       Address: &rpccheckout.Address{
          StreetAddress: "123 Main St",
          City:          "Beijing",
          State:         "Beijing",
          Country:       "China",
          ZipCode:       "0",
       },
       Email: "user@example.com",
       CreditCard: &rpcpayment.CreditCardInfo{
          CreditCardNumber:          "5302079249905900",
          CreditCardCvv:             123,
          CreditCardExpirationMonth: 12,
          CreditCardExpirationYear:  2025,
       },
    })
    if err != nil {
       klog.Errorf("checkout failed: %s", err)
       return "", err
    }

    // 获取下单后的订单信息
    orderResp, err := rpc.OrderClient.GetOrder(ctx, &rpcorder.GetOrderReq{
       UserId:  p.UserId,
       OrderId: checkoutResp.OrderId,
    })
    if err != nil {
       klog.Error(err)
       return "", err
    }

    // 将订单信息组合成SearchOrdersResult结构体返回
    orderItems := make([]OrderItem, 0)
    for _, item := range orderResp.Order.Order.OrderItems {
       product, err := rpc.ProductClient.GetProduct(ctx, &rpcproduct.GetProductReq{Id: item.Item.ProductId})
       if err != nil {
          klog.Error("get product name failed: %s", err)
          return "", err
       }
       orderItems = append(orderItems, OrderItem{
          ProductId:   item.Item.ProductId,
          ProductName: product.Product.Name,
          Quantity:    item.Item.Quantity,
          Cost:        item.Cost,
       })
    }
    order := SearchOrdersResult{
       OrderId:      orderResp.Order.Order.OrderId,
       UserId:       orderResp.Order.Order.UserId,
       UserCurrency: orderResp.Order.Order.UserCurrency,
       Email:        orderResp.Order.Order.Email,
       CreatedAt:    convertInt32ToTime(orderResp.Order.Order.CreatedAt),
       OrderItems:   orderItems,
       OrderState:   orderResp.Order.OrderState,
    }

    res, err := json.Marshal(order)
    if err != nil {
       klog.Error(err)
       return "", err
    }

    return string(res), nil
}

func convertInt32ToTime(timestamp int32) time.Time {
    seconds := int64(timestamp)

    return time.Unix(seconds, 0)
}
```

##### 3.3.2.3 **消息队列实现定时任务**

该代码实现使用RabbitMQ的延时队列和死信队列来处理订单超时，同时通过乐观锁保证了并发安全性。

1. **消费者(Consumer)**

- 负责处理死信队列中的超时订单消息
- 使用乐观锁机制处理订单，防止并发问题
- 实现了重试机制和消息确认机制
- 支持优雅关闭

```Go
// order/biz/dal/mq/consumer.go
// Consumer 消费者结构体
type Consumer struct {
    conn       *amqp.Connection
    ctx        context.Context
    channel    *amqp.Channel
    done       chan struct{}
    maxRetries int
    DB         *gorm.DB
    orderQuery model.OrderQuery
}

// NewConsumer 创建消费者
func NewConsumer(db *gorm.DB) (*Consumer, error) {
    channel, err := RabbitMQConn.Channel()
    if err != nil {
        return nil, err
    }
    ctx := context.Background()
    consumer := &Consumer{
        conn:       RabbitMQConn,
        ctx:        ctx,
        channel:    channel,
        done:       make(chan struct{}),
        maxRetries: conf.GetConf().RabbitMQ.MaxRetries,
        DB:         db,
        orderQuery: model.NewOrderQuery(ctx, db),
    }

    klog.CtxInfof(ctx, "RabbitMQ Consumer 初始化成功")
    return consumer, nil
}

// handleOrderWithOptimisticLock 使用乐观锁处理订单
func (c *Consumer) handleOrderWithOptimisticLock(orderMsg OrderMessage) error {
    var err error
    klog.CtxInfof(c.ctx, "正在处理订单: %d", orderMsg.OrderID)

    for i := 0; i < c.maxRetries; i++ {
        version, orderState, err := c.orderQuery.GetOrderVersionAndState(orderMsg.OrderID)
        if err != nil {
            klog.CtxErrorf(c.ctx, "获取订单版本号失败: %v", err)
            return err
        }
        // 如果订单状态不是已下单，不处理 -> 已被其他消费者处理过了 || 订单已取消、已完成
        if orderState != model.OrderStatePlaced {
            klog.CtxInfof(c.ctx, "订单 %d 状态不是已下单，不处理", orderMsg.OrderID)
            return nil
        }
        err = c.orderQuery.CancelOrderWithVersion(orderMsg.OrderID, model.CancelTypeTimeout, int32(time.Now().Unix()), version)
        if err == nil {
            klog.CtxInfof(c.ctx, "订单 %d 处理成功", orderMsg.OrderID)
            return nil
        }
        // 如果是乐观锁冲突，继续重试
        klog.CtxWarnf(c.ctx, "乐观锁冲突，正在重试 (%d/%d)", i+1, c.maxRetries)
    }

    return fmt.Errorf("达到最大重试次数，处理订单失败: %v", err)
}

// Consume 消费者消费消息
func (c *Consumer) Consume() error {
    // 设置预取计数，控制消费者同时处理的消息数量
    err := c.channel.Qos(1, 0, false)
    if err != nil {
        klog.CtxErrorf(c.ctx, "设置RabbitMQ Consumer预取计数失败: %v", err)
        return err
    }

    msgs, err := c.channel.Consume(
        DeadLetterQueue,
        "",    // 消费者标识
        false, // 自动确认
        false, // 独占
        false, // 非阻塞
        false, // 等待服务器确认
        nil,
    )
    if err != nil {
        klog.CtxErrorf(c.ctx, "RabbitMQ Consumer start failed: %v", err)
        return err
    }

    go func() {
        for msg := range msgs {
            var orderMsg OrderMessage
            if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
                klog.CtxErrorf(c.ctx, "Consumer解析订单消息失败: %v", err)
                _ = msg.Nack(false, false)
                continue
            }

            // 使用乐观锁处理订单
            err := c.handleOrderWithOptimisticLock(orderMsg)
            if err != nil {
                klog.CtxErrorf(c.ctx, "Consumer处理订单失败: %v", err)
                if err == gorm.ErrRecordNotFound {
                    _ = msg.Ack(false) // 订单不存在，直接确认
                    continue
                }
                _ = msg.Nack(false, true) // 重新入队
                continue
            }

            _ = msg.Ack(false)
        }
    }()

    klog.CtxInfof(c.ctx, "Consumer start successfully, listening dead letter queue: %s", DeadLetterQueue)
    <-c.done
    klog.CtxInfof(c.ctx, "Consumer Stopped!")
    return nil
}

// Stop 停止消费者
func (c *Consumer) Stop() {
    close(c.done)
    if c.channel != nil {
        _ = c.channel.Close()
    }
    if c.conn != nil {
        _ = c.conn.Close()
    }
}
```

2. **生产者(Producer)**

- 负责发送订单消息到延迟队列
- 实现了两组交换机和队列的配置：
  - 延迟交换机(order.delay.exchange)和延迟队列(order.delay.queue)
  - 死信交换机(order.dlx.exchange)和死信队列(order.dlx.queue)
- 消息在延迟队列中超时后会自动转发到死信队列

```Go
// order/biz/dal/mq/producer.go
const (
    DelayExchange      = "order.delay.exchange"
    DelayQueue         = "order.delay.queue"
    DeadLetterExchange = "order.dlx.exchange"
    DeadLetterQueue    = "order.dlx.queue"
)

var ProducerInstance *Producer // ProducerInstance 生产者实例

// OrderMessage 订单消息结构体
type OrderMessage struct {
    OrderID string `json:"order_id"`
}

// Producer 生产者结构体
type Producer struct {
    conn    *amqp.Connection
    ctx     context.Context
    channel *amqp.Channel
}

// NewProducer 创建生产者实例
func NewProducer(orderTimeout int) (*Producer, error) {
    ctx := context.Background()
    channel, err := RabbitMQConn.Channel()
    if err != nil {
        klog.CtxErrorf(ctx, "RabbitMQ Producer 初始化失败, err: %v", err)
        return nil, err
    }

    producer := &Producer{
        conn:    RabbitMQConn,
        ctx:     ctx,
        channel: channel,
    }

    err = producer.initializeQueue(orderTimeout)
    if err != nil {
        klog.CtxErrorf(ctx, "RabbitMQ Producer 初始化失败, 无法初始化队列, err: %v", err)
        return nil, err
    }
    klog.CtxInfof(ctx, "RabbitMQ Producer 初始化成功")
    return producer, nil
}

// initializeQueue 初始化交换机和队列
// NOTE！
// 如果使用完全相同的参数重复声明队列/交换机，RabbitMQ 会直接返回成功. 这是一个幂等操作，不会对现有队列造成任何影响
// 如果尝试用不同的参数重新声明一个已存在的队列/交换机，RabbitMQ 会返回错误
func (p *Producer) initializeQueue(orderTimeout int) error {
    // 声明死信交换机
    err := p.channel.ExchangeDeclare(
        DeadLetterExchange,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 初始化失败, 无法初始化死信交换机, err: %v", err)
        return err
    }

    // 声明死信队列
    _, err = p.channel.QueueDeclare(
        DeadLetterQueue,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 初始化失败, 无法初始化死信队列, err: %v", err)
        return err
    }

    // 绑定死信队列到死信交换机
    err = p.channel.QueueBind(
        DeadLetterQueue,
        DeadLetterQueue,
        DeadLetterExchange,
        false,
        nil,
    )
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 初始化失败, 无法绑定死信队列到死信交换机,err: %v", err)
        return err
    }

    // 声明延迟交换机
    err = p.channel.ExchangeDeclare(
        DelayExchange,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 初始化失败, 无法初始化延迟交换机,err: %v", err)
        return err
    }

    // 声明延迟队列，设置消息过期后转发到死信交换机
    args := amqp.Table{
        "x-dead-letter-exchange":    DeadLetterExchange,
        "x-dead-letter-routing-key": DeadLetterQueue,
        "x-message-ttl":             orderTimeout,
    }
    _, err = p.channel.QueueDeclare(
        DelayQueue,
        true,
        false,
        false,
        false,
        args,
    )
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 初始化失败, 无法初始化延迟队列,err: %v", err)
        return err
    }

    // 绑定延迟队列到延迟交换机
    err = p.channel.QueueBind(
        DelayQueue,
        DelayQueue,
        DelayExchange,
        false,
        nil,
    )
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 初始化失败, 无法绑定延迟队列到延迟交换机,err: %v", err)
        return err
    }
    return nil
}

// Stop 关闭连接
func (p *Producer) Stop() {
    if p.channel != nil {
        _ = p.channel.Close()
    }
    if p.conn != nil {
        _ = p.conn.Close()
    }
    klog.CtxInfof(p.ctx, "RabbitMQ Producer 关闭成功")
}

// SendDelayMessage 发送延迟消息
func (p *Producer) SendDelayMessage(orderID string) error {
    message := OrderMessage{
        OrderID: orderID,
    }

    body, err := json.Marshal(message)
    if err != nil {
        klog.CtxErrorf(p.ctx, "RabbitMQ Producer 发送消息失败, err: %v", err)
        return err
    }

    return p.channel.Publish(
        DelayExchange,
        DelayQueue,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
```

3. **启动**

- 在主程序启动时初始化生产者和消费者
- 提供了优雅关闭的钩子函数
- 消费者以异步方式运行，不阻塞主线程

```Go
// order/main.go
var consumer *mq.Consumer

func main() {
    ... ...
    existing code

    // 初始化MySQL和RabbitMQ
    dal.Init()
    // 初始化Kitex
    opts := kitexInit()

    startProducer()
    // defer mq.ProducerInstance.Stop()
    startConsumer(mysql.DB)
    // defer consumer.Stop()
    ... ...
    existing code
}

func kitexInit() (opts []server.Option) {
    ... ...
    existing code
    server.RegisterShutdownHook(func() {
        consumer.Stop()
        mq.ProducerInstance.Stop()
        _ = asyncWriter.Sync()
    })
    return
}

func startProducer() {
    klog.Info("Producer starting...")
    var err error
    mq.ProducerInstance, err = mq.NewProducer(conf.GetConf().RabbitMQ.OrderTimeout)
    if err != nil {
        klog.Fatalf("NewProducer failed, err: %v", err)
        panic(err)
    }
}

func startConsumer(db *gorm.DB) {
    klog.Info("Consumer starting...")
    var err error
    consumer, err = mq.NewConsumer(db)
    if err != nil {
        klog.Fatalf("NewConsumer failed, err: %v", err)
        panic(err)
    }
    go func() {
        _ = consumer.Consume()
    }()
}
```

##### 3.3.2.4 **部署配置**

以otel-collector和order service的部署为例，展示微服务部署到Kubernetes。

1. **otel-collector**

- 通过 ConfigMap 定义了 collector 的配置，支持 OTLP 接收器，配置了 Prometheus 和 Jaeger 导出器
- Service 配置暴露了 4317 端口用于 OTLP gRPC 通信
- Deployment 配置使用 otel-collector-contrib 镜像，设置了多个端口用于监控、指标收集和健康检查

```YAML
# otel-collector-deployment.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: otel-collector-config
data:
  otel-collector-config.yaml: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317

    exporters:
      prometheusremotewrite:
        endpoint: "http://victoriametrics-service:8428/api/v1/write"

      debug:

      otlp/jaeger:
        endpoint: jaeger-service:4317
        tls:
          insecure: true

    processors:
      batch:

    extensions:
      health_check:
      pprof:
        endpoint: :1888
      zpages:
        endpoint: :55679

    service:
      extensions: [ pprof, zpages, health_check ]
      pipelines:
        traces:
          receivers: [ otlp ]
          processors: [ batch ]
          exporters: [ debug, otlp/jaeger ]
        metrics:
          receivers: [ otlp ]
          processors: [ batch ]
          exporters: [ debug, prometheusremotewrite ]
---
apiVersion: v1
kind: Service
metadata:
  name: otel-collector-service
spec:
  selector:
    app: otel-collector
  ports:
  - name: otlp-grpc
    port: 4317
    targetPort: 4317
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: otel-collector
spec:
  selector:
    matchLabels:
      app: otel-collector
  template:
    metadata:
      labels:
        app: otel-collector
    spec:
      containers:
      - name: otel-collector
        image: otel/opentelemetry-collector-contrib-dev:latest
        args:
        - "--config=/etc/otel-collector-config.yaml"
        ports:
        - containerPort: 1888  # pprof extension
        - containerPort: 8888  # Prometheus metrics
        - containerPort: 8889  # Prometheus exporter
        - containerPort: 13133 # health check
        - containerPort: 4317  # OTLP gRPC
        - containerPort: 55679 # zpages
        volumeMounts:
        - name: otel-config
          mountPath: /etc/otel-collector-config.yaml
          subPath: otel-collector-config.yaml
      volumes:
      - name: otel-config
        configMap:
          name: otel-collector-config
```

2. **order service**

- Dockerfile 采用多阶段构建：
  - 第一阶段使用 golang:1.22-alpine 编译服务
  - 第二阶段使用 alpine 作为基础镜像运行服务
  - 设置了时区为上海，暴露了 gRPC(8085)、OpenTelemetry(4317) 和健康检查(8095) 端口
- Kubernetes 配置：
  - ConfigMap 存储环境变量
  - Service 暴露 gRPC 和 OpenTelemetry 端口
  - Deployment 配置了 2 个副本，设置了资源限制和健康检查探针
  - 包含了详细的存活和就绪探针配置
- 健康检查实现：
  - 提供了一个简单的 HTTP 健康检查接口
  - 返回服务状态和服务名称的 JSON 响应
  - 通过独立的 goroutine 运行健康检查服务

```Dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.work go.work.sum ./
COPY rpc_gen/ rpc_gen/
COPY common/ common/
COPY app/order/ app/order/
COPY app/auth/ app/auth/
COPY app/user/ app/user/
COPY app/cart/ app/cart/
COPY app/product/ app/product/
COPY app/payment/ app/payment/
COPY app/gateway/ app/gateway/
COPY app/checkout/ app/checkout/

WORKDIR /build/app/order
ENV GO111MODULE=on \
    # GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

RUN go mod tidy
RUN sh build.sh


FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/app/order/output/ /app/
RUN ls -al /app/

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

ENV GO_ENV=dev

# grpc
EXPOSE 8085
# opentelemetry
EXPOSE 4317
# health check
EXPOSE 8095

CMD ["sh", "bootstrap.sh"]
# order-deployment.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: order-config
  namespace: default
data:
  GO_ENV: "online"
---
apiVersion: v1
kind: Service
metadata:
  name: order-service
  namespace: default
  labels:
    app: order
spec:
  selector:
    app: order
  ports:
    - name: grpc
      port: 8085
      targetPort: 8085
    - name: opentelemetry
      port: 4317
      targetPort: 4317
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-deployment
  namespace: default
  labels:
    app: order
spec:
  replicas: 2
  selector:
    matchLabels:
      app: order
  template:
    metadata:
      labels:
        app: order
    spec:
      containers:
        - name: order
          image: order-service:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8085
              name: grpc
            - containerPort: 4317
              name: opentelemetry
            - containerPort: 8095
              name: health
          env:
            - name: GO_ENV
              valueFrom:
                configMapKeyRef:
                  name: order-config
                  key: GO_ENV
          resources:
            requests:
              cpu: "100m"
              memory: "128Mi"
            limits:
              cpu: "500m"
              memory: "512Mi"
          livenessProbe:
            httpGet:
              path: /health
              port: 8095
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8095
            initialDelaySeconds: 5
            periodSeconds: 10
// health check
type HealthResponse struct {
    Status      string `json:"status"`
    ServiceName string `json:"service"`
}

func StartHealthCheck(addr string, serviceName string) {
    handler := newHealthCheckHandler(serviceName)
    http.HandleFunc("/health", handler.healthCheckHandler)

    // 在后台启动健康检查服务
    go func() {
        if err := http.ListenAndServe(addr, nil); err != nil {
            klog.Errorf("Health check server failed to start: %v", err)
        }
    }()
}

type healthCheckHandler struct {
    serviceName string
}

func newHealthCheckHandler(serviceName string) *healthCheckHandler {
    return &healthCheckHandler{
        serviceName: serviceName,
    }
}

func (h *healthCheckHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{
        Status:      "ok",
        ServiceName: h.serviceName,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(response)
}
```

##### 3.3.2.5 casbin实现的rbac

```go
//gateway/biz/middleware/casbin.go
type CasbinMiddleware struct {
    enforcer *casbin.Enforcer
}

// NewCasbinEnforcer 创建并初始化 Casbin Enforcer
func NewCasbinEnforcer(db *gorm.DB) (*CasbinMiddleware, error) {
    // 创建 GORM 适配器
    adapter, err := gormadapter.NewAdapterByDB(db)
    if err != nil {
       hlog.Error("Casbin创建gorm适配器失败: %v", err)
       return nil, err
    }
    // 加载模型
    _, filename, _, _ := runtime.Caller(0)
    basePath := filepath.Dir(filepath.Dir(filename))
    modelPath := filepath.Join(basePath, "model", "rbac.conf")
    enforcer, err := casbin.NewEnforcer(modelPath, adapter)
    enforcer.AddFunction("RegexMatch", RegexMatch)
    if err != nil {
       hlog.Error("创建Casbin模型失败: %v", err)
       return nil, err
    }
    // 从数据库加载策略
    err = enforcer.LoadPolicy()
    if err != nil {
       hlog.Error("加载Casbin策略失败: %v", err)
       return nil, err
    }

    if err := initDefaultPolicies(enforcer); err != nil {
       hlog.Error("初始化默认策略失败: %v", err)
       return nil, err
    }

    return &CasbinMiddleware{enforcer: enforcer}, nil
}

func (cm *CasbinMiddleware) Middleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
       var role string
       fmt.Println(c.Get("user_id"))
       // 从上下文中获取角色
       roleVal, exists := c.Get("role")
       if !exists {
          role = "public" // 如果没有角色，则默认为 public
       } else {
          switch v := roleVal.(type) {
          case uint32:
             if v == 1 {
                role = "admin"
             } else if v == 2 {
                role = "user"
             } else {
                role = "public"
             }
          case string:
             role = v
          default:
             role = "public" // 默认设置为 public
          }
       }
       // 获取请求信息
       obj := string(c.Request.URI().Path())
       act := string(c.Request.Method())

       // 权限验证
       ok, err := cm.enforcer.Enforce(fmt.Sprint(role), obj, act)

       if err != nil {
          hlog.Error("Casbin 权限验证失败: %v", err)
          c.AbortWithStatus(500)
          return
       }

       if !ok {
          c.AbortWithStatus(403)
          return
       }
       c.Next(ctx)
    }
}

func initDefaultPolicies(enforcer *casbin.Enforcer) error {
    // 管理员权限
    if _, err := enforcer.AddPolicy("admin", ".*", ".*"); err != nil {
       return fmt.Errorf("添加管理员策略失败: %w", err)
    }

    // 公共访问权限
    policies := [][]string{
       {"public", "/auth/token", "POST"},
       {"public", "/auth/verify", "POST"},
       {"public", "/auth/renew", "POST"},
       {"public", "/user/register", "POST"},
       {"public", "/user/login", "POST"},
       {"user", "/user/logout", "POST"},
       {"user", "/user/update", "PUT"},
       {"user", "/user/info", "GET"},
       {"admin", "/user/delete", "DELETE"},
       {"admin", "/user/update_role", "PUT"},
       {"public", "/products", "GET"},
       {"public", "/product", "GET"},
       {"public", "/product/search", "GET"},
       {"admin", "/product", "POST"},
       {"admin", "/product/upload", "POST"},
       {"admin", "/product", "PUT"},
       {"admin", "/product", "DELETE"},
       {"user", "/orders", "POST"},
       {"user", "/orders", "GET"},
       {"user", "/orders/.*", "PUT"},
       {"user", "/orders/.*", "DELETE"},
       {"public", "/checkout", "POST"},
       {"public", "/checkout", "GET"},
       {"public", "/checkout/.*", "PUT"},
       {"public", "/checkout/.*", "DELETE"},
    }

    for _, p := range policies {
       if _, err := enforcer.AddPolicy(p[0], p[1], p[2]); err != nil {
          return fmt.Errorf("添加策略%v失败:%w", p, err)
       }
    }
    // 保存策略
    if err := enforcer.SavePolicy(); err != nil {
       return fmt.Errorf("保存Casbin策略失败:%w", err)
    }
    return nil
}

func RegexMatch(args ...any) (any, error) {
    key, ok := args[0].(string)
    if !ok {
       return nil, errors.New("key错误")
    }

    pattern, ok := args[1].(string)
    if !ok {
       return nil, errors.New("pattern错误")
    }
    matched, err := regexp.MatchString("^"+pattern+"$", key)
    if err != nil {
       return false, err
    }
    return matched, nil
}
```

##### 3.3.2.6 图片上传

```go
func (s *UploadImageService) Run(req *product.UploadImageReq) (resp *product.UploadImageResp, err error) {
    // 请求验证
    if len(req.ImageData) == 0 {
       return nil, apiErr.ImageDataRequiredErr
    }
    if req.FileName == "" {
       return nil, apiErr.FileNameRequiredErr
    }

    // 压缩图片处理
    tempFile := "data/temp_compressed.jpg"
    if err := img.ConvertAndCompressImage(s.ctx, req.ImageData, req.FileName, tempFile); err != nil {
       return nil, apiErr.ConvertErr(err)
    }
    defer func() {
       if err := os.Remove(tempFile); err != nil {
          klog.CtxErrorf(s.ctx, "删除临时文件失败: %s", err)
       }
    }()
    // 打开压缩后的图片文件
    compressedFile, err := os.Open(tempFile)
    if err != nil {
       return nil, apiErr.ConvertErr(err)
    }
    defer func() {
       if err := compressedFile.Close(); err != nil {
          klog.CtxErrorf(s.ctx, "关闭压缩文件失败: %s", err)
       }
    }()
    // 计算图片大小
    info, err := os.Stat(tempFile)
    if err != nil {
       return nil, err
    }
    // 上传图片到对象存储
    objectKey := img.GenerateObjectKey("image", ".jpeg")
    objectUrl, err := img.PutObject(objectKey, compressedFile, info.Size(), "image/jpeg")
    if err != nil {
       return nil, apiErr.ConvertErr(err)
    }
    // 返回上传成功的图片URL
    return &product.UploadImageResp{
       ImageUrl: objectUrl,
    }, nil
}

func ConvertAndCompressImage(ctx context.Context, imgData []byte, fileName string, dstPath string) error {
    // 确保目标目录存在
    dir := "data"
    if _, err := os.Stat(dir); os.IsNotExist(err) {
       // 创建目录
       if err := os.MkdirAll(dir, os.ModePerm); err != nil {
          return fmt.Errorf("创建目录失败: %w", err)
       }
    }
    tmpFilePath := fmt.Sprintf("data/%s", fileName)
    dst, err := os.Create(tmpFilePath)
    if err != nil {
       return fmt.Errorf("创建文件失败: %w", err)
    }
    defer func() {
       err := dst.Close()
       if err != nil {
          klog.CtxErrorf(ctx, "关闭文件失败: %s", err)
       }
       if err := os.Remove(tmpFilePath); err != nil {
          klog.CtxErrorf(ctx, "删除临时文件失败: %s", err)
       }
    }()

    // 将字节切片写入文件
    _, err = io.Copy(dst, bytes.NewReader(imgData))
    if err != nil {
       return fmt.Errorf("写入文件失败: %w", err)
    }
    srcFile, err := os.Open(tmpFilePath)
    if err != nil {
       return err
    }
    defer func() {
       err := srcFile.Close()
       if err != nil {
          klog.CtxErrorf(ctx, "关闭文件失败: %s", err)
       }
    }()
    // 解码图像
    srcImg, _, err := image.Decode(srcFile)
    if err != nil {
       return fmt.Errorf("failed to decode image: %w", err)
    }

    f, err := os.Create(dstPath)
    if err != nil {
       return err
    }
    defer func(f *os.File) {
       err := f.Close()
       if err != nil {
          klog.CtxErrorf(ctx, "关闭文件失败: %s", err)
       }
    }(f)

    // 压缩并保存图像为 JPEG
    if err := jpeg.Encode(f, srcImg, &jpeg.Options{Quality: 100}); err != nil {
       return fmt.Errorf("failed to encode JPEG: %w", err)
    }

    return nil
}

func GenerateObjectKey(uploadType string, fileExt string) string {
    return fmt.Sprintf("%s/%d/%s%s", uploadType, time.Now().Year(), uuid.NewV1().String(), fileExt)
}

// ms 是全局的 MinioService 实例
var ms = &minioDal.MinioService

// PutObject 用于上传对象
func PutObject(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
    opts := minio.PutObjectOptions{ContentType: contentType}
    _, err := (*ms).Client.PutObject(context.Background(), (*ms).Bucket, objectKey, reader, size, opts)
    if err != nil {
       return "", err
    }
    return (*ms).Domain + (*ms).Bucket + "/" + objectKey, nil
}

func GetObjectKeyFromUrl(fullUrl string) (objectKey string, ok bool) {
    objectKey = strings.TrimPrefix(fullUrl, (*ms).Domain+(*ms).Bucket+"/")
    if objectKey == fullUrl {
       return "", false
    }
    return objectKey, true
}

// DeleteObject 用于删除相应对象
func DeleteObject(objectKey string) error {
    err := (*ms).Client.RemoveObject(
       context.Background(),
       (*ms).Bucket,
       objectKey,
       minio.RemoveObjectOptions{ForceDelete: true},
    )
    if err != nil {
       return fmt.Errorf("failed to delete object: %w", err)
    }
    return nil
}

// DeleteObjectByUrlAsync 通过给定的 Url 异步删除对象
func DeleteObjectByUrlAsync(ctx context.Context, url string) {
    objectKey, ok := GetObjectKeyFromUrl(url)
    if ok {
       go func(objectKey string) {
          err := DeleteObject(objectKey)
          if err != nil {
             klog.CtxErrorf(ctx, "failed to delete object: %v", err)
          }
       }(objectKey)
    }
}
```

##### 3.3.2.7 部分handler、service——以product服务为例

```go
//handler
func UploadImage(ctx context.Context, c *app.RequestContext) {
    var req product.UploadImageReq
    if err := c.BindAndValidate(&req); err != nil {
       utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err) // 使用 400 错误码表示请求无效
       return
    }

    // 获取上传的文件
    file, err := c.FormFile("file")
    if err != nil {
       utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err) // 使用 500 错误码表示服务器错误
       return
    }

    // 打开文件并保证关闭
    src, err := file.Open()
    if err != nil {
       utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
       return
    }
    defer func() {
       if closeErr := src.Close(); closeErr != nil {
          hlog.CtxErrorf(ctx, "关闭文件失败: %s", closeErr)
       }
    }()

    // 读取文件内容到内存
    fileBytes, err := io.ReadAll(src)
    if err != nil {
       utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
       return
    }

    // 设置请求的图像数据
    req.Image = fileBytes
    req.Name = file.Filename

    // 调用服务层处理上传
    resp, err := service.NewUploadImageService(ctx, c).Run(&req)
    if err != nil {
       utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
       return
    }

    // 返回上传成功的响应
    utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

//service
func (h *UploadImageService) Run(req *product.UploadImageReq) (resp *product.UploadImageResp, err error) {
    result, err := rpc.ProductClient.UploadImage(h.Context, &rpcproduct.UploadImageReq{
       FileName:  req.Name,
       ImageData: req.Image,
       Target:    req.Target,
    })
    if err != nil {
       return nil, err
    }
    return &product.UploadImageResp{
       Url: result.ImageUrl,
    }, nil
}
```

##### 3.3.2.8 部分 RPC 调用

```go
/*
    Run

// 1. get cart
// 2. calculate cart
// 3. create order
// 4. empty cart
// 5. pay
// 6. change order result
// 7. finish
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
    fmt.Println("CheckoutService.Run")
    // Finish your business logic.

    // 1. get cart
    // Idempotent
    // get cart
    cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
    fmt.Println(req.UserId)
    // cartResult.Cart.Items
    if err != nil {
        return nil, kerrors.NewGRPCBizStatusError(5005001, err.Error())
    }
    if cartResult == nil || cartResult.Cart.Items == nil {
        return nil, kerrors.NewGRPCBizStatusError(5004001, "cart is empty")
    }

    // 2. calculate cart
    var total float32
    oi := make([]*order.OrderItem, 0)
    for _, cartItem := range cartResult.Cart.Items {
        productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
            Id: cartItem.ProductId,
        })

        if resultErr != nil {
            return nil, resultErr
        }

        if productResp.Product == nil {
            continue
        }

        p := productResp.Product.Price

        cost := p * float32(cartItem.Quantity)
        total += cost
        oi = append(oi, &order.OrderItem{
            Item: &cart.CartItem{
                ProductId: cartItem.ProductId,
                Quantity:  cartItem.Quantity,
            },
            Cost: cost,
        })
    }
    fmt.Println("total", total)

    // 3. create order
    orderReq := &order.PlaceOrderReq{
        UserId:       req.UserId,
        UserCurrency: "USD",
        OrderItems:   oi,
        Email:        req.Email,
    }
    if req.Address != nil {
        addr := req.Address
        zipCodeInt, _ := strconv.Atoi(addr.ZipCode)
        orderReq.Address = &order.Address{
            StreetAddress: addr.StreetAddress,
            City:          addr.City,
            Country:       addr.Country,
            State:         addr.State,
            ZipCode:       int32(zipCodeInt),
        }
    }
    orderResult, err := rpc.OrderClient.PlaceOrder(s.ctx, orderReq)
    if err != nil {
        err = fmt.Errorf("placeOrder.err:%v", err)
        return
    }
    klog.Info("orderResult", orderResult)
    fmt.Println("orderResult", orderResult)

    // 4. empty cart
    emptyResult, err := rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
    if err != nil {
        err = fmt.Errorf("emptyCart.err:%v", err)
        return
    }
    klog.Info("emptyResult")
    klog.Info(emptyResult)
    fmt.Println("emptyResult")
    fmt.Println(emptyResult)

    // 5. pay
    // ==charge
    var orderId string
    if orderResult != nil || orderResult.Order != nil {
        orderId = orderResult.Order.OrderId
    }

    payReq := &payment.ChargeReq{
        UserId:  req.UserId,
        OrderId: orderId,
        Amount:  total,
        CreditCard: &payment.CreditCardInfo{
            CreditCardNumber:          req.CreditCard.CreditCardNumber,
            CreditCardCvv:             req.CreditCard.CreditCardCvv,
            CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
            CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
        },
    }

    paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
    if err != nil {
        return nil, err
    }
    klog.Info("paymentResult")
    klog.Info(paymentResult)
    fmt.Println("paymentResult")
    fmt.Println(paymentResult)

    // 6. change order result
    // change order state
    _, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{UserId: req.UserId, OrderId: orderId})
    if err != nil {
        klog.Error(err)
        return
    }

    // 7. finish
    resp = &checkout.CheckoutResp{
        OrderId:       orderId,
        TransactionId: paymentResult.TransactionId,
    }
    fmt.Println("normal return ......")
    return
}
```

# 四、测试结果

### 4.1 功能测试

#### 4.1.1 auth service

##### 4.1.1.1 deliver token

| 用例描述      | 测试结果                                                     |
| ------------- | ------------------------------------------------------------ |
| 正常分发token | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NGI0ZDViMWNiMGQ1NzAyMGUyZDEzMGU1YTBmODE2MDdfT0lkbU5Telp0TFByVjE5U05PckVxdXJLczJLcWpiUklfVG9rZW46S3J5RGJGZXhtbzNVWld4OGZSRWN6RG1lbmV6XzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.1.2 renew token

| 用例描述      | 测试结果                                                     |
| ------------- | ------------------------------------------------------------ |
| 正常刷新token | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MDE5M2UxNjU2NmJjZDc4ZmE5ZmEzODA2ZjBjYzU1NmNfQlJaMmFsV2dwakhuWkM3Q2pZWURKakVWQlVFak04MExfVG9rZW46UnUxQWJmWHBNb1BPNnR4TEFlUGNJRU44bkdjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.1.3 verify token

| 用例描述      | 测试结果                                                     |
| ------------- | ------------------------------------------------------------ |
| 正常验证token | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NGQxMDQyOGRmZjc4MDhjMmE4YzZlZjNjMjI2MGRlYjBfcEJXU2lEMUQxb1FvRUR6U2I4SUZFb1ZnamFlVDRpNnVfVG9rZW46T3VHQ2J0QWl3b3ptT0h4eEpBMWM3dWsyblBnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

#### 4.1.2 user service

##### 4.1.2.1 /user/register (POST)

| 用例描述           | 测试结果                                                     |
| ------------------ | ------------------------------------------------------------ |
| 正常注册           | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=OTAwMzdmODE2ZTg5ZjVjYWIxMzlhYzMxMDNiNmI0YzJfaHlxQWxLeUR1MnptY3RIcTRrY0h2S3VsU2Jrc1ZPMlFfVG9rZW46UjFoSmJNOUZQb0xOUDh4VEF2dWMya2QzbkJmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 两次输入密码不一致 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YzI1NGM3YmJjZDZjMWJlNDdhMDI1MzVmYmFiYmJkZmVfeWFlcGJFN05RWlE0amNhU2xwaVBZN2dtTTR4Z2o3bGhfVG9rZW46QzVRd2JTdXRIb3NJR2R4YjFtcWN4WUhTbmZnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 邮箱为空           | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NjhmN2ExOTRkZmIyZGEzZTc2OTE4OTk3MzM2MjIxZWNfaHFOOThjZENHZ204V1hSRW9oU3lhVk50TkdybXYzT2lfVG9rZW46VXVDbGJhc1Fyb0QwbWt4UzRKcmMzSU5lbmpnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.2.2 /user/login (POST)

| 用例描述   | 测试结果                                                     |
| ---------- | ------------------------------------------------------------ |
| 正常登录   | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTMyMTExYmNjZWQxZjIyZjUxMjdkMzc2M2U1YWJiZWJfdzI4UzRZcjNPSDA1SUg5Zkc1NHFyRTF0YmV2bkxLV3NfVG9rZW46S1NjWWJlcnhsb2FSNnZ4TU02eWNKZElobmZjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 密码错误   | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YWU3YWQ0ZWIyMTk5YWJiMzM2MjYwNWJjMTQyMzRiMzhfaEFsbWVWWDVaUUpDcUk5N0NoVmU1TjdpamRBdjFKdlpfVG9rZW46VTRKTGJNR2xwb21mNlB4SDN0bmM2OHJXblBjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 邮箱为空   | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZGMyYjY2OWFiMWNjMmExZWVjZDdmZTZlZWJhZGNhMWNfdEtWbXI3dXdIbVdlVWNmZWY2bnFPZnl1RVpSM0dRRHZfVG9rZW46VTkwa2I5RDNDb0g0TUx4a1BzRGNnWHFVbldiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 用户被封禁 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Mjc2OTExNTg5YmM4ZDFmMWU2ZGYyZTgzZjg2NTVlNTVfUEFXdWRBcU9rOWxrWFRnSzBEeTlkZzUyUlJob0gxb2VfVG9rZW46RXNZUWI5cDk4b0lqU1p4YWd0cmNkdUQ5bmxlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.2.3 /user/logout (POST)

| 用例描述  | 测试结果                                                     |
| --------- | ------------------------------------------------------------ |
| 正常登出  | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YTk1YzY5MzkwZGVhZDFhNDQ4ZGE5ZmQyMzczNWEyMmFfT3B1Zkg1TkJKUENBaFZFZnhUbURJNlFKU2t6SGp2NWlfVG9rZW46WXZBMGJMajE1b0tybk54MGVIQWN2VmJxbnk0XzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 无效token | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YWE1ZjBiMmFhZTk0ZDczMTk4OGI3NTYwMTZhZTUxNGRfeExoZUlPOFZ5UTFnY3dDeHlEbmh3S3NWckNQMzdjRkhfVG9rZW46WkN5b2JwZ21ab3hpalN4UWdCcmNTaXVabk5hXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| token为空 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Zjc2YzY3NDcyYWIxYTdiYjk0M2FiMDk2YzMzMWNlNTVfMGdyVFFXVUhIZ2drQU56VE9CNUdhQUxSSUs2dHZNNFJfVG9rZW46S2lvRWJEc29Pb0p4RGl4OUxuemNwUnM2bnRlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.2.4 /user/delete (DELETE)

| 用例描述           | 测试结果                                                     |
| ------------------ | ------------------------------------------------------------ |
| 正常删除           | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTdiMjdjY2ViMTQ2YTFjYjgyMzI3YjU5YTFmMjIxNjVfakc4S0xkRlJmV0ZyejJWVFFGUlFRTnoxNkxhdzhyQkdfVG9rZW46UlBNQ2JyU2hyb0M4ekR4cXZZV2NwVVFXbkpnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 用户id为空或不存在 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Nzg5Mzc5MGVjZWU3ZGFlOGI0YTM0NzRiNDJjMjBmODJfZ3phRUQ3QVNDaHNiUWYwRHdaM1cxNnZFVmRDajNiTU1fVG9rZW46WWNLa2JTcXFWbzI1aU54akcwdmNLbnpNbjJnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.2.5 /user/update (PUT)

| 用例描述         | 测试结果                                                     |
| ---------------- | ------------------------------------------------------------ |
| 正常更新用户     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MzYyY2VhNWI5MjlhNmJmNGRmZGY2OGZkN2NiZDBiMjhfZ0llemIxb0JQZmJwM0R1MHdFU2dxS2JtN3UxSjQ2czRfVG9rZW46UkNNSGJXU0Rzb2dxakZ4VzF1VGNQZFlEbkJiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 更新的结构体为空 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTcyNWYzNGM1M2IyYzkwZmE3NzQ2NDkxODBjM2YzYTdfQ0V3cjdmVWhtdHpGd0trNmdNMjdkdFNiNWgzSEt4d0lfVG9rZW46UnRMZmJHeVFWb0EzQUl4am10QmNlendxbm5mXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.2.6 /user/info (GET)

| 用例描述         | 测试结果                                                     |
| ---------------- | ------------------------------------------------------------ |
| 正常获取用户信息 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTA5ZWZmMjgzNTI0MjY0NmJjZjQ4NDFlNTkwN2ZjZWNfcHdKbE5wbEllYkV6WFFya1dEcDRoNXgyem9XS3piWkpfVG9rZW46S0pjaWJtc0d0b1h6TGt4RGM5VmNkelBjbkhlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.2.7 /user/update_role (PUT)

| 用例描述         | 测试结果                                                     |
| ---------------- | ------------------------------------------------------------ |
| 正常修改用户角色 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTdhNWI2YjVhY2JhNWQ4MWM2M2ZiOWMwOGM0NTFiYTZfTVZIaW9oRDBmUFdvRFMyRVVSWTVDeWlCd2NwRGVWU2NfVG9rZW46VW44d2JPbFBNb2hxNWZ4dW5SemN4TWp0bkdmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 用户id错误       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjUwYzMwMWM0YTIxMTJhNWY2MWI1MTk5ZDNhMWExZWNfd2tacEFwMUlGdEp4VGR5aVU5YmVsMGJ3QmllaXRmTkhfVG9rZW46VjFlRWJOUEt4b0doWmd4bkcwN2NoOUdIbmpjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

#### 4.1.3 product service

##### 4.1.3.1 /product (POST) - 新建商品

| 用例描述     | 测试结果                                                     |
| :----------- | :----------------------------------------------------------- |
| 新建商品成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Y2YxNjQyMThhNWQxNjM0ZTFmNjc0N2FlM2RkMjY2NWVfUnRWTWwxUUE0cnRxclVuOHd2Qmdwa2hYYzJKZ21VRXNfVG9rZW46QkpPcWJtTWZqbzJIOEx4NzdseGNwWmhqbjJJXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 商品名称为空 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MWM4N2Y3OGM2YjQwMDVmMTM2Mzc1NTU0NTVhYmYzY2ZfZHJCWnRFU0JGQmx5TG5FZDB1QW1GYVA0MG82a29sZndfVG9rZW46VkdIWGIya0J2bzBpM1R4Snc1aGNjak15blE5XzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 商品价格为负 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MTViZGQ5ZmY5MmFkMjUwZmMzMjNkNWZhODEyMTMzNzVfdUZYVGNwVVRCeUVveVlkMk9TVVMxOTBhZENFZ1hnU2RfVG9rZW46TktNN2J4OEsyb0dnUmJ4em9WUmNnMkF0bjRkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.3.2 /product (GET) - 获取商品

| 用例描述         | 测试结果                                                     |
| :--------------- | :----------------------------------------------------------- |
| 获取商品详情成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NmUzZmRkNjY3NTk4N2ZlMDM4NWE5YWYxZmRiOWQ2MTVfTW1SWUYxNzNneEMxWVBJZmZLZmNrTFJ0Y3JrZ0NmNGlfVG9rZW46WVhzcWJ2SHdGb1hBeml4RzNWaGNqTGtjbmlmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 获取不存在的商品 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=N2I2MDk1MjFhNjlkNzgxYzk3ZjFkMjE4Mjk4MjNkYzRfM1NXMXpwNFFkRXR3MjVkOVJLb1ZMYW42ZmdiNWVPejNfVG9rZW46RW9aYWJacmZ5b3o1YkJ4WFU1emN0RThPbkdlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 商品ID为0        | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ODM4MDhhNGI0ZjU3ODdiNGI1ZjAzOGI4NWExYWNhNTFfOHZtak0wUDVTWWRtTUpZSHZnZEN3NTdqMml4RjJtQlVfVG9rZW46T0QxMGI0ZGVVb09vTnZ4N005Z2M2T0ZDbldiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.3.3 /product (PUT) - 修改商品

| 用例描述     | 测试结果                                                     |
| :----------- | :----------------------------------------------------------- |
| 修改商品成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NTljMDQ2ZDY4ZTQ0YjU2NTA5ZGRlNGZhYWYxOGMzODVfaVFkaW81VmZKSjgwZWhCRnJoMFhCOEhVTk5PaGF4U1hfVG9rZW46SkRzTWJIYkFhb21kVDB4Wk9SYWMySldQbndoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 商品ID为0    | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjU5ZGU3NGNhNDcxZTA2ZmJiOTk3MjAwNTRkMjAyNmNfcDFZQjV2cktTVDgyRGpoZVRkdzR3dHNwbkM5Tm83WkNfVG9rZW46SzlETWI3cjFBb2hoMzh4b2tIU2NhTVpObkRnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.3.4 /product (DELETE) - 删除商品

| 用例描述         | 测试结果                                                     |
| :--------------- | :----------------------------------------------------------- |
| 删除商品成功     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjUyYzY5ODIyMDY5N2Y1MzQ0YTEyNDNkYmZhZmYzOTVfZjVwQ2k4TVFReENKZjFaOEVRcnl3QWNTVkVmOE1iOUxfVG9rZW46QnZKYmI2a0w0b2xlM0x4N1EyMGNsdWtJbnY1XzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 删除不存在的商品 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YWU0MDkyZTY2MjAxYWQwZTI5OWViYzVkNWI1NGI2YThfdVJQRnROdDlYZHlDdm5ERUVuaVdRRjZqT2lxM3FINXNfVG9rZW46SWdpc2IxOEYyb3BCcXN4aHZEOGNLV1Jibm1jXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 商品ID为0        | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MjVmYzZmYjNlZjYwM2JiMjA5ZDBiYWQwMDkzZTUwODNfNHZwSG5uVzlxY1FyYW1YTk9OSG1pQ3U5TjdkblpObGxfVG9rZW46SlNseWJYVVZMb0g3Sjd4cDlRTGNNMTJBbmNiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.3.5/products (GET) - 获取商品列表

| 用例描述         | 测试结果                                                     |
| :--------------- | :----------------------------------------------------------- |
| 获取商品列表成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=OTdmNzMxNTQwYmFkZDcyMmEwZTA2OWIxZjAwZjJhMDBfM2VwaEUwNnFROWpUc3VFT25aNUhmU2R1ZHZ6YllRVXlfVG9rZW46VThrNGJaVkVIb2FISWd4Nk9GNmNveGxEbmVkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 每页数量为负数   | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MDg1MTJlMDM2YjU0MzM3NWE5NDY1M2ZkYjEwMjc0NmRfdkoxaW5OS3YxOFZjRVNLU2ozTmJITHhUSXZGVkRPOEpfVG9rZW46TjFmZ2JWZGFpb3BQM2h4MnpuWGNkd3JybjV2XzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 页码为负数       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YzAwZmI0YzQ1MTNiOWVhNDkyYjg2YWIxM2U2NThiMWJfVTFoR05LazRZT0pzQVV3VDlsMVlONU1nUG5DZXdVTEJfVG9rZW46UlV5RmJHdlNab3QxcVV4SU90QmM1VGk1bnpjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 每页数量为0      | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZmM1OTljOWE0MzFmYzdhYzIxZjAyYTBjZTM0ZTg3NWZfR1ViREpkN25Ddzk5d3lwUUVQWUtkUFBjaGd6Mm42allfVG9rZW46SWYwRWJGQnU1bzFKVkN4NndYb2NyNFU2bkpkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.3.6 /product/search (GET) - 搜寻商品(匹配所有字段)

| 用例描述             | 测试结果                                                     |
| :------------------- | :----------------------------------------------------------- |
| 搜寻商品列表成功     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MjBmOWQyN2EyMWE1ZmY5OTYyYWFjYzQxYjVlYWU4NmRfN0pEUjV0bk5iZXZ1UXhTd3NSVWM1QWxsdHJobk03a1VfVG9rZW46Rm9nTGJXVlVqb2tkbDl4SHYwQ2NIcVhsbkNiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 每页数量或页码为负数 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MjA2MWJjZmMxMDlhYmJjYjY4ZWZmYzA2Y2YwYjdlZThfZDV4Zk1vcVF1NVFjVWpveW9qWmFXc1FuS1luS0E1eWNfVG9rZW46SVJxUGIySDkwb2pjWDB4VlhQMmNaMUFabjVmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.3.7 /product/upload (POST) - 图片上传

| 用例描述     | 测试结果                                                     |
| :----------- | :----------------------------------------------------------- |
| 图片上传成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=M2QzYjNjZjBjOTJkN2Q4NGRmNThjMzViZWNkYmFiMzlfa0JhWFY4blNwa2JZME1NM3RtVnN5MHJqQVVSV2RxVklfVG9rZW46VTBvQWJ3Sm53b1FYMjh4aVBzR2NucEhXbktkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

#### 4.1.4 cart service

##### 4.1.4.1 /cart (GET) - 获取购物车

| 用例描述       | 测试结果                                                     |
| :------------- | :----------------------------------------------------------- |
| 正常获取购物车 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=OWJiZTcxMmU0ODA4MzUwOTUxZDI0ZTFiNGE5ODhmZjFfT2ZzM1hwUGFxNmVHSVB2c0E0bm51TjVTYkNRV1ZrMzRfVG9rZW46QW0zU2IzYUIxb0dGMlN4cTlCWGNKa0FkbjNlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 购物车为空     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Njc3NGI1YjA0ZTMxMzgzZTVkNGNlNmYyNTk2MjAyZWJfMUdDckpsTGNhMExYSHFLZ25UTWFZWHQ0S2YwaHBpcWpfVG9rZW46QVliaWJpeHA0b0Jyc2d4T2YxR2NURHk4bjNSXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.4.2 /cart (POST) - 添加商品到购物车

| 用例描述     | 测试结果                                                     |
| :----------- | :----------------------------------------------------------- |
| 添加商品成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ODQ2NzczZGZkZjFhYzc1MGM2YzU4MzFjNjAyOThjOThfamswOUVKV3hZemRxWUs1cjJZdlptbkxkeUdPRTlVYVFfVG9rZW46UkFwV2JzVmltb2FRczh4dDhDY2NWOEZwblBlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
|              |                                                              |

##### 4.1.4.3 /cart（DELETE）- 清空购物车 

| 用例描述       | 测试结果                                                     |
| :------------- | :----------------------------------------------------------- |
| 清空购物车成功 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZDA3YzE5ODA0Zjc0ZGY0NDU3YzlkN2EzYmRjZjM2YzhfZ0ZhZWt6aWs4UUJDVnRVeE93N1A1MXBKSXZTTXNURVdfVG9rZW46SUEwNWI5aDVFb0dHYkV4NURLZ2NxN2VXblRlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
|                |                                                              |

#### 4.1.5 order service

##### 4.1.5.1 **/orders (POST) - 创建订单接口**

| 用例描述     | 测试结果                                                     |
| ------------ | ------------------------------------------------------------ |
| 正常创建订单 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZmE4ZmE0ZDM1N2YyMjc5NzE2MWIzZDJmYTFiZTE1NDNfc0tnejlRaUlRR01MT2tKSGMyNGM1M1BaUFZ1ZTRjSUZfVG9rZW46RldFWGJsT2Fyb1BQaVR4cHpjZGNjTW5kbnFiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单项为空   | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZmUyMDlhYTBiYjY1MzcwZGU4MjJiZTBlZGNhYTliZGJfeUNHMldRWE1jNWRQOTlJeWxqZDdob0M0V21XRktuRzBfVG9rZW46SUhaN2JrY1FGb0ttQzl4WjFIWmN1dW9KbnNnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.5.2 **/orders/:orderId (GET) - 获取订单列表接口**

| 用例描述             | 测试结果                                                     |
| -------------------- | ------------------------------------------------------------ |
| 正常获取订单列表     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YzIyNTg1MDRhYThjNjgyYzBjMWVjZWU5NTFiNjE1N2VfbEo3WHhFRXdSQ215Qm1QQ3A3bGRORFc1bzdvNXZiNTdfVG9rZW46VVRVZWJla0l5bzlUZkh4V0ZMc2NlbGZxblVmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 获取不存在用户的订单 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Y2E5MTI5NDgzZTljZDdmOTE0NzZkZTdiNDAzNDNjZjhfM0J5YU1xczBYQ1J6ME5OMXJTd1V2ZGVYZlVmeVlMQnlfVG9rZW46VjlZdGJVcjd4b3JIdXd4MGVrTGNpRnhpbnFiXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.5.3 **/orders/:orderId/paid (PUT) - 标记订单支付接口**

| 用例描述           | 测试结果                                                     |
| ------------------ | ------------------------------------------------------------ |
| 正常标记支付       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Y2IzMGNmMWU5MzFlZmVkZmEzYjRhNmY0MmE1YTc1ZDdfcmhIRU00TFJMUHhkYUFvSXNDR3N3TFVLZGtpdlZsN29fVG9rZW46REs2RmJtSVFRb3hvNmp4ZkxMZGNmQ0Q0bkJoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单不存在         | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjI5ZWU4YmQwZDIxYWNlZjJlM2RiM2ViODQ0NDc0M2VfR1lEQlQ2ajI0VTRWNUUzNnIzOWZiQVlKbEdXbWdsTjBfVG9rZW46V29BRWJpR0Yxbzc0bU14MUlMTmNjSm9xbnJmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单不属于该用户   | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YWYxYWRkMDY5NWQ1NzVlN2M2NWNhY2FjZjVjNDQzZmZfdHBvdnRSOXhTNm92U1h1VGp3N1Izd2NkSlNzdmt4SGJfVG9rZW46VVhxR2J0SVhXb1h1aWN4Z0VsN2M5cjhYblZlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单已取消或已支付 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZmNlODUyYTMxNDFiYmE4NzY0ZjkyMDAwNTc4MzBkYThfUDBIQXg1aGVKNmRrV3R1SVB6N2lJMVJMWjk3Mmc2TFZfVG9rZW46R3FIM2Jkd3FPb242SFZ4Z3psa2Nra3lhblliXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.5.4 **/orders/:orderId (PUT) - 更新订单接口**

| 用例描述         | 测试结果                                                     |
| ---------------- | ------------------------------------------------------------ |
| 正常更新订单     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YjBhZGI0NWQ5ZWM4ZjkzNGRiOWVjM2ZhMGFiZTI4NmFfd0xoV0V2T21LNTB5TjBYS2ZYSTNxaEVhbjhYc1ZZQm5fVG9rZW46WURzMWJTV3JIb0ZlOTJ4TlpUT2NrdmZFbnJnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单不存在       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=OTRiY2ZjODdjMjY1YzA2ZThiNWRkYjNmMTM3MmFlZTlfeUhvQU5aT1UzRDlNbGQ4d25IOGVCUDFVMFJ3ejgxUmtfVG9rZW46UkxnU2Izbmxvb0luQ0l4YlJndGNUUWpzbnRnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单不属于该用户 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MTljZDU4NzJlMmExODljMzBmMjEwZmM5OWMzZjFmZGVfWWp2dkxVTUIzTk5qdTEzcmgzVG1MSEVwanE2YkI0cnZfVG9rZW46U3plQWJtUzZNb1o5YWx4ZzVKSGNXTnhpblJnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单已取消       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MWJkYWI5MjdhMWQzZjE2YmU5NjgxZTMyNTY5OTMyOGRfOGdzZEZIOTl6b0hXcHVnMjlKQWRuTlJLZzJxaEM1ZVZfVG9rZW46SmRIRWJ6NFl2b3EyajB4c0lEZGNLZTJFbkZlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 更新项为空       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NmJkMGIwOWIxOTkwYmE4ZDE4MjU5NDEzZjlkZDg5ZjhfcUpDZ0h0ZmNiS2Vtaldsa3FUdENLNXdMTXc1RExuR2xfVG9rZW46QzNRdWJrOUlhb2lxMDZ4ekJ6Y2NrcjlJbkFoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.5.5 **/orders/:orderId (DELETE) - 取消订单接口**

| 用例描述         | 测试结果                                                     |
| ---------------- | ------------------------------------------------------------ |
| 正常取消订单     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NGM3NzU4MjM3NjA3YzBjYjllNjI3MTY0Y2U1MmU4ZWJfeUR5czl4UXhvd2ZlWEhiZ0x5N2tKVEljakk5RXRhRmpfVG9rZW46T2FvN2JzYzMxb1NSb3R4OXFjOGNHeWNabmJjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单不存在       | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=OGRhMjA5ZDVkNjViMjdjZTNmNTE4ZjdiZDZhM2FkZjZfcWphVmN1RndONUtiekdQMkdxa21tODE5RjFkelFZRWhfVG9rZW46UFY1VWJ4aVFSb2MzYTd4Q1k1YmM1M0xFbmw3XzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 订单不属于该用户 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MGVkODYxY2FiYzdlOWIzZGRmMjFkY2EyOTUzZTc0OGRfWmV4NDYyQlg4enA3MEx5QkRwZTc0MkMxeXRySGZBS1JfVG9rZW46SUFXOWJoaFNqb29ueFd4cmpDNWN0OGp1bmhlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 重复取消订单     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NDMwNWFlYjNlZWRjOWUzMmQ5YTY5ZTRiODZlN2ViN2ZfeFkzbDhOVlVWVnU5R1kyRXpaS3d4dzM5QzVmcHkxMWVfVG9rZW46UHJPSWJ0cG5jbzN6MWd4NWR6cGN3SkxpbjNlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 取消时间为空     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NGYwOWQ4ZmE1ZjdkYjQ4YTAzNWI3NzA3NDFhZTNkMjNfdkxYN1I4Y0FEMlRmekMyZGJDU0w4aUc5d2ZObTlqRWlfVG9rZW46TUlSRWJXOTlnb1J2clV4OTZCQWNuVm1IbkRmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

#### 4.1.6 checkout service

##### 4.1.6.1 **/checkout (GET) - 结算接口**

| 用例描述 | 测试结果                                                     |
| :------- | :----------------------------------------------------------- |
| 结算     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NGM2NjQ4ZmU0YTk5NDRkNzJmNDAxMGQ5ODE0YjFiYzBfV3FaUHdrR3dac0hVYzZNdkg1Q0JBbHhhZjhCcXViZENfVG9rZW46RmlsWWJieFdIb21jeDN4dGU2emNSNnMxbnljXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YmZjOThkMTA4OGI4YzYzNTViOTg5ZmQxNzBlOGQ2MTVfUzliS1I1N01YdEtlNjNZQTB1blNreW5WN3V0amZyZ0ZfVG9rZW46WDYzR2I2clVJb1ppNmt4YXAxM2NBdVRjbjJlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YmVkNWUzMzNmZDRiODZmOTZkYmMwZTAyYjZkZWVhZTFfbUFKaGQ2UnVxeENzRWc2MXl6NnRDYmhRMHpoRXpoMkdfVG9rZW46RXVZV2JsNWFHb1I5bkJ4N092RWNTOUZLbklmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NDg2ZTJhYTZhOGU4ODQyNTBlODMzZTg5ZWI0MmMzMjRfWGpyRElNNkdrb05uOFl3N2R1a1AwbFJGQmZMTW5BTXBfVG9rZW46R0JTYWJCNGNGbzJ0VGN4SEhxQWNTTWl6bnloXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

#### 4.1.7 payment service

| 用例描述 | 测试结果                                                     |
| :------- | :----------------------------------------------------------- |
| 支付     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NjY3YzIzOTQ2MjhkZmRlNTFkMjk1MWUxNzQ2Mjg4OWFfcFhTcDhTYWZCaDZVVUhCMWxZQVVQcTgzU05TWXdpd1ZfVG9rZW46Qk1wQ2J1UUVGb0RodVN4eEhNemMwU0FwbnNkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

#### 4.1.8 ai service

##### 4.1.8.1 /ai/query (POST)

| 用例描述             | 测试结果                                                     |
| -------------------- | ------------------------------------------------------------ |
| 按照日期查询订单     | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YzU0ZGUwMDhkNjEwZWE4MDcyNTNkYjgxMGEwN2Q0ZjBfbHhJOElpYXdqZDV6TVRLSGNZQVYxMFhIeVNKQUY3ZVRfVG9rZW46WFk2ZmJ5c0t2b0lmZDZ4NzE0UmM0cFlsblZNXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 按照商品名称查询订单 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NzNmOWJiNzZlMWU4OWMxZGZlMTRmM2UzNDUyYzA5MjdfU3JseFR4a3BsM2lnYlRZcnAwQW9HUFN5N0RzNzc4QkZfVG9rZW46QjJJZWJkOHpMb0t5azB4ZHFLQWM2VmVybnNlXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 按照支付状态查询订单 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZWY0MTFlYTczMmJhZmM1NGJjZDQ5M2IyM2QzOGE2ZDRfNm9aTGdQdjBnZTlWcUs0OHJjQUJMdHNnYldNekxseTRfVG9rZW46VVpTdmJYUWg0b3IzOFZ4WVU5cWNvNW42bkplXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

##### 4.1.8.2 /ai/place (POST)

| 用例描述     | 测试结果                                                     |
| ------------ | ------------------------------------------------------------ |
| 下单一件商品 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NzhkODA4NTMzNjJlZGU1MjljNDEyMjVmMmRmMmI0MGZfc2ZwZlF6Z2x1OWFrNUk0OGF5ck90QVo0bUJrNXdtTFJfVG9rZW46U0QzTGIxNnM1b2daTUx4QTlxMGNqb2EzbldnXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 下单两件商品 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=MzFlNGQ4YzliOTI4YmViOGZhNDdiNjZlOTBjMzcxZGRfY1J6enBZbEhtWG1MTTdiVG04SGtoaDVLbTZlcUlFR3BfVG9rZW46VVhocGJzVkJVbzVEd0t4NzhncGNuenF2bmxoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |
| 下单多件商品 | ![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=YjQ5Yzk4ODg1NTc0ZjBhNjdhZjc2ZmJiYjk3NjZhY2VfeFdJWm4zVkUxVld4U092OGpleTI0c0oxMmhhQnJRMUhfVG9rZW46UjJiZGJwSDRnb2gzNkp4UTExM2M5M2JKbndoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA) |

### 4.2 性能测试

#### 4.2.1 用户登录

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZmE0YjNiZjMwZDgxZDc1MzY0ZDZmZTg2NzM5M2YzNTlfWGZWMXpkQXpCcmdIZlZGc1B3S2NXTXFGbWZ6T00wT1dfVG9rZW46RmxEWGJhWFJub0pJaWJ4b1JkcGNrNzR6bjVmXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=NGNkMGFmZWU0Y2M4MDE1NmUwMTQ3OTY1NmYyN2Q2ZWJfVUdRWWZ1U3FEZWtRZ2Q3eU5TNXZoaDBGN1VFUXRORkZfVG9rZW46UEZqcWJnU25Gb2hIenN4QlpRWWMzVWxIbndoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

#### 4.2.2 下单

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjQ3Njg4OTc3NjVjMWVkNWU1MTcwNGJhNDdmMzgwNjhfU09uRmZTYkV0Uk5obVZ0QU5oeWNuUmVZc1NwNUhUV3pfVG9rZW46VFVtS2JPM2s4bzZnTXB4VW04OGNJaWdZbkpjXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

#### 4.2.3 搜索商品

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=Y2IzZWM1YWE0NTYyYzJkZTVkYTk2ZDkxMTA2MjQxMzFfNjlnOGRLUDlaY2w2RE5mZG83cVFFR2VTY2VZMDFreG9fVG9rZW46TGtONGI0cGtZbzE3ZWh4Y3oxb2NIUERibjZkXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

![img](https://ebn7zhozsl.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjdlMDIxZDM2ZmJhYmNlYmJiZjVjYzk4MjNmN2QyMGFfM0JWUzFhbEtWVzdVOGxtYzZHSllDRlZnM1ZMOFFaQ3ZfVG9rZW46RUFhUGI5NDdob2xCa2t4YVNmY2NGdFRqbnVoXzE3NTk5OTczODY6MTc2MDAwMDk4Nl9WNA)

# 五、Demo 演示视频 （必填）

## https://www.bilibili.com/video/BV1G69BYBEun

# 六、项目总结与反思 

> 1. 目前仍存在的问题
>    1. **性能瓶颈**：在高并发场景下，部分微服务的响应时间较长，需要进一步优化性能。
>    2. **资源利用率**：某些时段资源使用不均衡，导致部分资源闲置或过载。
>    3. **数据一致性**：在分布式环境中，某些情况下数据一致性无法完全保证，需要改进数据同步机制。
> 2. 已识别出的优化项
>    1.  **安全性**
>    2. 将 JWT 证书更新为很短的有效期，并且使用更严格的加密算法。
>    3. RabbitMQ 重新设置用户权限，进一步添加身份认证策略，避免突发流量带来的拉起效果。
>    4. 在 MinIO 中建立更精简的访问权限控制，限制公共访问和需要鉴权的数据操作。
>    5.  **性能**
>    6. 优化 MySQL 索引管理，减少查询开销，提高查询速度。
>    7. 针对高并发场景，优化 表连接（JOIN）策略，避免全表扫描带来的性能下降。
>    8. 引入读写分离架构，使用 主从复制+Redis 缓存，减少数据库压力。
>    9. 在 OpenTelemetry 中针对高流量核心应用 进行链路数据采样和压缩，减少存储冗余，优化可观测性性能。
> 3. 架构演进的可能性
>    1. 考虑将 RabbitMQ 更换为 Kafka，实现更高效的消息队列处理。
>    2. 将部分代码重构为流水线处理，最大限度充分利用机器外部系统资源。
> 4. 项目过程中的反思与总结
>
> 本次项目基于 **Go + Kitex + Hertz** 构建了一个高并发电商平台，并通过 **Consul** 进行服务注册与发现。整体架构展现了良好的可扩展性和高效能，但仍有一些需要进一步优化的地方：
>
> - **API 网关安全性**：Hertz 在 API 端的高性能表现优秀，但在数据传输和访问控制上仍有优化空间，后续可加入流量限流、API 签名校验等机制。
> - **缓存策略优化**：Redis 在本项目中承担了高频查询的缓存功能，但面对热点数据仍可能出现访问瓶颈，未来可结合 **分布式缓存（如 Redis Cluster）** 或 **本地缓存（如 Caffeine）** 提升响应效率。
> - **高并发事务处理**：在订单支付和库存管理中，分布式事务仍是一个关键挑战，后续可引入 **TCC（Try-Confirm-Cancel）模式** 或 **Seata** 进行事务优化。
> - **团队协作提升**：本次项目锻炼了团队在 **微服务架构设计、DevOps、分布式计算** 等方面的能力，同时加强了协作开发和代码规范化管理，为后续更复杂的业务架构打下了坚实基础。
>
> 未来，我们将继续优化系统架构，引入更先进的技术和工具，不断提升系统的稳定性、可扩展性和安全性，以更好地支持高并发业务场景。

# 七、其他补充资料

[Git规范](https://kv8faq2pjwc.feishu.cn/docx/Tx9SdbusmoJYrBx3MCmcXM5tnbg?from=from_copylink)

[编码规范](https://ywn3zwhwg6s.feishu.cn/docx/H3AYdMOEGogyvXxawpYcBwTinbe?from=from_copylink)

[青训营项目安排](https://zjutjhwl.feishu.cn/wiki/ZctSwNpgbi56OGkNIEgc5ERPnYe?from=from_copylink)

[学习心得](https://ywn3zwhwg6s.feishu.cn/docx/RpKVdAKgmoJRpxxGiLIcR3TBnee?from=from_copylink)
