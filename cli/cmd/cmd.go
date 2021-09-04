package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	viperApiKey = "api-key"
	viperHost   = "host"
)

var (
	cfgFile string
	apiKey  string
	symbol  string
	value   string
	host    string
)

var cmdCreate = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr"},
	Run:     create,
}
var cmdConvert = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"conv"},
	Run:     convert,
}

var cmdRoot = &cobra.Command{
	Use: "",
}

func Execute() {

	cmdRoot.AddCommand(cmdCreate)

	cmdConvert.Flags().StringVarP(&apiKey, "api-key", "k", "", "API Key used for authentication")
	cmdConvert.Flags().StringVarP(&symbol, "symbol", "s", "", "Symbol for conversion (e.g. EUR/USD)")
	cmdConvert.Flags().StringVarP(&value, "value", "v", "", "Value to be converted (e.g 1.7658)")
	cmdConvert.MarkFlagRequired("symbol")
	cmdConvert.MarkFlagRequired("value")
	cmdRoot.AddCommand(cmdConvert)

	cmdRoot.PersistentFlags().StringVar(&host, "host", "http://127.0.0.1:8081", "Host for converter API")

	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	cmdRoot.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.converter.yaml)")

}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".converter")
		viper.SetConfigType("yaml")

		configPath := filepath.Join(home, ".converter.yaml")

		_, err = os.Stat(configPath)
		if !os.IsExist(err) {
			if _, err := os.Create(configPath); err != nil { // perm 0666

				fmt.Println("Can't create config:", err)
				os.Exit(1)

			}
		}

	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
