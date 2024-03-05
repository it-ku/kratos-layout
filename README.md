# Kratos Layout Project Template

> GO版本 ：  1.20

## How to use 
```
kratos new helloworld -r https://github.com/wx-micro/kratos-layout.git
```


## 对比官方的改动如下

- `GORM`  [官网直达](https://gorm.io/)
- `Redis` , Redis的版本最低需要7 ，如果使用6需要对 `github.com/go-redis/redis/v8` 降级
- `etcd`
- 修改了 `Makefile` 文件(不需要的手动删除)：
  
   + api 增加了 `--validate_out=paths=source_relative,lang=go:./api \` 参数校验
   + api 增加了 `--go-errors_out=paths=source_relative:./api \` 错误处理
