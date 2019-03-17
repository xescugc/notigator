package assets

import (
	"os"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

//go:generate go-bindata-assetfs -nometadata -ignore=bindata\.go -ignore=doc\.go -ignore=assets\.go -pkg assets ./...

// AssetFS it's a wrapper to make the assetFS public
func AssetFS() *assetfs.AssetFS {
	// Ideally I should only call this function but for some reason it did not work as expected so
	// I had to copy the content and do it correctly
	//return assetFS()
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo}
}
