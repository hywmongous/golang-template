package queries

type Query interface {
	Apply(handler QueryHandler) error
}
