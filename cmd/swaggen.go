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

func genSwaggerDocs(apicon cmdContext) (*exec.Cmd) {
	box = packr.NewBox("./swag")
	sbyte, _ := box.Find("swagger.exe")
	
	if runtime.GOOS == "windows" {
		swaggerPath = apicWinTempPath
	} else {
		swaggerPath = apicLinuxTempPath
	}

	if !fileDirExists(swaggerPath) {
		os.Mkdir(swaggerPath, 0755)
	}

	if !fileDirExists(swaggerExeTempPath) {
		ioutil.WriteFile(swaggerExeTempPath, sbyte,0755)
		//fmt.Println(ioerr)
	}

	genSwaggerYml(apicon)

	// swagArgs := fmt.Sprintf("serve -p %v -F=swagger %v", defaultSwaggerPort, swaggerYmlTempPath)
	// fmt.Println(swagArgs)
	swagexecCmd := exec.Command("swagger", "serve", "-p", "8090", "-F=swagger", "./gotemplate-swagger.yml") //exec.Command(swaggerExeTempPath, "serve", "-p 8090", "-F=swagger", swaggerYmlTempPath)
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

	return swagexecCmd
}

func genSwaggerYml(cmdCon cmdContext) (error) {
	swagYmlStr, _ := box.FindString("gotemplate-swagger.tpl")

	t, err := template.New("swagger").Parse(swagYmlStr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(swagYmlStr)

	 eerr := t.Execute(os.Stdout, cmdCon)
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