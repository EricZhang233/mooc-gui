package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/aoaostar/mooc/pkg/config"
	"github.com/aoaostar/mooc/pkg/yinghua"
	"github.com/aoaostar/mooc/pkg/task"
	"github.com/sirupsen/logrus"
	"strconv"
	"fmt"
	"errors"
	"sync"
)

// 用户管理视图
type UserManagementView struct {
	*walk.TabPage
	
	// 用户列表
	userListView    *walk.TableView
	userModel       *UserListModel
	
	// 用户操作按钮
	addButton       *walk.PushButton
	editButton      *walk.PushButton
	deleteButton    *walk.PushButton
	refreshButton   *walk.PushButton
	
	// 用户详情表单
	detailsGroup    *walk.GroupBox
	usernameEdit    *walk.LineEdit
	passwordEdit    *walk.LineEdit
	baseUrlEdit     *walk.LineEdit
	schoolIdEdit    *walk.NumberEdit
	nameEdit        *walk.LineEdit      // 新增：姓名
	platformEdit    *walk.LineEdit      // 新增：平台
	remarkEdit      *walk.TextEdit      // 新增：备注
	
	// 用户课程列表
	courseListView  *walk.TableView
	courseModel     *UserCourseModel
	
	// 课程操作按钮
	getCourseButton *walk.PushButton
	startButton     *walk.PushButton
	viewButton      *walk.PushButton
	
	// 保存取消按钮
	saveButton      *walk.PushButton
	cancelButton    *walk.PushButton
	
	// 数据
	currentUser     *config.User
	configManager   *ConfigManager
	
	mu              sync.Mutex
}

