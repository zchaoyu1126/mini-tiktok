package db

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

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
func getMysqlConf() *MysqlConfig {
	var c MysqlConfig

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

// connect mysql database
func InitMySql() *gorm.DB {
	conf := getMysqlConf()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// test the connect
	if sqlDB.Ping() != nil {
		sqlDB.Close()
		log.Fatal(err)
	}
	log.Println("mysql init success.")
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

// this method was removed in gormv2
// close the database
// func Close() {
// 	db.Close()
// }
