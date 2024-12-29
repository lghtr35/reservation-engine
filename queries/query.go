package queries

type Query interface {
	// Not sure if using generics better here
	Execute() (any, error)
}
