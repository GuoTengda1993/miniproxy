/*
 * @Description: message
 */
package proxy

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (p *Proxy) MessagePush(title, message string, isSucc bool) {
	useBox := false
	if title != "" {
		useBox = true
	}
	tp := "success"
	if !isSucc {
		tp = "error"
	}
	msg := &MsgInfo{
		Type:    tp,
		Message: message,
		Title:   title,
		UseBox:  useBox,
	}
	runtime.EventsEmit(p.ctx, EventNameMessage, msg)
}
