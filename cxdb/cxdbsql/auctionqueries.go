package cxdbsql

import (
	"database/sql"
	"fmt"

	"github.com/mit-dci/opencx/match"
)

// PlaceAuctionPuzzle puts a puzzle and ciphertext in the datastore.
func (db *DB) PlaceAuctionPuzzle(encryptedOrder *match.EncryptedAuctionOrder) (err error) {

	var tx *sql.Tx
	if tx, err = db.DBHandler.Begin(); err != nil {
		err = fmt.Errorf("Error when beginning transaction for NewAuction: %s", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("Error while placing puzzle order: \n%s", err)
			return
		}
		err = tx.Commit()
	}()

	// We don't really care about the result when trying to use a schema
	if _, err = tx.Exec("USE " + db.puzzleSchema + ";"); err != nil {
		err = fmt.Errorf("Error trying to use auction schema: %s", err)
		return
	}

	var puzzleBytes []byte
	if puzzleBytes, err = encryptedOrder.OrderPuzzle.Serialize(); err != nil {
		err = fmt.Errorf("Error serializing puzzle: %s", err)
		return
	}

	var concatBytes []byte
	concatBytes = append(puzzleBytes, encryptedOrder.OrderCiphertext...)

	// We concatenate ciphertext and puzzle and set "selected" to 1 by default
	placeAuctionPuzzleQuery := fmt.Sprintf("INSERT INTO %s VALUES (%x, %x, 'TRUE');", db.puzzleTable, encryptedOrder.IntendedAuction, concatBytes)
	if _, err = tx.Exec(placeAuctionPuzzleQuery); err != nil {
		err = fmt.Errorf("Error adding auction puzzle to puzzle orders: %s", err)
		return
	}

	// TODO
	return
}

// PlaceAuctionOrder places an order in the unencrypted datastore.
func (db *DB) PlaceAuctionOrder(*match.AuctionOrder) (err error) {

	// TODO
	return
}

// ViewAuctionOrderBook takes in a trading pair and auction ID, and returns auction orders.
func (db *DB) ViewAuctionOrderBook(tradingPair *match.Pair, auctionID [32]byte) (sellOrderBook []*match.AuctionOrder, buyOrderBook []*match.AuctionOrder, err error) {

	// TODO
	return
}

// ViewAuctionPuzzleBook takes in an auction ID, and returns encrypted auction orders, and puzzles.
// You don't know what auction IDs should be in the orders encrypted in the puzzle book, but this is
// what was submitted.
func (db *DB) ViewAuctionPuzzleBook(auctionID [32]byte) (orders []*match.EncryptedAuctionOrder, err error) {

	// TODO
	return
}

// NewAuction takes in an auction ID, and creates a new auction, returning the "height"
// of the auction.
func (db *DB) NewAuction(auctionID [32]byte) (height uint64, err error) {

	var tx *sql.Tx
	if tx, err = db.DBHandler.Begin(); err != nil {
		err = fmt.Errorf("Error when beginning transaction for NewAuction: %s", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("Error while creating new auction: \n%s", err)
			return
		}
		err = tx.Commit()
	}()

	// We don't really care about the result when trying to use a schema
	if _, err = tx.Exec("USE " + db.auctionOrderSchema + ";"); err != nil {
		err = fmt.Errorf("Error trying to use auction order schema: %s", err)
		return
	}

	auctionNumQuery := fmt.Sprintf("SELECT MAX(auctionNumber) FROM %s;", db.auctionOrderTable)
	if err = tx.QueryRow(auctionNumQuery).Scan(height); err != nil {
		err = fmt.Errorf("Could not find maximum auction number when creating new auction: %s", err)
		return
	}

	// Insert the new auction ID w/ current max height + 1
	height++
	insertNewAuctionQuery := fmt.Sprintf("INSERT INTO %s VALUES (%x, %d);", db.auctionOrderTable, auctionID, height)
	if _, err = tx.Exec(insertNewAuctionQuery); err != nil {
		err = fmt.Errorf("Error inserting new auction ID when creating new auction: %s", err)
		return
	}

	return
}
