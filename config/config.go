package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 用于存储所有配置项
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Server   ServerConfig   `mapstructure:"server"`
}

// DatabaseConfig 数据库相关配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`     // 数据库主机地址
	Port     int    `mapstructure:"port"`     // 数据库端口
	User     string `mapstructure:"user"`     // 数据库用户名
	Password string `mapstructure:"password"` // 数据库密码
	DBName   string `mapstructure:"dbname"`   // 数据库名称
}

// RedisConfig Redis相关配置
type RedisConfig struct {
	Host string `mapstructure:"host"` // Redis主机地址
	Port int    `mapstructure:"port"` // Redis端口
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"` // JWT密钥
	Expire int    `mapstructure:"expire"` // JWT过期时间（秒）
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

var Conf *Config

// InitConfig 初始化配置文件
func InitConfig() {
	// 设置配置文件的名称和路径
	viper.SetConfigName("config")   // 配置文件的名称（不带扩展名）
	viper.AddConfigPath("./config") // 配置文件所在路径
	viper.AddConfigPath("../../config") // 从cmd/server目录启动时的相对路径
	viper.AddConfigPath("../config") // 从其他子目录启动时的相对路径
	viper.SetConfigType("yaml")     // 配置文件类型

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件错误: %s", err)
	}

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

	fmt.Println("服务器端口:", Conf.Server.Port)
}
