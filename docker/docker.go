package docker

import (
	"os"
)

// CheckContainerRestart 检查容器是否有重启标记
func CheckContainerRestart() bool {
	_, err := os.Stat("/tmp/started.flag")
	if err != nil {
		file, _ := os.Create("/tmp/started.flag")
		if file != nil {
			file.Close()
		}
		return false
	}

	return true
}
