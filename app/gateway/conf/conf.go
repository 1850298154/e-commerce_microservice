package conf

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf     *Config
	once     sync.Once
	BasePath string
)

type Config struct {
	Env           string
	Hertz         Hertz         `yaml:"hertz"`
	MySQL         MySQL         `yaml:"mysql"`
	Redis         Redis         `yaml:"redis"`
	OpenTelemetry OpenTelemetry `yaml:"open_telemetry"`
	HealthCheck   HealthCheck   `yaml:"health_check"`
	Security      Security      `yaml:"security"`
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
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	DB       int    `yaml:"db"`
}

type OpenTelemetry struct {
	Endpoint string `yaml:"endpoint"`
}

type Hertz struct {
	Service         string `yaml:"service"`
	Address         string `yaml:"address"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
	RegistryAddr    string `yaml:"registry_addr"`
}

type HealthCheck struct {
	Addr string
}

type Security struct {
	PublicRoutes []string `yaml:"public_routes"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// 获取项目根目录
	_, filename, _, _ := runtime.Caller(0)
	BasePath = filepath.Join(filepath.Dir(filename), "..")

	prefix := "conf"
	var confFileRelPath string
	if env := GetEnv(); env != "online" {
		confFileRelPath = filepath.Join(BasePath, prefix, filepath.Join(env, "conf.yaml"))
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
		hlog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		hlog.Error("validate config error - %v", err)
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
