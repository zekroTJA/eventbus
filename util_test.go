package eventbus

import "testing"

func TestOpt(t *testing.T) {
	var r int

	r = opt([]int{})
	if r != 0 {
		t.Errorf("r was %d != expected 0", r)
	}

	r = opt([]int{2})
	if r != 2 {
		t.Errorf("r was %d != expected 2", r)
	}

	r = opt([]int{2, 3, 4})
	if r != 2 {
		t.Errorf("r was %d != expected 2", r)
	}

	r = opt([]int{}, 4)
	if r != 4 {
		t.Errorf("r was %d != expected 4", r)
	}

	r = opt([]int{5}, 4)
	if r != 5 {
		t.Errorf("r was %d != expected 5", r)
	}
}
