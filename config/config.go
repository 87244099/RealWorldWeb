package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// 把配置文件转成结构体
type config struct {
	Secret             string
	PublicKeyLocation  string
	PrivateKeyLocation string
	DefaultPortrait    string
	RedisAddr          string
}

var _config config

// 你可以改成大写，然后外部初始化
func init() {
	viper.SetConfigName("config")                          // name of config file (without extension)
	viper.SetConfigType("yaml")                            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("D:/software/GoDemo/RealWorldWeb") // optionally look for config in the working directory
	err := viper.ReadInConfig()                            // Find and read the config file
	if err != nil {                                        // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println("========")
	fmt.Println(viper.Get("secret"))
	fmt.Println("========")

	err = viper.Unmarshal(&_config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

}

func GetSecret() string {
	return _config.Secret
}

func GetPublicKeyLocation() string {
	return _config.PublicKeyLocation
}

func GetPrivateKeyLocation() string {
	return _config.PrivateKeyLocation
}

func GetDefaultPortrait() string {
	return _config.DefaultPortrait
}

func GetRedisAddr() string {
	return _config.RedisAddr
}
