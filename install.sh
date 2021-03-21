cp ./database/database.example.db ./database/database.db
cp ./config/config.example.yaml ./config/config.yaml
# 开启 go module 功能并加速下载
# 如果执行错误请检查你的 go 版本是否高于 1.13
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go build