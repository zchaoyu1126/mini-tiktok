package db

import (
	"io/ioutil"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fmt"

	"gopkg.in/yaml.v2"
)

// Mysql Configure Struct
type MysqlConfig struct {
	Host     string `yaml:"mysql_host"`
	UserName string `yaml:"mysql_username"`
	Password string `yaml:"mysql_password"`
	DbName   string `yaml:"mysql_dbname"`
	Port     string `yaml:"mysql_port"`
}

// Mysql Data Access Object
type MySQLConn struct {
	DB *gorm.DB
}

const DRIVER = "mysql"

var mysqlConn *MySQLConn
var mysqlOnce sync.Once

// Get Mysql configure from yaml.
// If read yaml file failed or unmarshal failed, system exit instantly.
func readMysqlConf() *MysqlConfig {
	config := new(MysqlConfig)
	file, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(file, config)
	return config
}

// connect mysql database
func InitMySql() *gorm.DB {
	conf := readMysqlConf()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("mysql database init failed")
	}

	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("mysql database init failed")
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// test the connect
	if sqlDB.Ping() != nil {
		sqlDB.Close()
		zap.L().Fatal("mysql database init failed")
	}
	zap.L().Info("mysql database init success")
	return db
}

// return mysql dao with singleten pattern
func NewMySQLConnInstance() *MySQLConn {
	mysqlOnce.Do(
		func() {
			db := InitMySql()
			mysqlConn = &MySQLConn{db}
		})
	return mysqlConn
}
