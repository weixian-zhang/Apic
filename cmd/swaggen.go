package cmd

import (
	"io/ioutil"
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"text/template"
	"path/filepath"
	"runtime"
	"github.com/gobuffalo/packr"
)

var swaggerPath string
var apicWinTempPath string = filepath.Join(os.Getenv("APPDATA"), "apic")
var swaggerExeTempPath string = filepath.Join(apicWinTempPath, "swagger.exe")
var swaggerYmlTempPath string = filepath.Join(apicWinTempPath, "gotemplate-swagger.yml")
var box packr.Box

func genSwaggerDocs(apiContext RestApiContext) (*exec.Cmd) {
	box = packr.NewBox("./swag")
	sbyte, _ := box.Find("swagger.exe")
	
	if runtime.GOOS == "windows" {
		swaggerPath = apicWinTempPath
	} else {
		swaggerPath = apicLinuxTempPath
	}

	if !fileDirExists(swaggerPath) {
		os.Mkdir(swaggerPath, 0755)
		//TODO: log, fail to create apic temp folder
	}

	if !fileDirExists(swaggerExeTempPath) {
		ioutil.WriteFile(swaggerExeTempPath, sbyte,0755)
		//TODO: log, log fail to create swagger yml at temp folder
	}

	genSwaggerYml(apiContext) //gens swagger yaml in apic temp folder

	swagexecCmd := exec.Command("swagger", "serve", "-p", apiContext.swaggerPort, "-F=swagger", "./gotemplate-swagger.yml")
	swagexecCmd.Dir = swaggerPath
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	swagexecCmd.Stdout = &out
	swagexecCmd.Stderr = &stderr

	// serr := swagexecCmd.Start()

	// if serr != nil {
	// 	fmt.Println(fmt.Sprint(serr) + ": " + stderr.String())
	// 	return nil
	// }
	// fmt.Println("Result: " + out.String())

	//TODO: log, swagger 
	return swagexecCmd
}

//swagger editor
//https://editor.swagger.io/
func genSwaggerYml(cmdCon RestApiContext) (error) {
	swagYmlStr, _ := box.FindString("gotemplate-swagger.tpl")

	t, err := template.New("swagger").Parse(swagYmlStr)
	if err != nil {
		fmt.Println(err)
		//TODO: log
		return err
	}
	fmt.Println(swagYmlStr)

	 eerr := t.Execute(os.Stdout, cmdCon)
	 //TODO: log
	 fmt.Println(eerr)

	return nil
}

func fileDirExists(filedirpath string) bool {
    if _, err := os.Stat(filedirpath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}