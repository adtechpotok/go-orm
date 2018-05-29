package go_orm

import (
	"testing"
	"fmt"
)

func TestBaseDbModelInit(t *testing.T) {
	var v BaseDbModel
	v = BaseDbModel{}
	v.GetCache()
	if !v.isInit {
		t.Error("Base Model should be inited before GetCache")
	}
}
func TestBaseDbModelGet(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	res := v.FindInCache(1)
	if !v.isInit {
		t.Error("Base Model should be inited before FindInCache")
	}

	if res != nil {
		t.Error("Not empty result")
	}

	v = BaseDbModel{}
	v.AddToCache(testPotok{true, 1, 2})
	if !v.isInit {
		t.Error("Base Model should be inited before GetCache")
	}
	res = v.FindInCache(1)
	if res == nil {
		t.Error("Empty result after put")
	}

}
func TestBaseDbModelGetInactive(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	v.AddToCache(testPotok{false, 1, 2})

	res := v.FindInCache(1)
	if res != nil {
		t.Error("Found inactive result")
	}

	if v.Len() != 1 {
		t.Error("Error result count")
	}
}

func TestBaseDbModelClear(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	v.AddToCache(testPotok{true, 1, 2})

	v.ClearCache()
	res := v.FindInCache(1)

	if res != nil {
		t.Error("Not empty result after clear")
	}
}

func TestAltIndex(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	v.RegisterIntIndex("test", AltIndexTest)
	v.RegisterStringIndex("testString", AltStringTest)

	v.AddToCache(testPotok{true, 1, 2})

	res := v.FindIndex("test", 2, true)
	if res == nil {
		t.Error("Empty result after put")
	}

	v.AddToCache(testPotok{true, 1, 3})

	res = v.FindIndex("test", 2, true)
	if res != nil {
		t.Error("Not empty result after put another alter index")
	}

	res = v.FindIndex("test", 3, true)
	if res == nil {
		t.Error("Empty result after one more")
	}
	res = v.FindIndex("testString", "3", true)
	if res == nil {
		t.Error("Empty result in string alter index")
	}
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	v := BaseDbModel{}
	v.FindIndex("test", true, true)

}

func AltIndexTest(p PotokOrm) int {
	return p.(testPotok).altId
}

func AltStringTest(p PotokOrm) string {
	return fmt.Sprintf("%d", p.(testPotok).altId)
}

type testPotok struct {
	active bool
	id     int
	altId  int
}

func (m testPotok) IsActive() bool { return m.active }
func (m testPotok) GetId() int     { return m.id }
