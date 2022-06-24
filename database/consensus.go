package database

import (
	"fmt"
	junotypes "github.com/forbole/juno/v3/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"time"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// GetLastBlock returns the last block stored inside the database based on the heights
func (db *Db) GetLastBlock() (*dbtypes.BlockRow, error) {
	stmt := `SELECT * FROM block ORDER BY height DESC LIMIT 1`

	var blocks []dbtypes.BlockRow
	if err := db.Sqlx.Select(&blocks, stmt); err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("cannot get block, no blocks saved")
	}

	return &blocks[0], nil
}

// GetLastBlockHeight returns the last block height stored inside the database
func (db *Db) GetLastBlockHeight() (int64, error) {
	block, err := db.GetLastBlock()
	if err != nil {
		return 0, err
	}
	if block == nil {
		return 0, fmt.Errorf("block table is empty")
	}
	return block.Height, nil
}

// -------------------------------------------------------------------------------------------------------------------

// getBlockHeightTime retrieves the block at the specific time
func (db *Db) getBlockHeightTime(pastTime time.Time) (dbtypes.BlockRow, error) {
	stmt := `SELECT * FROM block WHERE block.timestamp <= $1 ORDER BY block.timestamp DESC LIMIT 1;`

	var val []dbtypes.BlockRow
	if err := db.Sqlx.Select(&val, stmt, pastTime); err != nil {
		return dbtypes.BlockRow{}, err
	}

	if len(val) == 0 {
		return dbtypes.BlockRow{}, fmt.Errorf("cannot get block time, no blocks saved")
	}

	return val[0], nil
}

// GetBlockHeightTimeMinuteAgo return block height and time that a block proposals
// about a minute ago from input date
func (db *Db) GetBlockHeightTimeMinuteAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Minute * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeHourAgo return block height and time that a block proposals
// about a hour ago from input date
func (db *Db) GetBlockHeightTimeHourAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeDayAgo return block height and time that a block proposals
// about a day (24hour) ago from input date
func (db *Db) GetBlockHeightTimeDayAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -24)
	return db.getBlockHeightTime(pastTime)
}

// -------------------------------------------------------------------------------------------------------------------

// SaveAverageBlockTimePerMin save the average block time in average_block_time_per_minute table
func (db *Db) SaveAverageBlockTimePerMin(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_minute(average_time, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time, 
        height = excluded.height
WHERE average_block_time_per_minute.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per minute: %s", err)
	}

	return nil
}

// SaveAverageBlockTimePerHour save the average block time in average_block_time_per_hour table
func (db *Db) SaveAverageBlockTimePerHour(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_hour(average_time, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time,
        height = excluded.height
WHERE average_block_time_per_hour.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per hour: %s", err)
	}

	return nil
}

// SaveAverageBlockTimePerDay save the average block time in average_block_time_per_day table
func (db *Db) SaveAverageBlockTimePerDay(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_day(average_time, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time,
        height = excluded.height
WHERE average_block_time_per_day.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per day: %s", err)
	}

	return nil
}

