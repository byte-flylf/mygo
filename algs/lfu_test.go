package algs

import (
	"fmt"
	"testing"
)

func TestLFU(t *testing.T) {
	c := NewLFUCache()
	c.Set("a", "a")
	if v := c.Get("a"); v != "a" {
		t.Errorf("Value was not saved: %v != 'a'", v)
	}
	if l := c.Len(); l != 1 {
		t.Errorf("Length was not updated: %v != 1", l)
	}

	c.Set("b", "b")
	if v := c.Get("b"); v != "b" {
		t.Errorf("Value was not saved: %v != 'b'", v)
	}
	if l := c.Len(); l != 2 {
		t.Errorf("Length was not updated: %v != 2", l)
	}

	c.Get("a")
	evicted := c.Evict(1)
	if v := c.Get("a"); v != "a" {
		t.Errorf("Value was improperly evicted: %v != 'a'", v)
	}
	if v := c.Get("b"); v != nil {
		t.Errorf("Value was not evicted: %v", v)
	}
	if l := c.Len(); l != 1 {
		t.Errorf("Length was not updated: %v != 1", l)
	}
	if evicted != 1 {
		t.Errorf("Number of evicted items is wrong: %v != 1", evicted)
	}
}

func TestBoundsMgmt(t *testing.T) {
	c := NewLFUCache()
	c.UpperBound = 10
	c.LowerBound = 5

	for i := 0; i < 100; i++ {
		c.Set(fmt.Sprintf("%v", i), i)
	}
	if c.Len() > 10 {
		t.Errorf("Bounds management failed to evict properly: %v", c.Len())
	}
}
