## 环境准备
- 变量设置：`export GOPATH="/c/Users/faisco/go"`
- 安装模块：`go install github.com/gin-gonic/gin@latest`
- 镜像设置：`export GOPROXY=https://goproxy.cn`
  - goland设置：file->setting->go modules-> environment variables

## 项目初始化流程

- mkdir goquick_start
- cd goquick_start
- go mod init example.com/goquick_start 这里会出现go.sum的文件
- go get -u github.com/gin-gonic/gin
- 剩余参考：https://gin-gonic.com/docs/quickstart/

## 日志
- logrus: 最常用的日志库
- zap: 性能比较高

### logrus 使用指南
- docs: https://github.com/sirupsen/logrus
- install:
  - 如果文档里面没有，参考install go get
  - `go get github.com/sirupsen/logrus`
- usages
  - create instance by gin context: `log := logger.New(ctx)`
  - output with error: `log.WithError(err).Errorln("build json failed")`
  - output with custom info: `log.WithField("user", utils.JsonMarshal(body)).Infof("user login called")`
  - config custom log formater: `logrus.SetFormatter(&logrus.JSONFormatter{})`

## token生成

### 对称加密

#### jwt
- docs: https://github.com/golang-jwt/jwt
- install: `go get -u github.com/golang-jwt/jwt/v5`

### 非堆成加密

#### 生成密钥
- 电脑看是否有`openssl`
- 生成私钥:`openssl genrsa -out private.pem 2048`，当前路径下会出现一个`private.pem`的文件
- 由私钥生成公钥：`openssl rsa -in private.pem --outform PEM -pubout -out public.pem`，目录下会出现`public.pem`
#### 验证
- usage
```go
func VerifyJwt(token string) (bool, error) {
	var claim jwt.MapClaims
	claims, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {// 比如过期，内容被串改等都有提示
		return false, err
	}
	if claims.Valid {//最终看这个值处理
		return true, nil
	}
	return false, nil
}
```
- error type: 
  - 内容不对： `token signature is invalid: crypto/rsa: verification error`
  - 过期了：`token has invalid claims: token is expired`
- QA:
  - Q: 为什么jwt中间那一段用base64？
#### 生成token



## 配置文件
- background: 
  - 一些配置信息可能存放在go文件，配置文件或环境变量
  - 如果要统一这些方式，需要第三方库解决

### viper
- docs: https://github.com/spf13/viper
- usages
  - init
  
    ```go
     viper.SetConfigName("config") // name of config file (without extension)
  	  viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
  	  viper.AddConfigPath("..")     // optionally look for config in the working directory
  	  err := viper.ReadInConfig()   // Find and read the config file
  	  if err != nil {               // Handle errors reading the config file
  	  	panic(fmt.Errorf("fatal error config file: %w", err))
  	  }
    ```

  - read config: `viper.Get("secret")`
  - read struct: `viper.Unmarshal(&_config)`


## 数据库

这里使用docker去安装mysql

- install
  - 安装docker
  - docker hub找到mysql进行安装: `docker pull msyql`
  - run serve `docker run --name realworld_mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci`
  - check server: `docker ps`
  - 如果服务没有运行，可以查看对应实例的日志：`docker logs c0690cc9ff66
  - stop serve: `docker stop realworld_mysql`
  - rm instance: `docker rm realworld_mysql`
- connect
  - goland右侧有个database照着ip，port等配置下即可
  - 最后点击下[Test Connection]测试能否连通即可
- table design
  - database
    ```sql
        create database realworld;
    ```    
  - tables
    - users
      ```sql
      create table user (
      username varchar(64) not null primary key,
      password varchar(512) not null default '',
      email varchar(256) not null default '' unique,
      image varchar(1024) not null default ''
      bio varchar(1024) not null default ''
      ) engine=InnoDB default charset=utf8mb4
      ```
    - 
      `
### go操作数据库

#### sqlx

docs: https://github.com/jmoiron/sqlx
go: `go get github.com/jmoiron/sqlx`

