package syncer

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/siddontang/go-mysql/client"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type tableSyncer struct {
	// 字段集合
	tableColumnMap sync.Map
	// sign
	sign chan bool
}

/**
 * @Author lvxin0315@163.com
 * @Description //init
 * @Date 6:18 下午 2021/1/19
 * @Param
 * @return
 **/
func (syncer *tableSyncer) init() error {
	err := syncer.initTableColumn()
	if err != nil {
		return err
	}
	// 定时更新
	syncer.sign = make(chan bool)
	go syncer.updateTableColumn()
	return nil
}

/**
 * @Author lvxin0315@163.com
 * @Description 初始化表的字段信息
 * @Date 9:45 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *tableSyncer) initTableColumn() error {
	for _, syncerTable := range config.SyncerConfig.Tables {
		err := syncer.getTableFields(syncerTable.Name)
		if err != nil {
			logrus.Error("initTableColumn error: ", err)
			return err
		}
	}
	return nil
}

/**
 * @Author lvxin0315@163.com
 * @Description 获取表字段名称
 * @Date 3:22 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *tableSyncer) getTableFields(schemaTableName string) error {
	query := fmt.Sprintf("SHOW COLUMNS FROM %s", schemaTableName)
	c, err := client.Connect(fmt.Sprintf("%s:%d",
		config.MysqlConfig.Host,
		config.MysqlConfig.Port),
		config.MysqlConfig.User,
		config.MysqlConfig.Password,
		"")
	if err != nil {
		logrus.Error("getTableFields client.Connect error: ", err)
		return err
	}
	defer c.Close()
	rr, err := c.Execute(query)
	if err != nil {
		logrus.Error("getTableFields Execute error: ", err)
		return err
	}
	var fieldList []string

	for i := 0; i < rr.RowNumber(); i++ {
		colName, err := rr.GetString(i, 0)
		if err != nil {
			logrus.Error("getTableFields GetString error: ", err)
			return err
		}
		fieldList = append(fieldList, colName)
	}
	syncer.tableColumnMap.Store(fmt.Sprintf("%s", schemaTableName), fieldList)
	return nil
}

/**
 * @Author lvxin0315@163.com
 * @Description 获取table的字段列表
 * @Date 6:53 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *tableSyncer) getTableColumnList(name string) []string {
	columnNameList, ok := syncer.tableColumnMap.Load(name)
	if !ok {
		// TODO 为毛会没有
		fmt.Println("getTableColumnList: 为毛没有")
		err := syncer.getTableFields(name)
		if err != nil {
			panic(err)
		}
		return syncer.getTableColumnList(name)
	}
	return columnNameList.([]string)
}

/**
 * @Author lvxin0315@163.com
 * @Description 刷新并获取table的字段列表
 * @Date 7:00 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *tableSyncer) refreshAndGetTableColumnList(name string) []string {
	err := syncer.getTableFields(name)
	if err != nil {
		panic(err)
	}
	return syncer.getTableColumnList(name)
}

/**
 * @Author lvxin0315@163.com
 * @Description 更新table的字段
 * @Date 3:38 下午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *tableSyncer) updateTableColumn() {
	for {
		select {
		case <-syncer.sign:
			logrus.Info("tableSyncer is closed")
			return
		default:
			logrus.Info("updateTableColumn ...")
			time.Sleep(2 * time.Second)
			err := syncer.initTableColumn()
			if err != nil {
				logrus.Error("updateTableColumn error: ", err)
			}
		}
	}
}

/**
 * @Author lvxin0315@163.com
 * @Description 校验字段数与内存中的值，并返回
 * @Date 5:37 下午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *tableSyncer) checkColumnNumAndGetTableColumnList(tableName string, columnNum int) []string {
	columnNameList := syncer.getTableColumnList(tableName)
	//字段数量变化
	if len(columnNameList) != columnNum {
		logrus.Info("refreshAndGetTableColumnList")
		columnNameList = syncer.refreshAndGetTableColumnList(tableName)
	}
	return columnNameList
}

/**
 * @Author lvxin0315@163.com
 * @Description 判断是否被监控
 * @Date 6:03 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *tableSyncer) inTableList(name string) bool {
	for _, syncerTable := range config.SyncerConfig.Tables {
		if syncerTable.Name == name {
			return true
		}
	}
	return false
}

/**
 * @Author lvxin0315@163.com
 * @Description 关闭
 * @Date 5:53 下午 2021/1/19
 * @Param
 * @return
 **/
func (syncer *tableSyncer) close() {
	logrus.Info("tableSyncer is closing")
	syncer.sign <- false
	syncer.sign = nil
}
