package attachment

import (
	"time"

	"ez-admin-gin/server/internal/model"
)

// Entity 复用附件中心模型，统一模块内部命名。
type Entity = model.SystemAttachment

// View 表示附件中心列表与详情查询时的联合视图。
type View struct {
	ID           uint
	FileID       uint
	DisplayName  string
	Category     string
	BizType      string
	OriginalName string
	FileName     string
	Ext          string
	MimeType     string
	Size         int64
	URL          string
	UploaderID   uint
	Status       model.SystemAttachmentStatus
	Remark       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
