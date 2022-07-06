package distributed

import (
	"context"
	"encoding/binary"
	"io/ioutil"
)

var intObj = "atomicSequenceNumber"

func NextValue(ctx context.Context, proposed uint64) (uint64, error) {
	l, err := AwaitLock(ctx)
	defer l.Unlock(ctx)
	if err != nil {
		return 0, err
	}
	obj := client.Bucket(lockBucket).Object(intObj)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return 0, err
	}
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	lastUsed := binary.LittleEndian.Uint64(contents)

	if proposed > lastUsed {
		err = overwrite(ctx, proposed)
		if err != nil {
			return 0, err
		}
		return proposed, nil
	}
	next := lastUsed + 1
	err = overwrite(ctx, next)
	if err != nil {
		return 0, err
	}
	return next, nil
}

func overwrite(ctx context.Context, val uint64) error {
	obj := client.Bucket(lockBucket).Object(intObj)
	w := obj.NewWriter(ctx)
	defer w.Close()
	return binary.Write(w, binary.LittleEndian, val)
}
