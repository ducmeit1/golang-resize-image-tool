package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"path/filepath"
)

type GatewayHandler struct {
	initialized bool
	muxLambda   *gorillamux.GorillaMuxAdapter
	hw          *HelloWorldHandler
	rc          *ResizeCropHandler
}

func (gw *GatewayHandler) ServeHTTP(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request.Path = filepath.Clean(request.Path)
	gw.initRouter()
	return gw.muxLambda.Proxy(request)
}

func (gw *GatewayHandler) initRouter() {
	if !gw.initialized {
		r := mux.NewRouter().StrictSlash(false)
		gw.hw = new(HelloWorldHandler)
		gw.rc = new(ResizeCropHandler)
		r.HandleFunc("/go", gw.hw.ServeHTTP)
		r.HandleFunc("/{optional}", gw.rc.ServeHTTP).
			Methods("GET")
		gw.muxLambda = gorillamux.New(r)
		gw.initialized = true
	}
}
