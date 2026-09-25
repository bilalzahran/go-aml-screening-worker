package model

type Tenant struct {
	*BaseEntity `bson:",inline"`
	DbName      string `bson:"dbName"`
	Name        string `bson:"name"`
}
