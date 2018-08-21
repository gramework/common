package common

import (
	"crypto/rand"
	"io"
	"math/big"

	"github.com/pkg/errors"
)

var (
	secp256k1N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	keyBitLen = secp256k1N.BitLen()/8 + 8
)

// GenerateKey generates a private key using random source rand.
func GenerateKey() (b []byte, err error) {
	/* See Certicom's SEC1 3.2.1, pg.23 */
	/* See NSA's Suite B Implementer’s Guide to FIPS 186-3 (ECDSA) A.1.1, pg.18 */

	/* Select private key d randomly from [1, n) */

	/* Read N bit length random bytes + 64 extra bits  */
	b = make([]byte, keyBitLen)
	_, err = io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, errors.Errorf("Reading random reader: %v", err)
	}

	d := new(big.Int).SetBytes(b)

	/* Mod n-1 to shift d into [0, n-1) range */
	d.Mod(d, new(big.Int).Sub(secp256k1N, big.NewInt(1)))
	/* Add one to shift d to [1, n) range */
	d.Add(d, big.NewInt(1))

	return d.Bytes(), nil
}
