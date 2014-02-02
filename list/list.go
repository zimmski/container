package list

// Iterator is an iterator for a list
type Iterator interface {
	Next() bool
	Set(v interface{})
	Value() interface{}
}

// List holds a list
type List interface {
	Clear()
	Len() int

	First() Iterator
	Last() Iterator
	Get(i int) (interface{}, error)
	GetFunc(m func(v interface{}) bool) (interface{}, bool)
	Set(i int, v interface{}) error
	SetFunc(m func(v interface{}) bool, v interface{}) bool

	Copy() List
	ToArray() []interface{}

	InsertAt(i int, v interface{}) error
	RemoveAt(i int) (interface{}, error)
	RemoveFirstOccurrence(v interface{}) bool
	RemoveLastOccurrence(v interface{}) bool
	Pop() (interface{}, bool)
	Push(v interface{})
	PushList(l2 List)
	Shift() (interface{}, bool)
	Unshift(v interface{})
	UnshiftList(l2 List)

	Contains(v interface{}) bool
	IndexOf(v interface{}) (int, bool)
	LastIndexOf(v interface{}) (int, bool)
}
