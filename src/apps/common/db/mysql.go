package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 这个非常重要！ 否则 sql.Open 会出错
	"time"
)

/*
最主要代码在于
sqlDB, err := sql.Open("mysql", url)

*/

type _StoreUtil struct {
	storeDB         *sql.DB
	storeCmdTimeout time.Duration
}

type IStoreUtil interface {
	GetCmdTimeout() time.Duration
	GetSqlDB() *sql.DB
	GetConn(ctx context.Context) *sql.Conn

	// 感觉不是很能成为接口？

	InsertData(sqlStr string, args ...interface{}) error
	SelectData(sqlStr string, element interface{}, args ...interface{}) error
}

const (
	mysqlMaxPoolSize = 64               // 可有可无
	mysqlMinPoolSize = 32               // 可有可无
	mysqlMaxIdleTime = 30 * time.Second // 可有可无
	mysqlMaxLiftTime = 0 * time.Second  // 可有可无
)

func NewStoreUtil(mysqlUser, mysqlPwd, mysqlAddr, mysqlDataBase string, mysqlCmdTimeoutSec time.Duration) (IStoreUtil, error) {
	if mysqlCmdTimeoutSec == 0 {
		panic("NewStoreUtil mysqlCmdTimeoutSec is 0")
	}
	newMysql := func() (*sql.DB, error) {
		url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", mysqlUser, mysqlPwd, mysqlAddr, mysqlDataBase)
		sqlDB, err := sql.Open("mysql", url)
		if err != nil {
			return nil, fmt.Errorf("sql open err:%v", err)
		}
		sqlDB.SetMaxOpenConns(mysqlMaxPoolSize)
		sqlDB.SetMaxIdleConns(mysqlMinPoolSize)
		sqlDB.SetConnMaxIdleTime(mysqlMaxIdleTime)
		sqlDB.SetConnMaxLifetime(mysqlMaxLiftTime)

		conn, err := sqlDB.Conn(context.Background())
		if err != nil {
			fmt.Println("conn mysql err ", err)
			return nil, fmt.Errorf("sqlDB conn err:%v", err)
		}
		defer conn.Close()
		if err = conn.PingContext(context.Background()); err != nil {
			return nil, fmt.Errorf("conn PingContext err:%v", err)
		}
		return sqlDB, nil
	}

	sqlDB, err := newMysql()
	if err != nil {
		return nil, err
	}
	return &_StoreUtil{
		storeDB:         sqlDB,
		storeCmdTimeout: mysqlCmdTimeoutSec,
	}, nil
}

func (u *_StoreUtil) GetCmdTimeout() time.Duration {
	return u.storeCmdTimeout
}
func (u *_StoreUtil) GetSqlDB() *sql.DB {
	return u.storeDB
}
func (u *_StoreUtil) GetConn(ctx context.Context) *sql.Conn {
	ctx, cancel := context.WithTimeout(ctx, u.storeCmdTimeout)
	conn, err := u.storeDB.Conn(ctx)
	if err != nil {
		cancel() // 超时了err会不为nil吗？ 如果因为别的其他的err, 为啥还要调用cancel()?也可以不调吧
		panic(fmt.Sprintf("conn err:%v", err))
		return nil
	}
	return conn
}

func (u *_StoreUtil) InsertData(sqlStr string, args ...interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), u.GetCmdTimeout())
	conn := u.GetConn(ctx)
	defer conn.Close()
	_, err := conn.ExecContext(ctx, sqlStr, args...)
	return err
}

func (u *_StoreUtil) SelectData(sqlStr string, element interface{}, args ...interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), u.GetCmdTimeout())
	conn := u.GetConn(ctx)
	defer conn.Close()

	row := conn.QueryRowContext(ctx, sqlStr, args...)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.Scan(element)
	if err != nil {
		return err
	}
	return nil
}
