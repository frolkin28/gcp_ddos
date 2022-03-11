package lib

import (
	"errors"
	"flag"
)

const (
	CREATE = "create"
	DELETE = "delete"
	STOP   = "stop"
)

type IntputParams struct {
	Help       bool
	ProjectID  string
	ApiKeyPath string
	Url        string
	Duration   string
	Command    string
}

func getError() (IntputParams, error) {
	return IntputParams{}, errors.New("pass correct flags values, see --help")
}

func GetInputFlags() (IntputParams, error) {
	help := flag.Bool("help", false, "help")
	projectID := flag.String("pid", "", "GCP project_id")
	apiKey := flag.String("key", "", "GCP api key")
	url := flag.String("url", "", "GCP api key")
	duration := flag.String("d", "", "duration of attack")
	command := flag.String("command", "", "command to send action to GCP")
	flag.Parse()

	if *help {
		return IntputParams{Help: true}, nil
	} else if *command == DELETE {
		if *projectID != "" && *apiKey != "" {
			return IntputParams{Command: *command, ApiKeyPath: *apiKey, ProjectID: *projectID}, nil
		}
	} else if *command == STOP {
		if *projectID != "" && *apiKey != "" {
			return IntputParams{Command: *command, ApiKeyPath: *apiKey, ProjectID: *projectID}, nil
		}
	} else if *command == CREATE {
		if *projectID != "" && *apiKey != "" && *duration != "" && *url != "" {
			return IntputParams{
				Command:    *command,
				ApiKeyPath: *apiKey,
				Duration:   *duration,
				Url:        *url,
				ProjectID:  *projectID,
			}, nil
		}
	}

	return getError()
}
