package migration

import (
	"context"
	"fmt"
	"os"

	"github.com/SpectreFury/odin-book/backend/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// have a migration folder with cwith tables
// picks all the tables one after another
// has a single sql file inside which it parses
// runs it using conn.Query(context.Background(), query)


func loadFileNames(path string) []string {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	var fileNames []string
	for _, file := range dir {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames
}

func RunMigration(conn *pgxpool.Pool, path string) error {
	logger := logger.Logger {}
	sqlFiles := loadFileNames(path)

	for _, sqlFile := range sqlFiles {
		data, err := os.ReadFile(path + "/" + sqlFile)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		sqlQuery := string(data)

		rows, err := conn.Query(context.Background(), sqlQuery)
		if err != nil {

			return err
		}
		defer rows.Close()
	} 

	logger.Log("Migration completed")
	return nil
}
