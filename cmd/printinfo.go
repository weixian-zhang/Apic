package cmd

import (
	"github.com/gookit/color"
)


type pInfoApi struct {
	hostIP string
	hostName string
	path string
	querystring string
	headers []header
	resp string //json data or string
	swaggerPath string
}

type pInfoIngressRequest struct {
	clientIP string
}

func createApiInfo() {

}