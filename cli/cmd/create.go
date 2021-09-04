package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func create(cmd *cobra.Command, args []string) {

	apiKey, err := makeRequest(http.MethodPost, "create", nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	viper.Set(viperApiKey, string(apiKey))

	fmt.Println(apiKey)
}
