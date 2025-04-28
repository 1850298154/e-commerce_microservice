package conf

import (
	"log"
	"os"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"github.com/spf13/viper"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/validator.v2"
)

var (
	conf     *Config
	once     sync.Once
	BasePath string
)

type Config struct {
	Env           string
	Hertz         Hertz         `mapstructure:"hertz"`
	MySQL         MySQL         `mapstructure:"mysql"`
	Redis         Redis         `mapstructure:"redis"`
	OpenTelemetry OpenTelemetry `mapstructure:"open_telemetry"`
	HealthCheck   HealthCheck   `mapstructure:"health_check"`
	Security      Security      `mapstructure:"security"`
}

type MySQL struct {
	Host            string `mapstructure:"db_host"`
	Port            int    `mapstructure:"db_port"`
	User            string `mapstructure:"db_user"`
	Password        string `mapstructure:"db_password"`
	DBName          string `mapstructure:"db_name"`
	DSN             string `mapstructure:"dsn"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 秒
}

type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	DB       int    `mapstructure:"db"`
}

type OpenTelemetry struct {
	Endpoint string `mapstructure:"endpoint"`
}

type Hertz struct {
	Service         string `mapstructure:"service"`
	Address         string `mapstructure:"address"`
	EnablePprof     bool   `mapstructure:"enable_pprof"`
	EnableGzip      bool   `mapstructure:"enable_gzip"`
	EnableAccessLog bool   `mapstructure:"enable_access_log"`
	LogLevel        string `mapstructure:"log_level"`
	LogFileName     string `mapstructure:"log_file_name"`
	LogMaxSize      int    `mapstructure:"log_max_size"`
	LogMaxBackups   int    `mapstructure:"log_max_backups"`
	LogMaxAge       int    `mapstructure:"log_max_age"`
	RegistryAddr    string `mapstructure:"registry_addr"`
}

type HealthCheck struct {
	Addr string `mapstructure:"addr"`
}

type Security struct {
	PublicRoutes []string `mapstructure:"public_routes"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// 获取项目根目录
	// _, filename, _, _ := runtime.Caller(0)
	// BasePath = filepath.Join(filepath.Dir(filename), "..")
	//
	// prefix := "conf"
	// var confFileRelPath string
	// if env := GetEnv(); env != "online" {
	//     confFileRelPath = filepath.Join(BasePath, prefix, filepath.Join(env, "conf.yaml"))
	// } else {
	//     confFileRelPath = filepath.Join(prefix, filepath.Join(env, "conf.yaml"))
	// }
	//
	// content, err := os.ReadFile(confFileRelPath)
	// if err != nil {
	//     panic(err)
	// }
	//
	// conf = new(Config)
	// err = yaml.Unmarshal(content, conf)
	// if err != nil {
	//     hlog.Error("parse yaml error - %v", err)
	//     panic(err)
	// }
	// if err := validator.Validate(conf); err != nil {
	//     hlog.Error("validate config error - %v", err)
	//     panic(err)
	// }
	//
	// conf.Env = GetEnv()
	//
	// _, _ = pretty.Printf("%+v\n", conf)

	// 从etcd中获取配置信息
	_ = godotenv.Load()
	registryAddress := "http://" + os.Getenv("REGISTRY_ADDR")
	configPath := "/config/gateway/dev"
	runtimeConfig := viper.New()
	err := runtimeConfig.AddRemoteProvider("etcd3", registryAddress, configPath)
	if err != nil {
		return
	}
	runtimeConfig.SetConfigType("yaml")
	err = runtimeConfig.ReadRemoteConfig()
	if err != nil {
		log.Fatalln("viper read:", err)
	}
	err = runtimeConfig.WatchRemoteConfigOnChannel()
	if err != nil {
		log.Fatalln("viper watch err:", err)
	}

	conf = new(Config)
	err = runtimeConfig.Unmarshal(conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		klog.Error("validate config error - %v", err)
		panic(err)
	}
	conf.Env = GetEnv()

	_, _ = pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}

func LogLevel() hlog.Level {
	level := GetConf().Hertz.LogLevel
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
