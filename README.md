
## CloudHands
``` 
bot --自动化攻击机器人，本项目主要用于检测资产安全漏洞，至于用于不法的攻击行为请自负后果！！！

```

### 依赖软件
```
go version 1.15以上
redis --用来存储攻击结果
```
   
### 编译方法
```
1. cd GBWBot
2. sh build.sh
3. cd build
   sh install.sh
   
编译安装完成之后，目录结构：
/opt/bot/sbot/bin ---可执行程序，主要程序sbot 自动化攻击服务端控制程序，cbot自动化攻击客户端攻击程序
/opt/data/store/attack/script --攻击武器库
### 运行方法
启动服务端控制程序：/opt/bot/sbot/bin/sbot -cfile /opt/bot/sbot/conf/sbot.json



