package actorx

import (
	"ez-admin-gin/server/internal/platform/datascope"

	"github.com/gin-gonic/gin"
)

const currentActorKey = "current_actor"

func SetCurrentActor(c *gin.Context, actor datascope.Actor) {
	c.Set(currentActorKey, actor)
}

func CurrentActor(c *gin.Context) (datascope.Actor, bool) {
	value, ok := c.Get(currentActorKey)
	if !ok {
		return datascope.Actor{}, false
	}

	actor, ok := value.(datascope.Actor)
	return actor, ok
}
