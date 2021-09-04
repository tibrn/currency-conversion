package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func makeRequest(method, path string, body interface{}, modifiers ...func(*http.Request)) (string, error) {
	//Take host from viper or flag
	apiHost := viper.GetString(viperHost)
	isHostFromFlag := false

	if apiHost == "" {
		isHostFromFlag = true
		apiHost = host
	}

	buffBody := bytes.NewBuffer([]byte{})

	//Write body
	if body != nil {

		bytsBody, err := json.Marshal(body)

		if err != nil {
			return "", err
		}

		buffBody.Write(bytsBody)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, path), buffBody)

	if err != nil {
		return "", err
	}

	//Modify request before sending
	for _, modifier := range modifiers {
		modifier(req)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if isHostFromFlag {
		//Update persistent config from flags
		viper.Set(viperHost, host)
	}

	defer resp.Body.Close()

	//Read response
	respBody, err := ioutil.ReadAll(resp.Body)

	return string(respBody), err
}
