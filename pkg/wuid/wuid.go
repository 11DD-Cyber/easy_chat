package wuid

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

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
	return fmt.Sprintf("%d", w.Next())
}

func CombineId(aid, bid string) string {
	ids := []string{aid, bid}
	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b

	})
	return fmt.Sprintf("%s_%s", ids[0], ids[1])

}
