package wuid

import (
	"database/sql"
	"fmt"

	"github.com/edwingeng/wuid/mysql/wuid"
)

var w *wuid.WUID

func Init(dsn string) {
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}
	w = wuid.NewWUID("default", nil)
	// WUID 库的 LoadH28FromMysql 方法签名（核心）
	//func (w *WUID) LoadH28FromMysql(
	// newDB func() (*sql.DB, bool, error),  // 第一个参数：函数类型！
	// table string,                         // 第二个参数：表名
	//) error
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

func GenUid(dsn string) string {
	if w == nil {
		Init(dsn)
	}
	//w.Next() 返回的是一个 64 位整数（WUID 生成的全局唯一 ID）
	return fmt.Sprintf("%#016x", w.Next())
}
