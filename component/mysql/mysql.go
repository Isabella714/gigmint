package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Isabella714/gigmint/component"
)

var instance *MySQL

func init() {
	instance = NewMySQL()

	component.RegisterApplicationInitEventListener(instance)
	component.RegisterApplicationStopEventListener(instance)
}

type MySQLConfig struct {
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	DBName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifetime  int    `mapstructure:"max_life_time"`
}

type MySQL struct {
	ctx      *component.ApplicationContext
	config   *MySQLConfig
	instance *gorm.DB
}

func NewMySQL() *MySQL {
	return &MySQL{}
}

func (c *MySQL) Instantiate() error {
	err := c.ctx.UnmarshalKey("mysql", &c.config)
	if err != nil {
		return err
	}

	if c.config == nil {
		return errors.New("mysql config isn't found")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		c.config.User,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.DBName,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(c.config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(c.config.MaxLifetime))

	c.instance, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func Get(ctx context.Context) *gorm.DB {
	return instance.Get(ctx)
}

func (c *MySQL) Get(ctx context.Context) *gorm.DB {
	return c.instance
}

func (c *MySQL) Close() error {

	return nil
}

func (c *MySQL) BeforeInit() error {
	return nil
}

func (c *MySQL) AfterInit(applicationContext *component.ApplicationContext) error {
	c.ctx = applicationContext

	return c.Instantiate()
}

func (c *MySQL) BeforeStop() {
	return
}

func (c *MySQL) AfterStop() {
	_ = c.Close()

	return
}
