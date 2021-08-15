package iterator

type List interface {
	Iterator() Interface
}

type ArrayIntList struct {
	modifyVersion int
	data          []int
}

func NewArrayIntList(data []int) List {
	return &ArrayIntList{
		modifyVersion: 0,
		data:          data,
	}
}

func (a *ArrayIntList) Iterator() Interface {
	return &ArrayIntIterator{
		currentIndex: 0,
		data:         a,
	}
}
