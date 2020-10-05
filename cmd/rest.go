package cmd

import (
	
	"io"
	"fmt"
	"net/http"
	"strings"
	"time"
	"bufio"
	"os"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)


const apicLinuxTempPath string = "~/apic"
const defaultPort string = "8080"
const defaultRestPath string = "/api/new"
const defaultSwaggerPort string = "8090"
const defaultSwaggerUIPath string = "/docs"
const defaultSwaggerJsonPath string = "/swagger.json"

type cmdContext struct {
	port string
	hostName string
	configPath string
	swaggerPort string
	swaggerUIPath string
	swaggerJsonPath string
	apiCmds []apiCmd
}

type apiCmd struct {
	description string
	path string
	querystring string
	headerStr string
	headers []header
	cookieStr string
	cookies []cookie
	resp string
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

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Creates REST API",
	Long: 
	``, //TODO desc flags
	Run: restCmdExecute,
}


func init() {

	hostCmd.AddCommand(restCmd)

	restCmd.PersistentFlags().StringP("config", "", "", "config file to host series of APIs")

	restCmd.PersistentFlags().StringP("querystr", "q", "", "query string")
	
	restCmd.PersistentFlags().StringP("port", "p", "", "specify listening port, default:8080")

	restCmd.PersistentFlags().StringP("header", "d", "", "i.e: content-type=application/json custom-key=customvalue")

	restCmd.PersistentFlags().StringP("cookie", "k", "", "i.e: cookie1=value1 cookie2=value2")
	
	restCmd.PersistentFlags().StringP("resp", "r", "", "mock response (always json)")
}

func restCmdExecute(cmd *cobra.Command, args []string) {

	cmdContext := readCmds(cmd)

	createRest(cmdContext)

	printAPIsInfo(cmdContext)


	swagCmd := genSwaggerDocs(cmdContext)

	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()
	fmt.Println(r)
	fmt.Println(err.Error())

	swagCmd.Process.Kill()
}

func readCmds(cmd *cobra.Command) (cmdContext) {

	cmdContext :=newCmdContext()
	port := getPort(cmd)

	cmdContext.port = port
	configPath, _ := cmd.Flags().GetString("config")
	cmdContext.configPath = configPath

	if configPath != "" {
		cmdContext.apiCmds = readConfigFileCmds(configPath)
	} else {
		cmdContext.apiCmds = readCliCmd(cmd)
	}

	return cmdContext
}

func readCliCmd(cmd *cobra.Command) ([]apiCmd) {

	apis := []apiCmd{}

	apicmd := newAPICmd()
	apicmd.path = cmd.Flags().Arg(0)
	if apicmd.path == "" {
		apicmd.path = defaultRestPath
	} else {
		apicmd.path =  formatAPIPath(strings.TrimSpace(apicmd.path))
	}

	cmd.Flags().Visit(func(f *pflag.Flag) {

		switch f.Name {
			case "querystr":
				apicmd.querystring = strings.TrimSpace(f.Value.String())
			case "resp":
				apicmd.resp = f.Value.String()
			case "header":
				apicmd.headerStr = strings.TrimSpace(f.Value.String())
				apicmd.headers = createHeaderSlice(apicmd.headerStr)
			case "cookie":
				apicmd.cookieStr = strings.TrimSpace(f.Value.String())
				apicmd.cookies = createCookieSlice(apicmd.cookieStr)
		}
	})

	apis = append(apis, apicmd)

	return apis
}

func readConfigFileCmds(configPath string) ([]apiCmd) {

	//TODO: log err
	fmt.Println(configPath)

	if configPath == "" {
		return nil
	}

	return nil
}

func createRest(cmdCon cmdContext) {

	r := mux.NewRouter()

	for _, v := range cmdCon.apiCmds {

		resp := createResponse(v)

		createRestHandlers(r, v, resp)
	}

	go http.ListenAndServe(fmt.Sprintf(":%s", cmdCon.port), r)
}

func createRestHandlers(r *mux.Router, api apiCmd, resp response ) { //cmd apiCmd, cmdCon cmdContext) {
	
	r.HandleFunc(api.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("GET")
	
	r.HandleFunc(api.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("POST")

	r.HandleFunc(api.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("PUT")

	r.HandleFunc(api.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("DELETE")
}

func handleResponse(w http.ResponseWriter, r *http.Request, resp response) { //} api apiCmd, cmdCon cmdContext) {

	for _, h := range resp.headers {
		w.Header().Add(h.key, h.value)
	}

	for _, c := range resp.cookies {
		cookie := http.Cookie{
			Name: c.name,
			Value: c.value,
		}

		http.SetCookie(w,&cookie)
	}

	newResp := fmt.Sprintf(`%v:
%v`,  r.Method, resp.resp)
	
	io.WriteString(w, newResp)

	//TODO: print printinfo.createIngressRequestInfo
}

func createResponse(api apiCmd) (response) {

	resp := response{}
	resp.headers = api.headers
	resp.cookies = api.cookies
	resp.resp = strings.TrimSpace(api.resp)

	return resp
}

func createHeaderSlice(headerStr string) ([]header) {

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

func createCookieSlice(cookieStr string) ([]cookie) {
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


func getApiPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	} else {
		return "/api/new"
	}
}

func getPort(cmd *cobra.Command) (string) {
	var port string = defaultPort
	p, _ := cmd.Flags().GetString("port")
	if p != "" {
		port = p
	}
	return strings.TrimSpace(port)
}

func newCmdContext() (cmdContext) {
	host, _ := os.Hostname()

	return cmdContext {
		port: defaultPort,
		hostName: host,
		swaggerPort: defaultSwaggerPort,
		swaggerUIPath: defaultSwaggerUIPath,
		swaggerJsonPath: defaultSwaggerJsonPath,
	}
}

func newAPICmd() (apiCmd) {
	j, _ := json.Marshal("success")
	return apiCmd{
		path: "/api/new",
		resp: string(j),
		description: "default apic API returning success",
	}
}


