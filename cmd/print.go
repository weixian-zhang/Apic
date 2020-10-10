package cmd

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

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
	serverip := color.FgCyan.Sprintf("APIs running on: %v/%v", s, host)
	info += newLine

	info += serverip
	info += newLine

	
	port := strings.TrimSpace(context.Port)
	swagPort := strings.TrimSpace(context.swaggerPort)

	for _, api := range context.RestApis {

		qs := api.Querystring
		if api.Querystring != "" {
			qs = "?" + qs
		}

		apiAddr :=
			color.FgLightGreen.Sprintf("Api: http://%v:%v%v%v", host, port, formatAPIPath(api.Path), qs)
		swagUIAddr :=
			color.FgLightGreen.Sprintf("Swagger UI: http://%v:%v%v", host, swagPort, "/docs")

		swagJsonAddr :=
			color.FgLightGreen.Sprintf("Swagger Json: http://%v:%v%v", host, swagPort, "/swagger.json")

		info += apiAddr
		info += newLine

		info += swagUIAddr
		info += newLine

		info += swagJsonAddr
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

func printIngreReqInfo(r *http.Request) {
	clientIP := readSourceIP(r)
	fmt.Println(color.LightWhite.Sprintf("received at: %v, from client: %v", time.Now().Format("15:04"), clientIP))
}

func printErr(err error) {
	if err != nil {
		color.FgRed.Println(err.Error())
	}
}
func printInfo(msg string) {
	color.FgLightGreen.Println(msg)
}

func readSourceIP(r *http.Request) string {
    IPAddress := r.Header.Get("X-Real-Ip")
    if IPAddress == "" {
        IPAddress = r.Header.Get("X-Forwarded-For")
    }
    if IPAddress == "" {
        IPAddress = r.RemoteAddr
    }
    return IPAddress
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