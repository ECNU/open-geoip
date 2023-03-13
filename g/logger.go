package g

import (
	"fmt"
	"os"

	"github.com/toolkits/pkg/logger"
)

type LoggerSection struct {
	Dir       string `json:"dir"`
	Level     string `json:"level"`
	KeepHours uint   `json:"keepHours"`
}

func InitLog(l LoggerSection) {

	lb, err := logger.NewFileBackend(l.Dir)
	if err != nil {
		fmt.Println("cannot init logger:", err)
		os.Exit(1)
	}
	lb.SetRotateByHour(true)
	lb.SetKeepHours(l.KeepHours)

	logger.SetLogging(l.Level, lb)
}
