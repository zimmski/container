package list

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	l := New()

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	Nil(t, l.Pop())

	v := []interface{}{1, "a", 2, "b"}
	vLen := len(v)

	for i := 0; i < vLen; i++ {
		l.Push(v[i])

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value, v[0])
		Equal(t, l.Last().Value, v[i])
	}

	i := 0
	n := l.First()

	for i < vLen && n != nil {
		Equal(t, v[i], n.Value)

		i++
		n = n.Next()
	}

	Equal(t, i, vLen)
	Nil(t, n)

	i = vLen - 1
	n = l.Pop()

	for i > -1 && n != nil {
		Equal(t, v[i], n.Value)
		Nil(t, n.Next())
		Equal(t, l.Len(), i)

		i--
		n = l.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	Equal(t, l.Len(), 0)
}
