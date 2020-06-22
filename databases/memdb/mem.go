package memdb

import (
	"fmt"
)

/*
创建memDB
@param string dbName db名称
TODO 目前意义不大，为了之后多db提前准备
@return *MemDBSchema
*/
func NewMemDB(dbName string) (*MemDBSchema, error) {
	if dbName == "" {
		return nil, fmt.Errorf("dbName不能是空的")
	}
	return &MemDBSchema{
		Name:   dbName,
		Tables: make(map[string]*memTableSchema),
	}, nil
}

/**
删除memDB
@param string dbName db名称
@return error
*/
//func DropMemDB(dbName string) error {
//	if dbName == "" {
//		return fmt.Errorf("dbName不能是空的")
//	}
//	return nil
//}
