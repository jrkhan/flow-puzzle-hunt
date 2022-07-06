package distributed

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

var client *storage.Client

func init() {
	if _, has := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !has {
		fmt.Print("No GOOGLE_APPLICATION_CREDENTIALS found, this might still be okay")
	}
	ctx := context.Background()
	var err error
	fmt.Print()
	client, err = storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	b := client.Buckets(ctx, project())
	_, err = b.Next()
	if err != nil {
		panic(err)
	}
}

func project() string {
	val, has := os.LookupEnv("GCP_PROJECT")
	if !has {
		return "fuzzle-353902"
	}
	return val
}

var lockBucket = "puzzle-alley-mint-lock"
var lockObject = "lock-object"

var autoUnlock = true
var autoUnlockTime = time.Second * 10
var awaitTime = time.Second * 10

func DL(ctx context.Context, bucket string, obj string) (Unlockable, error) {

	b := client.Bucket(bucket)

	oh := b.Object(obj).If(storage.Conditions{DoesNotExist: true})
	w := oh.NewWriter(ctx)

	_, err := fmt.Fprintf(w, " ")
	if err == nil {
		// we have the lock
		closeErr := w.Close()
		if closeErr == nil {
			deleteHandle := b.Object(obj)
			return Unlockable{deleteHandle, true}, nil
		} else {
			err = closeErr
		}
	}

	// this is expected when we do not have the lock
	if !strings.Contains(err.Error(), "conditionNotMet") {
		return Unlockable{}, err
	}

	if autoUnlock {
		possibleDelete := b.Object(obj)
		attrs, err := possibleDelete.Attrs(ctx)
		if err != nil {
			fmt.Print(err)
			return Unlockable{}, err
		}

		aquired := attrs.Updated
		if aquired.Before(time.Now().Add(-autoUnlockTime)) {

			err = possibleDelete.Delete(ctx)
			if err != nil {
				return Unlockable{}, err
			}
		}
	}
	return Unlockable{}, nil
}

type Unlockable struct {
	oh       *storage.ObjectHandle
	Acquired bool
}

func (u *Unlockable) Unlock(ctx context.Context) error {
	return u.oh.Delete(ctx)
}

func AwaitLock(ctx context.Context) (Unlockable, error) {

	start := time.Now()
	hasLock := false

	for !hasLock && time.Now().Before(start.Add(awaitTime)) {
		ul, err := DL(ctx, lockBucket, lockObject)
		if err != nil {
			return Unlockable{}, err
		}
		if !ul.Acquired {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		return ul, nil
	}
	return Unlockable{}, errors.New("unable to aquire lock")
}
