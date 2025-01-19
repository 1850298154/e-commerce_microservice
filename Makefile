# 自定义模块的执行命令
.PHONY: all
all: help

default: help

.PHONY: help
help: ## 显示帮助信息，列出所有可用的目标命令。
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Initialize Project
.PHONY: init
init: ## 复制 `.env.example` 到 `.env` 文件，只需执行一次。
	@scripts/copy_env.sh

##@ Gen
.PHONY: gen
gen: ## 生成 {svc} 服务的客户端代码。例如：make gen svc=product
	@scripts/gen.sh ${svc}

.PHONY: gen-client
gen-client: ## 生成 {svc} 服务的客户端代码。例如：make gen-client svc=product
	@cd rpc_gen && cwgo client --type RPC --service ${svc} --module github.com/cloudwego/biz-demo/gomall/rpc_gen -I ../idl --idl ../idl/${svc}.proto

.PHONY: gen-server
gen-server: ## 生成 {svc} 服务的服务端代码。例如：make gen-server svc=product
	@cd app/${svc} && cwgo server --type RPC --service ${svc} --module github.com/cloudwego/biz-demo/gomall/app/${svc} --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/${svc}.proto


##@ Build
.PHONY: tidy
tidy: ## 执行 `go mod tidy` 清理 Go 模块的依赖。
	@scripts/tidy.sh

.PHONY: fmt
fmt: ## 格式化 Go 代码，使用 `gofmt`、`gofumpt` 和 `goimports`。
	@gofmt -l -w app
	@gofumpt -l -w app
	@goimports -l -w app

.PHONY: run
run: ## 运行指定的服务。例如：make run svc=product
	@scripts/run.sh ${svc}

##@ Development Env
.PHONY: env-start
env-start: ## 启动所有中间件服务作为 Docker 容器。
	@docker-compose up -d

.PHONY: env-stop
env-stop: ## 停止所有运行中的 Docker 容器。
	@docker-compose down

.PHONY: clean
clean: ## 清理所有临时文件和目录。
	@rm -r app/**/log/ app/**/tmp/

##@ Open
.PHONY: open-consul
open-consul: ## 在默认浏览器中打开 Consul UI。
	@open "http://localhost:8500/ui/"