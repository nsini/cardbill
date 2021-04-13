# 信用卡小账本

小账本是一个专门用来管理您的信用卡刷卡记录的小工具。

## 架构设计

该平台提供了一整套解决方案。

## 平台演示

演示地址: [https://cardbill.nsini.com/](https://cardbill.nsini.com/)

该服务部署在开普勒云平台上, [https://github.com/kplcloud/kplcloud](https://github.com/kplcloud/kplcloud)

开普勒平台演示地址: [https://kplcloud.nsini.com/about.html](https://kplcloud.nsini.com/about.html)

## 小程序

![](http://source.qiniu.cnd.nsini.com/images/2021/04/74/68/c1/20210412-4d75f32bc458108dc2f9c1f968574713.jpg?imageView2/2/w/1280/interlace/0/q/70)

## 安装说明

平台后端基于[go-kit](https://github.com/go-kit/kit)、前端基于[ant-design](https://github.com/ant-design/ant-design)(版本略老)框架进行开发。

后端所使用到的依赖全部都在[go.mod](go.mod)里，前端的依赖在`package.json`，详情的请看`yarn.lock`，感谢开源社区的贡献。

后端代码: [https://github.com/icowan/cardbill](https://github.com/icowan/cardbill)

前端代码: [https://github.com/icowan/cardbill-view](https://github.com/icowan/cardbill-view)

### 安装教程

[安装教程]

### 依赖

- Golang 1.13+ [安装手册](https://golang.org/dl/)
- MySQL 5.7+ (大多数据都存在mysql)

## 快速开始

1. 克隆

```
$ mkdir -p $GOPATH/src/github.com/icowan
$ cd $GOPATH/src/github.com/icowan
$ git clone https://github.com/icowan/cardbill.git
$ cd cardbill
```

2. 配置文件准备

    - app.cfg文件配置也放到该项目目录app.cfg配置请参考 [配置文件解析](https://docs.nsini.com/start/config.html)

3. docker-compose 启动

```
$ cd install/docker-compose
$ docker-compose up
```

4. make 启动

```
$ make run
```

## 文档

[文档]

### 视频教程

- [本地启动]
- [开普勒平台部署]
- [信用卡管理]
- [消费记录]

### 支持我

![](https://lattecake.oss-cn-beijing.aliyuncs.com/static%2Fimages%2Freward%2Fweixin-RMB-xxx.JPG)
![](https://lattecake.oss-cn-beijing.aliyuncs.com/static%2Fimages%2Freward%2Falipay-RMB-xxx.png)