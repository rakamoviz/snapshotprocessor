package streamprocessingstatus

type Enum uint8

const (
	Undefined Enum = iota
	Running
	Success
	Failed
	Partial
)

func (e Enum) String() string {
	return [...]string{"Undefined", "Running", "Success", "Failed", "Partial"}[e]
}

func (e Enum) Index() uint8 {
	return uint8(e)
}
