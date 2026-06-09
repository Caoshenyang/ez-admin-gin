package model

// CasbinRule 是 Casbin gorm-adapter 使用的策略表模型。
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Ptype string `gorm:"size:100;not null;default:''" json:"ptype"`
	V0    string `gorm:"size:100;not null;default:''" json:"v0"`
	V1    string `gorm:"size:100;not null;default:''" json:"v1"`
	V2    string `gorm:"size:100;not null;default:''" json:"v2"`
	V3    string `gorm:"size:100;not null;default:''" json:"v3"`
	V4    string `gorm:"size:100;not null;default:''" json:"v4"`
	V5    string `gorm:"size:100;not null;default:''" json:"v5"`
}

// TableName 固定 Casbin 策略表名。
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
