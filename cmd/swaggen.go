package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	//"time"
	//"syscall"
	"github.com/gobuffalo/packr/v2"
)

var swaggerPath string
const apicLinuxTempPath string = "~/apic"
var apicWinTempPath string = filepath.Join(os.Getenv("APPDATA"), "apic")
var swaggerTempPath string
var swaggerExeTempPath string = filepath.Join(apicWinTempPath, "swagger.exe")
var swaggerYmlTempPath string = filepath.Join(apicWinTempPath, "gotemplate-swagger.yml")
var box *packr.Box

//show swagger path in printApiinfo
//fix swagger not found
//remoteaddr to be ipv4

func initExeYmlSwagPath() {
	if runtime.GOOS == "windows" {
		swaggerPath = apicWinTempPath

		//err := exec.Command("setx", "path", swaggerPath).Run() //set swagger path to %PATH%
		//printErr(err)
	} else {
		swaggerPath = apicLinuxTempPath
	}

	swaggerTempPath = swaggerPath
	swaggerExeTempPath = filepath.Join(swaggerTempPath, "swagger.exe")
	swaggerYmlTempPath = filepath.Join(swaggerTempPath, "gotemplate-swagger.yml")
}

func serveSwaggerDocs(pexit chan bool, apiContext RestApiContext) (error) {

	initExeYmlSwagPath()

	box = packr.New("swagger", "./swag")

	if !fileDirExists(swaggerPath) {
		os.Mkdir(swaggerPath, 0755)
		//TODO: log, fail to create apic temp folder
	}

	err := genSwaggerYml(apiContext) //gens swagger yaml in apic temp folder
	if err != nil {
		return err
	}

	serr := execSwagger(pexit, apiContext)
	if serr != nil {
		printErr(serr)
		return serr
	}

	return nil
}

//swagger editor
//https://editor.swagger.io/
func genSwaggerYml(cmdCon RestApiContext) (error) {
	
	if fileDirExists(swaggerYmlTempPath) {
		err := os.Remove(swaggerYmlTempPath)
		if err != nil {
			printErr(err)
			return err
		}
	}

	//get template from packr packed file
	swagYmlStr, _ := box.FindString("gotemplate-swagger.tpl")

	t, err := template.New("swagger").Parse(swagYmlStr)
	if err != nil {
		printErr(err)
		return err
	}

	file, err := os.Create(swaggerYmlTempPath)
	if err != nil {
		printErr(err)
		return err
	}

	 eerr := t.Execute(file, cmdCon)
	 if eerr != nil {
		printErr(err)
		return err
	}

	return nil
}

func execSwagger(pexit chan bool, apiContext RestApiContext) (error) {
	sbyte, _ := box.Find("swagger.exe")

	if !fileDirExists(swaggerExeTempPath) {
		werr := ioutil.WriteFile(swaggerExeTempPath, sbyte, 0755)
		//TODO: log
		fmt.Println(werr)
	}

	swagexecCmd := exec.Command(swaggerExeTempPath, "serve", "-p", apiContext.swaggerPort, "-F=swagger", swaggerYmlTempPath)
	//swagexecCmd.Dir = swaggerPath

	var out bytes.Buffer
	var stderr bytes.Buffer
	swagexecCmd.Stdout = &out
	swagexecCmd.Stderr = &stderr

	serr := swagexecCmd.Start()

	if serr != nil {
		fmt.Println("swagstart: " + fmt.Sprint(serr.Error()) + ": " + stderr.String())
		return serr
	}

	printInfo(os.Getenv("PATH"))
	printInfo(swagexecCmd.String())

	go func() {
		//for {
		select {
			case <- pexit: //on cli exits kill swagger.exe
			swagexecCmd.Process.Kill()
		}
		//}
	}()

	return nil

	// fmt.Println(<-pexit)

	//proc, err := startProcess(swaggerExeTempPath, "serve", "-p", apiContext.swaggerPort, "-F=swagger", "./gotemplate-swagger.yml")
	
	// if err != nil {
	// 	printErr(err)
	// 	return err
	// }

	
}

func fileDirExists(filedirpath string) bool {
    if _, err := os.Stat(filedirpath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// func startProcess(exePath string, params ...string) (*os.Process, error) {
// 	const (
// 		UID = 501
// 		GUID = 100
// 		)

	// cmdToRun := swaggerPath
	// args := params
	// // The Credential fields are used to set UID, GID and attitional GIDS of the process
    // // You need to run the program as  root to do this
    //     var cred =  &syscall.Credential{ UID, GUID, []uint32{} }
    // // the Noctty flag is used to detach the process from parent tty
    // var sysproc = &syscall.SysProcAttr{  Credential:cred, Noctty:true }
    // var attr = os.ProcAttr{
    //     Dir: ".",
    //     Env: os.Environ(),
    //     Files: []*os.File{
    //         os.Stdin,
    //         nil,
    //         nil,
    //     },
    //         Sys:sysproc,

    // }

	// process, err := os.StartProcess(cmdToRun, args, procAttr)

	// if err != nil {
	// 	fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	// 	return nil, err
	// } //else {
	// // 	fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
	// // }

	// return process, nil
//}