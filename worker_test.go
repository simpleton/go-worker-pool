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

func TestHello(t *testing.T) {
	t.Log("Hello Test")
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
