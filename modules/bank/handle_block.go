package bank

import (
	app "github.com/ODIN-PROTOCOL/odin-core/app"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/bdjuno/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

func (m *Module) HandleBlock(
	_ *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, txs []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	var addresses []string
	addrMap := make(map[string]bool)

	for _, tx := range txs {
		for _, message := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err := m.cdc.UnpackAny(message, &stdMsg)
			if err != nil {
				return err
			}

			addresses, _ := m.messageParser(m.cdc, stdMsg)
			for _, address := range addresses {
				addrMap[address] = true
			}
		}
	}

	for addr, _ := range addrMap {
		if len(addr) < 4 || addr[:4] != app.Bech32MainPrefix {
			continue
		}
		addresses = append(addresses, addr)
	}

	balances, err := m.keeper.GetBalances(addresses, height)
	if err != nil {
		return err
	}

	bankBalances := make([]banktypes.Balance, len(balances))
	for i, balance := range balances {
		bankBalances[i] = convertBalance(balance)
	}

	return m.db.SaveAccountBalances(height, bankBalances)
}

func convertBalance(balance types.AccountBalance) banktypes.Balance {
	return banktypes.Balance{
		Address: balance.Address,
		Coins:   balance.Balance,
	}
}
