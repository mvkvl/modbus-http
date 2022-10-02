package queue

import (
	"fmt"
	"testing"
)

func TestQueueInsertRemove(t *testing.T) {

	q := CreateQueue()

	val := "item 1"
	q.Insert(val)

	exp := 1
	if q.Size() != exp {
		t.Errorf("expected size '%d', got '%d' instead", exp, q.Size())
	}

	v, err := q.Remove()
	if nil != err {
		t.Errorf("%q", err)
	}
	if v != val {
		t.Errorf("expected value '%s', got '%s' instead", val, v)
	}

	exp = 0
	if q.Size() != exp {
		t.Errorf("expected size '%d', got '%d' instead", exp, q.Size())
	}

}

func TestQueueInsert(t *testing.T) {

	q := CreateQueue()

	val1 := "item 1"
	val2 := "item 2"
	q.Insert(val1)
	q.Insert(val2)

	exp := 2
	if q.Size() != exp {
		t.Errorf("expected size '%d', got '%d' instead", exp, q.Size())
	}

}

func TestQueueRemove(t *testing.T) {

	q := CreateQueue()

	val1 := "item 1"
	val2 := "item 2"
	val3 := "item 3"
	val4 := "item 4"
	q.Insert(val1)
	q.Insert(val2)
	q.Insert(val3)
	q.Insert(val4)

	q.Remove()
	q.Remove()
	q.Remove()

	exp := 1
	if q.Size() != exp {
		t.Errorf("expected size '%d', got '%d' instead", exp, q.Size())
	}

}

func TestQueueRead(t *testing.T) {

	q := CreateQueue()

	q.Insert(0)
	q.Insert("item 2")
	q.Insert("item 3")
	q.Insert(1)

	v, _ := q.Remove()
	v, ok := v.(int) // ok - conversion type assertion result
	fmt.Printf("%d (%t)\n", v, ok)

	v, _ = q.Remove()
	fmt.Printf("%v\n", v)

	v, _ = q.Remove()
	fmt.Printf("%v\n", v)

	v, _ = q.Remove()
	fmt.Printf("%v\n", v)

	v, _ = q.Remove()
	fmt.Printf("%v\n", v)

}
