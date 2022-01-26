package interceptors

import (
	"clock_project_v2/app/models"

	"github.com/revel/revel"
)

func VerifyTopicId(c *revel.Controller) revel.Result {
	if topicId := c.Request.Header.GetAll("x-topic-id"); len(topicId) == 0 {
		c.Response.Status = 406
		return c.RenderJSON(models.Response{Content: "x-topic-id not found"})
	}
	return nil
}
