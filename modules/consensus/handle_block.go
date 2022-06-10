package consensus

import (
	"fmt"
	app "github.com/ODIN-PROTOCOL/odin-core/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/types"

	"github.com/rs/zerolog/log"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, txs []*types.Tx, _ *tmctypes.ResultValidators,
) error {
	err := m.updateBlockTimeFromGenesis(b)
	if err != nil {
		log.Error().Str("module", "consensus").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating block time from genesis")
	}

	err = m.db.SetAverageBlockSize(b)
	if err != nil {
		log.Error().Str("module", "consensus").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating average block size")
	}

	err = m.db.SetAverageBlockTime(b)
	if err != nil {
		log.Error().Str("module", "consensus").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating average block time")
	}

	if len(txs) > 0 {
		err = m.db.SetTxsPerDate(b)
		if err != nil {
			log.Error().Str("module", "consensus").Int64("height", b.Block.Height).
				Err(err).Msg("error while updating txs per day")
		}

		var blockFee int64
		for _, tx := range txs {
			blockFee += tx.GetFee().AmountOf("loki").Int64()
		}
		err = m.db.SetAverageFee(blockFee, b)
		if err != nil {
			log.Error().Str("module", "consensus").Int64("height", b.Block.Height).
				Err(err).Msg("error while updating block fee")
		}
	}

	return nil
}

// updateBlockTimeFromGenesis insert average block time from genesis
func (m *Module) updateBlockTimeFromGenesis(block *tmctypes.ResultBlock) error {
	log.Trace().Str("module", "consensus").Int64("height", block.Block.Height).
		Msg("updating block time from genesis")

	genesis, err := m.db.GetGenesis()
	if err != nil {
		return fmt.Errorf("error while getting genesis: %s", err)
	}
	if genesis == nil {
		return fmt.Errorf("genesis table is empty")
	}

	// Skip if the genesis does not exist
	if genesis == nil {
		return nil
	}

	newBlockTime := block.Block.Time.Sub(genesis.Time).Seconds() / float64(block.Block.Height-genesis.InitialHeight)
	return m.db.SaveAverageBlockTimeGenesis(newBlockTime, block.Block.Height)
}

func DecodeTX(txBytes []byte) (sdk.Tx, error) {
	encCfg := app.MakeEncodingConfig()
	decodedTx, err := encCfg.TxConfig.TxDecoder()(txBytes)
	return decodedTx, err
}
