package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func create(getHost func() (string, bool)) func(cmd *cobra.Command, args []string) {

	return func(cmd *cobra.Command, args []string) {
		host, isFromFlag := getHost()

		apiKey, err := makeRequest(host, http.MethodPost, "create", nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		if isFromFlag {
			viper.Set(viperHost, host)
		}

		viper.Set(viperApiKey, string(apiKey))

		fmt.Println(apiKey)
	}
}
