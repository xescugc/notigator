package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// RootCmd is the root command for the CLI
	RootCmd = &cobra.Command{
		Use:   "notigator",
		Short: "Aggregates Notifications",
		Long:  "Aggregates Notifications",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(serveCmd)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration path to the file (default is $HOME/.notigator.*)")
}

func initConfig() {
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

		// Search config in home directory with name ".notigator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".notigator")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
