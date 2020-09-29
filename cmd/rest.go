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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apic/cmd/stdout"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const defaultRestPath string = "api/new"

type apiCmd struct {
	path string
	port string
	querystring string
	headerStr string
	headers []header
	resp string
	configPath string
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
	hostName string
	ip string
	port string
	path string
	querystring string
	headers []header
	mockResp string
	requArrivedAt time.Time
}

type flagProp struct {
	
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

	restCmd.PersistentFlags().StringP("header", "d", "", "headers space-delimited: content-type=application/json custom-key=customvalue")
	
	restCmd.PersistentFlags().StringP("resp", "r", "", "mock response (always json)")
}

func restCmdExecute(cmd *cobra.Command, args []string) {

	apis := readCmds(cmd)

	createRest(cmd, apis)
}

func readCmds(cmd *cobra.Command) ([]apiCmd) {

	apiCmds, exist := readConfigFileCmds(cmd)

	if !exist {
		return readCliCmd(cmd)
	} else {
		return apiCmds
	}
}

func readCliCmd(cmd *cobra.Command) ([]apiCmd) {

	apis := []apiCmd{}

	apicmd := apiCmd{}
	apicmd.path = cmd.Flags().Arg(0)

	cmd.Flags().Visit(func(f *pflag.Flag) {

		switch f.Name {
			case "querystr":
				apicmd.querystring = f.Value.String()
			case "resp":
				apicmd.resp = f.Value.String()
		}
	})

	apis = append(apis, apicmd)

	return apis
}

func readConfigFileCmds(cmd *cobra.Command) ([]apiCmd, bool) {

	configPath, _ := cmd.Flags().GetString("config")
	//TODO: log err
	fmt.Println(configPath)

	if configPath == "" {
		return nil, false
	}

	return nil, true
}

func createRest(cmd *cobra.Command, apiCmds []apiCmd) {

	port := getPort(cmd)

	if len(apiCmds) == 0 {
		stdout.PInfo("cmd not found")
	}

	r := mux.NewRouter()

	for _, v := range apiCmds {

		createRestHandlers(r, v)
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func createRestHandlers(r *mux.Router, cmd apiCmd) {
	
	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd) }).Methods("GET")
	
	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd) }).Methods("POST")

	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd) }).Methods("PUT")

	r.HandleFunc(cmd.path, func(w http.ResponseWriter, r *http.Request){ handleResponse(w, r, cmd) }).Methods("DELETE")
}

func handleResponse(w http.ResponseWriter, r *http.Request, api apiCmd) {
	fmt.Println(api.resp)
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



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

