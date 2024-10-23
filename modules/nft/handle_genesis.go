package nft

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/x/nft"
	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbtypes "github.com/forbole/callisto/v4/database/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "nft").Msg("parsing genesis")

	// Read the genesis state
	var nftGenState nft.GenesisState
	err := m.cdc.UnmarshalJSON(appState[nft.ModuleName], &nftGenState)
	if err != nil {
		return fmt.Errorf("error while unmarshalling nft state: %s", err)
	}

	// Read the genesis state
	var onftGenState onfttypes.GenesisState
	err = m.cdc.UnmarshalJSON(appState[onfttypes.ModuleName], &onftGenState)
	if err != nil {
		return fmt.Errorf("error while unmarshalling onft state: %s", err)
	}

	classes := make([]*onfttypes.Class, 0, len(nftGenState.Classes))

	for _, class := range nftGenState.Classes {
		classes = append(classes, &onfttypes.Class{
			Id:          class.Id,
			Name:        class.Name,
			Symbol:      class.Symbol,
			Description: class.Description,
			Uri:         class.Uri,
			UriHash:     class.UriHash,
			Data:        class.Data,
			Owner:       onftGenState.ClassOwners[class.Id],
		})
	}

	err = m.db.SaveNFTClasses(classes, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving nft classes: %s", err)
	}

	for _, entry := range nftGenState.Entries {
		nfts := make([]*dbtypes.NFT, 0, len(entry.Nfts))
		for _, n := range entry.Nfts {
			nfts = append(nfts, &dbtypes.NFT{
				ClassId:    n.ClassId,
				Id:         n.Id,
				Uri:        n.Uri,
				UriHash:    n.UriHash,
				Owner:      entry.Owner,
				Data:       n.Data,
				MintTxHash: "",
			})
		}

		err = m.db.SaveNFTs(nfts, doc.InitialHeight)
		if err != nil {
			return fmt.Errorf("error while saving nfts: %s", err)
		}
	}

	return nil
}
