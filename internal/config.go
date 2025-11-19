/*
 * @Description: local config
 */
package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// local config
type localConfig struct {
	appName   string
	confDir   string
	fileName  string
	ProxyPort int      `json:"proxy_port"` // port, default: 8890
	ProxyHost []string `json:"proxy_host"`
}

const DefaultProxyPort = 8890

var LocalConf *localConfig = &localConfig{appName: "MiCarAppTestTools", fileName: "config.json"}

func (c *localConfig) getConfigDir() error {
	var baseDir string
	switch runtime.GOOS {
	case "darwin":
		// macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	case "windows":
		// Windows
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			return os.ErrNotExist
		}
	default:
		// Linux
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = filepath.Join(homeDir, ".config")
	}
	configDir := filepath.Join(baseDir, c.appName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	c.confDir = configDir
	return nil
}

func (c *localConfig) LoadConfig() error {
	defer func() {
		if LocalConf.ProxyPort <= 0 {
			LocalConf.ProxyPort = DefaultProxyPort
		}
	}()
	err := c.getConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(c.confDir, c.fileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, LocalConf)
}

func (c *localConfig) SaveConfig() error {
	configPath := filepath.Join(c.confDir, c.fileName)
	data, err := json.MarshalIndent(LocalConf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func (c *localConfig) SetProxyHost(hosts []string) {
	c.ProxyHost = hosts
}
