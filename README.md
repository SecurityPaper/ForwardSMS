# 项目介绍

4G模块，sim卡短信转发到企业微信软硬一体解决方案

## 使用硬件

使用移远EC20解决方案，[购买链接](https://detail.tmall.com/item.htm?spm=a1z10.5-b-s.w4011-23773508522.66.cd38a48eir4WR3&id=595437612613&skuId=4304510502236)

## 部署方案

1. 插入sim卡
2. 插入电脑usb口
3. 如果指示灯蓝绿交替常亮，代表识别卡信号成功，如果交替闪烁，代表没信号，请安装天线后尝试
4. 使用`lsusb`命令查看是否识别成功
5. 克隆项目到本地 `git clone https://github.com/SecurityPaper/ForwardSMS.git && cd ForwardSMS`
6. 注意，当前目录必须为项目内目录，然后使用命令 `docker run --privileged -v /dev/ttyUSB3:/dev/ttyUSB3 -v ./data:/data securitypaperorg/gammu-smsd:latest echo "a test sms from ec20" | /usr/bin/gammu -c /data/config/gammu-smsd.conf sendsms TEXT 133xxxxxxx`
7. 133xxxxxxx请替换为自己手机号
8. 如果成功发送短信，代表卡识别正确，如果返回`350`代表卡并未搜到信号。如果长时间没反应，请尝试另外几个`/dev/ttyUSBx`,直到发送短信成功。
9. 短信发送成功后配置`data/config/forward.yaml`文件
10. 根据第6条测试出来的usb端口号，配置`docker-compose.yaml`文件
11. 在当前目录执行`docker-compose up -d`

## 文件解释

`data/config/forward.yaml`
```yaml
# 如果有all这个配置，就是默认所有短信都会转发给这个机器人，建议发送给管理员，或者直接删除关闭
all:
  rule: all
  type: all
  url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxx

# 从上到下依次为 项目名称、规则（使用关键字匹配）、匹配方式（后续可能支持正则）、机器人url
测试:
  rule: 测试
  type: keyword
  url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxx
```

---

`data/config/gammu-smsd.conf`

```conf
[gammu]
# 配置短信转发端口，映射为USB3，但是要看清楚主机的端口地址，测试端口请查看文档
port = /dev/ttyUSB3

#连接方式和速率，保持默认即可
connection = at115200

# 配置smsd守护进程的配置
[smsd]

# 使用的存档模式
Service = sql

# 使用具体的数据库
Driver = sqlite3

# 数据库路径
DBDir = /data/db

# 数据库名称
Database = sms.db

# 日志存放位置
logfile = /data/log/gammu-smsd.log

# 开启debug，默认不开启
debuglevel = 0

```
---

`data/config/status.yaml`

```yaml
# 代表当前发送到第多少条短信，建议不要删除，否则会从数据库里第一条一直发到最后一条，文件会根据发送自动更新。
id: 0

```
---
`data/db/sms.db`
> 文件为sqlite3数据库，用来存储短信接收，如果有需要请定期备份，或者可以下载后查看。

---
`data/log/gammu-smsd.log`
> 这个文件是gammu-smsd服务产生的日志