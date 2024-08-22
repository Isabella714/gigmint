package component

import (
	"github.com/spf13/viper"
)

var (
	ctx *ApplicationContext
)

func init() {
	ctx = NewApplicationContext()
}

type ApplicationInitEventListener interface {
	BeforeInit() error
	AfterInit(applicationContext *ApplicationContext) error
}

func RegisterApplicationInitEventListener(listener ApplicationInitEventListener) {
	ctx.initEventListeners = append(ctx.initEventListeners, listener)
}

type ApplicationStopEventListener interface {
	BeforeStop()
	AfterStop()
}

func RegisterApplicationStopEventListener(listener ApplicationStopEventListener) {
	ctx.stopEventListeners = append(ctx.stopEventListeners, listener)
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
}

type ApplicationConfig struct {
	config       *viper.Viper //
	serverConfig *ServerConfig
}

type ApplicationContext struct {
	config             *ApplicationConfig
	initEventListeners []ApplicationInitEventListener
	stopEventListeners []ApplicationStopEventListener
}

func NewApplicationContext() *ApplicationContext {
	return &ApplicationContext{
		config:             &ApplicationConfig{},
		initEventListeners: []ApplicationInitEventListener{},
		stopEventListeners: []ApplicationStopEventListener{},
	}
}

func Init() error {
	return ctx.Init()
}

func (ctx *ApplicationContext) Init() error {
	for _, listener := range ctx.initEventListeners {
		err := listener.BeforeInit()
		if err != nil {
			return err
		}
	}

	err := ctx.loadConfig()
	if err != nil {
		return err
	}

	for _, listener := range ctx.initEventListeners {
		err := listener.AfterInit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx *ApplicationContext) loadConfig() error {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	ctx.config.config = viper.New()

	ctx.config.config.AutomaticEnv()
	ctx.config.config.AllowEmptyEnv(true)
	ctx.config.config.SetConfigName("config")
	ctx.config.config.AddConfigPath(".")

	err := ctx.config.config.ReadInConfig()
	if err != nil {
		return err
	}

	err = ctx.config.config.UnmarshalKey("server", &ctx.config.serverConfig)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *ApplicationContext) UnmarshalKey(key string, rawVal interface{}) (err error) {
	err = ctx.config.config.UnmarshalKey(key, rawVal)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigString(key string) (string, bool) {
	if ctx.config == nil {
		return "", false
	}

	if ctx.config.config.IsSet(key) {
		return ctx.config.config.GetString(key), true
	}

	return "", false
}
