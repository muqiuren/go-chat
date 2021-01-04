### About
A websocket-based chat program written in Golang, supports multi-room group chat and public channel chat.

### Features
* Support multi-room group chat
* Support public channel chat
* Heartbeat detection
* Show online user list

### Main Depends Library
The project uses module to manage dependencies, the following is the main dependency library
    github.com/gorilla/websocket

    github.com/gorilla/mux
    
    github.com/holys/initials-avatar
    
    gopkg.in/yaml.v2
    
    github.com/staori/go.uuid

### Directory Structure

```
├─config            // Configuration file directory
├─resource          // Resource catalog
│  ├─assert         // Static resource directory
│  │  ├─font
│  │  ├─image
│  │  ├─script
│  │  └─style
│  └─template       // Template file directory
│     ├─base.html   // Basic template
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
├─config.yaml        // Application configuration file
├─README.md            
├─server.go          // Service start
└─server_test.go   
```  

### Start
    // Download Project
    git clone https://www.github.com/muqiuren/go-chat
    
    // Enter the project root directory
    cd go-chat
    
    // Check dependencies
    go mod tidy
    
    // Start Application,Visiter http://localhost:8000
    go run server.go

### Result
![https://myblog.hatchblog.cn/uploads/20210103/51941357e2c42969a37142cc7bfc4f2c.png](https://myblog.hatchblog.cn/uploads/20210103/51941357e2c42969a37142cc7bfc4f2c.png)
![https://myblog.hatchblog.cn/uploads/20210103/247d42cea63c560f785470c98e45f0cc.png](https://myblog.hatchblog.cn/uploads/20210103/247d42cea63c560f785470c98e45f0cc.png)
![https://myblog.hatchblog.cn/uploads/20210103/c61a8dd2c8c62c2f6985c6cbae714380.png](https://myblog.hatchblog.cn/uploads/20210103/c61a8dd2c8c62c2f6985c6cbae714380.png)

### TODO
- [x] Frontend UI
- [x] Multi-room group chat
- [x] Public chat channel
- [x] Enter and leave room broadcast
- [x] Heartbeat detection
- [x] Show Online user list
- [ ] Turing access

### More
[Use Go to write a websocket-based chat program](https://myblog.hatchblog.cn/article-23.html)