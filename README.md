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
    - run
      serve `docker run --name realworld_mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci`
    - check server: `docker ps`
    - 如果服务没有运行，可以查看对应实例的日志：`docker logs c0690cc9ff66
    - stop serve: `docker stop realworld_mysql`
    - rm instance: `docker rm realworld_mysql`
    -
  导出全部数据：`docker exec realworld_mysql sh -c 'exec mysqldump -uroot -p123456 --all-databases' > ./sqls/realworld_mysql.sql`
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
instance.User(func (ctx){
ctx.Next()
}).PUT("/XXXX", func (ctx){

})
```

## 退出登录

把客户端的token给清理即可，res项目是清理localStorage里面的user

## ORM

### 基本用例

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

### 开启调试模式，查看完整sql

```go 
gormDB = gormDB.Debug()
```

然后终端段就会出现：

```shell
[107.243ms] [rows:1] DELETE FROM `article_comment` WHERE `article_comment`.`id` = 1
```

### 对象字段处理，比如Tag

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
return json.Marshal(j) //讲对象转换成字节数据
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

### 连表查询

article和user表连查，以article表为主进行左连接

- tables
    - article
      ```sql
        create table article
          (
              id              bigint(20)    not null primary key auto_increment,
              author_username varchar(255)  not null,
              title           varchar(4096) not null,
              slug            varchar(4096)  not null unique, /*这里扩容时会遇到：[42000][1071] Specified key was too long; max key length is 3072 bytes. 有2个方式：1.更换别的字符集，2.对这个列设置前缀索引*/
              `body`          text          not null,
              `description`   varchar(2048) not null default '',
              tag_list        varchar(1024) not null default '[]',
              created_at      datetime      not null default current_timestamp,
              updated_at      datetime      not null default current_timestamp
          ) engine = InnoDB
          default charset = utf8mb4;
      ```

    - user

      ```sql
      create table user (
        username varchar(64) not null primary key,
        password varchar(512) not null default '',
        email varchar(256) not null default '' unique,
        image varchar(1024) not null default ''
        bio varchar(1024) not null default ''
      ) engine=InnoDB default charset=utf8mb4;
      ```


- models

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

//如果数据库字段名为 author_user_name，那么它默认会被转换为 Go 结构体的字段名 AuthorUserName，符合大驼峰式命名的规则。如果你希望使用其他形式的命名规则，可以使用 column 标签来自定义数据库字段名。
AuthorUserEmail string `gorm:"->"` //AuthorUserEmail字段被标记为gorm:"->"，这表示它是一个关联字段，通过它可以访问到与当前模型关联的其他模型的Email字段。
AuthorUserImage string `gorm:"->"`
AuthorUserBio   string `gorm:"->"`
}
```

- handler

> gormDb.withContext(ctx).Model(tableA).Select(tableBSql).Joins(JoinSql).Order(orderFields).Find(&tableAList)

```go
err := gormDB.WithContext(ctx).Model(models.Article{}).
Select("article.*, user.email as author_user_name, user.bio as author_user_bio, user.image as author_user_image").
Joins("LEFT JOIN user on article.author_username = user.username").
Order("created_at desc").Limit(req.Limit).Offset(req.Offset).Find(&articles).Error
```

## web处理

- 获取queryString

```go
type ListArticleQuery struct {
Limit  int    `form:"limit"` //这些是go里面的结构体标签，如果你要读取它，一般要通过反射才行，类似java里面的注解
Offset int    `form:"offset"`
Tag    string `form:"tag"`
}

