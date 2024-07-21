package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// 把配置文件转成结构体
type config struct {
	Secret             string
	PublicKeyLocation  string
	PrivateKeyLocation string
	DefaultPortrait    string
	RedisAddr          string
	MysqlDSN           string `mapstructure:"mysqlDSN"`
}

var _config config

// 你可以改成大写，然后外部初始化
func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	currentPath, _ := os.Getwd()  //获取当前的工作目录
	//fmt.Println("currentPath=", currentPath)
	viper.AddConfigPath(currentPath) // optionally look for config in the working directory
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file 111 : %w; currentPath=%s", err, currentPath))
	}

	fmt.Println("========")
	fmt.Println(viper.Get("secret"))
	fmt.Println("========")

	err = viper.Unmarshal(&_config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file 222: %w", err))
	}

}

func GetSecret() string {
	return _config.Secret
}

func GetPublicKeyLocation() string {
	currentPath, _ := os.Getwd()
	return filepath.Join(currentPath, _config.PublicKeyLocation)
}

func GetPrivateKeyLocation() string {
	currentPath, _ := os.Getwd()
	return filepath.Join(currentPath, _config.PrivateKeyLocation)
}

func GetDefaultPortrait() string {
	return _config.DefaultPortrait
}

func GetRedisAddr() string {
	return _config.RedisAddr
}

func GetMysqlDSN() string {
	return _config.MysqlDSN
}
