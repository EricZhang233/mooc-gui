package bootstrap

import (
	"fmt"
	"github.com/aoaostar/mooc/pkg/config"
	"github.com/aoaostar/mooc/pkg/task"
	"github.com/aoaostar/mooc/pkg/util"
	"github.com/aoaostar/mooc/pkg/yinghua"
	"github.com/sirupsen/logrus"
)

// 任务状态观察者接口
type TaskStatusObserver interface {
	OnTaskStatusChanged(task task.Task, status string, progress float64)
	OnTaskCompleted(task task.Task)
	OnTaskError(task task.Task, err error)
}

// 日志观察者接口
type LogObserver interface {
	OnLogMessage(level, message string)
}

// 观察者列表
var taskObservers []TaskStatusObserver
var logObservers []LogObserver

// RegisterTaskObserver 注册任务观察者
func RegisterTaskObserver(observer TaskStatusObserver) {
	taskObservers = append(taskObservers, observer)
}

// RegisterLogObserver 注册日志观察者
func RegisterLogObserver(observer LogObserver) {
	logObservers = append(logObservers, observer)
}

// NotifyTaskStatus 通知任务状态变更
func NotifyTaskStatus(task task.Task, status string, progress float64) {
	for _, observer := range taskObservers {
		observer.OnTaskStatusChanged(task, status, progress)
	}
}

// NotifyTaskCompleted 通知任务完成
func NotifyTaskCompleted(task task.Task) {
	for _, observer := range taskObservers {
		observer.OnTaskCompleted(task)
	}
}

// NotifyTaskError 通知任务错误
func NotifyTaskError(task task.Task, err error) {
	for _, observer := range taskObservers {
		observer.OnTaskError(task, err)
	}
}

// Run 启动核心引擎
func Run() {
	InitLog()
	util.Copyright()
	err := InitConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	
	// 移除Web服务启动
	// go InitWeb()
	
	for _, user := range config.Conf.Users {
		send(user)
	}
	task.Start()
}

func send(user config.User) {
	instance := yinghua.New(user)
	err := instance.Login()
	if err != nil {
		logrus.Fatal(err)
	}
	instance.Output(fmt.Sprintf("登录成功"))
	err = instance.GetCourses()
	if err != nil {
		logrus.Fatal(err)
	}
	instance.Output(fmt.Sprintf("获取全部在学课程成功, 共计 %d 门\n", len(instance.Courses)))
	for _, course := range instance.Courses {
		task.Tasks = append(task.Tasks, task.Task{
			User:   user,
			Course: course,
			Status: false,
		})
	}
}
