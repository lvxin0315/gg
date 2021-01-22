package builder

const (
	MYSQL_TP_NULL      = "unknown"
	MYSQL_TP_INT       = "int"
	MYSQL_TP_TINYINT   = "tinyint"
	MYSQL_TP_SMALLINT  = "smallint"
	MYSQL_TP_MEDIUMINT = "mediumint"
	MYSQL_TP_BIGINT    = "bigint"
	MYSQL_TP_DECIMAL   = "decimal"
	MYSQL_TP_FLOAT     = "float"
	MYSQL_TP_DOUBLE    = "double"
	MYSQL_TP_TIMESTAMP = "timestamp"
	MYSQL_TP_TIME      = "time"
	MYSQL_TP_DATE      = "date"
	MYSQL_TP_YEAR      = "year"
	MYSQL_TP_ENUM      = "enum"
	MYSQL_TP_SET       = "set"
	MYSQL_TP_BLOB      = "blob"
	MYSQL_TP_VARCHAR   = "varchar"
	MYSQL_TP_CHAR      = "char"
	MYSQL_TP_JSON      = "json"
	MYSQL_TP_GEOMETRY  = "geometry"
)

type column struct {
	field string //字段名
	tp    string //类型
	null  bool   //是否允许空
	key   string //键
	def   string //默认值
}
