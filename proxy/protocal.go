/*
 * @Description: protocal
 */
package proxy

type FlowInfo struct {
	RemoteAddr       string            `json:"remote_addr"`
	ID               int64             `json:"id"`
	URL              string            `json:"url"`
	Scheme           string            `json:"scheme"`
	Host             string            `json:"host"`
	Path             string            `json:"path"`
	Method           string            `json:"method"`
	Headers          map[string]string `json:"headers"`
	Query            string            `json:"query"`
	Data             string            `json:"data"`
	ReqTime          int64             `json:"req_time"`
	ReqTimeStr       string            `json:"req_time_str"`
	ReqContentLength string            `json:"req_content_length"`

	Code              int               `json:"code"`
	ReturnCode        string            `json:"return_code"` // code in json response
	RespContentLength string            `json:"resp_content_length"`
	ResponseHeaders   map[string]string `json:"response_headers"`
	ResponseBody      string            `json:"response_body"`
	ResponseTime      int64             `json:"response_time"`
	ResponseTimeStr   string            `json:"response_time_str"`
	Curl              string            `json:"curl"`     // curl
	Duration          int64             `json:"duration"` // cost time
}

type MsgInfo struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	UseBox  bool   `json:"use_box"`
	Title   string `json:"title"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
