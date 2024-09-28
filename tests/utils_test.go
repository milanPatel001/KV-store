package tests

import (
	"math"
	"prac/utils"
	"testing"
)

func TestSetMaxValue(t *testing.T) {

	mv1 := utils.SetMaxValue[int]()

	if mv1 != math.MaxInt {
		t.Errorf("Max value expected %v, got %v", math.MaxInt, mv1)
	}

	mv2 := utils.SetMaxValue[int64]()

	if mv2 != math.MaxInt64 {
		t.Errorf("Max value expected %v, got %v", math.MaxInt64, mv2)
	}

	mv3 := utils.SetMaxValue[uint64]()

	if mv3 != uint64(math.MaxUint64) {
		t.Errorf("Max value expected %v, got %v", uint64(math.MaxUint64), mv3)
	}

	mv4 := utils.SetMaxValue[float64]()

	if mv4 != math.MaxFloat64 {
		t.Errorf("Max value expected %v, got %v", math.MaxFloat64, mv4)
	}

	mv5 := utils.SetMaxValue[string]()

	if mv5 != "INF" {
		t.Errorf("Max value expected %v, got %v", "INF", mv4)
	}

}
