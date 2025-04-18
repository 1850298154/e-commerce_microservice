package conf

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env           string
	Kitex         Kitex         `mapstructure:"kitex"`
	MySQL         MySQL         `mapstructure:"mysql"`
	Redis         Redis         `mapstructure:"redis"`
	Meili         Meili         `mapstructure:"meili"`
	Minio         Minio         `mapstructure:"minio"`
	Registry      Registry      `mapstructure:"registry"`
	OpenTelemetry OpenTelemetry `mapstructure:"open_telemetry"`
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

type Meili struct {
	Address string `mapstructure:"address"`
	APIKey  string `mapstructure:"api_key"`
}

type Minio struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Secure    bool   `mapstructure:"secure"`
	Address   string `mapstructure:"address"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
	TempDir   string `mapstructure:"temp_dir"`
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

type OpenTelemetry struct {
	Endpoint string `mapstructure:"endpoint"`
}

type HealthCheck struct {
	Addr string `mapstructure:"addr"`
}

type Registry struct {
	RegistryAddress []string `mapstructure:"registry_address"`
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// 从本地读取conf文件
	// _, filename, _, _ := runtime.Caller(0)
	// BasePath := filepath.Join(filepath.Dir(filename), "..")
	// prefix := "conf"
	// confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	// fmt.Println(confFileRelPath)
	// content, err := os.ReadFile(confFileRelPath)
	// if err != nil {
	//	 panic(err)
	// }
	// conf = new(Config)
	// err = yaml.Unmarshal(content, conf)
	// if err != nil {
	//	 klog.Error("parse yaml error - %v", err)
	//	 panic(err)
	// }
	// if err := validator.Validate(conf); err != nil {
	//	 klog.Error("validate config error - %v", err)
	//	 panic(err)
	// }
	// conf.Env = GetEnv()
	// _, _ = pretty.Printf("%+v\n", conf)

	// 从etcd中读取conf文件
	_ = godotenv.Load()
	registryAddress := "http://" + os.Getenv("REGISTRY_ADDR")
	configPath := "/config/product"
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