// 创建用户管理页面
func NewUserManagementPage(parent walk.Container, configManager *ConfigManager) (*UserManagementView, error) {
	view := new(UserManagementView)
	view.configManager = configManager
	
	// 创建选项卡页面
	if err := TabPage{
		AssignTo: &view.TabPage,
		Title:    "用户管理",
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			// 用户列表部分
			GroupBox{
				Title:  "用户列表",
				Layout: VBox{},
				Children: []Widget{
					TableView{
						AssignTo:         &view.userListView,
						CheckBoxes:       false,
						ColumnsOrderable: true,
						MultiSelection:   false,
						AlternatingRowBG: true,
						Columns: []TableViewColumn{
							{Title: "用户名", Width: 120},
							{Title: "姓名", Width: 100},
							{Title: "平台", Width: 120},
							{Title: "学校平台", Width: 200},
							{Title: "课程数", Width: 80},
							{Title: "状态", Width: 80},
						},
						OnCurrentIndexChanged: view.onUserSelected,
						StyleCell: func(style *walk.CellStyle) {
							if style.Row()%2 == 0 {
								style.BackgroundColor = walk.RGB(248, 248, 248)
							}
							
							// 根据状态设置颜色
							if style.Column() == 5 { // 状态列
								user := view.userModel.users[style.Row()]
								if user.Status == "在线" {
									style.TextColor = walk.RGB(0, 128, 0) // 绿色
								} else {
									style.TextColor = walk.RGB(128, 128, 128) // 灰色
								}
							}
						},
					},
				},
			},
			
			// 用户操作按钮
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &view.addButton,
						Text:     "添加用户",
						OnClicked: view.onAddUser,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.editButton,
						Text:     "编辑用户",
						OnClicked: view.onEditUser,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.deleteButton,
						Text:     "删除用户",
						OnClicked: view.onDeleteUser,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.refreshButton,
						Text:     "刷新",
						OnClicked: view.onRefresh,
						MinSize: Size{Width: 100, Height: 30},
					},
				},
			},
			
			// 用户详情部分
			GroupBox{
				AssignTo: &view.detailsGroup,
				Title:    "用户详情",
				Layout:   VBox{},
				Children: []Widget{
					// 基本信息
					GroupBox{
						Title:  "基本信息",
						Layout: Grid{Columns: 4},
						Children: []Widget{
							Label{Text: "用户名:"},
							LineEdit{
								AssignTo: &view.usernameEdit,
								MinSize:  Size{Width: 150},
							},
							
							Label{Text: "密码:"},
							LineEdit{
								AssignTo: &view.passwordEdit,
								MinSize:  Size{Width: 150},
								PasswordMode: true,
							},
							
							Label{Text: "姓名:"},
							LineEdit{
								AssignTo: &view.nameEdit,
								MinSize:  Size{Width: 150},
							},
							
							Label{Text: "平台:"},
							LineEdit{
								AssignTo: &view.platformEdit,
								MinSize:  Size{Width: 150},
							},
							
							Label{Text: "学校平台URL:"},
							LineEdit{
								AssignTo: &view.baseUrlEdit,
								MinSize:  Size{Width: 300},
								ColumnSpan: 3,
							},
							
							Label{Text: "学校ID:"},
							NumberEdit{
								AssignTo: &view.schoolIdEdit,
								MinSize:  Size{Width: 80},
								Value:    0,
								Decimals: 0,
							},
							
							Label{Text: "备注:"},
							TextEdit{
								AssignTo: &view.remarkEdit,
								MinSize:  Size{Width: 300, Height: 60},
								ColumnSpan: 3,
							},
						},
					},
					
					// 用户课程
					GroupBox{
						Title:  "用户课程",
						Layout: VBox{},
						Children: []Widget{
							TableView{
								AssignTo:         &view.courseListView,
								CheckBoxes:       false,
								ColumnsOrderable: true,
								MultiSelection:   true,
								AlternatingRowBG: true,
								Columns: []TableViewColumn{
									{Title: "课程名称", Width: 200},
									{Title: "进度", Width: 80},
									{Title: "状态", Width: 100},
								},
								StyleCell: func(style *walk.CellStyle) {
									if style.Row()%2 == 0 {
										style.BackgroundColor = walk.RGB(248, 248, 248)
									}
									
									// 根据状态设置颜色
									if style.Column() == 2 { // 状态列
										course := view.courseModel.courses[style.Row()]
										switch course.Status {
										case "已完成":
											style.TextColor = walk.RGB(0, 128, 0) // 绿色
										case "已结束":
											style.TextColor = walk.RGB(128, 128, 128) // 灰色
										case "进行中":
											style.TextColor = walk.RGB(0, 0, 255) // 蓝色
										}
									}
								},
							},
							
							// 课程操作按钮
							Composite{
								Layout: HBox{},
								Children: []Widget{
									PushButton{
										AssignTo: &view.getCourseButton,
										Text:     "获取课程",
										OnClicked: view.onGetCourses,
										MinSize: Size{Width: 100, Height: 30},
									},
									PushButton{
										AssignTo: &view.startButton,
										Text:     "开始学习",
										OnClicked: view.onStartCourse,
										MinSize: Size{Width: 100, Height: 30},
									},
									PushButton{
										AssignTo: &view.viewButton,
										Text:     "查看详情",
										OnClicked: view.onViewCourse,
										MinSize: Size{Width: 100, Height: 30},
									},
								},
							},
						},
					},
				},
			},
			
			// 保存取消按钮
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &view.saveButton,
						Text:     "保存",
						OnClicked: view.onSave,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.cancelButton,
						Text:     "取消",
						OnClicked: view.onCancel,
						MinSize: Size{Width: 100, Height: 30},
					},
				},
			},
		},
	}.Create(NewBuilder(parent)); err != nil {
		return nil, err
	}
	
	// 初始化用户列表模型
	view.userModel = NewUserListModel()
	view.userListView.SetModel(view.userModel)
	
	// 初始化课程列表模型
	view.courseModel = NewUserCourseModel()
	view.courseListView.SetModel(view.courseModel)
	
	// 初始化界面状态
	view.setEditMode(false)
	view.loadUsers()
	
	return view, nil
}

// 设置编辑模式
func (v *UserManagementView) setEditMode(editing bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	walk.MustDo(func() {
		v.usernameEdit.SetEnabled(editing)
		v.passwordEdit.SetEnabled(editing)
		v.baseUrlEdit.SetEnabled(editing)
		v.schoolIdEdit.SetEnabled(editing)
		v.nameEdit.SetEnabled(editing)
		v.platformEdit.SetEnabled(editing)
		v.remarkEdit.SetEnabled(editing)
		
		v.saveButton.SetEnabled(editing)
		v.cancelButton.SetEnabled(editing)
		
		v.addButton.SetEnabled(!editing)
		v.editButton.SetEnabled(!editing && v.userListView.CurrentIndex() >= 0)
		v.deleteButton.SetEnabled(!editing && v.userListView.CurrentIndex() >= 0)
		v.refreshButton.SetEnabled(!editing)
		
		v.getCourseButton.SetEnabled(!editing && v.userListView.CurrentIndex() >= 0)
		v.startButton.SetEnabled(!editing && v.courseListView.CurrentIndex() >= 0)
		v.viewButton.SetEnabled(!editing && v.courseListView.CurrentIndex() >= 0)
	})
	
	if !editing {
		v.clearForm()
	}
}

