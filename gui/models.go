package gui

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"sync"
)

// 任务列表项
type TaskItem struct {
	ID            int
	CourseName    string
	UserName      string
	Progress      float64
	Status        string
	StartTime     time.Time
	CurrentChapter string
	CurrentLesson  string
}

// 任务列表模型
type TaskListModel struct {
	walk.TableModelBase
	tasks []TaskItem
	mu    sync.Mutex
}

// 创建任务列表模型
func NewTaskListModel() *TaskListModel {
	return &TaskListModel{
		tasks: make([]TaskItem, 0),
	}
}

// 实现walk.TableModel接口
func (m *TaskListModel) RowCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.tasks)
}

func (m *TaskListModel) Value(row, col int) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if row < 0 || row >= len(m.tasks) {
		return nil
	}
	
	task := m.tasks[row]
	
	switch col {
	case 0:
		return task.CourseName
	case 1:
		return task.UserName
	case 2:
		return fmt.Sprintf("%.0f%%", task.Progress*100)
	case 3:
		return task.Status
	case 4:
		return task.StartTime.Format("15:04:05")
	}
	
	return nil
}

// 添加任务
func (m *TaskListModel) AddTask(task TaskItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.tasks = append(m.tasks, task)
	m.PublishRowsInserted(len(m.tasks)-1, len(m.tasks)-1)
}

// 清空任务
func (m *TaskListModel) ClearTasks() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.tasks = make([]TaskItem, 0)
	m.PublishRowsReset()
}

// 用户列表项
type UserItem struct {
	Username  string
	Name      string
	Platform  string
	BaseUrl   string
	CourseNum int
	Status    string
	User      config.User
}

// 用户列表模型
type UserListModel struct {
	walk.TableModelBase
	users []UserItem
	mu    sync.Mutex
}

// 创建用户列表模型
func NewUserListModel() *UserListModel {
	return &UserListModel{
		users: make([]UserItem, 0),
	}
}

// 实现walk.TableModel接口
func (m *UserListModel) RowCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.users)
}

func (m *UserListModel) Value(row, col int) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if row < 0 || row >= len(m.users) {
		return nil
	}
	
	user := m.users[row]
	
	switch col {
	case 0:
		return user.Username
	case 1:
		return user.Name
	case 2:
		return user.Platform
	case 3:
		return user.BaseUrl
	case 4:
		return user.CourseNum
	case 5:
		return user.Status
	}
	
	return nil
}

// 添加用户
func (m *UserListModel) AddUser(user UserItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.users = append(m.users, user)
	m.PublishRowsInserted(len(m.users)-1, len(m.users)-1)
}

// 清空用户
func (m *UserListModel) ClearUsers() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.users = make([]UserItem, 0)
	m.PublishRowsReset()
}

// 用户课程列表项
type UserCourseItem struct {
	Name     string
	Progress float64
	Status   string
	Course   yinghua.Course
}

// 用户课程列表模型
type UserCourseModel struct {
	walk.TableModelBase
	courses []UserCourseItem
	mu      sync.Mutex
}

// 创建用户课程列表模型
func NewUserCourseModel() *UserCourseModel {
	return &UserCourseModel{
		courses: make([]UserCourseItem, 0),
	}
}

// 实现walk.TableModel接口
func (m *UserCourseModel) RowCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.courses)
}

func (m *UserCourseModel) Value(row, col int) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if row < 0 || row >= len(m.courses) {
		return nil
	}
	
	course := m.courses[row]
	
	switch col {
	case 0:
		return course.Name
	case 1:
		return fmt.Sprintf("%.0f%%", course.Progress*100)
	case 2:
		return course.Status
	}
	
	return nil
}

// 添加课程
func (m *UserCourseModel) AddCourse(course UserCourseItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.courses = append(m.courses, course)
	m.PublishRowsInserted(len(m.courses)-1, len(m.courses)-1)
}

// 清空课程
func (m *UserCourseModel) ClearCourses() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.courses = make([]UserCourseItem, 0)
	m.PublishRowsReset()
}

// 创建构建器
func NewBuilder(parent walk.Container) walk.Builder {
	return walk.NewBuilder(parent)
}
