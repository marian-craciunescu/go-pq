package go_pq


type Data interface {
	Priority() uint64
	String() string
	//Compare(item Data) int
}
