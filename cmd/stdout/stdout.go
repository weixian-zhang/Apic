package stdout

import (
	"log"
	"github.com/gookit/color"
)

func PInfo(msgFormat string, params ...interface{}) {
	log.Println(gray(msgFormat), params)
}

func PQueryString(msgFormat string, params ...interface{}) {
	log.Println(yellow(msgFormat), params)
}

func PHeader(msgFormat string, params ...interface{}) {
	log.Println(green(msgFormat), params)
}

func PSysInfo(msgFormat string, params ...interface{}) {
	log.Println(cyan(msgFormat), params)
}

func cyan(msg string) (string) {
	return color.FgCyan.Render(msg)
}

func yellow(msg string) (string) {
	return color.FgLightYellow.Render(msg)
}

func green(msg string) (string) {
	return color.FgLightGreen.Render(msg)
}

func gray(msg string) (string) {
	return color.FgGray.Render(msg)
}