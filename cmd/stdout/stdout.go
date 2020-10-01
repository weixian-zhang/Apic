// package stdout

// import (
// 	"fmt"
// 	"github.com/gookit/color"
// )

// func Print(msgFormat string, params ...interface{}) {
// 	fmt.Println(msgFormat, params)
// }

// // func PQueryString(msgFormat string, params ...interface{}) {
// // 	log.Println(yellow(msgFormat), params)
// // }

// // func PHeader(msgFormat string, params ...interface{}) {
// // 	log.Println(green(msgFormat), params)
// // }

// // func PSysInfo(msgFormat string, params ...interface{}) {
// // 	log.Println(cyan(msgFormat), params)
// // }

// func Cyan(msg string) (string) {
// 	return color.FgCyan.Render(msg)
// }

// func Yellow(msg string) (string) {
// 	return color.FgLightYellow.Render(msg)
// }

// func Green(msg string) (string) {
// 	return color.FgLightGreen.Render(msg)
// }

// func Gray(msg string) (string) {
// 	return color.FgGray.Render(msg)
// }