package receipt

type Service interface {
	PrintReceipt() error
	GetReceipt() (string, error)
}
