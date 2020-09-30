/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"io"
	"fmt"
	"net/http"
	"net"
	"strings"
	"time"
	"os"
	"errors"
	"github.com/apic/cmd/stdout"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const defaultRestPath string = "api/new"

type cmdContext struct {
	port string
	configPath string
	apiCmds []apiCmd
}

type apiCmd struct {
	path string
	querystring string
	headerStr string
	headers []header
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

type response struct {
	srcIP string
	hostName string
	IP string
	port string
	path string
	querystring string
	headers []header
	mockResp string //json data or string
	swaggerPath string
	requArrivedAt time.Time
}

var defaultPort string = "8080"
var apicmd = apiCmd{}
//var port string = "8080"
var configPath string

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

	restCmd.PersistentFlags().StringP("config", "c", "", "config file to host series of APIs")

	restCmd.PersistentFlags().StringP("querystr", "q", "", "query string")
	
	restCmd.PersistentFlags().StringP("port", "p", "", "specify listening port, default:8080")

	restCmd.PersistentFlags().StringP("header", "d", "", "i.e: content-type=application/json custom-key=customvalue")
	
	restCmd.PersistentFlags().StringP("resp", "r", "", "mock response (always json)")
}

func restCmdExecute(cmd *cobra.Command, args []string) {

	cmdContext := readCmds(cmd)

	createRest(cmdContext)
}

func readCmds(cmd *cobra.Command) (cmdContext) {

	cmdContext := cmdContext{}
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

	cmd.Flags().Visit(func(f *pflag.Flag) {

		switch f.Name {
			case "querystr":
				apicmd.querystring = f.Value.String()
			case "resp":
				apicmd.resp = f.Value.String()
			case "header":
				apicmd.headerStr = f.Value.String()
				apicmd.headers = getHeaders(f.Value.String())
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

	if len(cmdCon.apiCmds) == 0 {
		stdout.Print("cmd not found")
	}

	r := mux.NewRouter()

	for _, v := range cmdCon.apiCmds {

		createRestHandlers(r, v, cmdCon)
	}

	http.ListenAndServe(fmt.Sprintf(":%s", cmdCon.port), r)
}

func createRestHandlers(r *mux.Router, cmd apiCmd, cmdCon cmdContext) {
	
	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd, cmdCon) }).Methods("GET")
	
	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd, cmdCon) }).Methods("POST")

	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd, cmdCon) }).Methods("PUT")

	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd, cmdCon) }).Methods("DELETE")
}

func handleResponse(w http.ResponseWriter, r *http.Request, api apiCmd, cmdCon cmdContext) {
	resp := createResponse(r, api, cmdCon)

	for _, h := range resp.headers {
		w.Header().Add(h.key, h.value)
	}

	io.WriteString(w, resp.mockResp)

	stdout.Print(stdout.Cyan(resp.mockResp))
}

func createResponse(r *http.Request, api apiCmd, cmdCon cmdContext) (response) {
	resp := response{}

	resp.srcIP = r.RemoteAddr

	host, _ := os.Hostname()
	resp.hostName = host

	resp.headers = api.headers

	resp.path = api.path
	resp.port = cmdCon.port

	resp.querystring = strings.TrimSpace(api.querystring)
	resp.mockResp = strings.TrimSpace(api.resp)
	ip, err := getLocalIP()

	if err == nil {
		resp.IP = ip
	}

	resp.requArrivedAt = time.Now()

	return resp
}


func getApiPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	} else {
		return "/api/new"
	}
}

func getHeaders(headerStr string) ([]header) {
	return nil
}

func getPort(cmd *cobra.Command) (string) {
	var port string = defaultPort
	p, _ := cmd.Flags().GetString("port")
	if p != "" {
		port = p
	}
	return strings.TrimSpace(port)
}

func getLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func newAPICmd() (apiCmd) {
	return apiCmd{
		path: "/api/new",
		resp: "mocked resp",
	}
}


