# beego-demo

A web demo using Beego framework, with MongoDB, MySQL and Redis support.

这是一个基于 [Beego](http://beego.me) 框架构建的应用 demo，后台数据库使用 [MongoDB](http://www.mongodb.org) 和 [MySQL](http://www.mysql.com)，并使用 [Redis](http://redis.io) 存储 session 和一些统计数据。

## API列表

### 第一部分

该部分使用的数据库是 MongoDB 和 Redis。

| 功能 | URL | Mode |
|------|:-----|------|
| 注册 | /v1/users/register | POST |
| 登录 | /v1/users/login    | POST |
| 登出    | /v1/users/logout   | POST |
| 修改密码   | /v1/users/passwd   | POST |
| 上传多个文件 | /v1/users/uploads   | POST |
| 下载文件 | /v1/users/downloads   | GET |

在 static/test 目录下有如下的测试表单，除了用于测试外，也可看出具体的数据通讯协议：

* register.html
* login.html
* logout.html
* passwd.html
* uploads.html

说明：

* 输入数据通过 form 表单提交，返回数据均为 json。
* 使用 Beego 的 ParseForm 功能将输入数据解析到 struct 中。
* 使用 Beego 的 Validation 功能对数据进行校验。
* 使用 [scrypt](https://godoc.org/golang.org/x/crypto/scrypt) 算法进行密码处理。
* 对于数据库返回的错误，单独区分“记录不存在”和“记录重复”两种错误。
* 由于 Beego 本身不支持多文件上传，故单独实现了 uploads API 来展示该功能，该功能与数据库无关。

### 第二部分

该部分使用的数据库是 MySQL。

| 功能 | URL | Mode |
|------|:-----|------|
| 获取一个角色信息 | /v1/roles/:id  | GET    |
| 获取所有角色信息 | /v1/roles      | GET    |
| 新增一个角色信息 | /v1/roles      | POST   |
| 修改一个角色信息 | /v1/roles/:id  | PUT    |
| 删除一个角色信息 | /v1/roles/:id  | DELETE |
| 认证           | /v1/roles/auth | POST |

roles表结构如下：

| Field    | Type         | Null | Key |
|----------|:-------------|------|-----|
| id       | bigint(20)   | NO   | PRI |
| name     | varchar(255) | YES  |
| password | varchar(255) | YES  |
| reg_date | datetime     | YES  |

初始建数据库表的脚本位于：scripts/sql/db.sql。

多记录 api，提供如下参数：

* query=col1:op1:val1,col2:op2:val2 ...
* order=col1:asc|desc,col2:asc|esc ...
* limit=n，缺省为 10
* offset=n，缺省为 0

query 的 op 值：

* eq，等于
* ne，不等于
* gt，大于
* ge，大于等于
* lt，小于
* le，小于等于

说明：

* 参考 RESTful 模式设计 API。
* 使用 JSON Web Token (JWT) 做认证手段。
* 输入数据采用 json，返回数据也是 json。
* 数据库操作使用原生 SQL，没有采用 ORM。
* 对可能的 NULL 值做了处理。
* 多记录查询通过拼接 SQL 语句实现，故对输入参数做了一些校验和处理。
* 同样单独区分“记录不存在”和“记录重复”两种数据库错误。

## 环境

### GO语言

包括安装 go，设置 $GOPATH 等，具体可参考：[How to Write Go Code](http://golang.org/doc/code.html)。

### MongoDB

在 conf/app.conf 中设置 MongoDB 参数，如：

```
[mongodb]
url = mongodb://127.0.0.1:27017/beego-demo
```

完整的 url 写法可参考：http://godoc.org/gopkg.in/mgo.v2#Dial

这里单独封装了一个 mymongo 包来实现数据库的初始化，以简化后续的数据库操作。

### MySQL

在 conf/app.conf 中设置 MySQL 参数，如：

```
[mysql]
url = root:root@/beego-demo?charset=utf8&parseTime=True&loc=Local
```

完整的 url 写法可参考：https://github.com/go-sql-driver/mysql#dsn-data-source-name

这里单独封装了一个 mymysql 包来实现数据库的初始化，以简化后续的数据库操作。

### Redis

在 conf/app.conf 中设置 Redis 参数，涉及两个地方，一个是 session，一个是 cache，两者可以不同：

```
sessionproviderconfig = 127.0.0.1:6379

[cache]
server = 127.0.0.1:6379
password =
```
这里单独封装了一个 myredis 包来实现数据库的初始化，以简化后续的数据库操作。

## 运行

将代码放在 $GOPATH/src 目录下：

```
$ cd $GOPATH/src/
$ git clone https://github.com/gwduan/beego-demo.git
```

从 go1.7 开始，使用 [glide](https://github.com/Masterminds/glide) 工具来管理依赖包，需要事先安装好 glide 。

安装依赖包：

```
$ cd $GOPATH/src/beego-demo/
$ glide install
```

如果不使用 glide ，也可手工安装依赖包，但可能会有包版本的兼容问题：

```
$ go get -u github.com/astaxie/beego
$ go get -u github.com/astaxie/beego/session/redis
$ go get -u gopkg.in/mgo.v2
$ go get -u github.com/garyburd/redigo/redis
$ go get -u github.com/go-sql-driver/mysql
$ go get -u golang.org/x/crypto/scrypt
$ go get -u github.com/dgrijalva/jwt-go
```

安装 bee 工具，调试阶段很好用：

```
$ go get -u github.com/beego/bee
```

开始运行：

```
$ cd $GOPATH/src/beego-demo/
$ bee run
```

当前版本：

```
$ bee version
______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.9.1

├── Beego     : 1.9.0
├── GoVersion : go1.9.1
├── GOOS      : darwin
├── GOARCH    : amd64
```

## 部署

正式部署时，可通过系统的 Init 服务来启动。在 scripts 目录下有 upstart 和 systemd 两套简易示例脚本，可参考使用。

例如，在 CentOS 6 下，复制 upstart/bdemo.conf 到 /etc/init/，相应修改后，执行：

```
# start bdemo
```

在 CentOS 7 下，复制 systemd/bdemo.service 到 /etc/systemd/system/，相应修改后，执行：

```
# systemctl daemon-reload
# systemctl enable bdemo.service
# systemctl start bdemo.service
```

由于 Init 是由 root 控制的，相应的服务缺省也具有 root 权限，故一般都应该做降权处理。可在 systemd 和 upstart 脚本中设置运行时的普通用户名和组名，具体可参考官方文档。

降权的问题在于普通用户无法绑定特权端口（如 80 ），不过实际环境下，还是建议在前面部署 Nginx 等成熟的 web 服务器，通过反向代理来访问应用。
