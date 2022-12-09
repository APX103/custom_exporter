# deploy_platform_exporter

> Prometheus exporter for embedded devices of mmdeploy platform

## Usage

``` sh
cd deploy_platform_exporter
go build -o deploy_platform_exporter main.go
nohup ./deploy_platform_exporter 2>&1 >/dev/null&
```
