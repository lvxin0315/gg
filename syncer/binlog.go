package syncer

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/sirupsen/logrus"
)

type BinlogSyncer struct {
	sign chan bool
	// Position
	position mysql.Position
	// 消息通道
	cs channelSyncer
	// 表结构
	ts tableSyncer
}

/**
 * @Author lvxin0315@163.com
 * @Description 启动
 * @Date 9:37 上午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) Start() {
	// 1. 初始化通讯通道
	syncer.cs = channelSyncer{}
	syncer.cs.initChannels()

	// 2. 表字段
	syncer.ts = tableSyncer{}
	err := syncer.ts.init()
	syncer.error(err)

	// 3. 起始pos
	err = syncer.getMasterPos()
	syncer.error(err)
	// 4. 启动监听
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
	syncer.cs.closeChannels()
	syncer.ts.close()
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
	if !syncer.ts.inTableList(tableName) {
		return
	}
	columnNameList := syncer.ts.checkColumnNumAndGetTableColumnList(tableName, int(rowsEv.ColumnCount))
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

	if !syncer.ts.inTableList(tableName) {
		return
	}
	columnNameList := syncer.ts.checkColumnNumAndGetTableColumnList(tableName, int(rowsEv.ColumnCount))
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

	if !syncer.ts.inTableList(tableName) {
		return
	}
	columnNameList := syncer.ts.checkColumnNumAndGetTableColumnList(tableName, int(rowsEv.ColumnCount))
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
 * @Description QUERY_EVENT
 * @Date 5:51 下午 2021/1/18
 * @Param
 * @return
 **/
func (syncer *BinlogSyncer) queryEvent(ev *replication.BinlogEvent) {
	queryEv := ev.Event.(*replication.QueryEvent)
	logrus.Info("QueryEvent Query:", string(queryEv.Query))
	// TODO 目前全字段更新
	_ = syncer.ts.initTableColumn()
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
