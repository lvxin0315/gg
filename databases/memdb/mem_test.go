package memdb

import (
	"fmt"
	"github.com/lvxin0315/gg/helper"
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
	_ = DropMemDB(TestDB)
	Convey("测试new内存DB", t, func() {
		db, err := NewMemDB(TestDB)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(len(db.Tables))
		testDB = db
	})
}

//建表
func TestMemDBSchema_CreateTableSchema(t *testing.T) {
	TestNewMemDB(t)
	Convey(fmt.Sprintf("测试创建表：%s", Table1), t, func() {
		createTable1, err := testDB.CreateTableSchema(Table1, reflect.TypeOf(new(testData)))
		if err != nil {
			fmt.Println(err)
			t.Fail()
			return
		}
		fmt.Println(createTable1.Data)
		table1 = createTable1
	})
}

//创建两张同名表
func TestMemDBSchema_ValidateTableName(t *testing.T) {
	TestMemDBSchema_CreateTableSchema(t)

	Convey(fmt.Sprintf("测试创建第二张表：%s", Table2), t, func() {
		table2, err := testDB.CreateTableSchema(Table2, reflect.TypeOf(new(testData)))
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(len(table2.Data))
	})

	Convey(fmt.Sprintf("测试创建第三张表：%s，应该报错误", Table2), t, func() {
		_, err := testDB.CreateTableSchema(Table2, reflect.TypeOf(new(testData)))
		if err == nil {
			fmt.Println(fmt.Errorf("未返回表同名错误"))
			t.Fail()
		}
		fmt.Println(err)
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
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(fmt.Sprintf("插入第一条数据后，长度是：%d", l))
	})
}

//insert 多条数据 1万
func TestMemTableSchema_Insert2(t *testing.T) {
	TestMemDBSchema_CreateTableSchema(t)
	Convey(fmt.Sprintf("开始插入"), t, func() {
		startTime := time.Now().Unix()
		limit := 10000
		for limit > 0 {
			_, err := table1.Insert(&testData{
				Name: fmt.Sprintf("name%d", 10000-limit),
				Age:  uint(10000 - limit),
			})
			if err != nil {
				fmt.Println(err)
				t.Fail()
				return
			}
			//fmt.Println(fmt.Sprintf("插入第%d条数据后，长度是：%d", 10000-limit, l))
			limit--
			//内存情况
			helper.PrintMemState()
		}
		//输出时间
		fmt.Println(fmt.Sprintf("消耗时间：%d", time.Now().Unix()-startTime))
		So(table1.Length(), ShouldEqual, 10000)
	})
}

//查询
func TestMemTableSchema_Select(t *testing.T) {
	//先插入10000
	TestMemTableSchema_Insert2(t)
	//长度
	fmt.Println("table1.Length():", table1.Length())
	selectFunc(t, 0, 100, 100)
	selectFunc(t, 100, 100, 100)
	selectFunc(t, 10000, 100, 0)
	selectFunc(t, 9999, 100, 1)
	selectFunc(t, 50000, 100, 0)
}

func selectFunc(t *testing.T, offset, limit uint, dataLength int) {
	Convey(fmt.Sprintf("查询用例: %d, %d", offset, limit), t, func() {
		data1, err := table1.Select(offset, limit)
		So(err, ShouldBeNil)
		fmt.Println("data1长度是：", len(data1))
		if dataLength > 0 {
			fmt.Println(fmt.Sprintf("查询用例: %d, %d 的最后一条数据是: %v", offset, limit, data1[len(data1)-1]))
			fmt.Println("Name:", data1[len(data1)-1].(*testData).Name)
			fmt.Println("Age:", data1[len(data1)-1].(*testData).Age)
		}
		So(len(data1), ShouldEqual, dataLength)
	})
}

//设置索引
func TestMemTableSchema_BaseIndex(t *testing.T) {
	//先插入10000
	TestMemTableSchema_Insert2(t)
	//长度
	fmt.Println("table1.Length():", table1.Length())
	Convey("简单的建立索引", t, func() {
		//定制callback
		callback := func(data interface{}) string {
			//年龄区分奇数偶数，被建立索引
			if data.(*testData).Age%2 == 0 {
				return "evenNumber"
			} else {
				return "oddNumber"
			}
		}
		length, err := table1.BaseIndex("number", callback)
		So(err, ShouldBeNil)
		fmt.Println("table1.BaseIndex number:", length)
		So(length, ShouldEqual, 2)
	})
}

//索引查询
func TestMemTableSchema_Index(t *testing.T) {
	TestMemTableSchema_BaseIndex(t)
	Convey("通过索引查询", t, func() {
		dataList, err := table1.Index("number", "evenNumber").Select(0, 100)
		So(err, ShouldBeNil)
		for i, v := range dataList {
			fmt.Println(fmt.Sprintf("dataList[%d]: %s, %d", i, v.(*testData).Name, v.(*testData).Age))
		}
	})
}
