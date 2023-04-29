package global

const Version = "1.1.1"
const Name = "readygo"
const Author = "zuolongxiao"
const Email = "zuolongxiao@gmail.com"

// Queryer Queryer interface
type Queryer interface {
	Query(string) string
}

// IContextWrapper
type IContextWrapper interface {
	GetUsername() string
}
