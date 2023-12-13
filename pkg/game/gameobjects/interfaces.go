package gameobjects

type DbReader interface {
	UpdateFromDb()
}

type DbWriter interface {
	WriteToDb()
}
