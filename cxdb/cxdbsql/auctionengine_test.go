package cxdbsql

import (
	"testing"

	"github.com/mit-dci/lit/coinparam"
	"github.com/mit-dci/opencx/match"
)

var (
	litereg, _       = match.AssetFromCoinParam(&coinparam.LiteRegNetParams)
	btcreg, _        = match.AssetFromCoinParam(&coinparam.RegressionNetParams)
	testAuctionOrder = &match.AuctionOrder{
		Pubkey:     [...]byte{0x02, 0xe7, 0xb7, 0xcf, 0xcf, 0x42, 0x2f, 0xdb, 0x68, 0x2c, 0x85, 0x02, 0xbf, 0x2e, 0xef, 0x9e, 0x2d, 0x87, 0x67, 0xf6, 0x14, 0x67, 0x41, 0x53, 0x4f, 0x37, 0x94, 0xe1, 0x40, 0xcc, 0xf9, 0xde, 0xb3},
		Nonce:      [2]byte{0x00, 0x00},
		AuctionID:  [32]byte{0xde, 0xad, 0xbe, 0xef},
		AmountWant: 100000,
		AmountHave: 10000,
		Side:       "buy",
		TradingPair: match.Pair{
			AssetWant: btcreg,
			AssetHave: litereg,
		},
		Signature: []byte{0x1b, 0xd6, 0x0f, 0xd3, 0xec, 0x5b, 0x73, 0xad, 0xa9, 0x8a, 0x92, 0x79, 0x82, 0x0f, 0x8e, 0xab, 0xf8, 0x8f, 0x47, 0x6e, 0xc3, 0x15, 0x33, 0x72, 0xd9, 0x90, 0x51, 0x41, 0xfd, 0x0a, 0xa1, 0xa2, 0x4a, 0x73, 0x75, 0x4c, 0xa5, 0x28, 0x4a, 0xc2, 0xed, 0x5a, 0xe9, 0x33, 0x22, 0xf4, 0x41, 0x1f, 0x9d, 0xd1, 0x78, 0xb9, 0x17, 0xd4, 0xe9, 0x72, 0x51, 0x7f, 0x5b, 0xd7, 0xe5, 0x12, 0xe7, 0x69, 0xb0},
	}
	testEncryptedOrder, _ = testAuctionOrder.TurnIntoEncryptedOrder(testStandardAuctionTime)
	testEncryptedBytes, _ = testEncryptedOrder.Serialize()
	// examplePubkeyOne =
)

func TestCreateEngineAllParams(t *testing.T) {
	var err error

	var killFunc func(t *testing.T)
	if killFunc, err = createUserAndDatabase(); err != nil {
		t.Skipf("Error creating user and database, skipping: %s", err)
	}

	var pairList []*match.Pair
	if pairList, err = match.GenerateAssetPairs(constCoinParams()); err != nil {
		t.Errorf("Error creating asset pairs from coin list: %s", err)
	}

	for _, pair := range pairList {
		if _, err = CreateAuctionEngineWithConf(pair, testConfig()); err != nil {
			t.Errorf("Error creating auction engine for pair: %s", err)
		}
	}

	killFunc(t)
}

func TestPlaceSingleAuctionOrder(t *testing.T) {
	var err error

	var killFunc func(t *testing.T)
	if killFunc, err = createUserAndDatabase(); err != nil {
		t.Skipf("Error creating user and database, skipping: %s", err)
	}

	var engine match.AuctionEngine
	if engine, err = CreateAuctionEngineWithConf(&testEncryptedOrder.IntendedPair, testConfig()); err != nil {
		t.Errorf("Error creating auction engine for pair: %s", err)
	}

	var idStruct *match.AuctionID = new(match.AuctionID)
	if err = idStruct.UnmarshalBinary(testEncryptedOrder.IntendedAuction[:]); err != nil {
		t.Errorf("Error unmarshalling auction ID: %s", err)
	}

	if _, err = engine.PlaceAuctionOrder(testAuctionOrder, idStruct); err != nil {
		t.Errorf("Error placing auction order: %s", err)
	}

	killFunc(t)
}
