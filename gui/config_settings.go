package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/aoaostar/mooc/pkg/config"
	"github.com/sirupsen/logrus"
	"sync"
)

// 配置设置视图
type ConfigSettingsView struct {
	*walk.TabPage
	
	// 全局设置
	limitEdit       *walk.NumberEdit
	serverEdit      *walk.LineEdit
	
	// 按钮
	saveButton      *walk.PushButton
	resetButton     *walk.PushButton
	importButton    *walk.PushButton
	exportButton    *walk.PushButton
	
	// 数据
	configManager   *ConfigManager
	
	mu              sync.Mutex
}

// 创建配置设置页面
func NewConfigSettingsPage(parent walk.Container, configManager *ConfigManager) (*ConfigSettingsView, error) {
	view := new(ConfigSettingsView)
	view.configManager = configManager
	
	// 创建选项卡页面
	if err := TabPage{
		AssignTo: &view.TabPage,
		Title:    "配置设置",
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			// 全局设置
			GroupBox{
				Title:  "全局设置",
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "并发任务数限制:"},
					NumberEdit{
						AssignTo: &view.limitEdit,
						MinValue: 1,
						MaxValue: 999999,
						Value:    3,
						Decimals: 0,
					},
					
					Label{Text: "服务器地址:"},
					LineEdit{
						AssignTo: &view.serverEdit,
						Text:     ":10086",
					},
				},
			},
			
			// 按钮区域
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &view.saveButton,
						Text:     "保存设置",
						OnClicked: view.onSave,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.resetButton,
						Text:     "重置为默认",
						OnClicked: view.onReset,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.importButton,
						Text:     "导入配置",
						OnClicked: view.onImport,
						MinSize: Size{Width: 100, Height: 30},
					},
					PushButton{
						AssignTo: &view.exportButton,
						Text:     "导出配置",
						OnClicked: view.onExport,
						MinSize: Size{Width: 100, Height: 30},
					},
				},
			},
		},
	}.Create(NewBuilder(parent)); err != nil {
		return nil, err
	}
	
	// 加载当前配置
	view.loadConfig()
	
	return view, nil
}

// 加载配置
func (v *ConfigSettingsView) loadConfig() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	conf := v.configManager.GetConfig()
	
	// 设置全局配置
	walk.MustDo(func() {
		v.limitEdit.SetValue(float64(conf.Global.Limit))
		v.serverEdit.SetText(conf.Global.Server)
	})
}

// 保存设置
func (v *ConfigSettingsView) onSave() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	conf := v.configManager.GetConfig()
	
	// 更新全局配置
	conf.Global.Limit = int(v.limitEdit.Value())
	conf.Global.Server = v.serverEdit.Text()
	
	// 保存配置
	err := v.configManager.SaveConfig(conf)
	if err != nil {
		walk.MsgBox(v.Form(), "错误", "保存配置失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	
	walk.MsgBox(v.Form(), "成功", "配置已保存", walk.MsgBoxIconInformation)
}

// 重置为默认
func (v *ConfigSettingsView) onReset() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	// 确认对话框
	if walk.MsgBox(v.Form(), "确认重置", "确定要将所有设置重置为默认值吗？", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
		return
	}
	
	// 设置默认值
	walk.MustDo(func() {
		v.limitEdit.SetValue(3)
		v.serverEdit.SetText(":10086")
	})
}

// 导入配置
func (v *ConfigSettingsView) onImport() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	dlg := new(walk.FileDialog)
	dlg.Title = "导入配置文件"
	dlg.Filter = "JSON文件 (*.json)|*.json|所有文件 (*.*)|*.*"
	
	if ok, err := dlg.ShowOpen(v.Form()); err != nil {
		walk.MsgBox(v.Form(), "错误", "打开文件对话框失败: "+err.Error(), walk.MsgBoxIconError)
		return
	} else if !ok {
		return
	}
	
	// 导入配置
	err := v.configManager.ImportConfig(dlg.FilePath)
	if err != nil {
		walk.MsgBox(v.Form(), "错误", "导入配置失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	
	// 重新加载配置
	v.loadConfig()
	
	walk.MsgBox(v.Form(), "成功", "配置已导入", walk.MsgBoxIconInformation)
}

// 导出配置
func (v *ConfigSettingsView) onExport() {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	dlg := new(walk.FileDialog)
	dlg.Title = "导出配置文件"
	dlg.Filter = "JSON文件 (*.json)|*.json|所有文件 (*.*)|*.*"
	
	if ok, err := dlg.ShowSave(v.Form()); err != nil {
		walk.MsgBox(v.Form(), "错误", "打开文件对话框失败: "+err.Error(), walk.MsgBoxIconError)
		return
	} else if !ok {
		return
	}
	
	// 导出配置
	err := v.configManager.ExportConfig(dlg.FilePath)
	if err != nil {
		walk.MsgBox(v.Form(), "错误", "导出配置失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	
	walk.MsgBox(v.Form(), "成功", "配置已导出到: "+dlg.FilePath, walk.MsgBoxIconInformation)
}
