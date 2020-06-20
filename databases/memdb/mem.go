package memdb

import (
	"fmt"
)

var memDB memDBSchema

/*
创建memDB
@param string dbName db名称
TODO 目前意义不大，为了之后多db提前准备
@return *memDBSchema
*/
func NewMemDB(dbName string) (*memDBSchema, error) {
	if dbName == "" {
		return &memDB, fmt.Errorf("dbName不能是空的")
	}
	if len(memDB.Tables) != 0 || memDB.Name != "" {
		return &memDB, fmt.Errorf("目前仅支持一个memdb")
	}
	memDB = memDBSchema{
		Name:   dbName,
		Tables: make(map[string]*memTableSchema),
	}
	return &memDB, nil
}

/**
删除memDB
@param string dbName db名称
@return error
*/
func DropMemDB(dbName string) error {
	if dbName == "" {
		return fmt.Errorf("dbName不能是空的")
	}
	memDB = memDBSchema{}
	return nil
}
