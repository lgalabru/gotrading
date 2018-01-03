package reporting

type Report interface {
	Encode() ([]byte, error)
}