// SaveAverageBlockTimeGenesis save the average block time in average_block_time_from_genesis table
func (db *Db) SaveAverageBlockTimeGenesis(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_from_genesis(average_time ,height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time, 
        height = excluded.height
WHERE average_block_time_from_genesis.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time since genesis: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveGenesis save the given genesis data
func (db *Db) SaveGenesis(genesis *types.Genesis) error {
	stmt := `
INSERT INTO genesis(time, chain_id, initial_height) 
VALUES ($1, $2, $3) ON CONFLICT (one_row_id) DO UPDATE 
    SET time = excluded.time,
        initial_height = excluded.initial_height,
        chain_id = excluded.chain_id`

	_, err := db.Sqlx.Exec(stmt, genesis.Time, genesis.ChainID, genesis.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis: %s", err)
	}

	return nil
}

// GetGenesis returns the genesis information stored inside the database
func (db *Db) GetGenesis() (*types.Genesis, error) {
	var rows []*dbtypes.GenesisRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM genesis;`)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows inside the genesis table")
	}

	row := rows[0]
	return types.NewGenesis(row.ChainID, row.Time, row.InitialHeight), nil
}

func (db *Db) SetBlockSize(size int, height int64) error {
	stmt := `UPDATE block SET size = $1 WHERE height = $2`

	_, err := db.Sql.Exec(stmt, size, height)
	if err != nil {
		return fmt.Errorf("error while setting block size: %s", err)
	}

	return nil
}

func (db *Db) SetAverageBlockSize(block *tmctypes.ResultBlock) error {
	var avgBlockSize []dbtypes.AverageBlockSize
	stmtSelect := `SELECT * FROM average_block_size WHERE date = $1`
	err := db.Sqlx.Select(&avgBlockSize, stmtSelect, TimeToUTCDate(block.Block.Time))
	if err != nil {
		return err
	}

	stmtInsert := `
INSERT INTO average_block_size (date, blocks_number, block_sizes, average_size)
VALUES ($1, $2, $3, $4)`
	if len(avgBlockSize) == 0 {
		_, err := db.Sqlx.Exec(stmtInsert, TimeToUTCDate(block.Block.Time), 1, block.Block.Size(), block.Block.Size())
		if err != nil {
			return fmt.Errorf("error while setting average block size: %s", err)
		}

		return nil
	}

	stmtUpdate := `
UPDATE average_block_size SET blocks_number = $1,
           					  block_sizes = $2,
                              average_size = $3
WHERE date = $4`
	avgBlockSize[0].BlockSizes += int64(block.Block.Size())
	avgBlockSize[0].BlocksNumber++
	avgBlockSize[0].AverageSize = avgBlockSize[0].BlockSizes / avgBlockSize[0].BlocksNumber

	_, err = db.Sqlx.Exec(
		stmtUpdate,
		avgBlockSize[0].BlocksNumber, avgBlockSize[0].BlockSizes,
		avgBlockSize[0].AverageSize, TimeToUTCDate(block.Block.Time),
	)
	if err != nil {
		return fmt.Errorf("error while setting average block size: %s", err)
	}

	return nil
}

func (db *Db) SetAverageBlockTime(block *tmctypes.ResultBlock) error {
	var avgBlockTime []dbtypes.AverageBlockTime
	stmtSelect := `SELECT * FROM average_block_time WHERE date = $1`
	err := db.Sqlx.Select(&avgBlockTime, stmtSelect, TimeToUTCDate(block.Block.Time))
	if err != nil {
		return err
	}

	stmtInsert := `
INSERT INTO average_block_time (date, last_timestamp, blocks_number, block_times, average_time)
VALUES ($1, $2, $3, $4, $5)`
	if len(avgBlockTime) == 0 {
		_, err := db.Sqlx.Exec(stmtInsert, TimeToUTCDate(block.Block.Time), block.Block.Time.Unix(), 1, 0, 0)
		if err != nil {
			return fmt.Errorf("error while setting average block time: %s", err)
		}

		return nil
	}

	stmtUpdate := `
UPDATE average_block_time SET last_timestamp = $1,
                              blocks_number = $2,
           					  block_times = $3,
                              average_time = $4
WHERE date = $5`
	avgBlockTime[0].BlockTimes += block.Block.Time.Unix() - avgBlockTime[0].LastTimestamp
	avgBlockTime[0].LastTimestamp = block.Block.Time.Unix()
	avgBlockTime[0].BlocksNumber++
	avgBlockTime[0].AverageTime = avgBlockTime[0].BlockTimes / avgBlockTime[0].BlocksNumber

	_, err = db.Sqlx.Exec(
		stmtUpdate, avgBlockTime[0].LastTimestamp,
		avgBlockTime[0].BlocksNumber, avgBlockTime[0].BlockTimes,
		avgBlockTime[0].AverageTime, TimeToUTCDate(block.Block.Time),
	)
	if err != nil {
		return fmt.Errorf("error while setting average block time: %s", err)
	}

	return nil
}

func (db *Db) SetTxsPerDate(block *tmctypes.ResultBlock) error {
	stmt := `
INSERT INTO txs_per_date (date, txs_number)
VALUES ($1, $2) ON CONFLICT (date) DO UPDATE
	SET txs_number = txs_per_date.txs_number + $2`

	_, err := db.Sqlx.Exec(stmt, TimeToUTCDate(block.Block.Time), len(block.Block.Txs))
	if err != nil {
		return fmt.Errorf("error while setting average block time: %s", err)
	}

	return nil
}

func (db *Db) SetAverageFee(blockFee int64, block *tmctypes.ResultBlock) error {
	var avgFee []dbtypes.AverageFee
	stmtSelect := `SELECT * FROM average_block_fee WHERE date = $1`
	err := db.Sqlx.Select(&avgFee, stmtSelect, TimeToUTCDate(block.Block.Time))
	if err != nil {
		return err
	}

	stmtInsert := `
INSERT INTO average_block_fee (date, blocks_number, block_fees, average_fee)
VALUES ($1, $2, $3, $4)`
	if len(avgFee) == 0 {
		_, err := db.Sqlx.Exec(stmtInsert, TimeToUTCDate(block.Block.Time), 1, blockFee, blockFee)
		if err != nil {
			return fmt.Errorf("error while setting average block size: %s", err)
		}

		return nil
	}

	stmtUpdate := `
UPDATE average_block_fee SET blocks_number = $1,
           					 block_fees = $2,
                             average_fee = $3
WHERE date = $4`
	avgFee[0].BlockFees += blockFee
	avgFee[0].BlocksNumber++
	avgFee[0].AverageFee = avgFee[0].BlockFees / avgFee[0].BlocksNumber

	_, err = db.Sqlx.Exec(
		stmtUpdate,
		avgFee[0].BlocksNumber, avgFee[0].BlockFees,
		avgFee[0].AverageFee, TimeToUTCDate(block.Block.Time),
	)
	if err != nil {
		return fmt.Errorf("error while resetting average block size: %s", err)
	}

	return nil
}

func (db *Db) SetTxSender(tx *junotypes.Tx) error {
	stmt := `UPDATE transaction SET sender = $1 WHERE hash LIKE $2`

	hash := tx.TxHash[:len(tx.TxHash)-1] + "%"

	_, err := db.Sqlx.Exec(stmt, tx.GetSigners()[0].String(), hash)
	if err != nil {
		return fmt.Errorf("error while setting tx senders: %s", err)
	}

	return nil
}

func TimeToUTCDate(t time.Time) time.Time {
	year, month, day := t.UTC().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
