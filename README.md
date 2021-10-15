# Plugin-Fortune
ZeroBot-Plugin-Dynamic 的动态库插件。

# 使用方法
- **每日运势** `import _ github.com/FloatTech/ZeroBot-Plugin/plugin_fortune`
    - [x] 运势|抽签
    - [x] 设置底图[车万 DC4 爱因斯坦 星空列车 樱云之恋 富婆妹 李清歌 公主连结 原神 明日方舟 碧蓝航线 碧蓝幻想 战双 阴阳师]
### 编译
#### 使用Actions编译
1. fork 本仓库
2. 上传修改后的`main.go`及其它新增文件。
3. 创建形如`v1.2.3`的`tag`，触发插件编译流程。
4. 编译好后前往`Release`页面下载即可。
#### 本地编译
```bash
# 本机架构
go build -ldflags "-s -w" -buildmode=plugin -o demo.so
# 交叉编译：详见 workflow 相关代码
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 CC=arm-linux-gnueabihf-gcc-9 CXX=g++-9-arm-linux-gnueabihf go build -ldflags="-s -w" -buildmode=plugin -o artifacts/zbpd-linux-armv6
```
### 开始使用
放置动态库到[ZeroBot-Plugin-Dynamic](https://github.com/FloatTech/ZeroBot-Plugin-Dynamic)的`plugins/`目录下，给机器人发送`/刷新插件`即可，或重启也可加载。
