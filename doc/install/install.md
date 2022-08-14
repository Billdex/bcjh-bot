# 下载安装

## 部署架构

![部署架构图](../media/部署架构.png)

## 下载

> 以下链接可能需要 `科学上网`，如果不能打开请自行搜索科学上网解决方案，如 [`iGuge`](https://iguge.xyz/)

爆炒江湖机器人查询程序 `bcjh-bot`：[下载地址](https://github.com/Billdex/bcjh-bot/releases)

第三方QQ客户端 `go-cqhttp`：[下载地址](https://github.com/Mrs4s/go-cqhttp/releases)

下载时请根据自己的 电脑/服务器 的系统版本选择合适的程序文件。

- windows/64位/pc 选择后缀为 `_windows_amd64`
- mac/intel芯片 选择后缀为 `_darwin_amd64`
- mac/m1、m2芯片 选择后缀为 `_darwin_arm64`
- linux 选择后缀为 `_linux`，具体子版本由所用的芯片架构决定

## 运行程序

!> 运行程序之前请先将前一步下载的压缩文件解压，请勿直接在压缩文件内运行程序！

### bcjh-bot

linux/mac 用户在解压路径使用命令行运行 `./bcjh-bot` 即可启动程序。

windows 用户双击`run.bat` 脚本启动程序，请勿直接双击 `bcjh-bot.exe`。

第一次运行时会生成默认配置文件 `config.ini`，根据自身需求修改配置文件，新手只需将 `admin` 选项改为自己用于程序管理员的QQ号即可。

修改完配置文件后重新运行程序。

各配置选项的详细说明请参考 [配置信息](./config.md)

### go-cqhttp

关于 `go-cqhttp` 的详细使用说明请参考 [go-cqhttp文档](https://docs.go-cqhttp.org/guide/quick_start.html)，这里仅列出一些关键步骤。

首次运行程序会提示选择通信方式，输入 `3` 并回车，选用反向 `websocket`，此时会在程序目录下自动生成配置文件 `config.yml`，需要编辑以下几项内容：

- **uin**: 此处数字改为你需要登录的机器人账号，建议使用小号登录机器人（用大号如果被封了就惨了）
- **password**: 登录密码。如果这项为空则需要在启动程序用扫码登录
- **universal**: 查询程序的服务地址，默认请使用 `ws://127.0.0.1:5800` 。如果修改了查询程序的 `port` 选项，则将结尾的数字改为相同的数值。

> 此处登录的机器人账号建议开启设备锁，否则会大大增加被风控的概率

修改完配置文件后重新运行程序，按照流程指引登录机器人账号即可。

## 获取数据

首次启动 `bcjh-bot` 与 `go-cqhttp` 程序时，数据库内是没有游戏相关数据的，暂时还无法进行游戏查询。

使用管理员账号向机器人登录的账号发送消息 `#更新`，机器人程序就会从图鉴网拉取游戏数据。取决于程序运行时的网络环境与计算机性能，整个更新过程可能会花费 10~20 分钟，请耐心等待。其中主要耗时在于绘制图鉴图片，一般在开始更新后的一分钟内即可查询文字数据。

更新时可以采用以下数据源：

- lgithub： L 图鉴网（[https://foodgame.github.io](https://foodgame.github.io)），可以最早获取到游戏最新数据，如没有科学上网可能访问很慢。
- lgitee：L 图鉴网的 gitee 版本（[https://foodgame.gitee.io](https://foodgame.gitee.io)），gitee政策更新后可能不能用了。
- 白菜菊花cf：小鱼部署在 CF Page 的上的白菜菊花网址（[https://bcjh.pages.dev](https://bcjh.pages.dev)），默认使用该数据源。
- 白菜菊花：  小鱼部署在自己服务器上的白菜菊花网址（[https://bcjh.xyz](https://bcjh.xyz)），上面的数据源如果用不了可以尝试用这个。

使用时作为参数添加到命令后面即可使用指定数据源更新，如 `#更新 白菜菊花`。

完成更新之后机器人账号将会回复一条关于本次数据更新耗时的消息，然后便可以尽情开始玩耍啦~

