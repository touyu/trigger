package cmd

import (
	"fmt"
	driver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

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

func mysqlBuildDSN(dbName string) string {
	c := driver.NewConfig()
	c.User = user
	c.Passwd = password
	c.DBName = dbName
	c.Net = "tcp"
	c.Addr = fmt.Sprintf("%s:%d", host, port)
	return c.FormatDSN()
}

func openDB(datasource string) *sqlx.DB {
	db, err := sqlx.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