// 清空表单
func (v *UserManagementView) clearForm() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	walk.MustDo(func() {
		v.usernameEdit.SetText("")
		v.passwordEdit.SetText("")
		v.baseUrlEdit.SetText("")
		v.schoolIdEdit.SetValue(0)
		v.nameEdit.SetText("")
		v.platformEdit.SetText("")
		v.remarkEdit.SetText("")
		v.courseModel.ClearCourses()
	})
}

// 加载用户列表
func (v *UserManagementView) loadUsers() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	// 获取配置中的用户列表
	conf := v.configManager.GetConfig()
	
	// 清空现有用户
	v.userModel.ClearUsers()
	
	// 添加用户到列表
	walk.MustDo(func() {
		for _, user := range conf.Users {
			v.userModel.AddUser(UserItem{
				Username:  user.Username,
				Name:      user.Name,
				Platform:  user.Platform,
				BaseUrl:   user.BaseURL,
				CourseNum: 0, // 需要获取课程后更新
				Status:    "离线",
				User:      user,
			})
		}
	})
}

// 用户选中事件
func (v *UserManagementView) onUserSelected() {
	index := v.userListView.CurrentIndex()
	if index < 0 {
		v.clearForm()
		v.setEditMode(false)
		return
	}
	
	// 获取选中的用户
	user := v.userModel.users[index].User
	v.currentUser = &user
	
	// 更新表单
	v.mu.Lock()
	defer v.mu.Unlock()
	
	walk.MustDo(func() {
		v.usernameEdit.SetText(user.Username)
		v.passwordEdit.SetText(user.Password)
		v.baseUrlEdit.SetText(user.BaseURL)
		v.schoolIdEdit.SetValue(float64(user.SchoolID))
		v.nameEdit.SetText(user.Name)
		v.platformEdit.SetText(user.Platform)
		v.remarkEdit.SetText(user.Remark)
		
		// 更新按钮状态
		v.editButton.SetEnabled(true)
		v.deleteButton.SetEnabled(true)
		v.getCourseButton.SetEnabled(true)
		
		// 清空课程列表
		v.courseModel.ClearCourses()
	})
}

// 添加用户
func (v *UserManagementView) onAddUser() {
	v.currentUser = nil
	v.clearForm()
	v.setEditMode(true)
}

// 编辑用户
func (v *UserManagementView) onEditUser() {
	if v.userListView.CurrentIndex() < 0 {
		return
	}
	
	v.setEditMode(true)
}

