package wasm

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	wasmtypes "github.com/ODIN-PROTOCOL/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/callisto/v4/types"
	"github.com/forbole/callisto/v4/utils"
	eventutils "github.com/forbole/callisto/v4/utils/events"
	juno "github.com/forbole/juno/v6/types"
)

var msgFilter = map[string]bool{
	"/cosmwasm.wasm.v1.MsgStoreCode":            true,
	"/cosmwasm.wasm.v1.MsgInstantiateContract":  true,
	"/cosmwasm.wasm.v1.MsgInstantiateContract2": true,
	"/cosmwasm.wasm.v1.MsgExecuteContract":      true,
	"/cosmwasm.wasm.v1.MsgMigrateContract":      true,
	"/cosmwasm.wasm.v1.MsgUpdateAdmin":          true,
	"/cosmwasm.wasm.v1.MsgClearAdmin":           true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	switch msg.GetType() {
	case "/cosmwasm.wasm.v1.MsgStoreCode":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgStoreCode{})
		err := m.HandleMsgStoreCode(index, tx, cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling MsgStoreCode: %s", err)
		}
	case "/cosmwasm.wasm.v1.MsgInstantiateContract":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgInstantiateContract{})
		err := m.HandleMsgInstantiateContract(index, tx, cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling MsgInstantiateContract: %s", err)
		}
	case "/cosmwasm.wasm.v1.MsgInstantiateContract2":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgInstantiateContract2{})
		err := m.HandleMsgInstantiateContract(index, tx, &wasmtypes.MsgInstantiateContract{
			Sender: cosmosMsg.Sender,
			Admin:  cosmosMsg.Admin,
			CodeID: cosmosMsg.CodeID,
			Label:  cosmosMsg.Label,
			Msg:    cosmosMsg.Msg,
			Funds:  cosmosMsg.Funds,
		})
		if err != nil {
			return fmt.Errorf("error while handling MsgInstantiateContract: %s", err)
		}
	case "/cosmwasm.wasm.v1.MsgExecuteContract":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgExecuteContract{})
		err := m.HandleMsgExecuteContract(index, tx, cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling MsgExecuteContract: %s", err)
		}
	case "/cosmwasm.wasm.v1.MsgMigrateContract":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgMigrateContract{})
		err := m.HandleMsgMigrateContract(index, tx, cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling MsgMigrateContract: %s", err)
		}
	case "/cosmwasm.wasm.v1.MsgUpdateAdmin":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgUpdateAdmin{})
		err := m.HandleMsgUpdateAdmin(cosmosMsg, tx)
		if err != nil {
			return fmt.Errorf("error while handling MsgUpdateAdmin: %s", err)
		}
	case "/cosmwasm.wasm.v1.MsgClearAdmin":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgClearAdmin{})
		err := m.HandleMsgClearAdmin(cosmosMsg, tx)
		if err != nil {
			return fmt.Errorf("error while handling MsgClearAdmin: %s", err)
		}
	}

	return nil
}

// HandleMsgStoreCode allows to properly handle a MsgStoreCode
// The Store Code Event is to upload the contract code on the chain, where a Code ID is returned
func (m *Module) HandleMsgStoreCode(index int, tx *juno.Transaction, msg *wasmtypes.MsgStoreCode) error {
	events := eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index)

	event, has := eventutils.FindEventByType(events, wasmtypes.EventTypeStoreCode)
	if !has {
		return errors.New("error while searching for EventTypeStoreCode")
	}

	codeIDKey, has := eventutils.FindAttributeByKey(event, wasmtypes.AttributeKeyCodeID)
	if !has {
		return errors.New("error while searching for AttributeKeyCodeID")
	}

	codeID, err := strconv.ParseUint(codeIDKey.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing code id to int64: %s", err)
	}

	return m.db.SaveWasmCode(
		types.NewWasmCode(
			msg.Sender, msg.WASMByteCode, msg.InstantiatePermission, codeID, int64(tx.Height),
		),
	)
}

