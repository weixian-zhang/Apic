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

	"github.com/spf13/cobra"
)

var configStr = `port: "8080"
rest:
- path: "/api/store/order/{orderid}"
  querystring: "" #optional
  response: "pickled cucumber" #simple strings or Json data
  headers: #optional
  - content-type=application/json
  - custom=customValue
  cookies: #default expiry 1hr, not configurable 

- path: "/api/order/inventory"
  querystring: "" #optional
  response: "{\"Product ID\":7615,\"SKU\":\"HEH-2245\",\"Name\":\"Simply Sweet Blouse\",\"Product URL\":\"https:\/\/www.domain.com\/product\/heh-2245\",\"Price\":42,\"Retail Price\":59.95,\"Thumbnail URL\":\"https:\/\/www.domain.com\/images\/heh-2245_600x600.png\",\"Search Keywords\":\"lorem, ipsum, dolor, ...\",\"Description\":\"Sociosqu facilisis duis ...\",\"Category\":\"Clothing>Tops>Blouses|Clearance|Tops On Sale\",\"Category ID\":\"285|512|604\",\"Brand\":\"Entity Apparel\",\"Child SKU\":\"HEH-2245-RSWD-SM|HEH-2245-RSWD-MD|HEH-2245-THGR-SM|EH-2245-THGR-MD|HEH-2245-DKCH-SM|HEH-2245-DKCH-MD\",\"Child Price\":\"42|59.99\",\"Color\":\"Rosewood|Thyme Green|Dark Charcoal\",\"Color Family\":\"Red|Green|Grey\",\"Color Swatches\":\"[{\\\"color\\\":\\\"Rosewood\\\", \\\"family\\\":\\\"Red\\\", \\\"swatch_hex\\\":\\\"#65000b\\\", \\\"thumbnail\\\":\\\"\/images\/heh-2245-rswd-sm_600x600.png\\\", \\\"price\\\":42}, {\\\"color\\\":\\\"Thyme Green\\\", \\\"family\\\":\\\"Green\\\", \\\"swatch_img\\\":\\\"\/swatches\/thyme_green.png\\\", \\\"thumbnail\\\":\\\"\/images\/heh-2245-thgr-sm_600x600.png\\\", \\\"price\\\":59.99}, {\\\"color\\\":\\\"Dark Charcoal\\\", \\\"family\\\":\\\"Grey\\\", \\\"swatch_hex\\\":\\\"#36454f\\\", \\\"thumbnail\\\":\\\"\/images\/heh-2245-dkch-sm_600x600.png\\\", \\\"price\\\":59.99}]\",\"Size\":\"Small|Medium\",\"Shoe Size\":\"\",\"Pants Size\":\"\",\"Occassion\":\"\",\"Season\":\"Summer|Spring\",\"Badges\":\"Exclusive|Clearance\",\"Rating Avg\":4.5,\"Rating Count\":10,\"Inventory Count\":8,\"Date Created\":\"2018-03-20 22:24:21\"}"
  headers: #optional
  - content-type=application/json
  - correlationId: 91a956ab-0e76-42b5-b506-981d73ac3e0d
  cookies: #default expiry 1hr, not configurable
  - approle: "user"

- path: "/api/user/list"
  querystring: "" #optional
  response: "{ \"email\": \"liam.walters@example.com\", \"gender\": \"male\", \"phone_number\": \"0438-376-652\", \"birthdate\": 826530877, \"location\": { \"street\": \"9156 dogwood ave\", \"city\": \"devonport\", \"state\": \"australian capital territory\", \"postcode\": 7374 }, \"username\": \"biglion964\", \"password\": \"training\", \"first_name\": \"liam\", \"last_name\": \"walters\", \"title\": \"mr\" }"
  headers: #optional
  - content-type=application/json
  - correlationId: 91a956ab-0e76-42b5-b506-981d73ac3e0d
  cookies: #default expiry 1hr, not configurable 
  - username: "Jonh Smith"
  - approle: "admin"`

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates sample config file",
	Long: `Generates a sample apic config file with APIs that is directly importable.
You can modify config to suit your needs and build up a readily available Microservice for testing`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	genCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
