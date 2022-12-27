# deploy_platform_exporter

> Prometheus exporter for embedded devices of mmdeploy platform

## Usage

``` sh
# CD to root of this repo
cd deploy_platform_exporter
# Build
go build -o deploy_platform_exporter main.go
# Run as deamon with no output
nohup ./deploy_platform_exporter 2>&1 >/dev/null&
```

## Extend

### Which functions can be extended

1. Metrics
2. Input Parameters

### Metrics Extend Steps

1. Create a file in [repo_path]/metrics/
2. Create a struct ending with "Metric"
3. Implement "Update() error" function for the struct
4. Register the struct in the same file in `func init()` like this.

``` go
func init() {
	registerCollector("cpu_metrics", defaultEnabled, &CPUMitrics{})
}
```

### Input Parameters Extend Steps

1. import this => `kingpin "gopkg.in/alecthomas/kingpin.v2"`
2. Add params any where in package

``` go
// Like this ↓
procPath = kingpin.Flag("path.procfs", "procfs mountpoint.")
                  .Default(procfs.DefaultMountPoint)
                  .String()
```
