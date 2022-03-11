package lib

import (
	"fmt"
	"regexp"
)

func getStartUpSript(url, duration string) string {
	template := fmt.Sprintf(
		"#! /bin/bash\n"+
			"sudo apt update\n"+
			"sudo apt install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common\n"+
			"sudo curl -fsSL https://download.docker.com/linux/debian/gpg |  apt-key add -\n"+
			"sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable\"\n"+
			"sudo apt update\n"+
			"sudo apt install -y docker-ce\n"+
			"sudo docker run -d alpine/bombardier -c 1000 -d %v -l %v", duration, url)
	return template
}

func getZonesList() [12]string {
	zones := [12]string{
		"europe-central2-a",
		"europe-central2-b",
		"europe-north1-a",
		"europe-north1-b",
		"asia-east2-a",
		"asia-east2-b",
		"asia-east2-c",
		"asia-east1-a",
		"asia-east1-b",
		"asia-east1-c",
		"asia-south1-a",
		"asia-southeast1-b",
	}
	return zones
}

func extractZoneFromUrl(zoneUrl string) string {
	r := regexp.MustCompile(
		`https://www.googleapis.com/compute/v\d+/projects/[0-9A-Za-z-_]+/zones/(?P<zone>[0-9A-Za-z-]+)`,
	)
	match := r.FindStringSubmatch(zoneUrl)
	index := r.SubexpIndex("zone")
	return match[index]
}
