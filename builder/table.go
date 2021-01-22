package builder

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/sirupsen/logrus"
)

/**
 * @Author lvxin0315@163.com
 * @Description 创建table
 * @Date 5:58 下午 2021/1/22
 * @Param
 * @return
 **/
func NewTable(name string) (Table, error) {
	t := Table{
		name: name,
	}
	err := t.init()
	if err != nil {
		logrus.Error("NewTable error: ", err)
		return t, err
	}
	return t, nil
}

/**
 * @Author lvxin0315@163.com
 * @Description 表信息
 * @Date 5:58 下午 2021/1/22
 * @Param
 * @return
 **/
type Table struct {
	name       string
	columns    []column
	columnMaps map[string]column
}

/**
 * @Author lvxin0315@163.com
 * @Description 初始化table 信息
 * @Date 5:57 下午 2021/1/22
 * @Param
 * @return
 **/
func (t *Table) init() error {
	query := fmt.Sprintf("SHOW COLUMNS FROM %s", t.name)
	c, err := client.Connect(fmt.Sprintf("%s:%d",
		config.MysqlConfig.Host,
		config.MysqlConfig.Port),
		config.MysqlConfig.User,
		config.MysqlConfig.Password,
		"")
	if err != nil {
		logrus.Error("Table.init client.Connect error: ", err)
		return err
	}
	defer c.Close()
	rr, err := c.Execute(query)
	if err != nil {
		logrus.Error("Table.init Execute error: ", err)
		return err
	}
	err = t.initColumns(rr)
	if err != nil {
		logrus.Error("Table.init initColumns error: ", err)
		return err
	}
	return nil
}

/**
 * @Author lvxin0315@163.com
 * @Description 初始化字段
 * @Date 5:58 下午 2021/1/22
 * @Param
 * @return
 **/
func (t *Table) initColumns(rr *mysql.Result) error {
	var columnList []column
	for i := 0; i < rr.RowNumber(); i++ {
		field, err := rr.GetStringByName(i, "Field")
		if err != nil {
			logrus.Error("NewColumnByMysqlResult Field error: ", err)
			return err
		}
		tp, err := rr.GetStringByName(i, "Type")
		if err != nil {
			logrus.Error("NewColumnByMysqlResult Type error: ", err)
			return err
		}
		null, err := rr.GetStringByName(i, "Null")
		if err != nil {
			logrus.Error("NewColumnByMysqlResult Null error: ", err)
			return err
		}
		key, err := rr.GetStringByName(i, "Key")
		if err != nil {
			logrus.Error("NewColumnByMysqlResult Key error: ", err)
			return err
		}
		def, err := rr.GetStringByName(i, "Default")
		if err != nil {
			logrus.Error("NewColumnByMysqlResult Default error: ", err)
			return err
		}
		columnList = append(columnList, column{
			field: field,
			tp:    tp,
			null:  "Yes" == null,
			key:   key,
			def:   def,
		})
	}
	t.columns = columnList
	// 初始化map
	t.columnMaps = make(map[string]column)
	for _, c := range t.columns {
		t.columnMaps[c.field] = c
	}
	return nil
}
