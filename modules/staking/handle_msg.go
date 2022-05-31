package staking

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return m.handleMsgCreateValidator(tx.Height, cosmosMsg)

	case *stakingtypes.MsgEditValidator:
		return m.handleEditValidator(tx.Height, cosmosMsg)

	case *stakingtypes.MsgDelegate:
		return m.handleDelegate(tx.Height, cosmosMsg)

	case *stakingtypes.MsgUndelegate:
		return m.handleUndelegate(tx.Height, cosmosMsg)

	case *stakingtypes.MsgBeginRedelegate:
		return m.handleBeginRedelegate(tx.Height, cosmosMsg)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func (m *Module) handleMsgCreateValidator(height int64, msg *stakingtypes.MsgCreateValidator) error {
	err := m.RefreshValidatorInfos(height, msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing validator from MsgCreateValidator: %s", err)
	}
	return nil
}

// handleEditValidator handles MsgEditValidator utils, updating the validator info
func (m *Module) handleEditValidator(height int64, msg *stakingtypes.MsgEditValidator) error {
	err := m.RefreshValidatorInfos(height, msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing validator from MsgEditValidator: %s", err)
	}

	return nil
}

func (m *Module) handleDelegate(height int64, msg *stakingtypes.MsgDelegate) error {
	err := m.RefreshDelegatorDelegations(height, msg.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing delegator from MsgDelegate: %s", err)
	}

	return nil
}

func (m *Module) handleUndelegate(height int64, msg *stakingtypes.MsgUndelegate) error {
	err := m.RefreshDelegatorDelegations(height, msg.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing delegator from MsgDelegate: %s", err)
	}

	return nil
}

func (m *Module) handleBeginRedelegate(height int64, msg *stakingtypes.MsgBeginRedelegate) error {
	err := m.RefreshDelegatorDelegations(height, msg.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing delegator from MsgDelegate: %s", err)
	}

	return nil
}
