/*
 * @Description: host
 */
package proxy

import (
	"miniproxy/internal"
	"miniproxy/utils"
	"sort"
	"strings"
)

func (p *Proxy) loadHostFromLocalConf() {
	if len(internal.LocalConf.ProxyHost) == 0 {
		utils.CopyMap2Sync(allowModifyHosts, p.allowHostMap)
		return
	}
	for _, h := range internal.LocalConf.ProxyHost {
		p.AddHost(h)
	}
}

// Save to the local configuration file before the service is closed.
func (p *Proxy) saveHostToLocalConf() {
	hosts := make([]string, 0)
	p.allowHostMap.Range(func(key, value any) bool {
		hosts = append(hosts, key.(string))
		return true
	})
	sort.Strings(hosts)
	reList := make([]string, 0)
	p.allowHostRe.Range(func(key, value any) bool {
		reList = append(reList, key.(string))
		return true
	})
	sort.Strings(reList)
	hosts = append(hosts, reList...)
	internal.LocalConf.SetProxyHost(hosts)
}

func (p *Proxy) HostList(isAllowed bool) []Option {
	if isAllowed {
		return p.AllowHosts()
	}
	return p.DisableHosts()
}

func (p *Proxy) AllowHosts() []Option {
	opts := make([]Option, 0)
	hosts := make([]string, 0)
	p.allowHostMap.Range(func(key, value any) bool {
		hosts = append(hosts, key.(string))
		return true
	})
	sort.Strings(hosts)
	for _, h := range hosts {
		opt := Option{Label: h, Value: h}
		opts = append(opts, opt)
	}
	reList := make([]string, 0)
	p.allowHostRe.Range(func(key, value any) bool {
		reList = append(reList, key.(string))
		return true
	})
	for _, k := range reList {
		k = "#" + k
		opt := Option{Label: k, Value: k}
		opts = append(opts, opt)
	}
	return opts
}

func (p *Proxy) DisableHosts() []Option {
	hosts := make([]string, 0)
	p.disableHostMap.Range(func(k, v any) bool {
		hosts = append(hosts, k.(string))
		return true
	})
	opts := make([]Option, 0)
	sort.Strings(hosts)
	for _, h := range hosts {
		opt := Option{Label: h, Value: h}
		opts = append(opts, opt)
	}
	return opts
}

func (p *Proxy) AddHost(h string) {
	h = strings.TrimSpace(h)
	isRe := false
	if strings.HasPrefix(h, "#") {
		isRe = true
		h = h[1:]
	}
	if isRe {
		p.allowHostRe.Store(h, true)
	} else {
		p.allowHostMap.Store(h, true)
	}
	p.disableHostMap.Delete(h)
}

func (p *Proxy) RemoveHost(h string) {
	h = strings.TrimSpace(h)
	isRe := false
	if strings.HasPrefix(h, "#") {
		isRe = true
		h = h[1:]
	}
	if isRe {
		p.allowHostRe.Delete(h)
	} else {
		p.allowHostMap.Delete(h)
	}
	p.disableHostMap.Store(h, true)
}

func (p *Proxy) ResetHost() {
	p.allowHostMap.Clear()
	utils.CopyMap2Sync(allowModifyHosts, p.allowHostMap)
	p.allowHostRe.Clear()
	p.disableHostMap.Clear()
}

func (p *Proxy) CatchHosts() []string {
	hosts := make([]string, 0)
	p.catchHostMap.Range(func(key, value any) bool {
		hosts = append(hosts, key.(string))
		return true
	})
	sort.Strings(hosts)
	return hosts
}
