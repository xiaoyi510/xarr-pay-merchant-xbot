# XArrPay 商户版机器人

[![Release](https://github.com/xiaoyi510/xarr-pay-merchant-xbot/actions/workflows/release.yml/badge.svg)](https://github.com/xiaoyi510/xarr-pay-merchant-xbot/actions/workflows/release.yml)

这是基于 [XBot](https://github.com/xiaoyi510/xbot) 框架开发的 XArrPay 商户版支付系统 OneBot 机器人，为 XArrPay 商户提供便捷的 QQ 机器人服务。

## 项目简介

XArrPay 商户版机器人是专为 XArrPay 支付系统打造的 QQ 机器人，支持商户账户管理、支付统计查询、渠道账户管理等功能，让商户能够通过 QQ 随时查看和管理自己的支付业务。

## 功能特性

- ✅ 基于 WebSocket 反向连接
- ✅ XArrPay 商户系统集成
- ✅ 商户账户绑定/解绑
- ✅ 账户余额查询
- ✅ 套餐信息查看
- ✅ 支付统计查询（今日/本周/本月/总计）
- ✅ 渠道账户管理
- ✅ 群聊白名单控制
- ✅ 超管权限管理
- ✅ 数据脱敏保护
- ✅ 灵活的配置管理
- ✅ 日志记录

## 项目结构

```
xarr-pay-merchant-xbot/
├── config/             # 配置文件目录
│   └── config.yaml    # 机器人配置文件
├── plugins/           # 插件目录
│   └── xarr-merchant/ # XArrPay 商户插件
│       ├── merchant.go    # 插件主文件
│       ├── handlers.go    # 消息处理器
│       ├── client.go      # API 客户端
│       ├── types.go       # 数据类型定义
│       └── utils.go       # 工具函数
├── data/              # 数据存储目录
├── logs/              # 日志目录
├── docker-compose.yml # Docker 编排文件
├── Dockerfile         # Docker 镜像构建文件
├── main.go            # 程序入口
├── go.mod             # Go 模块文件
└── README.md          # 项目说明
```

## 快速开始

### 环境要求

- Go 1.18 或更高版本
- OneBot 协议实现（如 go-cqhttp、Lagrange 等）
- XArrPay 商户版支付系统

### 安装步骤

1. **克隆项目**

```bash
git clone https://github.com/xiaoyi510/xarr-pay-merchant-xbot
cd xarr-pay-merchant-xbot
```

2. **安装依赖**

```bash
go mod tidy
```

3. **配置机器人**

编辑 `config/config.yaml` 文件，配置机器人连接信息：

```yaml
bot:
  nickname: ["商户机器人", "XArrPay"]
  super_users: [123456789]  # 超级用户 QQ 号
  command_prefix: "/"

drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8888"  # OneBot 实现的 WebSocket 地址
    access_token: ""
    reconnect_interval: 5
    max_reconnect: 0
    timeout: 30

log:
  level: "info"
  file: "logs/bot.log"

storage:
  type: "leveldb"
```

4. **配置商户系统**

启动机器人后，超级管理员需要私聊机器人设置商户系统配置：

```
/设置商户系统 https://your-xarrpay-api.com your-secret-key
/设置商户群聊 123456789,987654321
```

5. **运行机器人**

```bash
go run main.go
```

或使用 Docker：

```bash
docker-compose up -d
```

## 配置说明

### 机器人配置 (bot)

- `nickname`: 机器人的昵称列表，用于识别 @ 机器人
- `super_users`: 超级用户 QQ 号列表，拥有商户系统管理权限
- `command_prefix`: 命令前缀，默认为 "/"

### 驱动配置 (drivers)

- `type`: 驱动类型，支持 `ws_reverse`（反向 WebSocket）
- `url`: OneBot 实现的 WebSocket 地址
- `access_token`: 访问令牌（可选）
- `reconnect_interval`: 重连间隔（秒）
- `max_reconnect`: 最大重连次数，0 表示无限重连
- `timeout`: 连接超时时间（秒）

### 日志配置 (log)

- `level`: 日志级别（debug/info/warn/error）
- `file`: 日志文件路径

### 存储配置 (storage)

- `type`: 存储类型，支持 `leveldb`、`redis` 等

### 商户系统配置

商户系统配置通过机器人命令动态设置，不在配置文件中：

- **API 配置**: 通过 `/设置商户系统` 命令配置 XArrPay API 地址和密钥
- **群聊白名单**: 通过 `/设置商户群聊` 命令设置允许使用的群聊
- **权限控制**: 超级管理员拥有系统配置权限，普通用户只能查询自己的信息

## 功能使用说明

### 账户管理

#### 绑定商户账号
用户需要先在 XArrPay 商户系统中获取绑定 ticket，然后使用以下命令绑定：

```
/绑定 <ticket>
```

#### 解绑商户账号
```
/解绑
```

#### 查看个人信息
```
/我的信息
/个人信息
```

#### 查询账户余额
```
/余额
/查询余额
```

### 套餐管理

#### 查看套餐信息
```
/套餐信息
/我的套餐
```

显示内容包括：
- 套餐名称和到期时间
- 通道账号数限制
- 日限额和月限额
- 当前费率

### 统计查询

#### 今日统计
```
/今日统计
/今日
```

#### 完整统计
```
/统计
/支付统计
```

显示今日/本周/本月/总计的支付金额和订单数量。

#### 渠道账户列表
```
/渠道列表
/账户列表
```

查看所有渠道账户的状态、在线情况和今日额度使用情况。

### 超级管理员功能

#### 设置商户系统
```
/设置商户系统 <API地址> <Secret密钥>
```

#### 设置允许的群聊
```
/设置商户群聊 <群号1,群号2,...>
```

#### 查看系统配置
```
/查看商户配置
```

#### 帮助菜单
```
/商户帮助
```

## 数据安全

### 数据脱敏

为保护用户隐私，机器人在群聊中会对敏感信息进行脱敏处理：

- **用户ID**: `123****789`
- **用户名**: `张****`  
- **渠道账户名**: `支付****账户`

在私聊中显示完整信息，群聊中自动脱敏。

### 权限控制

- **超级管理员**: 可配置系统、查看配置、设置群聊白名单
- **普通用户**: 只能查询和管理自己的账户信息
- **群聊限制**: 只有白名单内的群聊可以使用商户功能

## 常见问题

### 1. 连接不上 OneBot 实现？

- 检查 `config.yaml` 中的 WebSocket 地址是否正确
- 确认 OneBot 实现（如 go-cqhttp）已正常运行
- 检查防火墙和网络设置

### 2. 商户功能无法使用？

- 确认超级管理员已设置商户系统配置（API 地址和密钥）
- 检查群聊是否在白名单中（群聊使用时）
- 确认用户已正确绑定商户账号

### 3. 绑定失败？

- 确认 ticket 是否有效且未过期
- 检查 XArrPay 商户系统是否正常运行
- 查看日志文件确认 API 连接状态

### 4. 查询信息失败？

- 确认用户已绑定商户账号
- 检查商户账号状态是否正常
- 确认网络连接和 API 配置正确

### 5. 如何调试？

- 将日志级别设置为 `debug`
- 查看 `logs/bot.log` 文件
- 检查商户系统 API 响应
- 使用 `/查看商户配置` 确认配置正确

## 部署说明

### Docker 部署

项目提供了 Docker 支持，可以使用以下方式快速部署：

```bash
# 构建镜像
docker build -t xarr-merchant-bot .

# 或使用 docker-compose
docker-compose up -d
```

### 生产环境建议

- 使用 Redis 作为存储后端（高并发场景）
- 配置日志轮转防止日志文件过大
- 设置合理的 API 超时时间
- 定期备份配置和数据
- 监控机器人运行状态

## 相关链接

- [XBot 框架](https://github.com/xiaoyi510/xbot)
- [OneBot 协议](https://github.com/botuniverse/onebot)
- [Go 官方文档](https://golang.org/doc/)

## 许可证

本项目基于 Apache 2.0 许可证开源。

## 贡献

欢迎提交 Issue 和 Pull Request！如果您有好的建议或发现了 Bug，请随时联系我们。

## 支持

如果您在使用过程中遇到问题，可以：

1. 查看本文档的常见问题部分
2. 提交 GitHub Issue
3. 联系项目维护者