// 删除用户
func (v *UserManagementView) onDeleteUser() {
	if v.userListView.CurrentIndex() < 0 {
		return
	}
	
	// 确认对话框
	if walk.MsgBox(v.Form(), "确认删除", "确定要删除选中的用户吗？", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
		return
	}
	
	// 删除用户
	index := v.userListView.CurrentIndex()
	conf := v.configManager.GetConfig()
	
	// 移除用户
	conf.Users = append(conf.Users[:index], conf.Users[index+1:]...)
	
	// 保存配置
	err := v.configManager.SaveConfig(conf)
	if err != nil {
		walk.MsgBox(v.Form(), "错误", "删除用户失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	
	// 刷新用户列表
	v.loadUsers()
}

// 刷新
func (v *UserManagementView) onRefresh() {
	v.loadUsers()
}

// 获取课程
func (v *UserManagementView) onGetCourses() {
	if v.currentUser == nil {
		return
	}
	
	// 显示进度对话框
	dlg, _ := walk.NewProgressDialog(v.Form())
	dlg.SetTitle("获取课程")
	dlg.SetValue(0)
	dlg.SetRange(0, 100)
	dlg.SetCancelable(true)
	
	// 在后台获取课程
	go func() {
		defer dlg.Close()
		
		// 创建英华学堂实例
		instance := yinghua.New(*v.currentUser)
		
		// 登录
		err := instance.Login()
		if err != nil {
			walk.MustDo(func() {
				walk.MsgBox(v.Form(), "错误", "登录失败: "+err.Error(), walk.MsgBoxIconError)
			})
			return
		}
		
		dlg.SetValue(50)
		
		// 获取课程
		err = instance.GetCourses()
		if err != nil {
			walk.MustDo(func() {
				walk.MsgBox(v.Form(), "错误", "获取课程失败: "+err.Error(), walk.MsgBoxIconError)
			})
			return
		}
		
		dlg.SetValue(100)
		
		// 更新课程列表
		v.mu.Lock()
		defer v.mu.Unlock()
		
		walk.MustDo(func() {
			v.courseModel.ClearCourses()
			
			for _, course := range instance.Courses {
				progress := 0.0
				if course.Progress1 != "" {
					progress, _ = strconv.ParseFloat(course.Progress1, 64)
				}
				
				status := "未开始"
				if course.Progress == 1 {
					status = "已完成"
				} else if course.State == 2 {
					status = "已结束"
				} else if progress > 0 {
					status = "进行中"
				}
				
				v.courseModel.AddCourse(UserCourseItem{
					Name:     course.Name,
					Progress: progress,
					Status:   status,
					Course:   course,
				})
			}
			
			// 更新用户状态
			index := v.userListView.CurrentIndex()
			if index >= 0 {
				v.userModel.users[index].CourseNum = len(instance.Courses)
				v.userModel.users[index].Status = "在线"
				v.userModel.PublishRowChanged(index)
			}
		})
	}()
}

// 开始学习
func (v *UserManagementView) onStartCourse() {
	if v.currentUser == nil || v.courseListView.CurrentIndex() < 0 {
		return
	}
	
	// 获取选中的课程
	course := v.courseModel.courses[v.courseListView.CurrentIndex()].Course
	
	// 确认对话框
	if walk.MsgBox(v.Form(), "确认", "确定要开始学习选中的课程吗？", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
		return
	}
	
	// 添加到任务队列
	task.Tasks = append(task.Tasks, task.Task{
		User:   *v.currentUser,
		Course: course,
		Status: false,
	})
	
	walk.MsgBox(v.Form(), "提示", "课程已添加到任务队列，请切换到进程监控页面查看进度。", walk.MsgBoxIconInformation)
}

// 查看课程详情
func (v *UserManagementView) onViewCourse() {
	if v.courseListView.CurrentIndex() < 0 {
		return
	}
	
	// 获取选中的课程
	course := v.courseModel.courses[v.courseListView.CurrentIndex()]
	
	// 显示课程详情对话框
	walk.MsgBox(v.Form(), "课程详情", 
		"课程名称: " + course.Name + "\n" +
		"进度: " + fmt.Sprintf("%.0f%%", course.Progress*100) + "\n" +
		"状态: " + course.Status + "\n" +
		"课程ID: " + strconv.Itoa(course.Course.ID) + "\n" +
		"开始时间: " + course.Course.StartTime + "\n" +
		"结束时间: " + course.Course.EndTime,
		walk.MsgBoxIconInformation)
}

// 保存
func (v *UserManagementView) onSave() {
	// 验证表单
	if v.usernameEdit.Text() == "" {
		walk.MsgBox(v.Form(), "错误", "用户名不能为空", walk.MsgBoxIconError)
		return
	}
	
	if v.passwordEdit.Text() == "" {
		walk.MsgBox(v.Form(), "错误", "密码不能为空", walk.MsgBoxIconError)
		return
	}
	
	if v.baseUrlEdit.Text() == "" {
		walk.MsgBox(v.Form(), "错误", "学校平台URL不能为空", walk.MsgBoxIconError)
		return
	}
	
	// 创建用户对象
	user := config.User{
		Username:  v.usernameEdit.Text(),
		Password:  v.passwordEdit.Text(),
		BaseURL:   v.baseUrlEdit.Text(),
		SchoolID:  int(v.schoolIdEdit.Value()),
		Name:      v.nameEdit.Text(),
		Platform:  v.platformEdit.Text(),
		Remark:    v.remarkEdit.Text(),
	}
	
	// 获取当前配置
	conf := v.configManager.GetConfig()
	
	// 添加或更新用户
	if v.currentUser == nil {
		// 添加新用户
		conf.Users = append(conf.Users, user)
	} else {
		// 更新现有用户
		index := v.userListView.CurrentIndex()
		if index >= 0 && index < len(conf.Users) {
			conf.Users[index] = user
		}
	}
	
	// 保存配置
	err := v.configManager.SaveConfig(conf)
	if err != nil {
		walk.MsgBox(v.Form(), "错误", "保存用户失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	
	// 刷新用户列表
	v.loadUsers()
	
	// 退出编辑模式
	v.setEditMode(false)
}

// 取消
func (v *UserManagementView) onCancel() {
	// 退出编辑模式
	v.setEditMode(false)
	
	// 恢复选中的用户
	if v.currentUser != nil {
		v.mu.Lock()
		defer v.mu.Unlock()
		
		walk.MustDo(func() {
			v.usernameEdit.SetText(v.currentUser.Username)
			v.passwordEdit.SetText(v.currentUser.Password)
			v.baseUrlEdit.SetText(v.currentUser.BaseURL)
			v.schoolIdEdit.SetValue(float64(v.currentUser.SchoolID))
			v.nameEdit.SetText(v.currentUser.Name)
			v.platformEdit.SetText(v.currentUser.Platform)
			v.remarkEdit.SetText(v.currentUser.Remark)
		})
	} else {
		v.clearForm()
	}
}
