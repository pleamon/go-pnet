package pnet

type Coding struct {
	Encode func([]byte) []byte
	Decode func([]byte) ([]byte, error)
}
