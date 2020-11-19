# ip-pools
IP 代理池收集工具

### 设置 golang 代理
```shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

go mod tidy 
go mod download     
```

### 查看数据库

```shell
boltbrowser my.db
```

