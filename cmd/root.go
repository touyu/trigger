package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "trigger",
	Short: "A automatic registering trigger tool",
	Long:  "A automatic registering trigger tool",
	Args:  cobra.MinimumNArgs(1),
	Run:   run,
}

var (
	user     string
	password string
	host     string
	port     uint64
	helpFlag bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "User")
	rootCmd.PersistentFlags().StringVarP(&password, "pass", "p", "", "Password")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "h", "", "Host")
	rootCmd.PersistentFlags().Uint64VarP(&port, "port", "P", 3306, "Port")
	rootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		panic("Not found database")
	}

	dbName := args[0]
	datasource := mysqlBuildDSN(dbName)
	db := openDB(datasource)
	defer db.Close()

	tables, err := getTables(db, dbName)
	if err != nil {
		log.Fatal(err)
	}

	triggers, err := getTriggers(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {
		exists := existsTrigger(table, triggers)
		if !exists {
			fmt.Println(fmt.Sprintf("Create set_logical_uniqueness_on_%s", table))
			err := createTrigger(db, table)
			if err != nil {
				log.Println("Error", err)
			}
		}
	}
}
