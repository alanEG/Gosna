package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func config_handling(action string) error {
	var err error
	var json_data []byte
	if action == "load" {
		config_data, err := ioutil.ReadFile(configFile)

		logger("error", "[Error]", err, err, 1)

		err = json.Unmarshal(config_data, &data)

		logger("error", "[Error]", err, err, 1)

	} else if action == "save" {
		json_data, err = json.MarshalIndent(data, "", "  ")

		ioutil.WriteFile(configFile, json_data, 0644)

		Execute_command("cd " + data.Config.Directory_work + ";git add .")

	}

	return err
}

//this for setup the config file and folder result and work directory
func setupFirst() {
	if data.Config.Is_first == true {
		Execute_command(fmt.Sprintf("mkdir %s;mkdir %s;cd %s;git init",
			data.Config.Directory_result,
			data.Config.Directory_work,
			data.Config.Directory_work,
		),
		)
		data.Config.Is_first = false
	}
}
