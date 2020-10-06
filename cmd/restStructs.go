package cmd

import(
	"time"
	"encoding/json"
	"strings"
	"os"
)

type cmdContext struct {
	Port string
	HostName string
	configPath string
	swaggerPort string
	swaggerUIPath string
	swaggerJsonPath string
	ApiCmds []apiCmd
}

type apiCmd struct {
	Description string
	Path string
	Querystring string
	headerStr string
	headers []header
	cookieStr string
	cookies []cookie
	Resp string
	data string
	auth apiAuth
}

type apiAuth struct {
	enableAuth bool
	mode string
	userName string
	password string
	jwtToken string
	publiCert string
}

type header struct {
	key string
	value string
}

type cookie struct {
	name string
	value string
	expiry time.Duration
}

type response struct {
	resp string
	headers []header
	cookies []cookie
}

func newResponse(api apiCmd) (response) {

	resp := response{}
	resp.headers = api.headers
	resp.cookies = api.cookies
	resp.resp = strings.TrimSpace(api.Resp)

	return resp
}

func newHeaderSlice(headerStr string) ([]header) {

	headers := []header{}

	hs := strings.Split(headerStr, " ")

	if len(hs) == 0 {
		return headers
	}

	for _, v := range hs {
		
		kvs := strings.Split(v, "=")

		if len(kvs) == 0 {
			return headers
		}

		headers = append(headers, header{
			key: kvs[0],value: kvs[1],
		})
	}

	return headers
}

func newCookieSlice(cookieStr string) ([]cookie) {
	cookies := []cookie{}

	cs := strings.Split(cookieStr, " ")

	if len(cs) == 0 {
		return cookies
	}

	for _, v := range cs {
		
		kvs := strings.Split(v, "=")

		if len(kvs) == 0 {
			return cookies
		}


		cookies = append(cookies, cookie{
			name: kvs[0],
			value: kvs[1],
			expiry: time.Hour * 1,
		})
	}

	return cookies
}

func newCmdContext() (cmdContext) {
	host, _ := os.Hostname()

	return cmdContext {
		Port: defaultPort,
		HostName: host,
		swaggerPort: defaultSwaggerPort,
		swaggerUIPath: defaultSwaggerUIPath,
		swaggerJsonPath: defaultSwaggerJsonPath,
	}
}

func newAPICmd() (apiCmd) {
	j, _ := json.Marshal("success")
	return apiCmd{
		Path: "/api/new",
		Resp: string(j),
		Description: "default apic API returning success",
	}
}