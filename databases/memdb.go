package databases

import (
	"github.com/lvxin0315/gg/databases/memdb"
	"github.com/lvxin0315/gg/etc"
	"github.com/sirupsen/logrus"
)

var memDB *memdb.MemDBSchema

func InitMemDB() *memdb.MemDBSchema {

	logrus.Info(etc.Config.MemDB)

	newMemDB, err := memdb.NewMemDB(etc.Config.MemDB.Name)
	if err != nil {
		panic(err)
	}
	memDB = newMemDB
	return memDB
}

//TODO 要不要copy，之后再斟酌
func NewMemBD() *memdb.MemDBSchema {
	return memDB
}
