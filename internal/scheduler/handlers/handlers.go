package handlers

type Enum string

const (
	StreamProcessing Enum = "stream:processing"
	ImageResizing    Enum = "image:resizing"
)

func (e Enum) String() string {
	return string(e)
}
