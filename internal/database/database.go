package database

import (
	"database/sql"
	"fmt"
	"time"

	"kitchen-api/internal/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrmDb struct {
	OrmInstance *gorm.DB
	DB          *sql.DB
}

func OpenPostgresORM(host string, port int, username string, password string, dbName string) (OrmDb, error) {
	sslMode := config.Default().GetString("db.sslmode")
	if sslMode == "" {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		username,
		password,
		dbName,
		sslMode,
	)

	if param := config.Default().GetString("db.param"); param != "" {
		dsn += " " + param
	}

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return OrmDb{}, fmt.Errorf("failed to open sql db: %w", err)
	}

	if connMaxLifetime := config.Default().GetInt("db.connMaxLifeTime"); connMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	}
	if maxOpenConn := config.Default().GetInt("db.maxOpenConn"); maxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConn)
	}
	if maxIdleConn := config.Default().GetInt("db.maxIdleConn"); maxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(maxIdleConn)
	}

	if err := sqlDB.Ping(); err != nil {
		return OrmDb{}, fmt.Errorf("failed to ping database: %w", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return OrmDb{}, fmt.Errorf("failed to open gorm db: %w", err)
	}

	return OrmDb{
		OrmInstance: gormDB,
		DB:          sqlDB,
	}, nil
}
