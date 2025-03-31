package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 用于存储所有配置项
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Server   ServerConfig   `mapstructure:"server"`
	SMS      SMSConfig      `mapstructure:"sms"`
}

// DatabaseConfig 数据库相关配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host" env:"DATABASE_HOST"`         // 数据库主机地址
	Port     int    `mapstructure:"port" env:"DATABASE_PORT"`         // 数据库端口
	User     string `mapstructure:"user" env:"DATABASE_USER"`         // 数据库用户名
	Password string `mapstructure:"password" env:"DATABASE_PASSWORD"` // 数据库密码
	DBName   string `mapstructure:"dbname" env:"DATABASE_NAME"`       // 数据库名称
}

// RedisConfig Redis相关配置
type RedisConfig struct {
	Host string `mapstructure:"host" env:"REDIS_HOST"` // Redis主机地址
	Port int    `mapstructure:"port" env:"REDIS_PORT"` // Redis端口
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret" env:"JWT_SECRET"` // JWT密钥
	Expire int    `mapstructure:"expire" env:"JWT_EXPIRE"` // JWT过期时间（秒）
}

type ServerConfig struct {
	Port int    `mapstructure:"port" env:"SERVER_PORT"`
	Mode string `mapstructure:"mode" env:"SERVER_MODE"` // 运行模式：debug或release
}

// SMSConfig 短信服务配置
type SMSConfig struct {
	Provider string          `mapstructure:"provider" env:"SMS_PROVIDER"` // 短信服务提供商，如aliyun
	Aliyun   AliyunSMSConfig `mapstructure:"aliyun"`                      // 阿里云短信配置
}

// AliyunSMSConfig 阿里云短信配置
type AliyunSMSConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id" env:"SMS_ALIYUN_ACCESS_KEY_ID"`         // 阿里云AccessKey ID
	AccessKeySecret string `mapstructure:"access_key_secret" env:"SMS_ALIYUN_ACCESS_KEY_SECRET"` // 阿里云AccessKey Secret
	RegionID        string `mapstructure:"region_id" env:"SMS_ALIYUN_REGION_ID"`                 // 地域ID
	SignName        string `mapstructure:"sign_name" env:"SMS_ALIYUN_SIGN_NAME"`                 // 短信签名
	TemplateCode    string `mapstructure:"template_code" env:"SMS_ALIYUN_TEMPLATE_CODE"`         // 短信模板ID
}

var Conf *Config

// InitConfig 初始化配置文件
func InitConfig() {
	// 创建一个新的Conf实例
	Conf = &Config{}

	// 加载.env文件中的环境变量
	loadEnvFile()

	// 设置配置文件的名称和路径
	viper.SetConfigName("config")       // 配置文件的名称（不带扩展名）
	viper.AddConfigPath("./config")     // 配置文件所在路径
	viper.AddConfigPath("../../config") // 从cmd/server目录启动时的相对路径
	viper.AddConfigPath("../config")    // 从其他子目录启动时的相对路径
	viper.SetConfigType("yaml")         // 配置文件类型

	// 读取配置文件作为默认值
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("警告: 读取配置文件错误: %s，将仅使用环境变量\n", err)
	}

	// 设置环境变量前缀和自动绑定环境变量
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	// 设置环境变量与配置项的映射关系
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 绑定环境变量
	// 数据库配置
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "DATABASE_NAME")

	// Redis配置
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")

	// JWT配置
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expire", "JWT_EXPIRE")

	// 服务器配置
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")

	// 短信服务配置
	viper.BindEnv("sms.provider", "SMS_PROVIDER")
	viper.BindEnv("sms.aliyun.access_key_id", "SMS_ALIYUN_ACCESS_KEY_ID")
	viper.BindEnv("sms.aliyun.access_key_secret", "SMS_ALIYUN_ACCESS_KEY_SECRET")
	viper.BindEnv("sms.aliyun.region_id", "SMS_ALIYUN_REGION_ID")
	viper.BindEnv("sms.aliyun.sign_name", "SMS_ALIYUN_SIGN_NAME")
	viper.BindEnv("sms.aliyun.template_code", "SMS_ALIYUN_TEMPLATE_CODE")

	// 将配置文件中的内容映射到结构体Config
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("配置文件映射到结构体时出错: %s", err)
	}

	// 检查配置文件中是否有必需的配置项
	// 如果缺少必要的数据库配置，直接报错
	if Conf.Database.Host == "" || Conf.Database.Port == 0 || Conf.Database.User == "" || Conf.Database.DBName == "" {
		log.Fatal("缺少必需的数据库配置项")
	}

	// 如果缺少Redis配置，直接报错
	if Conf.Redis.Host == "" || Conf.Redis.Port == 0 {
		log.Fatal("缺少必需的Redis配置项")
	}

	// 如果缺少JWT密钥配置，直接报错
	if Conf.JWT.Secret == "" || Conf.JWT.Expire <= 0 {
		log.Fatal("JWT配置不完整")
	}

	// 可选：打印配置项，方便调试
	fmt.Println("JWT配置:", Conf.JWT.Secret, "过期时间:", Conf.JWT.Expire)
	fmt.Println("数据库配置:", Conf.Database)
	fmt.Println("Redis配置:", Conf.Redis)
	if Conf.Server.Port == 0 {
		log.Fatal("缺少服务器端口配置")
	}

	// 检查服务器模式配置，如果未设置则默认为debug模式
	if Conf.Server.Mode == "" {
		Conf.Server.Mode = "debug"
	}

	fmt.Println("服务器端口:", Conf.Server.Port)
	fmt.Println("服务器模式:", Conf.Server.Mode)
}

// loadEnvFile 加载.env文件中的环境变量
func loadEnvFile() {
	// 尝试加载.env文件
	err := godotenv.Load()
	if err != nil {
		// 如果.env文件不存在，只打印警告，不终止程序
		fmt.Println("警告: .env文件未找到，将使用配置文件或系统环境变量")
	} else {
		fmt.Println("成功加载.env文件")
	}
}

// 注意：原overrideConfigWithEnv函数已被移除，因为现在使用viper的自动环境变量绑定功能
