package controllers

import (
	"clock_project_v2/app/models"

	"github.com/revel/revel"
)

type Health struct {
	*revel.Controller
}

func (c Health) Index() revel.Result {
	return c.RenderJSON(models.Response{Content: "running"})
}
