package global

// Queryer Queryer interface
type Queryer interface {
	Query(string) string
}
