package list

// Element holds one value coming from a list
type Element interface {
	Value() interface{}
}

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
	Get(i int) (Element, error)
	GetFunc(m func(v interface{}) bool) Element
	Set(i int, v interface{}) error
	SetFunc(m func(v interface{}) bool, v interface{})

	Copy() List
	ToArray() []interface{}

	InsertAt(i int, v interface{}) error
	RemoveAt(i int) (Element, error)
	RemoveFirstOccurrence(v interface{}) bool
	RemoveLastOccurrence(v interface{}) bool
	Pop() (Element, bool)
	Push(v interface{})
	PushList(l2 List)
	Shift() (Element, bool)
	Unshift(v interface{})
	UnshiftList(l2 List)

	Contains(v interface{}) bool
	IndexOf(v interface{}) (int, bool)
	LastIndexOf(v interface{}) (int, bool)
}
