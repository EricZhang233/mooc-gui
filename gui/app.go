package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/aoaostar/mooc/pkg/task"
	"github.com/aoaostar/mooc/bootstrap"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
	"fmt"
)

// 应用程序类
type App struct {
	MainWindow          *walk.MainWindow
	TabWidget           *walk.TabWidget
	ProcessMonitoringView *ProcessMonitoringView
	UserManagementView  *UserManagementView
	ConfigSettingsView  *ConfigSettingsView
	ConfigManager       *ConfigManager
	
	mu                  sync.Mutex
}

// 创建并运行应用程序
func RunApp() error {
	// 初始化配置管理器
	configPath := "./config.json"
	configManager := NewConfigManager(configPath)
	
	// 创建应用实例
	app := &App{
		ConfigManager: configManager,
	}
	
	// 创建并运行主窗口
	var mainWindow *walk.MainWindow
	var tabWidget *walk.TabWidget
	
	if _, err := (MainWindow{
		AssignTo: &mainWindow,
		Title:    "英华学堂网课助手 - GUI版",
		MinSize:  Size{Width: 900, Height: 700},
		Layout:   VBox{},
		Children: []Widget{
			TabWidget{
				AssignTo: &tabWidget,
			},
		},
		OnSizeChanged: func() {
			// 窗口大小变化时调整布局
			tabWidget.SetSize(walk.Size{
				Width:  mainWindow.ClientBounds().Width,
				Height: mainWindow.ClientBounds().Height,
			})
		},
	}).Create(); err != nil {
		logrus.Fatal(err)
		return err
	}
	
	app.MainWindow = mainWindow
	app.TabWidget = tabWidget
	
	// 创建选项卡页面
	processMonitoringPage, err := NewProcessMonitoringPage(tabWidget)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	app.ProcessMonitoringView = processMonitoringPage
	
	userManagementPage, err := NewUserManagementPage(tabWidget, configManager)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	app.UserManagementView = userManagementPage
	
	configSettingsPage, err := NewConfigSettingsPage(tabWidget, configManager)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	app.ConfigSettingsView = configSettingsPage
	
	// 设置图标
	icon, err := walk.Resources.Icon("ICON")
	if err == nil {
		mainWindow.SetIcon(icon)
	}
	
	// 注册为任务观察者
	bootstrap.RegisterTaskObserver(app)
	
	// 注册为日志观察者
	bootstrap.RegisterLogObserver(app)
	
	// 初始化核心引擎
	go func() {
		// 启动核心引擎，但不启动Web服务
		bootstrap.Run()
	}()
	
	// 窗口关闭事件处理
	mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		// 确认对话框
		if walk.MsgBox(mainWindow, "确认退出", "确定要退出应用程序吗？", 
			walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
			*canceled = true
			return
		}
		
		// 执行清理操作
		// ...
	})
	
	// 运行主窗口
	mainWindow.Run()
	
	return nil
}

// 实现TaskStatusObserver接口
func (app *App) OnTaskStatusChanged(task task.Task, status string, progress float64) {
	app.mu.Lock()
	defer app.mu.Unlock()
	
	// 确保在UI线程中执行
	walk.MustDo(func() {
		if app.ProcessMonitoringView != nil {
			app.ProcessMonitoringView.OnTaskStatusChanged(task, status, progress)
		}
	})
}

func (app *App) OnTaskCompleted(task task.Task) {
	app.mu.Lock()
	defer app.mu.Unlock()
	
	// 确保在UI线程中执行
	walk.MustDo(func() {
		if app.ProcessMonitoringView != nil {
			app.ProcessMonitoringView.OnTaskCompleted(task)
		}
	})
}

func (app *App) OnTaskError(task task.Task, err error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	
	// 确保在UI线程中执行
	walk.MustDo(func() {
		if app.ProcessMonitoringView != nil {
			app.ProcessMonitoringView.OnTaskError(task, err)
		}
	})
}

// 实现LogObserver接口
func (app *App) OnLogMessage(level, message string) {
	app.mu.Lock()
	defer app.mu.Unlock()
	
	// 确保在UI线程中执行
	walk.MustDo(func() {
		if app.ProcessMonitoringView != nil {
			app.ProcessMonitoringView.OnLogMessage(level, message)
		}
	})
}