```

```go
//读取特定的queryString到req上
var req request.ListArticleQuery
err2 := ctx.Bind(&req)
if err2 != nil {
ctx.AbortWithError(http.StatusBadRequest, err2)
return
}
```

## uuid

docs: https://github.com/google/uuid
install: go get github.com/google/uuid

## cookie

### 读取

```go 
ctx.Cookie("token")
```

### 写入

```go
ctx.SetCookie("token", token, 24*3600, "/", "localhost:8080", false, true)
```

## 文件上传

### 读取

```go
file, err := ctx.FormFile("file")
if err != nil {
log.WithError(err).Infof("get file failed")
}
```

### 保存

```go
err = ctx.SaveUploadedFile(file, "./"+file.Filename)
```

## redis

### redis

- install: docker pull redis
- run: docker run --name realworld_redis -p 6379:6379 -d redis
- docs: https://redis.io/docs/latest/get-started/
- usages

### go-redis

- install: go get github.com/redis/go-redis/v9
- docs: https://github.com/redis/go-redis
- setup:
- init

```go
func InitRedis() {
ctx := context.Background()
log := logger.New(ctx)
log.Infof("redisAddr=%s", config.GetRedisAddr())

rdb = redis.NewClient(&redis.Options{
Addr: config.GetRedisAddr(),
//Password: "", // no password set
//DB:       0,  // use default DB
})
ping := rdb.Ping(context.Background())
err := ping.Err()
if err != nil {
log.WithError(err).Infof("redis ping failed")
panic(err)
}
}
```

- getter

```go
rdb.Get(ctx, USER_PROFILE_KEY+userName).Result()
```

- setter

```go
    err = rdb.Set(ctx, USER_PROFILE_KEY+userName, string(js), time.Duration(ttl)*time.Second).Err()
```

- delete

```go
    err := rdb.Del(ctx, USER_PROFILE_KEY+userName).Err()

```

#### 分布式索

- docs: https://github.com/bsm/redislock
- install: go get github.com/bsm/redislock
- 应用场景：
    - 修改、依赖的资源比较多，不能单纯依赖数据库的锁
    - 定时任务修改数据

## SSE

- install
    - web
  ```go
      ticker := time.NewTicker(1 * time.Second)
      ctx.Stream(func(w io.Writer) bool {
          select {
          case <-ticker.C:
              log.Infof("send sse event")
              ctx.SSEvent("", "heartbeat: "+time.Now().String())
          }
          return true
      })
  ```
    - res：浏览器打开localhost，控制台输入如下代码，进行sse的客户端激活注册
  ```javascript
  let es = new EventSource("/api/sse")
  es.onmessage = function(event){
      console.log(event.data);
  }

  es.close() //主动结束
  ```

## WebSocket

- docs: https://github.com/gorilla/websocket
- install: go get github.com/gorilla/websocket
  - usage
    - web
    ```go
        conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
      if err != nil {
          log.WithError(err).Errorf("Upgrader ws failed")
          ctx.AbortWithStatus(http.StatusBadRequest)
          return
      }

      defer conn.Close()

      for {
          reqMsg := make(map[string]interface{})
          err := conn.ReadJSON(&reqMsg)
          if err != nil {
              log.WithError(err).Errorf("ReadJSON failed")
              return
          }
          log.Infof("read msg: %v\n", utils.JsonMarshal(reqMsg))

          if reqMsg["exit"] != nil {
              return
          }

          err = conn.WriteJSON(reqMsg)
          if err != nil {
              log.WithError(err).Errorf("WriteJSON failed")
              return
          }
      }
      ```
      - res
        - npm install -g wscat
        - wscat -c ws://localhost:8000/api/ws

## 容器化

### 前端项目

- 多步容器化
  - dockerfile
    ```dockerfile
    # stage1: build
    FROM node:18-alpine AS build
    # 下载node18
    
    # 设置当前工作文件夹
    WORKDIR /app
    
    # 把package.json拷贝到app下
    COPY package.json ./
    # 安装依赖
    RUN npm install
    # 复制啥来的？
    COPY . .
    
    RUN npm run build
    
    # stage2: nginx
    FROM nginx:latest
    ## copy --from= someStage fromSrc nowSrc
    COPY --from=build /app/dist /usr/share/nginx/html
    
    EXPOSE 80
    CMD ["nginx", "-g", "daemon off;"]
    ```
  - 浏览器输入: http://localhost:5432/
  - 

### docker集群-docker swarm
- 什么是swarm
  - 即docker的集群管理和编排功能
  - 节点：一个docker实例
  - 服务：
    - 创建服务的时候，你可以指定使用哪个镜像，或者对正在运行的容器做操作
- 如何使用swarm模式
- docker swarm init <manger-ip> 如果不输入，就以本机ip作为默认的管理员
然后你可以查看启动后的状态
```shell
$ docker node ls
ID                            HOSTNAME         STATUS    AVAILABILITY   MANAGER STATUS   ENGINE VERSION
ngi4l4rskb55rtgm4oj8xh47x *   docker-desktop   Ready     Active         Leader           26.1.1

