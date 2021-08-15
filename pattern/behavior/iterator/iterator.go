package iterator

type Interface interface {
	HasNext() bool
	Next()
	CurrentItem() interface{}
}

type ArrayIntIterator struct {
	currentIndex  int
	modifyVersion int
	data          *ArrayIntList
}

func (i *ArrayIntIterator) HasNext() bool {
	return i.currentIndex < len(i.data.data)
}

func (i *ArrayIntIterator) Next() {
	i.currentIndex++
}

func (i *ArrayIntIterator) CurrentItem() interface{} {
	return i.data.data[i.currentIndex]
}

func (i *ArrayIntIterator) CheckIsModify() bool {
	return i.data.modifyVersion != i.modifyVersion
}
