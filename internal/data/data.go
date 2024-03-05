package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"kratos-layout/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	etcdClient "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDiscovery,
	NewRegistrar,
	NewMysqlCmd,
	NewRedisCmd,
	NewGreeterRepo,
)

// Data .
type Data struct {
	cfg    *conf.Bootstrap
	logger *log.Helper
	db     *gorm.DB
	rdb    *redis.Client
}

// NewData .
func NewData(cfg *conf.Bootstrap, db *gorm.DB, redisCli *redis.Client, logger log.Logger) (*Data, func(), error) {
	logs := log.NewHelper(log.With(logger, "module", "layout/data"))
	cleanup := func() {
		logs.Info("closing the data resources")
	}
	return &Data{
		logger: logs,
		cfg:    cfg,
		db:     db,
		rdb:    redisCli,
	}, cleanup, nil
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	point := conf.Etcd.Address
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints: []string{point},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(client)
	return r
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	point := conf.Etcd.Address
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints: []string{point},
	})
	if err != nil {
		panic(err)
	}
	return etcd.New(client)
}

func NewMysqlCmd(conf *conf.Bootstrap, logger log.Logger) *gorm.DB {
	logs := log.NewHelper(log.With(logger, "module", "serviceName/data/mysql"))
	db, err := gorm.Open(mysql.Open(conf.Data.Database.Source), &gorm.Config{})
	if err != nil {
		logs.Fatalf("mysql connect error: %v", err)
	}
	// 如果是开发环境 打印sql
	if conf.Env.Mode == "dev" {
		db = db.Debug()
	}
	// todo: 数据迁移
	//db.AutoMigrate()
	return db
}

func NewRedisCmd(conf *conf.Data, logger log.Logger) *redis.Client {
	logs := log.NewHelper(log.With(logger, "module", "user-service/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		logs.Fatalf("redis connect error: %v", err)
	}
	return client
}
