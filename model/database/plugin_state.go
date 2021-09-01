package database

type PluginState struct {
	GroupId    int64  `xorm:"group_id pk"`
	PluginName string `xorm:"plugin_name pk"`
	State      bool   `xorm:"state"`
}

func (PluginState) TableName() string {
	return "plugin_state"
}
