package model

import "time"

// UserPost 保存用户和岗位的绑定关系。
type UserPost struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 固定用户岗位关系表名。
func (UserPost) TableName() string {
	return "sys_user_post"
}
