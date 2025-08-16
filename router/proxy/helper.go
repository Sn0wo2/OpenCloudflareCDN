package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go/http3"
)

func Request(ctx *gin.Context) {
	remote, err := url.Parse(config.Instance.OriginalServer)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var rt http.RoundTripper

	h3 := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if res, err := h3.RoundTrip(&http.Request{
		Method: http.MethodHead,
		URL:    remote,
	}); err == nil {
		_ = res.Body.Close()
		rt = h3
	} else {
		rt = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = rt
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = ctx.Request.URL.Path
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
