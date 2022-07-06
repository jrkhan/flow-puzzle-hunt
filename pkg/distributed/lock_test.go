package distributed

import (
	"context"
	"testing"
	"time"
)

func TestGetLock(t *testing.T) {
	ctx := context.Background()
	ul, err := AwaitLock(ctx)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ul.Unlock(ctx)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestLockContention(t *testing.T) {
	ctx := context.Background()
	ul, err := AwaitLock(ctx)
	defer ul.Unlock(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	// get the lock again - it should time out
	awaitTime = time.Millisecond * 200
	_, err = AwaitLock(ctx)
	if err == nil {
		t.Log("Expected second lock was not acquired")
		t.FailNow()
	}
	_, err = AwaitLock(ctx)
	if err == nil {
		t.Log("Expected third lock was not acquired")
		t.FailNow()
	}

	// get the lock again - now we can aquire, old one timed out
	awaitTime = time.Second * 2
	autoUnlockTime = time.Millisecond * 50
	v, err := AwaitLock(ctx)
	defer v.Unlock(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

}
