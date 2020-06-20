package memdb

import (
	"fmt"
	"reflect"
)

type memDBSchema struct {
	Name string
	//表
	Tables map[string]*memTableSchema
}

/*
创建表，表名唯一
*/
func (db *memDBSchema) CreateTableSchema(tableName string, dataType reflect.Type) (*memTableSchema, error) {
	err := db.ValidateTableName(tableName)
	if err != nil {
		return nil, err
	}
	//创建table
	newTableSchema := new(memTableSchema)
	newTableSchema.Name = tableName
	newTableSchema.dataType = dataType
	newTableSchema.Indexes = make(map[string]*memIndexSchema)
	db.Tables[tableName] = newTableSchema
	return newTableSchema, nil
}

/*
验证TableName唯一
@param string tableName table名称
@return error
*/
func (db *memDBSchema) ValidateTableName(tableName string) error {
	for name := range db.Tables {
		if name == tableName {
			return fmt.Errorf("「%s」表名已经存在", name)
		}
	}
	return nil
}

//表
type memTableSchema struct {
	Name string
	//索引
	Indexes  map[string]*memIndexSchema
	Data     []interface{}
	dataType reflect.Type
}

/*
添加数据到table
@param interface{} data 数据，对应table初始化类型
@return int length 表的数据长度
*/
func (table *memTableSchema) Insert(data interface{}) (length int, err error) {
	//验证类型
	err = table.ValidateDataType(data)
	if err != nil {
		return
	}
	table.Data = append(table.Data, data)
	length = len(table.Data)
	return length, err
}

/*
查询数据
@param uint offset 其实索引位
@param uint limit 查询长度
@return []interface{}
*/
func (table *memTableSchema) Select(offset, limit uint) ([]interface{}, error) {
	//data 空
	if table.IsEmpty() {
		//fmt.Println("data 空")
		return nil, nil
	}
	//索引大于长度
	if offset > uint(table.Length()+1) {
		//fmt.Println("索引大于长度")
		return nil, nil
	}
	//长度验证
	if (offset + limit) > uint(table.Length()+1) {
		//fmt.Println("长度验证")
		return table.Data[offset:], nil
	}
	//正常
	//fmt.Println("正常", offset, limit)
	return table.Data[offset:][:limit], nil
}

//获取长度
func (table *memTableSchema) Length() int {
	return len(table.Data)
}

//获取长度
func (table *memTableSchema) IsEmpty() bool {
	if len(table.Data) == 0 || table.Data == nil {
		return true
	}
	return false
}

/*
验证table的数据类型
*/
func (table *memTableSchema) ValidateDataType(data interface{}) error {
	if table.dataType != reflect.TypeOf(data) {
		return fmt.Errorf("table的数据类型是「%v」,参数的类型是「%v」", table.dataType, reflect.TypeOf(data))
	}
	return nil
}

//索引
type memIndexSchema struct {
	Value interface{}
	Data  interface{}
}