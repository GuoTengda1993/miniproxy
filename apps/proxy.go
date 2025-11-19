package apps

import (
	"miniproxy/proxy"
	"miniproxy/utils"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ProxyInfo() map[string]any {
	return a.proxy.Info()
}

func (a *App) ProxyStart() {
	a.proxy.Run()
}

func (a *App) ProxyStop() {
	a.proxy.Stop()
}

func (a *App) ProxyShowCert() string {
	return a.proxy.ShowCert()
}

func (a *App) ProxyCloseCert() {
	a.proxy.CloseCert()
}

func (a *App) ProxySaveCert() string {
	saveFileOption := runtime.OpenDialogOptions{
		Title:                      "Select Folder",
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		TreatPackagesAsDirectories: false,
	}

	saveFilepath, err := runtime.OpenDirectoryDialog(a.ctx, saveFileOption)
	if err != nil || saveFilepath == "" {
		return "Folder not selected"
	}
	err = a.proxy.SaveCertLocal(saveFilepath)
	return utils.StringfyError(err)
}

func (a *App) ProxyFlowListByPage(filter, remoteAddr, host string, page int, onlyError bool) map[string]any {
	now := time.Now().UnixMilli()
	defer a.cache.Store("proxy-flow-query", now)
	if loadTime, ok := a.cache.Load("proxy-flow-query"); ok {
		if now-loadTime.(int64) <= 100 {
			return map[string]any{
				"list":  nil,
				"total": -1,
			}
		}
	}

	list, total := a.proxy.FlowListByPage(filter, remoteAddr, host, page, onlyError)
	res := map[string]any{
		"list":  list,
		"total": total,
	}
	return res
}

func (a *App) ProxyFlowClear() {
	a.proxy.FlowClear()
}

func (a *App) ProxyHostList(isAllowed bool) []proxy.Option {
	return a.proxy.HostList(isAllowed)
}

func (a *App) ProxyAddHost(h string) {
	a.proxy.AddHost(h)
}

func (a *App) ProxyRemoveHost(h string) {
	a.proxy.RemoveHost(h)
}

func (a *App) ProxyResetHost() {
	a.proxy.ResetHost()
}

func (a *App) ProxyRemoteAddrs() []string {
	return a.proxy.RemoteAddrs()
}

func (a *App) ProxyCatchHosts() []string {
	return a.proxy.CatchHosts()
}
