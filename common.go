package go_orm

type SchemaPotok interface {
	SchemaName() string
	TableName() string
}

func FullTableName(entity SchemaPotok) string {
	return entity.SchemaName() + "." + entity.TableName()
}
