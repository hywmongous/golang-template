package commands

type Command interface {
	Apply(handler CommandHandler) error
}
