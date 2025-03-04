package conf

import (
	"os"
	"path/filepath"
	"runtime"
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
	OpenTelemetry OpenTelemetry `yaml:"open_telemetry"`
	MySQL         MySQL         `yaml:"mysql"`
	Redis         Redis         `yaml:"redis"`
	Registry      Registry      `yaml:"registry"`
	HealthCheck   HealthCheck   `yaml:"health_check"`
}

type MySQL struct {
	Host     string `yaml:"db_host"`
	Port     int    `yaml:"db_port"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	DBName   string `yaml:"db_name"`

	DSN string `yaml:"dsn"`
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `yaml:"max_idle_conns"`
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `yaml:"max_open_conns"`
	// ConnMaxLifetime 连接最大存活时间
	ConnMaxLifetime int `yaml:"conn_max_lifetime"` // 秒
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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

type HealthCheck struct {
	Addr string
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// 获取项目根目录
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(filename), "..")
	prefix := "conf"
	var confFileRelPath string
	if env := GetEnv(); env != "online" {
		confFileRelPath = filepath.Join(basePath, prefix, filepath.Join(env, "conf.yaml"))
	} else {
		confFileRelPath = filepath.Join(prefix, filepath.Join(env, "conf.yaml"))
	}
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
		return "dev"
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
