/*
app
*/
package apps

import (
	"context"
	"fmt"
	"sync"

	"miniproxy/internal"
	"miniproxy/proxy"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	CacheKeySavePath = "save-path"
)

// App struct
type App struct {
	ctx   context.Context
	cache *sync.Map

	proxy *proxy.Proxy
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	internal.LocalConf.LoadConfig()
	a.ctx = ctx
	a.cache = &sync.Map{}
	a.proxy = proxy.NewProxy(ctx, internal.LocalConf.ProxyPort)
	a.proxy.InitFromLocalConf()
}

func (a *App) Shutdown(ctx context.Context) {
	a.proxy.SaveToLocalConf()
	internal.LocalConf.SaveConfig()

	a.proxy.Stop()
	runtime.EventsOffAll(a.ctx)
}

func (a *App) LogPrint(data any) {
	fmt.Printf("[P] >>> >>> >>> %+v\n", data)
}

func (a *App) WindowGetHeight() (height int) {
	_, height = runtime.WindowGetSize(a.ctx)
	return
}

func (a *App) WindowGetWidth() (width int) {
	width, _ = runtime.WindowGetSize(a.ctx)
	return
}

// front inactive
func (a *App) FrontPause() {
	a.proxy.FrontPause()
}

// front active
func (a *App) FrontResume() {
	a.proxy.FrontResume()
}
