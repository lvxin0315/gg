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
	newTableSchema.initIndex()
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
	Indexes  map[string]memIndexSchema
	Data     []interface{}
	dataType reflect.Type
}

//索引
type memIndexSchema map[string]*memTableSchema

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

//判断是空的
func (table *memTableSchema) IsEmpty() bool {
	if len(table.Data) == 0 || table.Data == nil {
		return true
	}
	return false
}

//初始化索引
func (table *memTableSchema) initIndex() {
	if table.Indexes == nil {
		table.Indexes = make(map[string]memIndexSchema)
	}
}

//初始化索引table的基础信息
func (table *memTableSchema) initIndexTable(tableName, index string, dataType reflect.Type) *memTableSchema {
	if table.Name == "" {
		table.Name = fmt.Sprintf("index_%s_%s", tableName, index)
	}
	if table.dataType == nil {
		table.dataType = dataType
	}
	return table
}

/*
使用索引获取结果，配合Select
@param string index 索引名称
@return this
*/
func (table *memTableSchema) Index(index string, value string) *memTableSchema {
	return table.Indexes[index][value]
}

/*
简版索引
通过callback遍历data，建立
@param string index 索引名称
@param func(interface{}) string callback 回调方法
@return int 索引内容长度
*/
func (table *memTableSchema) BaseIndex(index string, callback func(interface{}) string) (length int, err error) {
	table.initIndex()
	mis := memIndexSchema{}
	//遍历data
	for _, item := range table.Data {
		if mis[callback(item)] == nil {
			mis[callback(item)] = new(memTableSchema)
			mis[callback(item)].initIndexTable(table.Name, index, table.dataType)
		}
		_, err = mis[callback(item)].Insert(item)
		if err != nil {
			break
		}
	}
	//刷新
	table.Indexes[index] = mis
	length = len(mis)
	return
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
