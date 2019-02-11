package cxserver

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/mit-dci/lit/btcutil/base58"

	"github.com/mit-dci/lit/wallit"
)

// NewAddressLTC returns a new address based on the keygen retrieved from the wallet
func (server *OpencxServer) NewAddressLTC(username string) (string, error) {
	// No really what is this
	return server.getLTCAddrFunc()(username)
}

// NewAddressBTC returns a new address based on the keygen retrieved from the wallet
func (server *OpencxServer) NewAddressBTC(username string) (string, error) {
	// What is this
	return server.getBTCAddrFunc()(username)
}

// NewAddressVTC returns a new address based on the keygen retrieved from the wallet
func (server *OpencxServer) NewAddressVTC(username string) (string, error) {
	// Is this currying
	return server.getVTCAddrFunc()(username)
}

// getVTCAddrFunc is used by NewAddressVTC as well as UpdateAddresses to call the address closure
func (server *OpencxServer) getVTCAddrFunc() func(string) (string, error) {
	return GetAddrFunction(server.OpencxVTCWallet)
}

// getBTCAddrFunc is used by NewAddressBTC as well as UpdateAddresses to call the address closure
func (server *OpencxServer) getBTCAddrFunc() func(string) (string, error) {
	return GetAddrFunction(server.OpencxBTCWallet)
}

// getLTCAddrFunc is used by NewAddressLTC as well as UpdateAddresses to call the address closure
func (server *OpencxServer) getLTCAddrFunc() func(string) (string, error) {
	return GetAddrFunction(server.OpencxLTCWallet)
}

// GetAddrFunction returns a function that can safely be called by the DB
func GetAddrFunction(wallet *wallit.Wallit) func(string) (string, error) {
	pubKeyHashAddrID := wallet.Param.PubKeyHashAddrID
	coinType := wallet.Param.HDCoinType
	return func(username string) (addr string, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("Problem with address closure: \n%s", err)
			}
		}()

		sha := sha256.New()
		if _, err = sha.Write([]byte(username)); err != nil {
			return
		}

		// TODO: figure out if necessary
		// We mod by 0x80000000 to make sure it's not hardened
		data := binary.BigEndian.Uint32(sha.Sum(nil)[:]) % 0x80000000

		// create a new keygen based on the data
		keygen := wallit.GetWalletKeygen(data, coinType)

		// Register the address with the wallet
		if err = wallet.AddPorTxoAdr(keygen); err != nil {
			return
		}

		// get the pubkeyhash from the keygen
		nAdr160 := wallet.PathPubHash160(keygen)

		// Encode the pkhash that we're given in base58check, return as a string
		addr = base58.CheckEncode(nAdr160[:], pubKeyHashAddrID)

		return
	}
}

// UpdateAddresses updates all the addresses in the DB with the address functions defined.
func (server *OpencxServer) UpdateAddresses() error {

	// Lock ingest so they wait for the db
	server.LockIngests()

	// Call DB method with functions
	if err := server.OpencxDB.UpdateDepositAddresses(server.getLTCAddrFunc(), server.getBTCAddrFunc(), server.getVTCAddrFunc()); err != nil {
		return err
	}

	server.UnlockIngests()

	return nil
}
