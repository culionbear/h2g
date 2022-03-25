package h2g

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kataras/iris/v12"
)

type Manager struct {
	handler		map[string]*module
	service		string
	method		string
}

func New(c Config) *Manager {
	m := &Manager{
		handler: make(map[string]*module),
		service: c.Service,
		method: c.Method,
	}
	for k, f := range c.Handler {
		m.handler[k] = newModule(f())
	}
	if m.service == "" {
		m.service = "service"
	}
	if m.method == "" {
		m.method = "method"
	}
	return m
}

func (m *Manager) Service(ctx iris.Context) {
	service, method := ctx.Params().Get(m.service), ctx.Params().Get(m.method)
	if service == "" || method == "" {
		m.writeError(ctx, errors.New("param is error"))
		return
	}
	buf, err := ctx.GetBody()
	if err != nil {
		m.writeError(ctx, err)
		return
	}
	buf, err = m.execute(service, method, buf)
	if err != nil {
		m.writeError(ctx, err)
		return
	}
	m.write(ctx, http.StatusOK, buf)
}

func (m *Manager) Path(uri string) string {
	return fmt.Sprintf("%s/{%s}/{%s}", uri, m.service, m.method)
}

func (m *Manager) execute(service, method string, body []byte) ([]byte, error) {
	h, ok := m.handler[service]
	if !ok {
		return nil, fmt.Errorf("handler[%s] is not exists", service)
	}
	return h.call(method, body)
}

func (m *Manager) write(ctx iris.Context, code int, buf []byte) {
	ctx.StatusCode(code)
	ctx.Write(buf)
}

func (m *Manager) writeJson(ctx iris.Context, code int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		m.writeError(ctx, err)
	}
	m.write(ctx, code, buf)
}

func (m *Manager) writeError(ctx iris.Context, err error) {
	m.writeJson(ctx, http.StatusInternalServerError, map[string]interface{}{"msg": err.Error()})
}
