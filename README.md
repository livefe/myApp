# 项目名称

该项目是一个基于Go语言、Gin框架、MySQL数据库、Redis缓存的后端API，提供用户管理和社区互动功能，支持用户认证、信息查询以及社区帖子的管理、评论和点赞功能。后端使用JWT进行身份验证，采用分层架构设计。

## 项目目录结构

```
project/
├── cmd/                              # 命令行工具目录
│   ├── server/                       # HTTP服务器
│   │   └── main.go                   # 服务器入口文件
│   ├── migrate/                      # 数据库迁移工具
│   │   └── migrate.go                # 数据库迁移入口文件
│   └── seed/                         # 测试数据生成工具
│       └── seed.go                   # 测试数据生成入口文件
├── config/                           # 配置管理
│   ├── config.go                     # 配置读取与管理
│   └── config.yaml                   # 配置文件，区分开发与生产环境
├── database/                         # 数据库管理
│   └── database.go                   # 数据库连接与初始化
├── handler/                          # 路由处理层
│   ├── user/                         # 用户模块相关请求处理
│   ├── community/                    # 社区模块相关请求处理
│   └── order/                        # 订单模块相关请求处理
├── service/                          # 业务层，封装业务逻辑
│   ├── user.go                       # 用户服务逻辑
│   ├── community.go                  # 社区模块业务逻辑
│   └── order.go                      # 订单模块业务逻辑
├── repository/                       # 数据存取层，负责与数据库的交互
│   ├── user.go                       # 用户数据存取
│   ├── community.go                  # 社区数据存取
│   └── order.go                      # 订单数据存取
├── model/                            # 数据模型层
│   ├── user.go                       # 用户模型
│   ├── community.go                  # 社区模型
│   └── order.go                      # 订单模型
├── router/                           # 路由管理层
│   ├── router.go                     # 总路由，初始化所有模块的路由
│   ├── user.go                       # 用户模块路由
│   ├── community.go                  # 社区模块路由
│   └── order.go                      # 订单模块路由
├── middleware/                       # 中间件
│   ├── jwt.go                        # JWT验证中间件
│   ├── cors.go                       # 跨域中间件
│   ├── logger.go                     # 请求日志中间件
│   └── rate_limiter.go               # 请求限流中间件
├── test/                             # 测试目录
│   └── api_test.go                   # API测试文件
├── README.md                         # 项目说明文件
├── go.mod                            # Go模块管理
└── go.sum                            # Go模块管理
```

## 项目概述

这个项目实现了一个后端API，支持以下主要功能：
- 用户注册、登录和信息管理
- 使用JWT进行用户认证
- 社区模块，支持创建帖子、评论和点赞

该项目采用了分层架构设计，按照功能划分了多个模块，以增强可维护性和可扩展性。

## 技术栈

- **Go**: 编程语言
- **Gin**: Web框架
- **GORM**: ORM框架，用于数据库操作
- **MySQL**: 数据库管理系统
- **Redis**: 缓存系统
- **JWT**: JSON Web Token，用于用户认证

## 安装与配置

### 1. 克隆代码库

```bash
git clone https://your-repository-url.git
cd your-project-directory
```

### 2. 安装依赖

使用Go模块管理安装依赖：

```bash
go mod tidy
```

### 3. 配置文件

在 `config/config.yaml` 文件中，配置数据库连接、Redis配置和JWT密钥。你可以选择开发环境 (`dev`) 或生产环境 (`prod`) 配置。

**示例 config.yaml**

```yaml
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  dbname: "your_db_name"

redis:
  host: "localhost"
  port: 6379

jwt:
  secret: "your-secret-key"
```

### 4. 数据库初始化

确保数据库已创建并配置正确，使用以下命令运行数据库迁移和测试数据生成：

```bash
# 运行数据库迁移
go run cmd/migrate/migrate.go

# 生成测试数据（可选）
go run cmd/seed/seed.go
```

### 5. 启动服务

启动Go应用：

```bash
go run cmd/server/main.go
```

服务默认会在 `localhost:8080` 启动。

## 功能说明

### 用户模块

- **POST /api/user/login**: 用户登录，返回JWT Token
- **GET /api/user/info**: 获取当前登录用户信息
- **PUT /api/user/profile**: 更新用户资料

### 社区模块

- **POST /api/community/post**: 发布帖子
- **GET /api/community/posts**: 获取所有帖子
- **POST /api/community/comment**: 对帖子发表评论
- **POST /api/community/like**: 点赞帖子

## 中间件

- **JWT验证**: 所有需要登录的接口都要求携带有效的JWT Token。
- **跨域支持 (CORS)**: 支持跨域请求。
- **请求日志**: 所有请求会记录日志，便于调试与监控。
- **限流**: 对高频请求进行限制，防止滥用。

## 代码结构

### `config/` - 配置管理

- `config.go`: 负责加载和解析配置文件。
- `config.yaml`: 配置文件，包含数据库、Redis和JWT相关配置。

### `handler/` - 路由处理层

处理请求的逻辑，每个模块有独立的处理文件。

- `user/`: 用户模块请求处理文件，包括登录、用户信息管理等。
- `community/`: 社区模块请求处理文件，包括帖子、评论、点赞等。

### `service/` - 业务逻辑层

封装具体的业务逻辑。

- `user.go`: 用户相关的服务逻辑。
- `community.go`: 社区相关的服务逻辑。

### `repository/` - 数据存取层

负责与数据库的交互，数据存取操作。

- `user.go`: 用户相关的数据存取。
- `community.go`: 社区相关的数据存取。

### `model/` - 数据模型层

定义数据库模型，映射到相应的数据库表。

- `user.go`: 用户模型。
- `community.go`: 社区模型。

### `router/` - 路由管理

- `router.go`: 初始化所有路由。
- `user.go`: 用户模块的路由设置。
- `community.go`: 社区模块的路由设置。

### `middleware/` - 中间件

- `jwt.go`: JWT验证中间件，验证每个请求的JWT Token。
- `cors.go`: 支持跨域请求的中间件。
- `logger.go`: 记录请求日志的中间件。
- `rate_limiter.go`: 限制请求频率的中间件。

### `pkg/` - 公共工具层

存放公共工具类。

- `jwt.go`: 用于JWT生成和验证的工具类。
- `logger.go`: 日志工具类。
- `redis.go`: Redis工具类，用于连接和操作Redis。
- `error.go`: 错误处理工具类。
- `validator.go`: 用于请求参数验证的工具类。

## 开发与贡献

### 如何贡献

我们欢迎任何形式的贡献！如果你发现bug或有任何建议，欢迎提交Issue或Pull Request。

1. Fork 本仓库。
2. 创建新的分支 (`git checkout -b feature/your-feature`).
3. 提交修改 (`git commit -am 'Add new feature'`).
4. 推送到远程仓库 (`git push origin feature/your-feature`).
5. 创建Pull Request。

### 常见问题

1. **如何配置开发和生产环境？**

   在 `config/config.yaml` 中分别配置开发环境 (`dev`) 和生产环境 (`prod`) 配置信息。根据环境变量或配置文件切换不同的数据库、Redis配置等。

2. **如何扩展功能？**

   在 `handler/` 和 `service/` 目录下分别增加新的处理和服务逻辑文件。

## 许可证

该项目使用 MIT 许可证，详细信息请参见 LICENSE 文件。

---

© 2025 项目版权所有。
