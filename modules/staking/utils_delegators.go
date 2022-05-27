package staking

import "github.com/forbole/bdjuno/v3/types"

func (m *Module) RefreshDelegatorDelegations(height int64, address string) error {
	delegations, err := m.source.GetDelegationsTotal(height, address)
	if err != nil {
		return err
	}

	return m.db.SaveStakingDelegator(types.Delegator{
		Address:     address,
		Delegations: delegations,
		Height:      height,
	})
}
