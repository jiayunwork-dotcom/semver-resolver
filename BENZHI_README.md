Go 语义化版本约束求解命令行：resolve 从候选列表挑出满足 ^/~/>= 约束的最高版本，check 判断单版本是否命中；无参数或 serve 在 :8080 提供 /api/resolve、/api/check、/api/parse。

## 构建 / 运行 / 测试

```text
go build ./...
go run . resolve --constraint "^1.2.3" --versions 1.2.0,1.2.3,1.9.0,2.0.0
go run . serve
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
