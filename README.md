# ccuser
CCUser CLI tool. Written in Go.


## download
https://github.com/mozillazg/ccuser/releases

## usage

### 基本用法

```
$ ccuser --help
$ ccuser status
$ ccuser -u username -p password login
$ ccuser -u username -p password logout
$
$ export CCUSER_USERNAME='username'
$ export CCUSER_PASSWORD='password'
$ ccuser login
$ ccuser logout
```

### 在公司

```
$ ccuser login
```

### 在家或 VPN

```
$ ccuser -b login     # heartbeat 模式，定期发送心跳请求，防止被下线
```
