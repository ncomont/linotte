package csv

const (
	LineTypeHeader  = "HEADER"
	LineTypeContent = "CONTENT"
)

type Line struct {
	Index    int
	Type     string
	Elements []string
	Error    error
}
