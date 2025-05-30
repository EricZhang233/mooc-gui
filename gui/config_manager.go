package gui

import (
	"github.com/aoaostar/mooc/pkg/config"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// 配置管理器
type ConfigManager struct {
	configPath string
	config     config.Config
	mu         sync.RWMutex
}

// 创建配置管理器
func NewConfigManager(configPath string) *ConfigManager {
	manager := &ConfigManager{
		configPath: configPath,
	}
	
	// 加载配置
	manager.LoadConfig()
	
	return manager
}

// 加载配置
func (cm *ConfigManager) LoadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// 检查配置文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// 创建默认配置
		cm.config = config.Config{
			Global: config.Global{
				Server: ":10086",
				Limit:  3,
			},
			Users: []config.User{},
		}
		
		// 保存默认配置
		return cm.saveConfigInternal()
	}
	
	// 读取配置文件
	data, err := ioutil.ReadFile(cm.configPath)
	if err != nil {
		return err
	}
	
	// 解析JSON
	err = json.Unmarshal(data, &cm.config)
	if err != nil {
		return err
	}
	
	return nil
}

// 获取配置
func (cm *ConfigManager) GetConfig() config.Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	return cm.config
}

// 保存配置
func (cm *ConfigManager) SaveConfig(conf config.Config) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.config = conf
	return cm.saveConfigInternal()
}

// 内部保存配置方法
func (cm *ConfigManager) saveConfigInternal() error {
	// 序列化为JSON
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}
	
	// 写入文件
	return ioutil.WriteFile(cm.configPath, data, 0644)
}

// 导入配置
func (cm *ConfigManager) ImportConfig(path string) error {
	// 读取文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	
	// 解析JSON
	var conf config.Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}
	
	// 保存配置
	return cm.SaveConfig(conf)
}

// 导出配置
func (cm *ConfigManager) ExportConfig(path string) error {
	// 获取当前配置
	conf := cm.GetConfig()
	
	// 序列化为JSON
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	
	// 写入文件
	return ioutil.WriteFile(path, data, 0644)
}
