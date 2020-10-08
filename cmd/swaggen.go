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

	"github.com/gobuffalo/packr"
)

var swaggerPath string
var apicWinTempPath string = filepath.Join(os.Getenv("APPDATA"), "apic")
var swaggerExeTempPath string = filepath.Join(apicWinTempPath, "swagger.exe")
var swaggerYmlTempPath string = filepath.Join(apicWinTempPath, "gotemplate-swagger.yml")
var box packr.Box

func serveSwaggerDocs(pexit <-chan bool, apiContext RestApiContext) (error) {


	box = packr.NewBox("./swag")
	
	if runtime.GOOS == "windows" {
		swaggerPath = apicWinTempPath
	} else {
		swaggerPath = apicLinuxTempPath
	}

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
			//TODO: log
			fmt.Println(err)
			return err
		}
	}

	//get template from packr packed file
	swagYmlStr, _ := box.FindString("gotemplate-swagger.tpl")

	t, err := template.New("swagger").Parse(swagYmlStr)
	if err != nil {
		fmt.Println(err)
		//TODO: log
		return err
	}
	fmt.Println(swagYmlStr)

	file, err := os.Create(swaggerYmlTempPath)
	if err != nil {
		//TODO: log
		fmt.Println(err)
		return err
	}

	 eerr := t.Execute(file, cmdCon)
	 //TODO: log
	 fmt.Println(eerr)

	return nil
}

func execSwagger(pexit <-chan bool, apiContext RestApiContext) (error) {
	sbyte, _ := box.Find("swagger.exe")

	if !fileDirExists(swaggerExeTempPath) {
		ioutil.WriteFile(swaggerExeTempPath, sbyte, 0755)
	}

	swagexecCmd := exec.Command("swagger", "serve", "-p", apiContext.swaggerPort, "-F=swagger", "./gotemplate-swagger.yml")
	swagexecCmd.Dir = swaggerPath
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	swagexecCmd.Stdout = &out
	swagexecCmd.Stderr = &stderr

	serr := swagexecCmd.Start()

	if serr != nil {
		fmt.Println(fmt.Sprint(serr) + ": " + stderr.String())
		return serr
	}
	var val bool = <-pexit
	fmt.Println(val)

	go func() {
		//for {
		select {
			case <- pexit: //on cli exits kill swagger.exe
				swagexecCmd.Process.Kill()
				os.Exit(1)

			//time.Sleep(2 * time.Second)
		}
		//}
	}()

	
	
	return nil
}

func fileDirExists(filedirpath string) bool {
    if _, err := os.Stat(filedirpath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}