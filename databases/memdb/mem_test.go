package memdb

import (
	"fmt"
	"github.com/lvxin0315/gg/helper"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"time"
)

var testDB *memDBSchema

type testData struct {
	Name string
	Age  uint
}

const (
	TestDB = "testDB"
	Table1 = "test_table_1"
	Table2 = "test_table_2"
	Table3 = "test_table_3"
)

var table1 *memTableSchema

//建库
func TestNewMemDB(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	_ = DropMemDB(TestDB)
	Convey("测试new内存DB", t, func() {
		db, err := NewMemDB(TestDB)
		if err != nil {
			logrus.Error(err)
			t.Fail()
		}
		logrus.Println(len(db.Tables))
		testDB = db
	})
}

//建表
func TestMemDBSchema_CreateTableSchema(t *testing.T) {
	TestNewMemDB(t)
	Convey(fmt.Sprintf("测试创建表：%s", Table1), t, func() {
		createTable1, err := testDB.CreateTableSchema(Table1, reflect.TypeOf(new(testData)))
		if err != nil {
			logrus.Error(err)
			t.Fail()
			return
		}
		logrus.Println(createTable1.Data)
		table1 = createTable1
	})
}

//创建两张同名表
func TestMemDBSchema_ValidateTableName(t *testing.T) {
	TestMemDBSchema_CreateTableSchema(t)

	Convey(fmt.Sprintf("测试创建第二张表：%s", Table2), t, func() {
		table2, err := testDB.CreateTableSchema(Table2, reflect.TypeOf(new(testData)))
		if err != nil {
			logrus.Error(err)
			t.Fail()
		}
		logrus.Println(len(table2.Data))
	})

	Convey(fmt.Sprintf("测试创建第三张表：%s，应该报错误", Table2), t, func() {
		_, err := testDB.CreateTableSchema(Table2, reflect.TypeOf(new(testData)))
		if err == nil {
			logrus.Error(fmt.Errorf("未返回表同名错误"))
			t.Fail()
		}
		logrus.Println(err)
	})
}

//insert 一条数据
func TestMemTableSchema_Insert(t *testing.T) {
	TestMemDBSchema_CreateTableSchema(t)
	Convey(fmt.Sprintf("插入第一条数据"), t, func() {
		l, err := table1.Insert(&testData{
			Name: "name1",
			Age:  1,
		})
		if err != nil {
			logrus.Error(err)
			t.Fail()
		}
		logrus.Println(fmt.Sprintf("插入第一条数据后，长度是：%d", l))
	})
}

//insert 多条数据 1万
func TestMemTableSchema_Insert2(t *testing.T) {
	TestMemDBSchema_CreateTableSchema(t)
	Convey(fmt.Sprintf("开始插入"), t, func() {
		startTime := time.Now().Unix()
		limit := 10000
		for limit >= 0 {
			l, err := table1.Insert(&testData{
				Name: fmt.Sprintf("name%d", 10000-limit),
				Age:  uint(10000 - limit),
			})
			if err != nil {
				logrus.Error(err)
				t.Fail()
			}
			logrus.Println(fmt.Sprintf("插入第%d条数据后，长度是：%d", 10000-limit, l))
			limit--
			//内存情况
			helper.PrintMemState()
		}
		//输出时间
		logrus.Println(fmt.Sprintf("消耗时间：%d", time.Now().Unix()-startTime))
	})
}
