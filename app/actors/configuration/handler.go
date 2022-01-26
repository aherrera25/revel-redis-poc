package configuration

import (
	"fmt"

	"github.com/revel/revel"
)

// Factory method building the dependencies to access configuration data in conf/app.conf
func Build() func(varName string) (string, bool) {
	revel.Config.SetSection(revel.Config.StringDefault("env", "dev"))
	return func(varName string) (string, bool) {
		entityName, found := revel.Config.String(varName)
		if !found {
			return fmt.Sprintf("configuration variable %s not found", varName), found
		}
		return entityName, found
	}
}
