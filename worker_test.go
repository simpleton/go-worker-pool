package workerpool

import (
	"testing"
  "time"
  "log"
)

type simpleTask struct {
  name string
}

func (s *simpleTask) Run() {
  log.Println(s.name)
  time.Sleep(time.Second)
}

func TestNew(t *testing.T) {
	t.Log("Hello New")
  {
    pool := New(3)
    names := []string {"a", "aa", "aaa", "bbb"}
    for _, name := range names {
      np := simpleTask{
        name: name,
      }
      pool.Submit(&np)
    }
  }
}

func TestNewDefault(t *testing.T) {
	t.Log("Hello New")
  {
    pool := NewDefault()
    names := []string {"a", "aa", "aaa", "bbb"}
    for _, name := range names {
      np := simpleTask{
        name: name,
      }
      pool.Submit(&np)
    }
  }
}
