# GoMall

--- 



## 贡献者

--- 


## 项目结构

--- 
- apis: 网关定义
- common: 通用模块
  - config: 配置文件解析
    consts: 常量定义
    middleware: 中间件
    response: 响应格式
    utils: 工具包
- construct: Docker Compose 文件
  - depend: 基础依赖
  - observability: 可观性
- dal: 数据访问层
- services: 微服务实例
- test: 测试文件
  - rpc: rpc测试
  - web: web测试

## 服务依赖

---
### 强制依赖
- mysql
- redis
- rabbitmq
- consul

### 非强制依赖
- Jaeger



