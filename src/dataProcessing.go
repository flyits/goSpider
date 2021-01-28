package goSpider

type DataHandler struct {
	Data DataItem

	ops uint64
}
type DataItem interface {
	GetHandler() string
}

func (DataHandler) getHandlerFunc() string {
	return "default"
}
