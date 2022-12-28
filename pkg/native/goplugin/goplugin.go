package goplugin

import (
	"fmt"
	"plugin"

	"github.com/ariaghora/zmol/pkg/val"
)

func Z_load_module(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return val.ERROR("load_module takes 2 arguments")
	}
	pluginPath := args[0]
	moduleName := args[1]

	// ensure both are strings
	if pluginPath.Type() != val.ZSTRING || moduleName.Type() != val.ZSTRING {
		return val.ERROR("load_module takes 2 strings: plugin path and exported module name")
	}

	plug, err := plugin.Open(pluginPath.(*val.ZString).Value)
	if err != nil {
		return val.ERROR(fmt.Sprintf("cannot open plugin: %s", err))
	}
	zModule, err := plug.Lookup(moduleName.(*val.ZString).Value)
	if err != nil {
		return val.ERROR(fmt.Sprintf("cannot find module %s in plugin: %s", moduleName.(*val.ZString).Value, err))
	}
	zModuleRef, ok := zModule.(**val.ZModule)
	if !ok {
		return val.ERROR(fmt.Sprintf("module %s is not a ZModule", moduleName))
	}

	return *zModuleRef
}

var GoPluginModule = val.MODULE(
	"goplugin",
	&val.Env{
		SymTable: map[string]val.ZValue{
			"load_module": &val.ZNativeFunc{Fn: Z_load_module},
		},
	},
)
