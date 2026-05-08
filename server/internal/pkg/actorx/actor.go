// Package actorx 在 Gin 上下文中存取当前登录人的数据权限上下文。
package actorx

import (
	"ez-admin-gin/server/internal/platform/datascope"

	"github.com/gin-gonic/gin"
)

const currentActorKey = "current_actor"

// SetCurrentActor 将当前登录人的数据权限上下文写入 Gin 上下文。
func SetCurrentActor(c *gin.Context, actor datascope.Actor) {
	c.Set(currentActorKey, actor)
}

// CurrentActor 从 Gin 上下文中读取当前登录人的数据权限上下文。
func CurrentActor(c *gin.Context) (datascope.Actor, bool) {
	value, ok := c.Get(currentActorKey)
	if !ok {
		return datascope.Actor{}, false
	}

	actor, ok := value.(datascope.Actor)
	return actor, ok
}
