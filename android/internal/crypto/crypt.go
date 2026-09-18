package crypto

import (
	"github.com/xtaci/kcp-go/v5"
)

// NewBlockCrypt creates an AES BlockCrypt instance using the provided secret key.
func NewBlockCrypt(key []byte) (kcp.BlockCrypt, error) {
	return kcp.NewAESBlockCrypt(key)
}
