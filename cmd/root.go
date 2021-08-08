package cmd

import (
	"fmt"
	"log"
	"os"
	"quant/internal/app/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "quant",
		Short: "A binance quant tool",
		Long:  `Quant is a tool of binance.`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("root Execute err: %+v", errors.WithStack(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config/config.yaml or $HOME/quant/config.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <ADDRESS>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	InitConfig(cfgFile)
}

// InitConfig returns the configuration extracted from env variables or config file.
func InitConfig(configFile string) {
	// path, _ := os.Getwd()
	// log.Printf("InitConfig path => %s", path)
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(fmt.Sprintf("%s%c%s", home, os.PathSeparator, "quant"))
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
	log.Printf("Using config file: %s", viper.ConfigFileUsed())

	err = viper.Unmarshal(&config.C)
	cobra.CheckErr(err)
	log.Printf("InitConfig: %#v", config.C)
}
