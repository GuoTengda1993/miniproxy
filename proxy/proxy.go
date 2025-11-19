/*
proxy
*/
package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"miniproxy/utils"

	"github.com/elazarl/goproxy"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Proxy struct {
	ctx context.Context

	port         int
	server       *http.Server
	isRunning    bool
	ticker       *time.Ticker
	tickerStatus bool
	stopChan     chan bool

	// cert download server
	certDownloadServer *CertServer

	allowHostMap   *sync.Map
	allowHostRe    *sync.Map
	disableHostMap *sync.Map

	flowRecords      *utils.FixedQueue[*FlowInfo]
	flowInfoMap      *sync.Map
	flowUnfinished   *sync.Map
	remoteAddrMap    *sync.Map
	remoteDisableMap *sync.Map
	catchHostMap     *sync.Map

	// Tmp cache
	// flowChangeMark(bool): mark new flow
	// frontPause(bool): mark front is inactive
	cache *sync.Map
}

func NewProxy(ctx context.Context, port int) *Proxy {
	p := &Proxy{
		ctx:  ctx,
		port: port,

		allowHostMap:   &sync.Map{},
		allowHostRe:    &sync.Map{},
		disableHostMap: &sync.Map{},

		flowRecords:      utils.NewFixedQueue[*FlowInfo](1 << 10), // max 1024
		flowInfoMap:      &sync.Map{},
		flowUnfinished:   &sync.Map{},
		remoteAddrMap:    &sync.Map{},
		remoteDisableMap: &sync.Map{},
		catchHostMap:     &sync.Map{},

		cache: &sync.Map{},
	}
	utils.CopyMap2Sync(allowModifyHosts, p.allowHostMap)
	return p
}

func (p *Proxy) InitFromLocalConf() {
	p.loadHostFromLocalConf()
}

func (p *Proxy) SaveToLocalConf() {
	p.saveHostToLocalConf()
}

// proxy info
func (p *Proxy) Info() map[string]any {
	ip := utils.GetLocalIp()
	result := map[string]any{
		"isRunning": p.isRunning,
		"addr":      fmt.Sprintf("%s:%d", ip, p.port),
	}
	return result
}

func (p *Proxy) setCA() error {
	goproxyCa, err := tls.X509KeyPair([]byte(CERTIFICATE), []byte(RSA_PRIVATE_KEY))
	if err != nil {
		runtime.LogErrorf(p.ctx, ">>> load ca error: %s", err.Error())
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		runtime.LogErrorf(p.ctx, ">>> parse ca error: %s", err.Error())
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

func (p *Proxy) reqHostFilter() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		if remoteAddr := strings.Split(req.RemoteAddr, ":")[0]; remoteAddr != "" {
			if _, ok := p.remoteDisableMap.Load(remoteAddr); ok {
				return false
			}
		}

		host := req.Host
		if strings.HasSuffix(host, ":443") {
			host = strings.ReplaceAll(host, ":443", "")
		}
		if _, ok := p.allowHostMap.Load(host); ok {
			return true
		}
		reList := make([]string, 0)
		p.allowHostRe.Range(func(key, value any) bool {
			reList = append(reList, key.(string))
			return true
		})
		for _, pattern := range reList {
			re := regexp.MustCompile(pattern)
			if re.MatchString(req.Host) {
				return true
			}
		}
		p.disableHostMap.LoadOrStore(host, true)
		return false
	}
}

func (p *Proxy) Restart() {
	p.Stop()
	p.Run()
}

// Start proxy
func (p *Proxy) Run() {
	if p.isRunning {
		return
	}
	if !utils.PortCheck(p.port) {
		p.MessagePush("", fmt.Sprintf("Port(%d) is occupied, cannot start", p.port), false)
		return
	}
	p.stopChan = make(chan bool)
	pxy := goproxy.NewProxyHttpServer()
	p.setCA()
	// pxy.Verbose = true
	pxy.OnRequest(p.reqHostFilter()).HandleConnect(goproxy.AlwaysMitm)
	pxy.OnRequest(p.reqHostFilter()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			remoteAddr := strings.Split(r.RemoteAddr, ":")[0]
			// record request
			defer func() {
				p.flowRecordReq(r, ctx)
			}()
			if _, ok := p.remoteAddrMap.LoadOrStore(remoteAddr, true); !ok {
				runtime.EventsEmit(p.ctx, EventNameProxyRemoteAddr, true)
			}
			if _, ok := p.catchHostMap.LoadOrStore(r.Host, true); !ok {
				runtime.EventsEmit(p.ctx, EventNameProxyHost, true)
			}
			return r, nil
		})
	pxy.OnResponse(p.reqHostFilter()).DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if resp == nil {
				return resp
			}
			// record response
			p.flowRecordResp(resp, ctx)
			return resp
		})

	p.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", p.port),
		Handler: pxy,
	}
	go p.server.ListenAndServe()
	p.isRunning = true
	p.startTicker()
}

func (p *Proxy) startTicker() {
	if p.tickerStatus {
		return
	}
	if p.ticker == nil {
		p.ticker = time.NewTicker(time.Millisecond * 357)
	} else {
		p.ticker.Reset(time.Millisecond * 357)
	}
	p.tickerStatus = true
	go func() {
		for {
			select {
			case <-p.ticker.C:
				if p.frontIsPause() {
					continue
				}
				if tmp, loaded := p.cache.Load(cacheKeyFlowChangeMark); loaded {
					if tmp.(bool) {
						runtime.EventsEmit(p.ctx, EventNameReloadFlow, true)
					}
				}
			case <-p.stopChan:
				return
			}
		}
	}()
}

func (p *Proxy) stopTicker() {
	if p.tickerStatus {
		p.tickerStatus = false
		p.stopChan <- true
	}
}

// stop proxy
func (p *Proxy) Stop() {
	if p.isRunning {
		p.server.Close()
		defer close(p.stopChan)
	}
	if p.certDownloadServer != nil {
		p.certDownloadServer.Stop()
	}
	p.stopTicker()
	if p.ticker != nil {
		p.ticker.Stop()
	}
	p.isRunning = false
}

func (p *Proxy) frontIsPause() bool {
	if _, loaded := p.cache.Load(cacheKeyFrontPause); loaded {
		return true
	}
	return false
}

func (p *Proxy) FrontPause() {
	p.cache.Store(cacheKeyFrontPause, true)
}

func (p *Proxy) FrontResume() {
	p.cache.Delete(cacheKeyFrontPause)
	p.cache.Store(cacheKeyFlowChangeMark, true)
}

func (p *Proxy) RemoteAddrs() []string {
	addrs := make([]string, 0)
	p.remoteAddrMap.Range(func(key, value any) bool {
		addr := key.(string)
		if value.(string) != "" {
			addr = fmt.Sprintf("%s-%s", value.(string), addr)
		}
		addrs = append(addrs, addr)
		return true
	})
	sort.Strings(addrs)
	return addrs
}
