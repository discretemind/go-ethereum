package merkle

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zmitton/go-merklemountainrange/db"
	"github.com/zmitton/go-merklemountainrange/digest"
	"github.com/zmitton/go-merklemountainrange/mmr"
	"golang.org/x/crypto/sha3"
	"sync"
)

var Mmr *mmrTree

type mmrTree struct {
	lock     sync.Mutex
	instance *mmr.Mmr
}

func (h *mmrTree) Hash(value interface{}) (root common.Hash) {
	fmt.Println("Calculate header hash: ", value)
	//buf := &bytes.Buffer{}
	//rlp.Encode(buf, value)
	//var leafHash common.Hash
	var dataHash [64]byte
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, value)
	hw.Sum(dataHash[:32])

	h.lock.Lock()
	index := h.instance.GetLeafLength()
	fmt.Printf("Lead index %d %x", index, dataHash)
	h.instance.Append(dataHash[:], index)
	//TODO: modify GetRoot() to return (may be include MMR implementation into the project without reference to 3rdParty lib)
	hash := h.instance.GetRoot()
	index = h.instance.GetLeafLength()
	fmt.Printf("Root hash %d %x", index, hash)
	h.lock.Unlock()
	copy(root[:32], hash[:32])
	return
}

//TODO: init just for test
func init() {
	//TODO: move data path to config
	d := db.CreateFilebaseddb("./data/mmr/data.mmr", 64)
	//d := db.NewMemorybaseddb(0, map[int64][]byte{})
	//fileDb := db.OpenFilebaseddb("./data/mmr/data.mmr")
	Mmr = &mmrTree{
		instance: mmr.New(digest.Keccak256FlyHash, d),
	}
}
