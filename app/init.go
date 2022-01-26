package app

import (
	"clock_project_v2/app/actors/configuration"
	sessionHandler "clock_project_v2/app/actors/redis"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
	_ "github.com/revel/modules"
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	Pool redis.Pool

	GetFromConfig func(varName string) (string, bool)
)

// Initialize redis connection pool
func InitDB() {
	Pool = redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// initialize settings state
func initializingSettings() {
	// initialize multienv configuration set up by config var 'env'
	GetFromConfig = configuration.Build()
	// get configuration based on specific env section settings from conf/app.conf
	content, found := GetFromConfig("entity.name")
	if !found {
		revel.AppLog.Fatal(content)
		return
	}
	// Build dependency inversion
	datasourceSetter, client := sessionHandler.SessionDependenciesFactory("MSET", Pool)
	defer client.Close()
	// Call composable function
	id := uuid.New().String()
	currMessage, status := datasourceSetter(
		//building redis composed keys
		"SettingsId", id,
	)

	if status >= 400 {
		revel.AppLog.Fatal(fmt.Sprintf("%s", currMessage))
		return
	}

	currMessage, status = datasourceSetter(
		//building redis composed keys
		fmt.Sprintf("%s.%s.%s", content, id, "Hours"), "tick",
		fmt.Sprintf("%s.%s.%s", content, id, "Minutes"), "tack",
		fmt.Sprintf("%s.%s.%s", content, id, "Seconds"), "toe",
	)

	if status >= 400 {
		revel.AppLog.Fatal(fmt.Sprintf("%s", currMessage))
		return
	}
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// Run default base security headers initilization
	revel.OnAppStart(StartupScript)
	// Run datastore connection pool initialization
	revel.OnAppStart(InitDB)
	// Run the configuration set up process
	revel.OnAppStart(initializingSettings)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func StartupScript() {
	// revel.DevMod and revel.RunMode work here
	// Use this script to check for dev mode and set dev/prod startup scripts here!
	executionMode := "production"
	if revel.DevMode {
		executionMode = "development"
	}
	revel.AppLog.Info(fmt.Sprintf("Running on %s mode", executionMode))
}
