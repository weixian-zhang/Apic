package cmd

import (
	
	"io"
	"fmt"
	"net/http"
	"bufio"
	"os"
	"strings"
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

	restApiContext := readCmds(cmd)

	createRest(restApiContext)

	printAPIsInfo(restApiContext)


	swagCmd := genSwaggerDocs(restApiContext)

	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()
	fmt.Println(r)
	fmt.Println(err.Error())

	swagCmd.Process.Kill()
}

func readCmds(cmd *cobra.Command) (RestApiContext) {

	restApiContext :=newRestApiContext()
	port := getPort(cmd)

	restApiContext.Port = port
	configPath, _ := cmd.Flags().GetString("config")
	restApiContext.configPath = configPath

	if configPath != "" {
		restApiContext.RestApis = readConfigFileCmds(configPath)
	} else {
		restApiContext.RestApis = readCliCmd(cmd)
	}

	return restApiContext
}

func readCliCmd(cmd *cobra.Command) ([]RestApi) {

	apis := []RestApi{}

	api := newRestApi()
	api.Path = cmd.Flags().Arg(0)
	if api.Path == "" {
		api.Path = defaultRestPath
	} else {
		api.Path =  formatAPIPath(strings.TrimSpace(api.Path))
	}

	cmd.Flags().Visit(func(f *pflag.Flag) {

		switch f.Name {
			case "querystr":
				api.Querystring = strings.TrimSpace(f.Value.String())
			case "resp":
				api.Resp = f.Value.String()
			case "header":
				api.headerStr = strings.TrimSpace(f.Value.String())
				api.headers = newHeaderSlice(api.headerStr)
			case "cookie":
				api.cookieStr = strings.TrimSpace(f.Value.String())
				api.cookies = newCookieSlice(api.cookieStr)
		}
	})

	apis = append(apis, api)

	return apis
}

func readConfigFileCmds(configPath string) ([]RestApi) {

	//TODO: log err
	fmt.Println(configPath)

	if configPath == "" {
		return nil
	}

	return nil
}

func createRest(apiContext RestApiContext) {

	r := mux.NewRouter()

	for _, v := range apiContext.RestApis {

		resp := newResponse(v)

		createRestHandlers(r, v, resp)
	}

	go http.ListenAndServe(fmt.Sprintf(":%s", apiContext.Port), r)
}

func createRestHandlers(r *mux.Router, api RestApi, resp response ) { //cmd RestApi, cmdCon RestApiContext) {
	
	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("GET")
	
	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("POST")

	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("PUT")

	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("DELETE")
}

func handleResponse(w http.ResponseWriter, r *http.Request, resp response) { //} api RestApi, cmdCon RestApiContext) {

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


