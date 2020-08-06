package cmd

import (
	"fmt"
	"os"
)

func GetInfo(infoType string) {
	switch infoType {
	case "packages":
		packageList, _ := getPackages()
		for _, p := range packageList {
			InfoPrint(p)
		}
	case "users":
		userList, _ := getLocalUsers()
		for _, u := range userList {
			InfoPrint(fmt.Sprint(u))
		}
	default:
		if infoType == "" {
			FailPrint("No info type provided.")
		} else {
			FailPrint("No info for \"" + infoType + "\" found.")
		}
		os.Exit(1)
	}
}
