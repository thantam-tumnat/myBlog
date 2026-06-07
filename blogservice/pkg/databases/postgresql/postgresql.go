package postgresql

import (
	"fmt"
	"blogservice/configs"
	"blogservice/modules/logs"
	"blogservice/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQLDBConnection(cfg *configs.Configs) (*gorm.DB, error) {
	postgresUrl, _, err := utils.ConnectionUrlBuilder("postgresql", cfg)
	if err != nil {
		return nil, err
	}
	dial := postgres.Open(postgresUrl)
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		logs.Error(fmt.Sprintf("error, can't connect to database, %s", err.Error()))
		return nil, err
	}

	table := fmt.Sprintf("set search_path=%v", cfg.PostgreSQL.Schema)
	db.Exec(table)

	logs.Info("postgreSQL database has been connected 🐘")
	return db, nil
}
