package cmd

import "github.com/spf13/viper"

func getHost() (string, bool) {
	//Take host from viper or flag
	apiHost := viper.GetString(viperHost)

	if apiHost == "" {
		return host, true
	}

	return apiHost, false
}
