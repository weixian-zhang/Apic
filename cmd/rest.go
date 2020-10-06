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

	cmdContext.Port = port
	configPath, _ := cmd.Flags().GetString("config")
	cmdContext.configPath = configPath

	if configPath != "" {
		cmdContext.ApiCmds = readConfigFileCmds(configPath)
	} else {
		cmdContext.ApiCmds = readCliCmd(cmd)
	}

	return cmdContext
}

func readCliCmd(cmd *cobra.Command) ([]apiCmd) {

	apis := []apiCmd{}

	apicmd := newAPICmd()
	apicmd.Path = cmd.Flags().Arg(0)
	if apicmd.Path == "" {
		apicmd.Path = defaultRestPath
	} else {
		apicmd.Path =  formatAPIPath(strings.TrimSpace(apicmd.Path))
	}

	cmd.Flags().Visit(func(f *pflag.Flag) {

		switch f.Name {
			case "querystr":
				apicmd.Querystring = strings.TrimSpace(f.Value.String())
			case "resp":
				apicmd.Resp = f.Value.String()
			case "header":
				apicmd.headerStr = strings.TrimSpace(f.Value.String())
				apicmd.headers = newHeaderSlice(apicmd.headerStr)
			case "cookie":
				apicmd.cookieStr = strings.TrimSpace(f.Value.String())
				apicmd.cookies = newCookieSlice(apicmd.cookieStr)
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

	for _, v := range cmdCon.ApiCmds {

		resp := newResponse(v)

		createRestHandlers(r, v, resp)
	}

	go http.ListenAndServe(fmt.Sprintf(":%s", cmdCon.Port), r)
}

func createRestHandlers(r *mux.Router, api apiCmd, resp response ) { //cmd apiCmd, cmdCon cmdContext) {
	
	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("GET")
	
	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("POST")

	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("PUT")

	r.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, resp) }).Methods("DELETE")
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


