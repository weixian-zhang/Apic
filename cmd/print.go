package cmd

import (
	"errors"
	"net"
	"os"
	"strings"
	"fmt"
	"github.com/gookit/color"
)


// type pInfoApi struct {
// 	hostIP string
// 	hostName string
// 	path string
// 	querystring string
// 	headers []header
// 	resp string //json data or string
// 	swaggerPath string
// }

// type pInfoIngressRequest struct {
// 	clientIP string
// }

func printAPIsInfo(context RestApiContext) {

	var info string
	var newLine string = `
`
	host, _ := os.Hostname()
	s, _ := getLocalIP()
	serverip := color.FgCyan.Sprintf("APIs running at: %v / %v", s, host)
	info += newLine

	info += serverip
	info += newLine

	
	port := strings.TrimSpace(context.Port)

	for _, api := range context.RestApis {

		qs := api.Querystring
		if api.Querystring != "" {
			qs = "?" + qs
		}

		fqdn :=
			color.FgLightGreen.Sprintf("FQDN: http://%v:%v%v%v", host, port, formatAPIPath(api.Path), qs)

		info += fqdn
		info += newLine

		headerStr := "not specified"
		if api.headerStr != "" {
			headerStr = api.headerStr
		}
		headers :=  color.FgLightMagenta.Sprintf("Response Headers: %v", headerStr)
		info += headers
		info += newLine

		cookieStr := "not specified"
		if api.cookieStr != "" {
			cookieStr = api.cookieStr
		}
		cookies := color.FgLightYellow.Sprintf("Response Cookies: %v", cookieStr)
		info += cookies
		info += newLine
	}

	methods := color.FgGray.Sprintf("Methods: GET,POST,PUT,DELETE")

	info += methods
	info += newLine
	info += newLine

	fmt.Println(info)
}

func printIngressRequestInfo() {

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

func formatAPIPath(path string) (string) {
	var newPath string
	newPath = strings.TrimSpace(path)

	if fc := newPath[0:1]; fc != "/" {
		newPath = "/" + newPath
	}

	return newPath
}