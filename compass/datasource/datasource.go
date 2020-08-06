package datasource

import "compass/util"

type DataSource struct {
	util.BaseModel
	Name     string      `json:"name"`
	PluginID string      `json:"pluginId"`
	Health   bool        `json:"health"`
	Data     interface{} `json:"data"`
}