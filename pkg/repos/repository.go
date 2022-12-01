package repo

type Repository interface {
	Store(interface{})
	Find()
}
