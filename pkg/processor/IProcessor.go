package processor

type IProcessor interface {
	getBalance(address string)
}
