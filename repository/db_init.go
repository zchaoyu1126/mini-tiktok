package repository

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"fmt"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
)

//配置参数映射结构体
type DBConfig struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Port     string `yaml:"post"`
}

type MysqlDao struct {
	db *gorm.DB
}

//指定驱动
const DRIVER = "mysql"

var mysqlDao *MysqlDao
var once sync.Once

//获取配置参数数据
func (c *DBConfig) getConf() *DBConfig {
	// 读取config/application.yaml文件
	yamlFile, err := ioutil.ReadFile("config/application.yaml")
	// 若出现错误，打印错误提示
	if err != nil {
		log.Fatal(err)
	}
	//将读取的字符串转换成结构体conf
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

//初始化连接数据库，生成可操作基本增删改查结构的变量
func InitMySql() *gorm.DB {
	var c DBConfig
	//获取yaml配置参数
	conf := c.getConf()
	//将yaml配置参数拼接成连接数据库的url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Url,
		conf.Port,
		conf.DbName,
	)
	//连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// 建立一个连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(20)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 测试是否能访问数据库
	if sqlDB.Ping() != nil {
		sqlDB.Close()
		log.Fatal(err)
	}
	return db
}

func NewDaoInstance() *MysqlDao {
	once.Do(
		func() {
			db := InitMySql()
			mysqlDao = &MysqlDao{db}
		})
	return mysqlDao
}

// 数据库关闭操作被gormv2废除了
// 关闭数据库连接
// func Close() {
// 	db.Close()
// }
