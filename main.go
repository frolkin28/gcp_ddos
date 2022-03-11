package main

import (
	"ddos/lib"
	"fmt"
)

func main() {
	params, err := lib.GetInputFlags()
	if err != nil {
		fmt.Println(err)
		return
	} else if params.Help {
		getHelp()
		return
	}
	switch {
	case params.Command == lib.CREATE:
		lib.CreateInstances(params)
	case params.Command == lib.DELETE:
		lib.DeleteAllInstances(params)
	case params.Command == lib.STOP:
		lib.StopAllInstances(params)
	default:
		fmt.Println("No such command")
	}
}

func getHelp() {
	template := `Args:
	command - one of create/stop/delete
	pid - your GCP project id
	fkey - path to GCP api key json file
	url - url to ddos
	d - duration of attack e.g. 3600s`
	fmt.Println(template)
}
