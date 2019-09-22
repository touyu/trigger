package cmd

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register trigger to Database",
	Long:  "Register trigger to Database",
	Run:   register,
}

func register(cmd *cobra.Command, args []string) {
	dbConfig := viper.GetStringMapString(Env)
	datasource := dbConfig["datasource"]
	fmt.Println("Register trigger for", Env)
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
			fmt.Println(fmt.Sprintf("Create set_logical_uniqueness_on_%s", table))
			err := createTrigger(db, table)
			if err != nil {
				log.Println("Error", err)
			}
		}
	}
}

func getDatabaseName(db *sqlx.DB) (string, error) {
	type Database struct {
		Name string `db:"name"`
	}

	rows, err := db.Queryx("SELECT database() as name")
	if err != nil {
		log.Fatal(err)
	}

	var database Database
	for rows.Next() {
		err := rows.StructScan(&database)
		if err != nil {
			return "", err
		}
		return database.Name, nil
	}

	return "", nil
}

func getTables(db *sqlx.DB, databaseName string) ([]string, error) {
	type Table struct {
		Name string `db:"name"`
	}

	rows, err := db.Queryx(fmt.Sprintf("SELECT TABLE_NAME name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'", databaseName))
	if err != nil {
		log.Fatal(err)
	}

	var tables = make([]string, 0)
	for rows.Next() {
		var table Table
		err := rows.StructScan(&table)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table.Name)
	}

	return tables, nil
}

func getTriggers(db *sqlx.DB) ([]string, error) {
	type Trigger struct {
		Name string `db:"name"`
	}

	rows, err := db.Queryx("select trigger_name as name from information_schema.triggers")
	if err != nil {
		log.Fatal(err)
	}

	var tables = make([]string, 0)
	for rows.Next() {
		var trigger Trigger
		err := rows.StructScan(&trigger)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, trigger.Name)
	}

	return tables, nil
}

func existsTrigger(table string, triggers []string) bool {
	triggerName := fmt.Sprintf("set_logical_uniqueness_on_%s", table)
	for _, trigger := range triggers {
		if triggerName == trigger {
			return true
		}
	}

	return false
}

func createTrigger(db *sqlx.DB, tableName string) error {
	sql := fmt.Sprintf(`
CREATE TRIGGER set_logical_uniqueness_on_%s BEFORE UPDATE ON %s FOR EACH ROW
BEGIN
  IF NEW.deleted_at IS NULL THEN
    SET NEW.logical_uniqueness = true;
  ELSE
    SET NEW.logical_uniqueness = NULL;
  END IF;
END
`, tableName, tableName)

	_, err := db.Exec(sql)
	return err
}

func connectDB(datasource string) *sqlx.DB {
	db, err := sqlx.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
