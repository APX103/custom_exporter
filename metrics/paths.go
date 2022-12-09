package metrics

import (
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	procPath   = kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).String()
	sysPath    = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	rootfsPath = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
)

func procFilePath(name string) string {
	return filepath.Join(*procPath, name)
}

func sysFilePath(name string) string {
	return filepath.Join(*sysPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(*rootfsPath, name)
}

func rootfsStripPrefix(path string) string {
	if *rootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, *rootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}
