package databases

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//plugin:gg_after_create
func ggAfterCreate(scope *gorm.Scope) {
	logrus.Println("ggAfterCreate:")
	for _, v := range scope.Fields() {
		logrus.Println("field:", v.Name, "value:", v.Field.String())
	}
	logrus.Println("SQL:", scope.SQL)
	for _, v := range scope.SQLVars {
		logrus.Println("SQLVars:", v)
	}
	logrus.Println("CombinedConditionSql:", scope.CombinedConditionSql())
	logrus.Println("TableName:", scope.TableName())
}

//plugin:gg_before_query_destination
func ggBeforeQueryDestination(scope *gorm.Scope) {
	logrus.Println("ggBeforeQueryDestination:")
	logrus.Println("TableName:", scope.TableName())
	for _, v := range scope.SQLVars {
		logrus.Println("SQLVars:", v)
	}
	logrus.Println("SQL:", scope.SQL)
}
