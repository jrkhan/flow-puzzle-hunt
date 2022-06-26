package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

var assets = []string{
	"https://ipfs.io/ipfs/QmU4SndLCLdbGtnthGm7ZbwKSw6cUs6WRUPdJ18pUrePQc?filename=piece-1.png",
}

type Piece struct {
	Index int
	IPFS  string `json:"ipfs"`
	CSS   string `json:"css"`
}

func main() {
	mp := map[string]int{}
	for i := 0; i < 20; i++ {
		uuidWithHyphen := uuid.New()
		mp[uuidWithHyphen.String()] = i + 1

		//fmt.Printf("\"%v\": %v,\n", uuidWithHyphen, i+1)

		//uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
		//fmt.Println(uuid)
	}
	rs, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		panic("unable to build result")
	}
	fmt.Println(string(rs))
}
