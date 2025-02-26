# go-mall


## 运行脚本使用
项目目录下的run.go为服务运行脚本，使用如下：
```shell
go run run.go --services service1,service2
```
example1:
```shell
go run run.go --services user,goods
```

example2:
不进行指定services默认启动所有服务

```shell
go run run.go
```