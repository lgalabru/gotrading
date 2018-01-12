package reporting

type Report interface {
	Encode() ([]byte, error)
	Description() string
}
