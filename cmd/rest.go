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
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

const defaultRestPath string = "api/new"

type apiCmd struct {
	path string
	port string
	querystring string
	headerStr string
	headers []header
	resp string
	dataPath string
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

var api apiCmd = apiCmd{}
var apis []apiCmd
// var path string
// var port string
var configPath string

// var singleApi api = api{
// 	path: "api/new",
// 	port: "8080",
// 	data: "{}",
// 	querystring: "",
// }

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

	restCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file to host series of APIs")

	restCmd.Flags().StringVarP(&api.querystring, "querystring", "q", "", "query string")
	
	restCmd.Flags().StringVarP(&api.port, "port", "p", "", "specify listening port, default:8080")

	restCmd.Flags().StringVarP(&api.headerStr, "header", "h", "", "headers space-delimited: content-type=application/json custom-key=customvalue")
	
	restCmd.Flags().StringVarP(&api.resp, "resp", "r", "", "mock response (always json)")
	
}

func restCmdExecute(cmd *cobra.Command, args []string) {

	

	//singleApi.path = path
 
	//fmt.Println(singleApi.path)

	// if len(args) == 0 {
	// 	createRest(defaultRestPath)
	// } else {
	// 	path := args[0]
	// 	createRest(path)
	// }

}

func createRest() {
	r := mux.NewRouter()

	http.ListenAndServe(":80", r)
}

func createRestHandlers(r *mux.Router, path string) {
	

	r.HandleFunc(path, httpHandler).Methods("GET")
	r.HandleFunc(path, httpHandler).Methods("POST")
	r.HandleFunc(path, httpHandler).Methods("PUT")
	r.HandleFunc(path, httpHandler).Methods("DELETE")
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	
}


func httpResponse(w http.ResponseWriter, r *http.Request, api apiCmd, resp response) {

}

func readCmd(apis *[]apiCmd) {

	exist := readConfig(apis)

	if !exist {
		readSingleCmd(apis)
	}
}

func readSingleCmd(apis *[]apiCmd) {
	
}

func readConfig(apis *[]apiCmd) (bool) {
	if configPath == "" {
		return false
	}

	return true
}

func getPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	} else {
		return "/api/new"
	}
}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

