package main

import (
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {}
func init() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
}

func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{}
}

type httpContext struct {
	types.DefaultHttpContext
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	method, err := proxywasm.GetHttpRequestHeader(":method")
	if err != nil {
		proxywasm.LogWarnf("failed to get :method header: %v", err)
		method = "UNKNOWN"
	}

	path, err := proxywasm.GetHttpRequestHeader(":path")
	if err != nil {
		proxywasm.LogWarnf("failed to get :path header: %v", err)
		path = "UNKNOWN"
	}

	host, err := proxywasm.GetHttpRequestHeader(":authority")
	if err != nil {
		host, err = proxywasm.GetHttpRequestHeader("host")
		if err != nil {
			proxywasm.LogWarnf("failed to get host header: %v", err)
			host = "UNKNOWN"
		}
	}

	proxywasm.LogInfof("HTTP Request: method=%s path=%s host=%s", method, path, host)

	return types.ActionContinue
}
