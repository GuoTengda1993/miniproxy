/*
 * @Description: cert download service
 */
package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"miniproxy/utils"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type CertServer struct {
	ctx       context.Context
	server    *http.Server
	router    *gin.Engine
	port      int
	isStarted bool
}

func NewCertServer(ctx context.Context, port int) *CertServer {
	return &CertServer{
		ctx:    ctx,
		port:   port,
		router: gin.Default(),
	}
}

func (s *CertServer) Start() {
	if s.isStarted {
		return
	}
	s.download()
	gin.SetMode(gin.ReleaseMode)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}
	go s.server.ListenAndServe()
	s.isStarted = true
}

func (s *CertServer) Stop() {
	if s.isStarted {
		s.server.Close()
		s.isStarted = false
	}
}

func (s *CertServer) download() {
	filename := "cartools-ca-cert.pem"
	s.router.GET(CertDownloadUrl, func(c *gin.Context) {
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Cache-Control", "no-cache")
		c.Header("Pragma", "no-cache")

		reader := bytes.NewReader([]byte(CERTIFICATE))
		response := c.Writer
		b, err := io.ReadAll(reader)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = response.Write(b)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func (p *Proxy) ShowCert() string {
	if p.certDownloadServer == nil {
		p.certDownloadServer = NewCertServer(p.ctx, p.port+1)
	}
	p.certDownloadServer.Start()
	return fmt.Sprintf("http://%s:%d%s", utils.GetLocalIp(), p.port+1, CertDownloadUrl)
}

func (p *Proxy) CloseCert() {
	if p.certDownloadServer == nil {
		return
	}
	p.certDownloadServer.Stop()
}

func (p *Proxy) SaveCertLocal(dir string) error {
	filepath := path.Join(dir, "cartools-ca-cert.pem")
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	if err != nil {
		return fmt.Errorf("文件创建失败: %s", err.Error())
	}
	defer f.Close()
	_, err = f.WriteString(CERTIFICATE)
	if err != nil {
		return fmt.Errorf("文件写入失败: %s", err.Error())
	}
	return nil
}
