package infra

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DatabaseConfig struct {
	Database   string        `mapstructure:"database"`
	Host       string        `mapstructure:"host"`
	Port       string        `mapstructure:"port"`
	Username   string        `mapstructure:"username"`
	Password   string        `mapstructure:"password"`
	Timeout    time.Duration `mapstructure:"timeout"`
	Connection Connection    `mapstructure:"connection"`
}

type Connection struct {
	Maxidle int `json:"maxidle"`
	Maxopen int `json:"maxopen"`
}

var TransactionDB *gorm.DB

func InitDb() *gorm.DB {
	cfg := NewDBConfig("db.transaction")
	dsn := genDSN(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connection db failed")
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.Connection.Maxidle)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.Connection.Maxopen)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(cfg.Timeout)
	return db
}

func genDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)
}

func NewDBConfig(db string) *DatabaseConfig {
	var cfg DatabaseConfig
	err := viper.UnmarshalKey(db, &cfg)
	if err != nil {
		fmt.Println("errors when call configs")
		panic(err)
	}
	return &cfg
}
