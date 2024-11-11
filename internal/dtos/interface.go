package dtos

type Dto interface {
	ToDatabaseModel()
	FromDatabaseModel()
}