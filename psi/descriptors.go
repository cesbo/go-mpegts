package psi

import "errors"

type Descriptors []byte

var (
	ErrDescriptorFormat = errors.New("descriptor: invalid format")
)

func (d Descriptors) Check() error {
	end := len(d)
	skip := 0
	next := 0

	for skip < end {
		next = skip + 2
		if next > end {
			return ErrDescriptorFormat
		}
		next += int(d[skip+1])
		if next > end {
			return ErrDescriptorFormat
		}
		skip = next
	}

	return nil
}

func (d Descriptors) Next() Descriptors {
	next := 2 + int(d[1])

	if len(d) == next {
		return nil
	} else {
		return d[next:]
	}
}
