package httppool

import (
	"cache"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const DefaultBasePath = "/_gee_cache_/"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: DefaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v))
}

func (p *HTTPPool) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + req.URL.Path)
	}
	p.Log("%s %s", req.Method, req.URL.Path)
	parts := strings.SplitN(req.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		resp.WriteHeader(http.StatusBadRequest)
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := cache.GetGroup(groupName)
	if group == nil {
		resp.WriteHeader(http.StatusNotFound)
		http.Error(resp, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.WriteHeader(http.StatusOK)
	resp.Write(view.ByteSlice())
}
