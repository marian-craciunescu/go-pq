package go_pq


import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
)


// IndexList a sorted list of fetch entries
type PriorityQueue struct {
	items []Data

	sync.RWMutex
}

func newPQ(size int) *PriorityQueue {
	return &PriorityQueue{items: make([]Data, 0, size)}
}

func (pq *PriorityQueue) Len() int {
	pq.RLock()
	defer pq.RUnlock()

	return len(pq.items)
}

func (pq *PriorityQueue) Merge(other *PriorityQueue) {
	pq.Insert(other.toSliceArray()...)
}

//Insert  adds in the sorted list a new element
func (pq *PriorityQueue) Insert(items ...Data) {
	for _, elem := range items {
		pq.InsertElem(elem)
	}
}

func (pq *PriorityQueue) InsertElem(elem Data) {
	pq.Lock()
	defer pq.Unlock()

	// first element on list just append at the end
	if len(pq.items) == 0 {
		pq.items = append(pq.items, elem)
		return
	}

	// if the first element in list have a bigger id...Insert new element on the start of list
	if (pq.items[0]).Priority() >= (elem).Priority(){
		pq.items = append([]Data{elem}, pq.items...)
		return
	}

	if (pq.items[len(pq.items)-1]).Priority() <= (elem).Priority() {
		pq.items = append(pq.items, elem)
		return
	}

	//found the correct position to make an insertion sort
	for i := 1; i <= len(pq.items)-1; i++ {
		if (pq.items[i]).Priority() > (elem).Priority() {
			pq.items = append(pq.items[:i], append([]Data{elem}, pq.items[i:]...)...)
			return
		}
	}
}

// Clear empties the current list
func (pq *PriorityQueue) Clear() {
	pq.items = make([]Data, 0)
}

// GetIndexEntryFromID performs a binarySearch retrieving the
// true, the position and list and the actual entry if found
// false , -1 ,nil if position is not found
// search performs a binary search returning:
// - `true` in case the item was found
// - `position` position of the item
// - `bestIndex` the closest index to the searched item if not found.
// - `index` the index if found
func (pq *PriorityQueue) Search(searchID uint64) (bool, int, int, Data) {
	pq.RLock()
	defer pq.RUnlock()

	if len(pq.items) == 0 {
		return false, -1, -1, nil
	}

	h := len(pq.items) - 1
	f := 0
	bestIndex := f

	for f <= h {
		mid := (h + f) / 2
		if (pq.items[mid]).Priority()== searchID {
			return true, mid, bestIndex, pq.items[mid]
		} else if (pq.items[mid]).Priority()< searchID {
			f = mid + 1
		} else {
			h = mid - 1
		}

		if abs((pq.items[mid]).Priority(), searchID) <= abs((pq.items[bestIndex]).Priority(), searchID) {
			bestIndex = mid
		}
	}

	return false, -1, bestIndex, nil
}

//Back retrieves the element with the biggest id or nil if list is empty
func (pq *PriorityQueue) Back() Data {
	pq.RLock()
	defer pq.RUnlock()

	if len(pq.items) == 0 {
		return nil
	}

	return pq.items[len(pq.items)-1]
}

//Front retrieves the element with the smallest id or nil if list is empty
func (pq *PriorityQueue) Front() Data {
	pq.RLock()
	defer pq.RUnlock()

	if len(pq.items) == 0 {
		return nil
	}

	return pq.items[0]
}

func (pq *PriorityQueue) toSliceArray() []Data {
	pq.RLock()
	defer pq.RUnlock()

	return pq.items
}

//Front retrieves the element at the given index or nil if position is incorrect or list is empty
func (pq *PriorityQueue) Get(pos int) Data {
	pq.RLock()
	defer pq.RUnlock()

	if len(pq.items) == 0 || pos < 0 || pos >= len(pq.items) {
		log.WithFields(log.Fields{
			"Len": len(pq.items),
			"pos": pos,
		}).Info("Empty list or invalid index")
		return nil
	}

	return pq.items[pos]
}

func (pq *PriorityQueue) MapWithPredicate(predicate func(elem Data, i int) error) error {
	pq.RLock()
	defer pq.RUnlock()

	for i, elem := range pq.items {
		if err := predicate(elem, i); err != nil {
			return err
		}
	}

	return nil
}

func (pq *PriorityQueue) String() string {
	pq.RLock()
	defer pq.RUnlock()

	s := ""
	for i, elem := range pq.items {
		s += fmt.Sprintf("[%d:%d %s] ", i, (elem).Priority(), (elem).String())
	}
	return s
}

// Contains returns true if given ID is between first and last item in the list
func (pq *PriorityQueue) Contains(id uint64) bool {
	pq.RLock()
	defer pq.RUnlock()

	if len(pq.items) == 0 {
		return false
	}

	return (pq.items[0]).Priority() <= id && id <= (pq.items[len(pq.items)-1]).Priority()
}

func abs(m1, m2 uint64) uint64 {
	if m1 > m2 {
		return m1 - m2
	}

	return m2 - m1
}

