package nft

import (
	"fmt"

	"cosmossdk.io/x/nft"
	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	wasmtypes "github.com/ODIN-PROTOCOL/wasmd/x/wasm/types"
	"github.com/forbole/callisto/v4/utils"
	juno "github.com/forbole/juno/v6/types"
)

var msgFilter = map[string]bool{
	"/onft.v1.MsgCreateNFTClass":           true,
	"/onft.v1.MsgTransferClassOwnership":   true,
	"/onft.v1.MsgMintNFT":                  true,
	"/nft.v1beta1.MsgSend":                 true,
	"/cosmwasm.wasm.v1.MsgExecuteContract": true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

func (m *Module) HandleMsg(index int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	switch msg.GetType() {
	case "/onft.v1.MsgCreateNFTClass":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &onfttypes.MsgCreateNFTClass{})
		return m.handleMsgCreateNFTClass(index, tx, cosmosMsg)
	case "/onft.v1.MsgTransferClassOwnership":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &onfttypes.MsgTransferClassOwnership{})
		return m.handleMsgTransferClassOwnership(tx, cosmosMsg)
	case "/onft.v1.MsgMintNFT":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &onfttypes.MsgMintNFT{})
		return m.handleMsgMintNFT(index, tx, cosmosMsg)
	case "/nft.v1beta1.MsgSend":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &nft.MsgSend{})
		return m.handleMsgSend(tx, cosmosMsg)
	case "/cosmwasm.wasm.v1.MsgExecuteContract":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &wasmtypes.MsgExecuteContract{})
		return m.handleMsgExecuteContract(index, tx, cosmosMsg)
	}

	return nil
}

func (m *Module) handleMsgCreateNFTClass(index int, tx *juno.Transaction, msg *onfttypes.MsgCreateNFTClass) error {
	event, err := tx.FindEventByType(index, onfttypes.EventTypeCreateNFTClass)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeCreateNFTClass: %s", err)
	}

	classID, err := tx.FindAttributeByKey(event, onfttypes.AttributeKeyClassID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyClassID: %s", err)
	}

	return m.db.SaveNFTClass(&onfttypes.Class{
		Id:          classID,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        msg.Data,
		Owner:       msg.Sender,
	}, int64(tx.Height))
}

func (m *Module) handleMsgTransferClassOwnership(tx *juno.Transaction, msg *onfttypes.MsgTransferClassOwnership) error {
	return m.db.ChangeNFTClassOwner(msg.ClassId, msg.NewOwner, int64(tx.Height))
}

func (m *Module) handleMsgMintNFT(index int, tx *juno.Transaction, msg *onfttypes.MsgMintNFT) error {
	event, err := tx.FindEventByType(index, onfttypes.EventTypeMintNFT)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeMintNFT: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, onfttypes.AttributeKeyNFTID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyNFTID: %s", err)
	}

	return m.db.SaveNFT(&onfttypes.NFT{
		Id:      id,
		ClassId: msg.ClassId,
		Uri:     msg.Uri,
		UriHash: msg.UriHash,
		Data:    msg.Data,
		Owner:   msg.Sender,
	}, int64(tx.Height))
}

func (m *Module) handleMsgSend(tx *juno.Transaction, msg *nft.MsgSend) error {
	return m.db.ChangeNFTOwner(msg.ClassId, msg.Id, msg.Receiver, int64(tx.Height))
}

func (m *Module) handleMsgExecuteContract(index int, tx *juno.Transaction, msg *wasmtypes.MsgExecuteContract) error {
	event, err := tx.FindEventByType(index, "wasm")
	if err != nil {
		return fmt.Errorf("error while searching for wasm: %s", err)
	}

	action, err := tx.FindAttributeByKey(event, "action")
	if err != nil {
		return nil //fmt.Errorf("error while searching for action: %s", err)
	}
	if action != "mint_nft_success" {
		return nil
	}

	classID, err := tx.FindAttributeByKey(event, "class_id")
	if err != nil {
		return fmt.Errorf("error while searching for class_id: %s", err)
	}

	nftID, err := tx.FindAttributeByKey(event, "nft_id")
	if err != nil {
		return fmt.Errorf("error while searching for class_id: %s", err)
	}

	n, err := m.source.NFT(int64(tx.Height), nftID, classID)
	if err != nil {
		return fmt.Errorf("failed to fetch nft: %s", err)
	}

	return m.db.SaveNFT(&onfttypes.NFT{
		Id:      n.Id,
		ClassId: n.ClassId,
		Uri:     n.Uri,
		UriHash: n.UriHash,
		Data:    n.Data,
		Owner:   n.Owner,
	}, int64(tx.Height))
}
