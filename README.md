# beego-demo

A web demo using Beego framework, with MongoDB, MySQL and Redis support.

这是一个基于[Beego](http://beego.me)框架构建的应用demo，后台数据库使用[MongoDB](http://www.mongodb.org)和[MySQL](http://www.mysql.com)，并使用[Redis](http://redis.io)存储session和一些统计数据。

## API列表1

该部分使用的数据库是MongoDB。

| 功能 | URL | Mode |
|------|:-----|------|
| 注册 | /v1/users/register | POST |
| 登录 | /v1/users/login    | POST |
| 登出 | /v1/users/logout   | POST |
| 修改密码 | /v1/users/passwd   | POST |
| 上传多个文件 | /v1/users/uploads   | POST |

输入数据通过form表单提交，返回数据均为json。

在static/test/目录下有如下的测试表单，除了用于测试外，也可看出具体的数据通讯协议：
* register.html
* login.html
* logout.html
* passwd.html
* uploads.html

## API列表2

该部分使用的数据库是MySQL。

| 功能 | URL | Mode |
|------|:-----|------|
| 获取一个角色信息 | /v1/roles/:id | GET    |
| 获取所有角色信息 | /v1/roles     | GET    |
| 新增一个角色信息 | /v1/roles     | POST   |
| 修改一个角色信息 | /v1/roles/:id | PUT    |
| 删除一个角色信息 | /v1/roles/:id | DELETE |

表结构如下：

| Field    | Type         | Null | Key |
|----------|:-------------|------|-----|
| id       | bigint(20)   | NO   | PRI |
| name     | varchar(255) | YES  |
| password | varchar(255) | YES  |
| reg_date | datetime     | YES  |

建数据库表的脚本：scripts/sql/db.sql。

多记录api，提供如下参数：
* query=col1:op1:val1,col2:op2:val2 ...
* order=col1:asc|desc,col2:asc|esc ...
* limit=n，缺省为10
* offset=n，缺省为0

query的op值：
* eq，等于
* ne，不等于
* gt，大于
* ge，大于等于
* lt，小于
* le，小于等于

## 环境

### GO语言

包括安装go，设置$GOPATH等，具体可参考：[How to Write Go Code](http://golang.org/doc/code.html)。

### MongoDB

在conf/app.conf中设置MongoDB参数，如：

```
[mongodb]
url = mongodb://127.0.0.1:27017/beego-demo
```

完整的url写法可参考：http://godoc.org/gopkg.in/mgo.v2#Dial

### MySQL

在conf/app.conf中设置MySQL参数，如：

```
[mysql]
url = root:root@/beego-demo?charset=utf8&parseTime=True&loc=Local
```

完整的url写法可参考：https://github.com/go-sql-driver/mysql#dsn-data-source-name

### Redis

在conf/app.conf中设置Redis参数，涉及两个地方，一个是session，一个是cache，两者可以不同：

```
sessionsavepath = 127.0.0.1:6379

[cache]
server = 127.0.0.1:6379
password =
```

### Beego

安装/升级所有依赖包：

```
$ go get -u github.com/astaxie/beego
$ go get -u github.com/beego/bee
$ go get -u github.com/astaxie/beego/session/redis
$ go get -u gopkg.in/mgo.v2
$ go get -u github.com/garyburd/redigo/redis
$ go get -u github.com/go-sql-driver/mysql
$ go get -u golang.org/x/crypto/scrypt
```

当前版本：

```
$ bee version
bee   :1.2.4
beego :1.4.3
Go    :go version go1.4.2 darwin/amd64
```

## 运行

将代码放在$GOPATH/src/目录下，运行（开发模式）：

```
$ cd $GOPATH/src/beego-demo/
$ bee run
```