#### go-sql-driver

docs: https://github.com/go-sql-driver/mysql
go: `go get -u github.com/go-sql-driver/mysql`



## 密码

docs: https://dropbox.tech/security/how-dropbox-securely-stores-your-passwords

### 原理
- 散列hash
- 加盐
- 验证

### 流程

- 注册流程
  - 客户端发送明文
  - 服务端对密码进行加盐
  - 然后进行加密
  - 将加密文本入库
- 登录流程
  - 客户端发送明文
  - 服务端对密码进行加盐
  - 然后db密码和加盐后的客户端密码做对比

### usages
- generate
```go
func HashPassword(password string) (string, error) {
	password += salt
	//deal with password
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10) //加盐的强度
	if err != nil {
		return "", err
	}

	//TODO 为什么存储的时候最好密码base64一次？
	return base64.StdEncoding.EncodeToString([]byte(bcryptPassword)), nil
}
```
- compare
```go
func CheckPassword(plain, hash string) bool {
	plain += salt
	bcryptedPassword, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(bcryptedPassword, []byte(plain))
	return err == nil
}

```

## 中间件
```go
instance.User(func(ctx){
	ctx.Next()
}).PUT("/XXXX", func(ctx){
	
})
```

## 退出登录

把客户端的token给清理即可，res项目是清理localStorage里面的user

## orm数据库

docs: https://gorm.io/docs
install: `go get -u gorm.io/gorm`
initialize:
    
  - 声明变量

    ```go 
        var db *sqlx.DB //为什么这个会被引用？
        var err error
        var gormDB *gorm.DB
    ```

  - 链接数据库：

    ```go
       
        db, err = sqlx.Open("mysql", "root:123456@(localhost:3306)/realworld")

    ```
    
  - 绑定orm：

    ```go
      gormDB, err = gorm.Open(gorm_mysql.New(gorm_mysql.Config{
          Conn: db,
      }), &gorm.Config{})
    ```

声明model
```go
type Article struct {
	Id             int64  `db:"id"` //依赖db的字段自增
	AuthorUsername string `gorm:"column:author_username"`
	Title          string
	Slug           string
	Body           string
	Description    string
	TagList        TagList `gorm:"type:string"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
```
实现获取表名的接口：
```go
func (a Article) TableName() string {
    return "article" //todo what is syntax sugar?
}
```

对象字段处理，比如Tag
- 声明TagList结构体
```go
type TagList []string //难道结合下面的scan可以实现面向对象语法？
```
- 实现结构体的Scan接口
```go
// Scan将数据库中存储的字节数据转TagList类型，判断能否转换成功
func (j *TagList) Scan(value interface{}) error {
	// 将value转换为[]byte类型
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	// 创建一个json.RawMessage类型的变量
	result := json.RawMessage{}

	// 将字节解析为json.RawMessage类型的对象
	err := json.Unmarshal(bytes, &result)
	return err
}
```
- 将TagList类型的数据，转成数据库可以存储的
```go
// Value return json value, implement driver.Valuer interface
func (j TagList) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)//讲对象转换成字节数据
}
```

创建curd

  ```go
    func CreateArticle(ctx context.Context, article *models.Article) error {
		return gormDB.WithContext(ctx).Create(article).Error
    }
  ```

调用curl
  ```
	ctx := context.TODO()
    err = CreateArticle(ctx, &models.Article{
        AuthorUsername: article["author"].(map[string]interface{})["username"].(string), //这个拿到的是一个对象
        Title:          article["title"].(string),
        Slug:           article["slug"].(string),
        Body:           article["body"].(string),
        Description:    article["description"].(string),
        TagList:        tagList, //数组要怎么转成json？
        CreatedAt:      createdAt,
        UpdatedAt:      updatedAt,
    })
    if err != nil {
        t.Fatal(err)
    }
  ```

## FQA

- Q: varchar扩容出现3072的限制
  - 字段改成前缀索引  

- 