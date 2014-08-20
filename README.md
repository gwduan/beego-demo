# beego-demo

A web demo using Beego framework, with MongoDB and Redis support.

这是一个基于[Beego](http://beego.me)框架构建的应用demo，后台数据库使用[MongoDB](http://www.mongodb.org)，并使用[Redis](http://redis.io)存储session和一些统计数据。

## API列表

| 功能 | URL | Mode |
|------|:-----|------|
| 注册 | /v1/users/register | POST |
| 登录 | /v1/users/login    | POST |
| 登出 | /v1/users/logout   | POST |
| 改密 | /v1/users/passwd   | POST |

输入数据通过form表单提交，返回数据均为json。

在static/test/目录下有如下的测试表单，除了用于测试外，也可看出具体的数据通讯协议：
* register.html
* login.html
* logout.html
* passwd.html

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

### Redis

在conf/app.conf中设置Redis参数，涉及两个地方，一个是session，一个是cache，两者可以不同：

```
sessionsavepath = 127.0.0.1:6379

[cache]
server = 127.0.0.1:6379
password =
```

### Beego

安装所有依赖包：

```
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
$ go get github.com/astaxie/beego/session/redis
$ go get gopkg.in/mgo.v2
$ go get github.com/garyburd/redigo/redis
```

当前版本：

```
$ bee version
bee   :1.2.1
beego :1.4.0
Go    :go version go1.3.1 darwin/amd64
```

## 运行

将代码放在$GOPATH/src/目录下，运行（开发模式）：

```
$ cd $GOPATH/src/beego-demo/
$ bee run
```

