package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/aoaostar/mooc/pkg/task"
	"github.com/aoaostar/mooc/bootstrap"
	"fmt"
	"time"
	"sync"
)

// 进程监控视图
type ProcessMonitoringView struct {
	*walk.TabPage
	
	// 任务列表
	taskListView    *walk.TableView
	taskModel       *TaskListModel
	
	// 任务详情
	detailsGroup    *walk.GroupBox
	courseNameLabel *walk.Label
	userNameLabel   *walk.Label
	progressBar     *walk.ProgressBar
	statusLabel     *walk.Label
	chapterLabel    *walk.Label
	lessonLabel     *walk.Label
	
	// 日志显示
	logView         *walk.TextEdit
	
	// 控制按钮
	startAllButton  *walk.PushButton
	pauseAllButton  *walk.PushButton
	stopAllButton   *walk.PushButton
	globalStatusLabel *walk.Label
	
	// 数据
	mu              sync.Mutex
}

// 创建进程监控页面
func NewProcessMonitoringPage(parent walk.Container) (*ProcessMonitoringView, error) {
	view := new(ProcessMonitoringView)
	
	// 创建选项卡页面
	if err := TabPage{
		AssignTo: &view.TabPage,
		Title:    "进程监控",
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			// 任务列表部分
			GroupBox{
				Title:  "任务列表",
				Layout: VBox{},
				Children: []Widget{
					TableView{
						AssignTo:         &view.taskListView,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						MultiSelection:   true,
						AlternatingRowBG: true,
						Columns: []TableViewColumn{
							{Title: "课程名称", Width: 150},
							{Title: "用户", Width: 80},
							{Title: "进度", Width: 80},
							{Title: "状态", Width: 80},
							{Title: "开始时间", Width: 120},
						},
						OnCurrentIndexChanged: view.onTaskSelected,
						StyleCell: func(style *walk.CellStyle) {
							if style.Row()%2 == 0 {
								style.BackgroundColor = walk.RGB(248, 248, 248)
							}
							
							// 根据状态设置颜色
							task := view.taskModel.tasks[style.Row()]
							if style.Column() == 3 { // 状态列
								switch task.Status {
								case "已完成":
									style.TextColor = walk.RGB(0, 128, 0) // 绿色
								case "错误":
									style.TextColor = walk.RGB(255, 0, 0) // 红色
								case "进行中":
									style.TextColor = walk.RGB(0, 0, 255) // 蓝色
								}
							}
						},
					},
				},
			},
			
			// 任务详情部分
			GroupBox{
				Title:  "任务详情",
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "课程:"},
					Label{AssignTo: &view.courseNameLabel},
					
					Label{Text: "用户:"},
					Label{AssignTo: &view.userNameLabel},
					
					Label{Text: "进度:"},
					ProgressBar{AssignTo: &view.progressBar, MaxValue: 100},
					
					Label{Text: "状态:"},
					Label{AssignTo: &view.statusLabel},
					
					Label{Text: "当前章节:"},
					Label{AssignTo: &view.chapterLabel},
					
					Label{Text: "当前课时:"},
					Label{AssignTo: &view.lessonLabel},
				},
			},
			
			// 日志输出部分
			GroupBox{
				Title:  "日志输出",
				Layout: VBox{},
				Children: []Widget{
					TextEdit{
						AssignTo:    &view.logView,
						ReadOnly:    true,
						VScroll:     true,
						MaxSize:     Size{Width: 0, Height: 200},
					},
				},
			},
			
			// 控制按钮部分
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &view.startAllButton,
						Text:     "开始全部",
						OnClicked: view.onStartAll,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.pauseAllButton,
						Text:     "暂停全部",
						OnClicked: view.onPauseAll,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.stopAllButton,
						Text:     "停止全部",
						OnClicked: view.onStopAll,
						MinSize: Size{Width: 100, Height: 30},
					},
					HSpacer{},
					Label{Text: "状态:"},
					Label{AssignTo: &view.globalStatusLabel, Text: "就绪"},
				},
			},
		},
	}.Create(NewBuilder(parent)); err != nil {
		return nil, err
	}
	
	// 初始化任务列表模型
	view.taskModel = NewTaskListModel()
	view.taskListView.SetModel(view.taskModel)
	
	// 注册为任务观察者
	bootstrap.RegisterTaskObserver(view)
	
	// 注册为日志观察者
	bootstrap.RegisterLogObserver(view)
	
	return view, nil
}

