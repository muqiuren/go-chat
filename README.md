### 简介
使用Golang编写的基于websocket的聊天程序，支持多房间群聊，公共频道聊天。

### 功能特性
* 支持多房间群聊
* 支持公共频道聊天
* 心跳检测
* 显示在线用户列表

### 主要依赖库
项目使用module管理依赖，下面是主要依赖库

    github.com/gorilla/websocket

    github.com/gorilla/mux
    
    github.com/holys/initials-avatar
    
    gopkg.in/yaml.v2
    
    github.com/staori/go.uuid

### 目录结构

```
├─config            // 配置文件目录
├─resource          // 资源目录
│  ├─assert         // 静态资源目录
│  │  ├─font
│  │  ├─image
│  │  ├─script
│  │  └─style
│  └─template       // 模板文件目录
│     ├─base.html   // 基础模板
│     ├─room.html   
│     └─home.html
├─service           
│  ├─app.go
│  ├─avatar.go
│  ├─client.go
│  ├─helper.go
│  ├─helper_test.go
│  ├─interrupt.go
│  ├─loader.go
│  ├─room.go
│  └─router.go
├─config.yaml        // 应用配置文件
├─README.md            
├─server.go          // 服务启动
└─server_test.go   
```  

### 运行
    // 下载项目
    git clone https://www.github.com/muqiuren/go-chat
    
    // 进入项目根目录
    cd go-chat
    
    // 检查依赖
    go mod tidy
    
    // 启动应用,访问http://localhost:8000
    go run server.go

### 效果
![](https://github.com/muqiuren/go-chat/blob/master/resource/assert/image/1.png)
![](https://github.com/muqiuren/go-chat/blob/master/resource/assert/image/2.png)
![](https://github.com/muqiuren/go-chat/blob/master/resource/assert/image/3.png)

### TODO
- [x] 前端聊天界面
- [x] 多房间
- [x] 公共聊天频道
- [x] 进入离开房间广播
- [x] 心跳检测
- [x] 在线用户列表
- [ ] 图灵接入

### 更多
[使用Go编写基于websocket聊天程序详解](https://myblog.hatchblog.cn/article-23.html)
