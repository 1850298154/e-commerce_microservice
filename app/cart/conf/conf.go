package conf

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
)

var (
	conf     *Config
	once     sync.Once
	BasePath string
)

type Config struct {
	Env           string        `mapstructure:"env"`
	Kitex         Kitex         `mapstructure:"kitex"`
	MySQL         MySQL         `mapstructure:"mysql"`
	Redis         Redis         `mapstructure:"redis"`
	Registry      Registry      `mapstructure:"registry"`
	Opentelemetry OpenTelemetry `mapstructure:"opentelemetry"`
}

type MySQL struct {
	DSN string `mapstructure:"dsn"`
}

type Redis struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Kitex struct {
	Service       string `mapstructure:"service"`
	Address       string `mapstructure:"address"`
	LogLevel      string `mapstructure:"log_level"`
	LogFileName   string `mapstructure:"log_file_name"`
	LogMaxSize    int    `mapstructure:"log_max_size"`
	LogMaxBackups int    `mapstructure:"log_max_backups"`
	LogMaxAge     int    `mapstructure:"log_max_age"`
}

type Registry struct {
	RegistryAddress []string `mapstructure:"registry_address"`
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
}

type OpenTelemetry struct {
	Endpoint string `mapstructure:"endpoint"`
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
	// 	confFileRelPath = filepath.Join(BasePath, prefix, filepath.Join(env, "conf.yaml"))
	// } else {
	// 	confFileRelPath = filepath.Join(prefix, filepath.Join(env, "conf.yaml"))
	// }
	//
	// fmt.Println(confFileRelPath)
	// content, err := os.ReadFile(confFileRelPath)
	// if err != nil {
	// 	panic(err)
	// }
	// conf = new(Config)
	// err = yaml.Unmarshal(content, conf)
	// if err != nil {
	// 	klog.Error("parse yaml error - %v", err)
	// 	panic(err)
	// }
	// if err := validator.Validate(conf); err != nil {
	// 	klog.Error("validate config error - %v", err)
	// 	panic(err)
	// }
	// conf.Env = GetEnv()
	// _, _ = pretty.Printf("%+v\n", conf)

	// 从etcd中读取配置
	_ = godotenv.Load()
	registryAddress := "http://" + os.Getenv("REGISTRY_ADDR")
	configPath := "/config/cart/dev"
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
