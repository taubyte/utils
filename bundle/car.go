package bundle

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-libipfs/blocks"
	"github.com/ipfs/go-unixfsnode/data/builder"
	"github.com/ipld/go-car/v2"
	"github.com/ipld/go-car/v2/blockstore"
	dagpb "github.com/ipld/go-codec-dagpb"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
)

// reference: https://github.com/ipld/go-car/blob/master/cmd/car/create.go
func Car(ctx context.Context, source, target string) error {
	// Make sure source exists
	_, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create car file
	carFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("Create `%s`: %s", target, err)
	}
	defer carFile.Close()

	// make a cid with the right length that we eventually will patch with the root.
	hasher, err := multihash.GetHasher(multihash.SHA2_256)
	if err != nil {
		return err
	}

	digest := hasher.Sum([]byte{})

	hash, err := multihash.Encode(digest, multihash.SHA2_256)
	if err != nil {
		return err
	}

	_cid := cid.NewCidV1(uint64(multicodec.DagPb), hash)

	blockReadWriter, err := blockstore.OpenReadWriteFile(carFile, []cid.Cid{_cid})
	if err != nil {
		return err
	}

	// Write the unixfs blocks into the store.
	root, err := writeFiles(ctx, blockReadWriter, source)
	if err != nil {
		return err
	}

	if err := blockReadWriter.Finalize(); err != nil {
		return err
	}

	// re-open/finalize with the final root.
	err = car.ReplaceRootsInFile(carFile.Name(), []cid.Cid{root})
	if err != nil {
		return err
	}

	return nil
}

func writeFiles(ctx context.Context, bs *blockstore.ReadWrite, paths ...string) (cid.Cid, error) {
	// paths is the files
	ls := cidlink.DefaultLinkSystem()
	ls.TrustedStorage = true
	ls.StorageReadOpener = func(_ ipld.LinkContext, l ipld.Link) (io.Reader, error) {
		cl, ok := l.(cidlink.Link)
		if !ok {
			return nil, fmt.Errorf("not a cidlink")
		}

		blk, err := bs.Get(ctx, cl.Cid)
		if err != nil {
			return nil, fmt.Errorf("getting cid `%s` failed with: %s", cl.Cid, err)
		}

		return bytes.NewBuffer(blk.RawData()), nil
	}
	ls.StorageWriteOpener = func(_ ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
		buf := bytes.NewBuffer(nil)
		return buf, func(l ipld.Link) error {
			cl, ok := l.(cidlink.Link)
			if !ok {
				return fmt.Errorf("not a cidlink")
			}

			blk, err := blocks.NewBlockWithCid(buf.Bytes(), cl.Cid)
			if err != nil {
				return fmt.Errorf("creating new block with cid `%s` failed with: %s", cl.Cid, err)
			}

			err = bs.Put(ctx, blk)
			if err != nil {
				return fmt.Errorf("putting block failed with: %s", err)
			}

			return nil
		}, nil
	}

	topLevel := make([]dagpb.PBLink, 0, len(paths))
	for _, p := range paths {
		l, size, err := builder.BuildUnixFSRecursive(p, &ls)
		if err != nil {
			return cid.Undef, err
		}

		name := path.Base(p)
		entry, err := builder.BuildUnixFSDirectoryEntry(name, int64(size), l)
		if err != nil {
			return cid.Undef, err
		}

		topLevel = append(topLevel, entry)
	}

	// make a directory for the file(s).
	root, _, err := builder.BuildUnixFSDirectory(topLevel, &ls)
	if err != nil {
		return cid.Undef, nil
	}

	rcl, ok := root.(cidlink.Link)
	if !ok {
		return cid.Undef, fmt.Errorf("could not interpret %s", root)
	}

	return rcl.Cid, nil
}
