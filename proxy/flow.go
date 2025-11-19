package proxy

import (
	"bytes"
	"fmt"
	"io"
	"miniproxy/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/tidwall/gjson"
	"moul.io/http2curl"
)

const pageSize = 50

func (p *Proxy) flowRecordReq(req *http.Request, ctx *goproxy.ProxyCtx) {
	reqTime := time.Now()
	reqTimeInt := reqTime.UnixMilli()
	reqTimeStr := reqTime.Format("2006-01-02 15:04:05.99")

	reqBody := ""
	if strings.ToUpper(req.Method) == "POST" {
		if buf, err := io.ReadAll(req.Body); err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(buf))
			reqBody = string(buf)
		}
	}

	curlCommandStr := ""
	if curlCommand, err := http2curl.GetCurlCommand(req); err == nil {
		curlCommandStr = curlCommand.String()
	}

	newHeaders := make(map[string]string)
	for k, v := range req.Header {
		if len(v) == 0 {
			continue
		}
		newHeaders[k] = strings.Join(v, " ")
	}

	flow := &FlowInfo{
		RemoteAddr:       strings.Split(req.RemoteAddr, ":")[0],
		ID:               ctx.Session,
		URL:              fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path),
		Scheme:           req.URL.Scheme,
		Host:             req.Host,
		Path:             req.URL.Path,
		Method:           req.Method,
		Headers:          newHeaders,
		Query:            req.URL.RawQuery,
		Data:             reqBody,
		ReqTime:          reqTimeInt,
		ReqTimeStr:       reqTimeStr,
		ReqContentLength: utils.Bytes2String(req.ContentLength),
		Curl:             curlCommandStr,
	}
	p.flowUnfinished.Store(ctx.Session, flow)
	p.flowRecords.Enqueue(flow)
	p.flowChange(true)
}

func (p *Proxy) flowRecordResp(r *http.Response, ctx *goproxy.ProxyCtx) {
	f, ok := p.flowUnfinished.LoadAndDelete(ctx.Session)
	if !ok {
		return
	}
	flow := f.(*FlowInfo)
	now := time.Now()
	respTime := now.UnixMilli()

	respContentLength := ""
	if r.ContentLength > 0 {
		respContentLength = fmt.Sprintf("%d (%s)", r.ContentLength, utils.Bytes2String(r.ContentLength))
	}

	respHeaders := make(map[string]string)
	for k, v := range r.Header {
		if len(v) == 0 {
			continue
		}
		respHeaders[k] = strings.Join(v, " ")
	}

	flow.ResponseHeaders = respHeaders
	flow.Duration = respTime - flow.ReqTime
	flow.Code = r.StatusCode
	flow.ResponseTime = respTime
	flow.ResponseTimeStr = now.Format("2006-01-02 15:04:05.99")
	flow.RespContentLength = respContentLength

	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err == nil {
		flow.ResponseBody = string(resp)
		getCode := gjson.Get(flow.ResponseBody, "code")
		if getCode.Exists() {
			flow.ReturnCode = utils.ToStr(getCode.Int())
		}
	} else {
		flow.ResponseBody = fmt.Sprintf("{\"ERROR\":\"响应解析失败: %s\"}", err.Error())
	}
	r.Body = io.NopCloser(bytes.NewReader(resp))
	p.flowChange(true)
}

func (p *Proxy) FlowList(filter, remoteAddr, host string, reverse bool, onlyError bool) []*FlowInfo {
	defer p.flowChange(false)
	filter = strings.TrimSpace(filter)
	allFlow := p.flowRecords.ToSlice(reverse)
	if filter == "" && remoteAddr == "" && host == "" && !onlyError {
		return allFlow
	}
	if remoteAddr != "" && strings.Contains(remoteAddr, "-") {
		remoteAddr = strings.Split(remoteAddr, "-")[1]
	}

	filter = strings.ToLower(filter)
	filterCode, _ := strconv.Atoi(filter)

	result := make([]*FlowInfo, 0)
	for _, flow := range allFlow {
		if remoteAddr != "" && flow.RemoteAddr != remoteAddr {
			continue
		}
		if host != "" && flow.Host != host {
			continue
		}

		if onlyError && (flow.Code == 200 && (flow.ReturnCode == "0" || flow.ReturnCode == "200")) {
			continue
		}
		if filter != "" {
			if filterCode > 100 && filterCode < 600 {
				if flow.Code == filterCode {
					result = append(result, flow)
				}
			} else {
				if strings.Contains(strings.ToLower(flow.Path), filter) {
					result = append(result, flow)
				}
			}
		} else {
			result = append(result, flow)
		}
	}
	return result
}

func (p *Proxy) FlowListByPage(filter, remoteAddr, host string, page int, onlyError bool) ([]*FlowInfo, int) {
	if page <= 0 {
		page = 1
	}
	all := p.FlowList(filter, remoteAddr, host, true, onlyError)
	total := len(all)
	if len(all)/pageSize+1 <= 1 {
		return all, total
	}
	start := (page - 1) * pageSize
	if start > total {
		start = 0
		page = 1
	}
	end := page * pageSize
	if end > total {
		end = total
	}
	return all[start:end], total
}

func (p *Proxy) FlowClear() {
	p.flowRecords.Flush(true)
	p.flowChange(true)
	p.flowUnfinished.Clear()
	p.remoteAddrMap.Clear()
	p.remoteDisableMap.Clear()
	p.catchHostMap.Clear()
}

func (p *Proxy) flowChange(mark bool) {
	p.cache.Store(cacheKeyFlowChangeMark, mark)
}
