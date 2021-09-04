package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func convert(cmd *cobra.Command, args []string) {

	apiKey := viper.Get(viperApiKey)

	if apiKey == "" {
		fmt.Println("You need to input api-key flag or to create a new project!")
		return
	}

	val, err := strconv.ParseFloat(value, strconv.IntSize)

	if err != nil {
		fmt.Println(val)
		return
	}

	body := map[string]interface{}{
		"symbol": symbol,
		"value":  val,
	}

	value, err := makeRequest(http.MethodGet, "convert", body, func(r *http.Request) {
		r.Header.Add("Authorization", apiKey)
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)
}
