package bootstrap

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// 自定义日志Hook
type GuiLogHook struct{}

func (h *GuiLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *GuiLogHook) Fire(entry *logrus.Entry) error {
	// 将日志转发给所有观察者
	for _, observer := range logObservers {
		observer.OnLogMessage(entry.Level.String(), entry.Message)
	}
	return nil
}

func InitLog() {
	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	
	// 创建日志目录
	os.MkdirAll("./logs", 0755)
	
	// 创建日志文件
	file, err := os.OpenFile("./logs/aoaostar.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// 同时输出到文件，但不再输出到控制台
		logrus.SetOutput(file)
	} else {
		logrus.Warn("Failed to log to file, using default stderr")
	}
	
	// 添加GUI日志Hook
	logrus.AddHook(&GuiLogHook{})
}
