# 🛍️ Go-Mall - 基于Go-zero的微服务电商系统
[![Go Version](https://img.shields.io/badge/go-1.20%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## 🌟 项目简介

基于Go语言与Go-zero框架开发的轻量级抖音电商系统，支持高并发场景与全链路监控。通过微服务架构实现高扩展性，集成AI查询、分布式事务、缓存预热等企业级功能。

## 🚀 技术栈

### 核心框架

| 类别     | 技术选型     |
|--------|----------|
| 开发语言   | Go 1.20+ |
| RPC 框架 | Go-Zero  |
| 服务治理   | Consul   |
| 消息队列   | RabbitMQ |
| 事务管理   | DTM      |

### 数据存储

| 存储类型 | 技术方案              |
|------|-------------------|
| 关系型  | MySQL 8.0         |
| 缓存   | Redis 6.0         |
| 搜索   | Elasticsearch 8.x |

### 运维体系

| 领域    | 工具链                     |
|-------|-------------------------|
| 容器化   | Docker/K8s              |
| 监控    | Prometheus+Grafana      |
| 日志    | EFK Stack               |
| CI/CD | GitHub Actions + ArgoCD |

## 📁 项目结构

```mermaid

```

## ✨ 核心功能

1. AI 智能查询：自然语言解析商品搜索条件
2. 分布式事务：通过 DTM 实现 TCC/SAGA 事务模式
3. 防超卖系统：Redis Lua 脚本原子库存扣减
4. 智能推荐：基于用户行为的商品推荐引擎
5. 全链路监控：APM + 日志 + 指标三维监控体系
6. 弹性扩容：K8s 自动扩缩容应对流量洪峰

## 🚦 项目亮点

1. 🤖 AI 代码审查：大模型辅助代码质量检查
2. ⚙️ 动态缓存预热：基于访问热度自动缓存加载
3. 🕵️ 操作审计：敏感操作全生命周期追踪
4. 🛡️ 安全防护：JWT 指纹验证 + SQL 注入防护
5. 🎨 个性头像：MD5 哈希生成唯一头像 URL

## 📦 启动与安装部署

### 本地启动

1. 克隆项目

```shell
git clone https://github.com/jijizhazha1024/go-mall.git
```

2. 安装依赖

```shell
go mod tidy
```

3. 启动服务

```shell
go run run.go # 启动所有服务
```

### docker-compose 部署

```shell
docker-compose up -d
```

### K8s 部署

```shell
kubectl apply -f manifests/
```

> 注意：
> 以上操作部署方案都需要进行确保基础环境满足要求，如：
> - [基础设施服务](construct/depend/docker-compose.yaml)
> - [.env](.env.example)

## 🛠️ 开发规范

1. 代码提交遵循 Conventional Commits 规范
2. 采用 Git Flow 工作流管理分支
3. 所有 PR 需通过 AI 辅助代码审查
4. 单元测试覆盖率不低于 80%