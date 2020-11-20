package database

var TablesName = []string{
	new(Chef).TableName(),
	new(Equip).TableName(),
	new(Recipe).TableName(),
	new(Guest).TableName(),
	new(Material).TableName(),
	new(Skill).TableName(),
}
