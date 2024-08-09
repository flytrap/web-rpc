# web-rpc

通过 web 接口远程执行命令， 用于解决容器内部需要调用宿主机命令的情况

## 启动命令

```bash
go run main.go web -c config.json
```

## 调用

```bash
curl --location --request POST 'localhost:8080/api/rpc/v1/run/reloadNginx'
```

## 配置

commands

- code: 命令代码
- exec: 具体命令字符串

## 编译

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/web-rpc-linux .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/web-rpc-darwin .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/web-rpc.exe .
```
