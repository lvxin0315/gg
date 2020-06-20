package helper

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

var m runtime.MemStats

//调试用，输出内存使用情况
func PrintMemState() {
	runtime.ReadMemStats(&m)
	logrus.Debugf("%d Kb\n", m.Alloc/1024)
}
