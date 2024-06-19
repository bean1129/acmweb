package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"acmweb/system/config"
)

type ZDatabase struct {
	Inst *sql.DB
}

func NewDB() *ZDatabase {
	// 获取配置实例
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=3s&parseTime=true", config.CONFIG.Mysql.Username, config.CONFIG.Mysql.Password, config.CONFIG.Mysql.Host, config.CONFIG.Mysql.Port, config.CONFIG.Mysql.Database, config.CONFIG.Mysql.Charset)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Printf("数据库连接错误:%v", err.Error())
		return nil
	}

	// 通过engine.Ping()来进行数据库的连接测试是否可以连接到数据库。
	err = db.Ping()
	if err == nil {
		fmt.Println("数据库连接成功")
		// 关闭连接
		// defer ZDB.Close()
	} else {
		fmt.Printf("数据库连接错误:%v", err.Error())
		return nil
	}

	// 设置连接池的空闲数大小
	db.SetMaxIdleConns(config.CONFIG.Mysql.MaxIdleCons)
	// 设置最大打开连接数
	db.SetMaxOpenConns(config.CONFIG.Mysql.MaxOpenCons)
	db.SetConnMaxLifetime(time.Minute * 60)

	// 开启调试模式和打印日志,会在控制台打印执行的sql
	if config.CONFIG.Mysql.Debug {
		fmt.Println("调试模式")
	}
	inst := new(ZDatabase)
	inst.Inst = db
	return inst
}

func (c *ZDatabase) Query(sql string, args ...any) (map[string]string, error) {
	rows, err := c.Inst.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		log.Printf("查询出错,SQL语句:%s\n错误详情:%s\n", sql, err.Error())
		return nil, err
	}

	// 获取列名cols
	cols, _ := rows.Columns()
	if len(cols) > 0 {
		buff := make([]interface{}, len(cols))       // 创建临时切片buff
		data := make([][]byte, len(cols))            // 创建存储数据的字节切片2维数组data
		dataKv := make(map[string]string, len(cols)) // 创建dataKv, 键值对的map对象
		for i, _ := range buff {
			buff[i] = &data[i] // 将字节切片地址赋值给临时切片,这样data才是真正存放数据
		}

		for rows.Next() {
			rows.Scan(buff...) // ...是必须的,表示切片
		}

		for k, col := range data {
			dataKv[cols[k]] = string(col)
			// fmt.Printf("%30s:\t%s\n", cols[k], col)
		}
		return dataKv, nil
	} else {
		return nil, nil
	}
}

func (c *ZDatabase) QueryRows(sql string, args ...any) ([]map[string]string, error) {
	rows, err := c.Inst.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("查询出错:\nSQL:\n%s, 错误详情:%s\n", sql, err.Error())
		return nil, err
	}
	// 获取列名cols
	cols, _ := rows.Columns()
	if len(cols) > 0 {
		var ret []map[string]string
		for rows.Next() {
			buff := make([]interface{}, len(cols))
			data := make([][]byte, len(cols)) // 数据库中的NULL值可以扫描到字节中
			for i, _ := range buff {
				buff[i] = &data[i]
			}
			rows.Scan(buff...) // 扫描到buff接口中，实际是字符串类型data中
			// 将每一行数据存放到数组中
			dataKv := make(map[string]string, len(cols))
			for k, col := range data { // k是index，col是对应的值
				// fmt.Printf("%30s:\t%s\n", cols[k], col)
				dataKv[cols[k]] = string(col)
			}
			ret = append(ret, dataKv)
		}
		return ret, nil
	} else {
		return nil, nil
	}
}

func (c *ZDatabase) QueryResultRows(sql string, args ...any) (resultData []map[string]interface{}, err error) {
	rows, err := c.Inst.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("查询出错:\nSQL:\n%s, 错误详情:%s\n", sql, err.Error())
		return nil, err
	}
	// 获取列名cols
	cols, _ := rows.Columns()
	columnCount := len(cols)
	values, valuesPoints := make([]interface{}, columnCount), make([]interface{}, columnCount)
	for rows.Next() {
		for i := 0; i < columnCount; i++ {
			valuesPoints[i] = &values[i]
		}

		rows.Scan(valuesPoints...)
		row := make(map[string]interface{})

		for i, val := range values {
			key := cols[i]
			// 判断val的值的类型
			var v interface{}
			b, ok := val.([]byte) //判断是否为[]byte
			if ok {
				v = string(b)
			} else {
				v = val
			}
			// 列名与值对应
			row[key] = v
		}
		resultData = append(resultData, row)
	}
	return
}

func (c *ZDatabase) StartTrans() (*sql.Tx, error) {
	tx, err := c.Inst.Begin()
	return tx, err
}

func (c *ZDatabase) CommitTrans(tx *sql.Tx) bool {
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

func (c *ZDatabase) Rollback(tx *sql.Tx) bool {
	err := tx.Rollback()
	if err != nil {
		return false
	}
	return true
}

func (c *ZDatabase) Exec(sql string, args ...any) (int64, error) {
	res, err := c.Inst.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return rows, nil
}

func (c *ZDatabase) ExecTx(tx *sql.Tx, sql string, args ...any) (int64, error) {
	res, err := tx.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return rows, err
}
