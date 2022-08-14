# 配置信息

默认生成的配置文件 `config.ini` 内容如下所示：

```ini
[server]
port = 5800   # 程序运行的端口号

[bot]
private_msg_max_len = 20   # 私聊消息有多条数据时的最大长度，用于查询列表时候的分页
group_msg_max_len = 10   # 群聊消息有多条数据时的最大长度，用于查询列表时候的分页
exchange_msg_max_len = 3   # 兑换码消息可查询的最大消息数
admin = 123456789   # 管理员QQ号，要配置多个时用半角逗号分隔

# 数据库相关配置
[database]
use_local = true  # 是否开启本地数据库，开启后将在程序目录自动创建一个数据库文件。
# use_local 为false时会读取以下配置连接 MySQL
host = 127.0.0.1:3306
database = bcjh
user = bcjh
password =

# 静态资源路径
[resource]
image = ./resource/image/
font = ./resource/font/
shortcut = ./resource/shortcut/

# 日志相关配置
[log]
# 日志级别: DEBUG, INFO, WARN, ERROR, PANIC, FATAL
level = INFO
# 日志输出路径
out_path = ./logs/bcjh-bot.log
```

新手使用时一般只需要改动 `admin` 选项即可，将这里改为自己希望作为管理员的QQ号。

如果有Mysql数据库使用经验，并且希望让机器人查询程序连接到自己的数据库，可以将 `use_local` 选项设为 `false`，并将 `[database]` 中剩余选项修改为你自己的数据库连接配置。