package types

type StakingDelegatorRow struct {
	Address string `db:"address"`
	Stake   string `db:"stake"`
}

func NewStakingDelegatorRow(address, stake string) StakingDelegatorRow {
	return StakingDelegatorRow{
		Address: address,
		Stake:   stake,
	}
}
