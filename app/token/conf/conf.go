package conf

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/validator.v2"
)

var (
	conf *Config
	once sync.Once
)

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

type Config struct {
	Env           string
	Kitex         Kitex         `mapstructure:"kitex"`
	MySQL         MySQL         `mapstructure:"mysql"`
	OpenTelemetry OpenTelemetry `mapstructure:"open_telemetry"`
	Redis         Redis         `mapstructure:"redis"`
	Registry      Registry      `mapstructure:"registry"`
	JWT           JWT           `mapstructure:"jwt"`
	HealthCheck   HealthCheck   `mapstructure:"health_check"`
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
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Kitex struct {
	Service        string `mapstructure:"service"`
	Address        string `mapstructure:"address"`
	LogLevel       string `mapstructure:"log_level"`
	LogFileName    string `mapstructure:"log_file_name"`
	LogMaxSize     int    `mapstructure:"log_max_size"`
	LogMaxBackups  int    `mapstructure:"log_max_backups"`
	LogMaxAge      int    `mapstructure:"log_max_age"`
	MaxConnections int    `mapstructure:"max_connections"`
	MaxQPS         int    `mapstructure:"max_qps"`
}

type Registry struct {
	RegistryAddress []string `mapstructure:"registry_address"`
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
}

type JWT struct {
	SigningKey string `mapstructure:"key"`
}

type OpenTelemetry struct {
	Endpoint string `mapstructure:"endpoint"`
}

type HealthCheck struct {
	Addr string `mapstructure:"addr"`
}

func initConf() {
	// _, filename, _, _ := runtime.Caller(0)
	// basePath := filepath.Join(filepath.Dir(filename), "..")
	//
	// prefix := "conf"
	// var confFileRelPath string
	// if env := GetEnv(); env != "online" {
	//     confFileRelPath = filepath.Join(basePath, prefix, filepath.Join(env, "conf.yaml"))
	// } else {
	//     confFileRelPath = filepath.Join(prefix, filepath.Join(env, "conf.yaml"))
	// }
	// content, err := os.ReadFile(confFileRelPath)
	// if err != nil {
	//     panic(err)
	// }
	// conf = new(Config)
	// err = yaml.Unmarshal(content, conf)
	// if err != nil {
	//     klog.Error("parse yaml error - %v", err)
	//     panic(err)
	// }
	// if err := validator.Validate(conf); err != nil {
	//     klog.Error("validate config error - %v", err)
	//     panic(err)
	// }
	// conf.Env = GetEnv()
	// if _, err := pretty.Printf("%+v\n", conf); err != nil { // 检查 pretty.Printf 的错误返回值
	//     klog.Error("pretty.Printf error - %v", err)
	// }

	_ = godotenv.Load()
	registryAddress := "http://" + os.Getenv("REGISTRY_ADDR")
	configPath := "/config/auth/dev"
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

func LogLevel() klog.Level {
	level := GetConf().Kitex.LogLevel
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
