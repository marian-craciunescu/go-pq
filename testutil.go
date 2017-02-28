package go_pq

import "fmt"

type IntWrapper struct {
	value int
}

func (i IntWrapper) Priority() uint64 {
	return uint64(i.value)
}

func (i IntWrapper) String() string {
	return fmt.Sprintf("%d", i.value)
}
