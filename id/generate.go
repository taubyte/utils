package id

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	mh "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
)

var (
	RandomnessLenth   = 32
	RandomnessTimeout = 1 * time.Second
)

// best effort to get random bytes
func getRandom() []byte {
	ctx, ctx_cancel := context.WithTimeout(context.Background(), RandomnessTimeout)
	defer ctx_cancel()
	for {
		select {
		case <-ctx.Done():
			return []byte{}
		default:
			r := make([]byte, RandomnessLenth)
			n, err := io.ReadFull(rand.Reader, r)
			if n == len(r) && err == nil {
				return r
			}
		}
	}
}

// Generate hash ID
// provided parameters + current timestamp + randomness
func Generate(args ...interface{}) string {
	// check https://github.com/ipsn/go-ipfs/blob/master/gxlibs/github.com/libp2p/go-libp2p-peer/peer.go#L154 and https://github.com/ipsn/go-ipfs/blob/master/gxlibs/github.com/libp2p/go-libp2p-peer/peer.go#L39
	args = append(args, getRandom())
	args = append(args, time.Now().Unix())
	hash, _ := mh.Sum([]byte(fmt.Sprint(args...)), mh.SHA2_256, -1)
	return hash.B58String()
}
