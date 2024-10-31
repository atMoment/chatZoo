package main

/*
import (
	"ChatZoo/common/db"
	"context"
	"fmt"
	"time"
)

const (
	MysqlDataBase      = "chatZoo" //数据库名字
	MysqlUser          = "root"
	MysqlPwd           = "111111"
	MysqlAddr          = "127.0.0.1:3306"
	MysqlCmdTimeoutSec = 3 * time.Second

	UserTable = "User"
)

var StoreUtil db.IStoreUtil

func init() {
	var err error
	StoreUtil, err = db.NewStoreUtil(MysqlUser, MysqlPwd, MysqlAddr, MysqlDataBase, MysqlCmdTimeoutSec)
	if err != nil {
		panic(" init err " + err.Error())
	}
}

func insertUser(id string, userData []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), StoreUtil.GetCmdTimeout())
	conn := StoreUtil.GetConn(ctx)
	defer conn.Close()
	sqlStr := fmt.Sprintf("insert into %s(ID, Data) values (?,?)  ", UserTable)
	_, err := conn.ExecContext(ctx, sqlStr, id, userData)
	return err
}

func deleteUser(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), StoreUtil.GetCmdTimeout())
	conn := StoreUtil.GetConn(ctx)
	defer conn.Close()

	sqlStr := fmt.Sprintf("delete from %s where ID = ?", UserTable)
	_, err := conn.ExecContext(ctx, sqlStr, id)
	return err
}

func updateUserData(id string, userData []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), StoreUtil.GetCmdTimeout())
	conn := StoreUtil.GetConn(ctx)
	defer conn.Close()

	sqlStr := fmt.Sprintf("update %s set Data = ? where ID = ?", UserTable)
	_, err := conn.ExecContext(ctx, sqlStr, id, userData)
	return err
}

func getUserData(id string) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), StoreUtil.GetCmdTimeout())
	conn := StoreUtil.GetConn(ctx)
	defer conn.Close()

	sqlStr := fmt.Sprintf("select Data from %s where ID = ? ", UserTable)

	row := conn.QueryRowContext(ctx, sqlStr, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var userData []byte
	err := row.Scan(&userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}
*/
