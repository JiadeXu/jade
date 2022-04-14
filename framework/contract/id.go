package contract

const IDKey = "jade:id"

type IDService interface {
	NewID() string
}
