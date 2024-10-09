package tests

import (
	"prac/utils"
	"testing"
)

func TestDoesExist(t *testing.T) {
	bl := utils.CreateBloomFilter(10, 0.01)
	var err error

	if err = bl.Set("ok1"); err != nil {
		t.Error(err)
	}
	if err = bl.Set("doksds"); err != nil {
		t.Error(err)
	}
	if err = bl.Set("oops"); err != nil {
		t.Error(err)
	}
	if err = bl.Set("wyvern"); err != nil {
		t.Error(err)
	}

	if bl.DoesExist("ok1") == false {
		t.Error("ok1 was supposed to exist but not found !!!")
	}

	if bl.DoesExist("doksds") == false {
		t.Error("doksds was supposed to exist but not found !!!")
	}

	if bl.DoesExist("oops") == false {
		t.Error("oops was supposed to exist but not found !!!")
	}
	if bl.DoesExist("wyvern") == false {
		t.Error("wyvern was supposed to exist but not found !!!")
	}

	if bl.DoesExist("ok2") == true {
		t.Error("ok2 was supposed to not exist but found in bloom filter !!!")
	}

}

func TestSetGetBit(t *testing.T) {
	bl := utils.CreateBloomFilter(10, 0.01)
	var err error

	if err = bl.SetBit(10); err != nil {
		t.Error(err)
	}

	if bl.GetBit(10) == false {
		t.Error("Expected 10th bit to be set to 1 but found 0")
	}

	if err = bl.SetBit(20); err != nil {
		t.Error(err)
	}

	if bl.GetBit(20) == false {
		t.Error("Expected 20th bit to be set to 1 but found 0")
	}

	if bl.GetBit(5) == true {
		t.Error("Expected 5th bit to be set to 0 but found 1")
	}
}
