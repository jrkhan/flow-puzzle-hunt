package distributed

import (
	"context"
	"testing"
)

func TestIncreasing(t *testing.T) {
	intObj = "atomicSequenceNumber-test"
	ctx := context.Background()
	overwrite(ctx, 0)

	next, err := NextValue(ctx, 0)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if next != 1 {
		t.Logf("expected 1 but got %v", next)
		t.FailNow()
	}

	next, err = NextValue(ctx, 0)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if next != 2 {
		t.Logf("expected 2 but got %v", next)
		t.FailNow()
	}

	next, err = NextValue(ctx, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if next != 10 {
		t.Logf("expected 10 but got %v", next)
		t.FailNow()
	}
}