```
- 创建一个服务 
  - docker service create --name helloworld alpine ping docker.com
- 查看服务实例
  - docker service ps helloworld
- 扩容服务
  - docker service scale helloworld=5
- 删除服务
  - docker service rm helloworld
- 对现有服务进行更新
  - docker service update --image alpine --args "ping baidu.com" helloworld 
### 把前端项目放到swarm上
- 创建镜像仓库: 本质是从docker hub拉一个特定镜像下来，他作为镜像仓库的服务运行
  - docker pull registry
  - docker run -d -p 5000:5000 --restart always --name registry registry
  - 本机访问：http://localhost:5000/v2/
  - 打包镜像
    - docker tag vue3-realworld:latest localhost:5000/vue3-realworld:v1.0.0
  - 发布镜像
    - docker push localhost:5000/vue3-realworld:v1.0.0
  - docker service create --env BASE-URL=/ --env VITE_API_HOST=http://realworld-docker-compose-frontend:8080 --publish published=80,target=80 --replicas 3 --name realworld-frontend localhost:5000/vue3-realworld:v1.0.0
- 


- 

## FQA

### 数据库

- Q: varchar扩容出现3072的限制
    - 字段改成前缀索引：一般是字段设置了unique导致的

- Q: 连表查询出现：Error 1267 (HY000): Illegal mix of collations (utf8mb4_0900_ai_ci,IMPLICIT) and (
  utf8mb4_general_ci,IMPLICIT) for operation '='
    - 先确定你要用什么字符集，基于这个标准去查
    - 问题原因一般是：库，表，列的字符集三者没有保持一致导致的
    - 导出数据库的初始化sql`：mysqldump -u your_username -p your_database_name your_table_name > output_file.sql`
        - 如果你是docker里面执行，按如下操作
            - 查看运行的镜像：docker ps
            - 把对应的文件复制出来：docker cp <container_name_or_id>:<container_path> <host_path>
    - 查看output_file.sql里面，库，表和列的字符集是否存在不一致的情况下，全部改成一样即可
    -
- Q: request returned Internal Server Error for API route and version http://%2F%2F.%2Fpipe%2Fdocker_engine/_ping, check if the server supports the requested API version
  - launch docker desktop, pls
- Q: Error response from daemon: pull access denied for vue3-realworld, repository does not exist or may require 'docker login': denied: requested access to the resource is denied
  - docker hub确实没有这个镜像，要你自己打包下
  - 去对象项目基于dockerfile打包下：`docker build -t vue3-realworld .`
- Q: 在docker compose up启动服务的时候出现：panic: dial tcp 172.23.0.3:3306: connect: connection refused
  - 目前排查情况：
    - backend容器可以ping mysql的ip和service：证明网络是正常的
  - 采取措施：
    - 将backend的启动延迟若干秒，比如20s：服务恢复正常，说明最然backend依赖其他server started的状态，但是不代表server可以被连接
  - 常用命令：
    - docker-compose down 
    - docker-compose up --build
- Q: 后端是8080端口，前端老是请求8000端口
  - 前端是域名指定是针对nginx的，而静态资源只能写死域名，无法通过环境变量运行时注入，所以只能写localhost然后修改为后端服务对应的端口

- 
