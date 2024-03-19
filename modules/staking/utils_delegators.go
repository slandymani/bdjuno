package staking

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

func (m *Module) RefreshDelegatorDelegations(height int64, address string) error {
	delegations, err := m.source.GetDelegationsTotal(height, address)
	if err != nil {
		return err
	}

	pool, err := m.source.GetPool(height)
	if err != nil {
		return err
	}

	delegationsPercent := float32(delegations.AmountOf("loki").Int64()) * 100 / float32(pool.BondedTokens.Int64())

	return m.db.SaveStakingDelegator(types.Delegator{
		Address:               address,
		Delegations:           delegations.AmountOf("loki").Int64(),
		DelegationsPercentage: fmt.Sprintf("%f", delegationsPercent),
		Height:                height,
	})
}
