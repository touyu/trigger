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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var (
	Env string
)

func init() {
	viper.SetConfigName("dbconfig")
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