// 任务选中事件处理
func (v *ProcessMonitoringView) onTaskSelected() {
	// 获取选中的任务
	index := v.taskListView.CurrentIndex()
	if index < 0 {
		return
	}
	
	task := v.taskModel.tasks[index]
	
	// 更新详情视图
	v.courseNameLabel.SetText(task.CourseName)
	v.userNameLabel.SetText(task.UserName)
	v.progressBar.SetValue(int(task.Progress * 100))
	v.statusLabel.SetText(task.Status)
	v.chapterLabel.SetText(task.CurrentChapter)
	v.lessonLabel.SetText(task.CurrentLesson)
}

// 开始全部任务
func (v *ProcessMonitoringView) onStartAll() {
	// 调用核心逻辑启动所有任务
	v.globalStatusLabel.SetText("运行中")
	// TODO: 实现任务启动逻辑
}

// 暂停全部任务
func (v *ProcessMonitoringView) onPauseAll() {
	// 调用核心逻辑暂停所有任务
	v.globalStatusLabel.SetText("已暂停")
	// TODO: 实现任务暂停逻辑
}

// 停止全部任务
func (v *ProcessMonitoringView) onStopAll() {
	// 调用核心逻辑停止所有任务
	v.globalStatusLabel.SetText("已停止")
	// TODO: 实现任务停止逻辑
}

// 实现TaskStatusObserver接口
func (v *ProcessMonitoringView) OnTaskStatusChanged(task task.Task, status string, progress float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	// 在UI线程中执行
	walk.MustDo(func() {
		// 查找任务是否已存在
		found := false
		for i, t := range v.taskModel.tasks {
			if t.ID == task.Course.ID {
				// 更新现有任务
				v.taskModel.tasks[i].Status = status
				v.taskModel.tasks[i].Progress = progress
				v.taskModel.PublishRowChanged(i)
				found = true
				break
			}
		}
		
		// 如果任务不存在，添加新任务
		if !found {
			v.taskModel.AddTask(TaskItem{
				ID:         task.Course.ID,
				CourseName: task.Course.Name,
				UserName:   task.User.Username,
				Progress:   progress,
				Status:     status,
				StartTime:  time.Now(),
				CurrentChapter: "",
				CurrentLesson: "",
			})
		}
		
		// 如果当前选中的是这个任务，更新详情视图
		if v.taskListView.CurrentIndex() >= 0 {
			currentTask := v.taskModel.tasks[v.taskListView.CurrentIndex()]
			if currentTask.ID == task.Course.ID {
				v.progressBar.SetValue(int(progress * 100))
				v.statusLabel.SetText(status)
			}
		}
	})
}

func (v *ProcessMonitoringView) OnTaskCompleted(task task.Task) {
	v.OnTaskStatusChanged(task, "已完成", 1.0)
}

func (v *ProcessMonitoringView) OnTaskError(task task.Task, err error) {
	v.OnTaskStatusChanged(task, "错误", task.Course.Progress)
	v.AppendLog("error", err.Error())
}

// 实现LogObserver接口
func (v *ProcessMonitoringView) OnLogMessage(level, message string) {
	v.AppendLog(level, message)
}

// 添加日志
func (v *ProcessMonitoringView) AppendLog(level, message string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	// 在UI线程中执行
	walk.MustDo(func() {
		// 根据日志级别设置颜色
		var textColor walk.Color
		switch level {
		case "error", "fatal", "panic":
			textColor = walk.RGB(255, 0, 0) // 红色
		case "warn", "warning":
			textColor = walk.RGB(255, 165, 0) // 橙色
		case "info":
			textColor = walk.RGB(0, 0, 0) // 黑色
		default:
			textColor = walk.RGB(128, 128, 128) // 灰色
		}
		
		// 格式化日志时间
		timestamp := time.Now().Format("15:04:05")
		logText := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
		
		// 添加到日志视图
		v.logView.SetTextColor(textColor)
		v.logView.AppendText(logText + "\r\n")
		
		// 滚动到底部
		v.logView.SendMessage(win.EM_SCROLLCARET, 0, 0)
		
		// 限制日志行数
		if v.logView.LineCount() > 1000 {
			// 删除前100行
			v.logView.SetText(v.logView.Text()[v.logView.LineLength(0)*100:])
		}
	})
}
