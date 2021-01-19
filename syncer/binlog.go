package syncer

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type BinlogSyncer struct {
	sign chan bool
	// 字段集合
	tableColumnMap sync.Map
	// Position
	position mysql.Position
}

/**
 * @Author lvxin0315@163.com
 * @Description 启动
 * @Date 9:37 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) Start() {
	// 1. 表字段
	err := syncer.initTableColumn()
	syncer.error(err)
	// 2. 起始pos
	err = syncer.getMasterPos()
	syncer.error(err)
	go syncer.updateTableColumn()
	// 3. 启动监听
	syncer.sign = make(chan bool)
	go func() {
		err = syncer.listenBinlog()
		syncer.error(err)
	}()

}

/**
 * @Author lvxin0315@163.com
 * @Description 终止
 * @Date 9:44 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) Close() {
	syncer.sign <- false
	syncer.sign = nil
}

/**
 * @Author lvxin0315@163.com
 * @Description 初始化表的字段信息
 * @Date 9:45 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) initTableColumn() error {
	for _, schemaTableName := range config.SyncerConfig.Tables {
		err := syncer.getTableFields(schemaTableName)
		if err != nil {
			logrus.Error("initTableColumn error: ", err)
			syncer.Close()
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
func (syncer *BinlogSyncer) getTableFields(schemaTableName string) error {
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
 * @Description binlog 起点
 * @Date 3:23 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) getMasterPos() error {
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
	rr, err := c.Execute("SHOW MASTER STATUS")
	if err != nil {
		logrus.Error("getMasterPos Execute error: ", err)
		return err
	}
	name, _ := rr.GetString(0, 0)
	pos, _ := rr.GetInt(0, 1)
	syncer.position = mysql.Position{Name: name, Pos: uint32(pos)}
	return nil
}

func (syncer *BinlogSyncer) error(err error) {
	if err == nil {
		return
	}
	if config.CommonConfig.AppDebug {
		logrus.Error(err)
	}
	panic(err)
}

/**
 * @Author lvxin0315@163.com
 * @Description 监听解析binlog
 * @Date 3:23 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) listenBinlog() error {
	cfg := replication.BinlogSyncerConfig{
		ServerID:       uint32(config.SyncerConfig.ServerID),
		Flavor:         config.MysqlConfig.Flavor,
		Host:           config.MysqlConfig.Host,
		Port:           uint16(config.MysqlConfig.Port),
		User:           config.MysqlConfig.User,
		Password:       config.MysqlConfig.Password,
		RawModeEnabled: config.SyncerConfig.Raw,
		UseDecimal:     true,
	}
	b := replication.NewBinlogSyncer(cfg)
	defer b.Close()
	s, err := b.StartSync(syncer.position)
	if err != nil {
		logrus.Error("Start sync error: ", err)
		return err
	}

	for {
		select {
		case <-syncer.sign:
			logrus.Info("正常退出")
			return nil
		default:
			for _, ev := range s.DumpEvents() {
				logrus.Info("ev.Header.EventType:", ev.Header.EventType)
				// 处理日志
				syncer.dumpEvent(ev)
			}
		}

	}
}

/**
 * @Author lvxin0315@163.com
 * @Description 获取table的字段列表
 * @Date 6:53 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) getTableColumnList(name string) []string {
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
func (syncer *BinlogSyncer) refreshAndGetTableColumnList(name string) []string {
	err := syncer.getTableFields(name)
	if err != nil {
		panic(err)
	}
	return syncer.getTableColumnList(name)
}

/**
 * @Author lvxin0315@163.com
 * @Description 写event
 * @Date 11:43 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) writeEvent(ev *replication.BinlogEvent) {
	rowsEv := ev.Event.(*replication.RowsEvent)
	schema := rowsEv.Table.Schema
	table := rowsEv.Table.Table
	tableName := fmt.Sprintf("%s.%s", schema, table)

	if !syncer.inTableList(tableName) {
		return
	}
	columnNameList := syncer.checkColumnNumAndGetTableColumnList(tableName, int(rowsEv.ColumnCount))
	if len(columnNameList) == 0 {
		logrus.Error(tableName, " - 字段长度为0")
		return
	}
	for _, dataList := range rowsEv.Rows {
		for index, data := range dataList {
			//TODO
			logrus.Info(fmt.Sprintf("%s : %v", columnNameList[index], data))
		}
	}
}

/**
 * @Author lvxin0315@163.com
 * @Description 判断是否被监控
 * @Date 6:03 下午 2021/1/15
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) inTableList(name string) bool {
	for _, t := range config.SyncerConfig.Tables {
		if t == name {
			return true
		}
	}
	return false
}

/**
 * @Author lvxin0315@163.com
 * @Description 更新table的字段
 * @Date 3:38 下午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) updateTableColumn() {
	for {
		time.Sleep(60 * time.Second)
		logrus.Info("updateTableColumn ...")
		err := syncer.initTableColumn()
		syncer.error(err)
	}
}

/**
 * @Author lvxin0315@163.com
 * @Description 更新event
 * @Date 11:43 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) updateEvent(ev *replication.BinlogEvent) {
	rowsEv := ev.Event.(*replication.RowsEvent)
	schema := rowsEv.Table.Schema
	table := rowsEv.Table.Table
	tableName := fmt.Sprintf("%s.%s", schema, table)

	if !syncer.inTableList(tableName) {
		return
	}
	columnNameList := syncer.checkColumnNumAndGetTableColumnList(tableName, int(rowsEv.ColumnCount))
	if len(columnNameList) == 0 {
		logrus.Error(tableName, " - 字段长度为0")
		return
	}
	for _, dataList := range rowsEv.Rows {
		for index, data := range dataList {
			//TODO
			logrus.Info(fmt.Sprintf("%s : %v", columnNameList[index], data))
		}
	}
}

/**
 * @Author lvxin0315@163.com
 * @Description 删除event
 * @Date 5:34 下午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) deleteEvent(ev *replication.BinlogEvent) {
	rowsEv := ev.Event.(*replication.RowsEvent)
	schema := rowsEv.Table.Schema
	table := rowsEv.Table.Table
	tableName := fmt.Sprintf("%s.%s", schema, table)

	if !syncer.inTableList(tableName) {
		return
	}
	columnNameList := syncer.checkColumnNumAndGetTableColumnList(tableName, int(rowsEv.ColumnCount))
	if len(columnNameList) == 0 {
		logrus.Error(tableName, " - 字段长度为0")
		return
	}
	for _, dataList := range rowsEv.Rows {
		for index, data := range dataList {
			//TODO
			logrus.Info(fmt.Sprintf("%s : %v", columnNameList[index], data))
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
func (syncer *BinlogSyncer) checkColumnNumAndGetTableColumnList(tableName string, columnNum int) []string {
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
 * @Description QUERY_EVENT
 * @Date 5:51 下午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) queryEvent(ev *replication.BinlogEvent) {
	queryEv := ev.Event.(*replication.QueryEvent)
	logrus.Info("QueryEvent Query:", string(queryEv.Query))
	// TODO 目前全字段更新
	syncer.initTableColumn()
}

/**
 * @Author lvxin0315@163.com
 * @Description event 处理
 * @Date 10:26 上午 2021/1/19
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) dumpEvent(ev *replication.BinlogEvent) {
	switch ev.Header.EventType {
	case replication.WRITE_ROWS_EVENTv0:
	case replication.WRITE_ROWS_EVENTv1:
	case replication.WRITE_ROWS_EVENTv2:
		syncer.writeEvent(ev)

	case replication.UPDATE_ROWS_EVENTv0:
	case replication.UPDATE_ROWS_EVENTv1:
	case replication.UPDATE_ROWS_EVENTv2:
		syncer.updateEvent(ev)

	case replication.DELETE_ROWS_EVENTv0:
	case replication.DELETE_ROWS_EVENTv1:
	case replication.DELETE_ROWS_EVENTv2:
		syncer.deleteEvent(ev)

	case replication.QUERY_EVENT: //包含表结构变化
		syncer.queryEvent(ev)

	default:
		logrus.Info(ev.Header.EventType)
	}
}
