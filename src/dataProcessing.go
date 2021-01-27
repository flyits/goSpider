package goSpiderFarmwork

type DataHandler struct {
	Data DataItem

	Ops uint64
}
type DataItem interface {
	GetHandler() string
}

func (DataHandler) getHandlerFunc() string {
	return "default"
}
