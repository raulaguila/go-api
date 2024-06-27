package datatransferobject

type Lookup uint8

const (
	Body Lookup = iota
	Query
	Params
	Cookie
)