// HandleMsgInstantiateContract allows to properly handle a MsgInstantiateContract
// Instantiate Contract Event instantiates an executable contract with the code previously stored with Store Code Event
func (m *Module) HandleMsgInstantiateContract(index int, tx *juno.Transaction, msg *wasmtypes.MsgInstantiateContract) error {
	events := eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index)

	event, has := eventutils.FindEventByType(events, wasmtypes.EventTypeInstantiate)
	if !has {
		return errors.New("error while searching for EventTypeInstantiate")
	}

	contractAddress, has := eventutils.FindAttributeByKey(event, wasmtypes.AttributeKeyContractAddr)
	if !has {
		return errors.New("error while searching for AttributeKeyContractAddr")
	}

	err := m.db.UpdateMsgInvolvedAccountsAddresses(contractAddress.Value, tx.TxHash)
	if err != nil {
		return fmt.Errorf("error while saving contract address inside involved accounts addresses: %s", err)
	}

	// Get result data
	resultData, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if err != nil {
		resultData = ""
	}
	resultDataBz, err := base64.StdEncoding.DecodeString(resultData)
	if err != nil {
		return fmt.Errorf("error while decoding result data: %s", err)
	}

	// Get the contract info
	contractInfo, err := m.source.GetContractInfo(int64(tx.Height), contractAddress.Value)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	// Get contract info extension
	var contractInfoExt string
	if contractInfo.Extension != nil {
		var extentionI wasmtypes.ContractInfoExtension
		err = m.cdc.UnpackAny(contractInfo.Extension, &extentionI)
		if err != nil {
			return fmt.Errorf("error while getting contract info extension: %s", err)
		}
		contractInfoExt = extentionI.String()
	}

	// Get contract states
	contractStates, err := m.source.GetContractStates(int64(tx.Height), contractAddress.Value)
	if err != nil {
		return fmt.Errorf("error while getting genesis contract states: %s", err)
	}

	contract := types.NewWasmContract(
		msg.Sender, msg.Admin, msg.CodeID, msg.Label, msg.Msg, msg.Funds,
		contractAddress.Value, string(resultDataBz), timestamp,
		contractInfo.Creator, contractInfoExt, contractStates, int64(tx.Height),
	)
	return m.db.SaveWasmContracts(
		[]types.WasmContract{contract},
	)
}

// HandleMsgExecuteContract allows to properly handle a MsgExecuteContract
// Execute Event executes an instantiated contract
func (m *Module) HandleMsgExecuteContract(index int, tx *juno.Transaction, msg *wasmtypes.MsgExecuteContract) error {
	events := eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index)

	event, has := eventutils.FindEventByType(events, wasmtypes.EventTypeExecute)
	if !has {
		return errors.New("error while searching for EventTypeExecute")
	}

	resultData, has := eventutils.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if !has {
		resultData.Value = ""
	}
	resultDataBz, err := base64.StdEncoding.DecodeString(resultData.Value)
	if err != nil {
		return fmt.Errorf("error while decoding result data: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	err = m.db.UpdateMsgInvolvedAccountsAddresses(msg.Contract, tx.TxHash)
	if err != nil {
		return fmt.Errorf("error while saving contract address inside the involved addresses in message table: %s", err)
	}

	return m.db.SaveWasmExecuteContract(
		types.NewWasmExecuteContract(
			msg.Sender, msg.Contract, msg.Msg, msg.Funds,
			string(resultDataBz), timestamp, int64(tx.Height),
		),
	)
}

// HandleMsgMigrateContract allows to properly handle a MsgMigrateContract
// Migrate Contract Event upgrade the contract by updating code ID generated from new Store Code Event
func (m *Module) HandleMsgMigrateContract(index int, tx *juno.Transaction, msg *wasmtypes.MsgMigrateContract) error {
	events := eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index)

	event, has := eventutils.FindEventByType(events, wasmtypes.EventTypeMigrate)
	if !has {
		return errors.New("error while searching for EventTypeMigrate")
	}

	resultData, has := eventutils.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if !has {
		resultData.Value = ""
	}
	resultDataBz, err := base64.StdEncoding.DecodeString(resultData.Value)
	if err != nil {
		return fmt.Errorf("error while decoding result data: %s", err)
	}

	err = m.db.UpdateMsgInvolvedAccountsAddresses(msg.Contract, tx.TxHash)
	if err != nil {
		return fmt.Errorf("error while saving contract address inside the involved addresses in message table: %s", err)
	}

	return m.db.UpdateContractWithMsgMigrateContract(msg.Sender, msg.Contract, msg.CodeID, msg.Msg, string(resultDataBz))
}

// HandleMsgUpdateAdmin allows to properly handle a MsgUpdateAdmin
// Update Admin Event updates the contract admin who can migrate the wasm contract
func (m *Module) HandleMsgUpdateAdmin(msg *wasmtypes.MsgUpdateAdmin, tx *juno.Transaction) error {
	err := m.db.UpdateMsgInvolvedAccountsAddresses(msg.Contract, tx.TxHash)
	if err != nil {
		return fmt.Errorf("error while saving contract address inside the involved addresses in message table: %s", err)
	}

	return m.db.UpdateContractAdmin(msg.Sender, msg.Contract, msg.NewAdmin)
}

// HandleMsgClearAdmin allows to properly handle a MsgClearAdmin
// Clear Admin Event clears the admin which make the contract no longer migratable
func (m *Module) HandleMsgClearAdmin(msg *wasmtypes.MsgClearAdmin, tx *juno.Transaction) error {
	err := m.db.UpdateMsgInvolvedAccountsAddresses(msg.Contract, tx.TxHash)
	if err != nil {
		return fmt.Errorf("error while saving contract address inside the involved addresses in message table: %s", err)
	}

	return m.db.UpdateContractAdmin(msg.Sender, msg.Contract, "")
}
