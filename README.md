# proglog

一个用 Go 实现的分布式提交日志（commit log）服务。通过 gRPC 提供写入/读取与流式接口，使用 Raft 做一致性复制，Serf 做节点发现与成员变更，并支持 TLS 双向认证与 Casbin ACL 鉴权。

**功能特性**
- gRPC API: `Produce/Consume` 及流式 `ProduceStream/ConsumeStream`。
- Raft 一致性复制与日志快照。
- Serf 成员发现与自动加入/退出。
- gRPC 负载均衡：写入走 Leader，读取按 Follower 轮询。
- TLS 与 ACL：双向 TLS + Casbin 策略控制。
- 可观测性：zap 日志与 OpenCensus 指标。

**目录结构**
- `cmd/proglog` 可执行入口（CLI）。
- `api/v1` gRPC 协议与生成代码。
- `internal/log` 本地日志与分布式日志实现。
- `internal/server` gRPC 服务端实现。
- `internal/discovery` Serf 成员管理。
- `internal/loadbalance` gRPC resolver 与 balancer。
- `internal/auth` Casbin 鉴权。
- `deploy/proglog` Helm Chart。
- `test` 示例 CA/证书与 ACL 策略样例。

**快速开始（本地单节点）**
1. 初始化配置目录与测试证书（依赖 `cfssl`、`cfssljson`）。
```
make init
make gencert
```

1. 启动单节点（示例）。
```
go run ./cmd/proglog \
  --bootstrap \
  --node-name node-1 \
  --bind-addr 127.0.0.1:8401 \
  --rpc-port 8400 \
  --data-dir /tmp/proglog-1 \
  --acl-model-file ~/.proglog/model.conf \
  --acl-policy-file ~/.proglog/policy.csv \
  --server-tls-cert-file ~/.proglog/server.pem \
  --server-tls-key-file ~/.proglog/server-key.pem \
  --server-tls-ca-file ~/.proglog/ca.pem \
  --peer-tls-cert-file ~/.proglog/root-client.pem \
  --peer-tls-key-file ~/.proglog/root-client-key.pem \
  --peer-tls-ca-file ~/.proglog/ca.pem
```

**多节点启动（示例）**
- 第一个节点加 `--bootstrap`。
- 其余节点增加 `--start-join-addrs` 指向已存在节点的 `bind-addr`。
```
go run ./cmd/proglog \
  --node-name node-2 \
  --bind-addr 127.0.0.1:8402 \
  --rpc-port 8403 \
  --data-dir /tmp/proglog-2 \
  --start-join-addrs 127.0.0.1:8401 \
  --acl-model-file ~/.proglog/model.conf \
  --acl-policy-file ~/.proglog/policy.csv \
  --server-tls-cert-file ~/.proglog/server.pem \
  --server-tls-key-file ~/.proglog/server-key.pem \
  --server-tls-ca-file ~/.proglog/ca.pem \
  --peer-tls-cert-file ~/.proglog/root-client.pem \
  --peer-tls-key-file ~/.proglog/root-client-key.pem \
  --peer-tls-ca-file ~/.proglog/ca.pem
```

**配置与安全**
- ACL 模型与策略示例位于 `test/model.conf`、`test/policy.csv`，`make init` 后默认路径为 `~/.proglog`。
- TLS 示例证书由 `make gencert` 生成，默认同样放在 `~/.proglog`。
- 支持通过 `CONFIG_DIR` 环境变量修改默认配置目录。

**开发**
- 运行测试：
```
make test
```
- 生成 proto 代码：
```
make compile
```

**部署**
- Helm Chart 位于 `deploy/proglog`。

**协议**
- gRPC 定义见 `api/v1/log.proto`。
