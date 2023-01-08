package filequeue_test

import (
	"github.com/logxxx/utils/filequeue"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileQueue2(t *testing.T) {
	type Foo struct {
		Name      string `json:"name"`
		Age       int    `json:'age'`
		IsTeacher bool   `json:"is_teacher"`
	}

	queue := filequeue.NewFileQueue("test_queue1")

	req1 := &Foo{
		Name:      "11",
		Age:       22,
		IsTeacher: true,
	}

	err := queue.Push(req1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFileQueue(t *testing.T) {

	type Foo struct {
		Name      string `json:"name"`
		Age       int    `json:'age'`
		IsTeacher bool   `json:"is_teacher"`
	}

	queue := filequeue.NewFileQueue("test_queue1")

	req1 := &Foo{
		Name:      "11",
		Age:       22,
		IsTeacher: true,
	}

	{
		err := queue.Push(req1)
		if err != nil {
			t.Fatal(err)
		}

		resp1 := &Foo{}
		err = queue.Pop(resp1)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, req1.Name, resp1.Name)
		assert.Equal(t, req1.Age, resp1.Age)
		assert.Equal(t, req1.IsTeacher, resp1.IsTeacher)

		resp2 := &Foo{}
		err = queue.Pop(resp2)
		if err != nil {
			t.Fatal(err)
		}

		assert.Empty(t, resp2)

		queue.Clean()

	}

	{
		for i := 0; i <= 1000; i++ {
			req1.Age = i
			err := queue.Push(req1)
			if err != nil {
				t.Fatal(err)
			}
		}

		for i := 0; i <= 1000; i++ {
			resp1 := &Foo{}
			err := queue.MustPop(resp1)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, req1.Name, resp1.Name)
			assert.Equal(t, i, resp1.Age)
			assert.Equal(t, req1.IsTeacher, resp1.IsTeacher)
		}
		resp2 := &Foo{}
		err := queue.MustPop(resp2)
		assert.Equal(t, filequeue.ErrEmpty, err)

		queue.Clean()

	}

}
