package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show triggers of database",
	Long:  "Show triggers of database",
	Run:   showList,
}

func showList(cmd *cobra.Command, args []string) {
	dbConfig := viper.GetStringMapString(Env)
	datasource := dbConfig["datasource"]
	fmt.Println("Show triggers of", Env)
	db := connectDB(datasource)
	defer db.Close()

	databaseName, err := getDatabaseName(db)
	if err != nil {
		log.Fatal(err)
	}

	tables, err := getTables(db, databaseName)
	if err != nil {
		log.Fatal(err)
	}

	triggers, err := getTriggers(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {
		exists := existsTrigger(table, triggers)
		if exists {
			fmt.Println(fmt.Sprintf("Already exists set_logical_uniqueness_on_%s", table))
		} else {
			fmt.Println(fmt.Sprintf("Not exists set_logical_uniqueness_on_%s", table))
		}
	}
}
