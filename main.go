package main

import (
	"github.com/aoaostar/mooc/gui"
	"github.com/sirupsen/logrus"
)

func main() {
	// 运行GUI应用程序
	err := gui.RunApp()
	if err != nil {
		logrus.Fatal(err)
	}
}
