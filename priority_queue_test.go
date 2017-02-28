package go_pq

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"math/rand"
)

func TestPriorityQueue_Add(t *testing.T) {

	a:= assert.New(t)
	pq :=newPQ(10)

	pq.InsertElem(IntWrapper{1})
	pq.InsertElem(IntWrapper{2})
	pq.InsertElem(IntWrapper{3})
	pq.InsertElem(IntWrapper{5})
	pq.InsertElem(IntWrapper{4})

	st:= pq.String()
	fmt.Println(st)
	a.Equal(st, "[0:1 1] [1:2 2] [2:3 3] [3:4 4] [4:5 5] ")
}


func Test_SortedListSanity(t *testing.T) {

	a := assert.New(t)
	pq	 := newPQ(1000)

	generatedIds := make([]int, 0, 11)

	for i := 0; i < 11; i++ {
		msgID := rand.Intn(50)
		generatedIds = append(generatedIds, msgID)


		pq.Insert(IntWrapper{msgID})
	}
	min := 200
	max := 0

	for _, id := range generatedIds {
		if max < id {
			max = id
		}
		if min > id {
			min = id
		}
		found, pos, _, foundEntry := pq.Search(uint64(id))
		a.True(found)
		a.Equal(foundEntry.Priority(), uint64(id))
		a.True(pos >= 0 && pos <= len(generatedIds))
	}


	a.Equal(uint64(min), pq.Front().Priority())
	a.Equal(uint64(max), pq.Back().Priority())

	found, pos, _, foundEntry := pq.Search(uint64(46))
	a.False(found, "Element should not be found since is a number greater than the random generated upper limit")
	a.Equal(pos, -1)
	a.Nil(foundEntry)

	a.Equal(pq.Front().Priority(), pq.Get(0).Priority(), "First element should contain the smallest element")
	a.Nil(pq.Get(-1), "Trying to get an invalid index will return nil")

	pq.Clear()
	a.Nil(pq.Front())
	a.Nil(pq.Back())

}