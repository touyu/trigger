package cmd

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "trigger",
	Short: "",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("AAAAAAAAAAAAAAAAAAA")
	},
}

var (
	Env string
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVarP(&Env, "env", "e", "development", "Environment")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}