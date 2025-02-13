package conf

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env           string
	Kitex         Kitex         `yaml:"kitex"`
	MySQL         MySQL         `yaml:"mysql"`
	Redis         Redis         `yaml:"redis"`
	Meili         Meili         `yaml:"meili"`
	Minio         Minio         `yaml:"minio"`
	Registry      Registry      `yaml:"registry"`
	OpenTelemetry OpenTelemetry `yaml:"openTelemetry"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Meili struct {
	Address string `yaml:"address"`
	APIKey  string `yaml:"api_key"`
}

type Minio struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Secure    bool   `yaml:"secure"`
	Address   string `yaml:"address"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
	TempDir   string `yaml:"temp_dir"`
}

type Kitex struct {
	Service        string `yaml:"service"`
	Address        string `yaml:"address"`
	LogLevel       string `yaml:"log_level"`
	LogFileName    string `yaml:"log_file_name"`
	LogMaxSize     int    `yaml:"log_max_size"`
	LogMaxBackups  int    `yaml:"log_max_backups"`
	LogMaxAge      int    `yaml:"log_max_age"`
	MaxConnections int    `yaml:"max_connections"`
	MaxQPS         int    `yaml:"max_qps"`
}

type OpenTelemetry struct {
	Endpoint string `yaml:"endpoint"`
}

type Registry struct {
	RegistryAddress []string `yaml:"registry_address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	prefix := "conf"
	confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	content, err := os.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
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
