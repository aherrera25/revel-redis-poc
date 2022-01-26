package controllers

import (
	sessionHandler "clock_project_v2/app/actors/redis"
	"clock_project_v2/app/interceptors"
	"clock_project_v2/app/models"
	"fmt"

	"clock_project_v2/app"

	"github.com/google/uuid"
	"github.com/revel/revel"
)

type Settings struct {
	*revel.Controller
}

func (c Settings) Set(settings models.Settings) revel.Result {
	idHeaders := c.Request.Header.GetAll("x-content-id")
	var id string
	if len(idHeaders) == 0 {
		id = uuid.New().String()
	} else {
		id = idHeaders[0]
	}

	// get configuration based on specific env section settings from conf/app.conf
	content, _ := app.GetFromConfig("entity.name")

	// Build dependency inversion and Initialize the connection client session from a connection pool
	datasourceSetter, client := sessionHandler.SessionDependenciesFactory("MSET", app.Pool)
	// Assert redis session will be closed
	defer client.Close()

	// Call composable function
	currMessage, status := datasourceSetter(
		//building redis composed keys as follows: <entityName>.<uuid>.<attribute>
		fmt.Sprintf("%s.%s.%s", content, id, "Hours"), settings.Hours,
		fmt.Sprintf("%s.%s.%s", content, id, "Minutes"), settings.Minutes,
		fmt.Sprintf("%s.%s.%s", content, id, "Seconds"), settings.Seconds,
	)
	// setting up a contexto to access the stored information
	responseWithId := models.ResponseWithId{
		ContentId: id,
	}
	responseWithId.Content = currMessage

	c.Response.Status = status
	return c.RenderJSON(responseWithId)
}

func (c Settings) Get() revel.Result {
	// Initialize the connection client session from a connection pool
	id := c.Request.Header.GetAll("x-content-id")[0]
	// TODO: x-content-id validation
	// get configuration based on specific env section settings from conf/app.conf
	content, _ := app.GetFromConfig("entity.name")
	// Build dependency inversion
	composableFn, client := sessionHandler.SessionDependenciesFactory("MGET", app.Pool)
	// Assert redis session will be closed
	defer client.Close()
	// Call composable function
	currMessage, status := composableFn(
		//building redis composed keys as follows: <entityName>.<uuid>.<attribute>
		fmt.Sprintf("%s.%s.%s", content, id, "Hours"),
		fmt.Sprintf("%s.%s.%s", content, id, "Minutes"),
		fmt.Sprintf("%s.%s.%s", content, id, "Seconds"),
	)
	// multiresponse as interface{} to []interface{} in order to build the response for a MGET operation
	jsonMessage := currMessage.([]interface{})

	c.Response.Status = status
	return c.RenderJSON(models.Settings{
		Hours:   string(jsonMessage[0].([]byte)),
		Minutes: string(jsonMessage[1].([]byte)),
		Seconds: string(jsonMessage[2].([]byte)),
	})
}

func init() {
	// controller level interceptor validation
	revel.InterceptFunc(interceptors.VerifyTopicId, revel.BEFORE, &Settings{})
}
