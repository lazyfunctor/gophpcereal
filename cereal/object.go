package cereal

type PropertyType int

const (
    Public PropertyType = iota
    Private PropertyType = iota
    Protected PropertyType = iota
)

type Property struct {
    PropValue interface{}
    PropType PropertyType
}

type PHPObject struct {
    ClassName string
    Properties map[string]Property
}