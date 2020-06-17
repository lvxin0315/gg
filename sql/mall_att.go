package model

type MallArticle struct {
	Id            uint32 `gorm:"type:INT(10) UNSIGNED;AUTO_INCREMENT;NOT NULL"`
	Cid           string `gorm:"type:VARCHAR(255);"`
	Title         string `gorm:"type:VARCHAR(255);NOT NULL"`
	Author        string `gorm:"type:VARCHAR(255);"`
	ImageInput    string `gorm:"type:VARCHAR(255);NOT NULL"`
	Synopsis      string `gorm:"type:VARCHAR(255);"`
	ShareTitle    string `gorm:"type:VARCHAR(255);"`
	ShareSynopsis string `gorm:"type:VARCHAR(255);"`
	Visit         string `gorm:"type:VARCHAR(255);"`
	Sort          uint32 `gorm:"type:INT(10) UNSIGNED;NOT NULL"`
	Url           string `gorm:"type:VARCHAR(255);"`
	Status        uint8  `gorm:"type:TINYINT(1) UNSIGNED;NOT NULL"`
	AddTime       string `gorm:"type:VARCHAR(255);NOT NULL"`
	Hide          uint8  `gorm:"type:TINYINT(1) UNSIGNED;NOT NULL"`
	AdminId       uint32 `gorm:"type:INT(10) UNSIGNED;NOT NULL"`
	MerId         uint32 `gorm:"type:INT(10) UNSIGNED;"`
	ProductId     int32  `gorm:"type:INT(10);NOT NULL"`
	IsHot         uint8  `gorm:"type:TINYINT(1) UNSIGNED;NOT NULL"`
	IsBanner      uint8  `gorm:"type:TINYINT(1) UNSIGNED;NOT NULL"`
}
