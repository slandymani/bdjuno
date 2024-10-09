package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	shell "github.com/ipfs/go-ipfs-api"

	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
)

func (db *Db) SaveNFTClass(class *onfttypes.Class, height int64) error {
	return db.SaveNFTClasses([]*onfttypes.Class{class}, height)
}

func (db *Db) SaveNFTClasses(classes []*onfttypes.Class, height int64) error {
	stmt := `INSERT INTO nft_class(id, name, symbol, description, uri, uri_hash, data, owner, metadata, height) VALUES`

	var params []interface{}
	for i, class := range classes {
		metadata, _ := fetchMetadataFromIPFS(class.Uri)

		vi := i * 10
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10)
		params = append(params, class.Id,
			class.Name,
			class.Symbol,
			class.Description,
			class.Uri,
			class.UriHash,
			class.Data,
			class.Owner,
			metadata,
			height,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
ON CONFLICT (id) DO UPDATE 
    SET name = excluded.name,
    	symbol = excluded.symbol,
    	description = excluded.description,
        uri = excluded.uri,
    	uri_hash = excluded.uri_hash,
    	data = excluded.data,
    	owner = excluded.owner,
    	metadata = excluded.metadata,
    	height = excluded.height
WHERE nft_class.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing classes: %s", err)
	}
	return nil
}

func (db *Db) ChangeNFTClassOwner(classID, newOwner string, height int64) error {
	stmt := `UPDATE nft_class SET owner=$1, height=$2 WHERE id=$3 and height<=$2`

	_, err := db.SQL.Exec(stmt, newOwner, height, classID)
	if err != nil {
		return fmt.Errorf("error while updating nft owner: %s", err)
	}
	return nil
}

func (db *Db) SaveNFT(nft *onfttypes.NFT, height int64) error {
	return db.SaveNFTs([]*onfttypes.NFT{nft}, height)
}

func (db *Db) SaveNFTs(nfts []*onfttypes.NFT, height int64) error {
	stmt := `INSERT INTO nft(id, class_id, uri, uri_hash, data, owner, metadata, height) VALUES`

	var params []interface{}
	for i, nft := range nfts {
		metadata, _ := fetchMetadataFromIPFS(nft.Uri)

		vi := i * 8
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8)
		params = append(params, nft.Id,
			nft.ClassId,
			nft.Uri,
			nft.UriHash,
			nft.Data,
			nft.Owner,
			metadata,
			height,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
ON CONFLICT (id, class_id) DO UPDATE 
    SET uri = excluded.uri,
    	uri_hash = excluded.uri_hash,
    	data = excluded.data,
    	owner = excluded.owner,
    	metadata = excluded.metadata,
    	height = excluded.height
WHERE nft.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing nfts: %s", err)
	}

	return nil
}

func (db *Db) ChangeNFTOwner(classID, nftID, newOwner string, height int64) error {
	stmt := `UPDATE nft SET owner=$1, height=$2 WHERE id=$3 and class_id=$4 and height<=$2`

	_, err := db.SQL.Exec(stmt, newOwner, height, nftID, classID)
	if err != nil {
		return fmt.Errorf("error while updating nft owner: %s", err)
	}

	return nil
}

func fetchMetadataFromIPFS(cid string) (string, error) {
	if cid == "" {
		return "", nil
	}

	parts := strings.Split(cid, "/")
	cid = parts[len(parts)-1]

	// Connect to the local IPFS node
	sh := shell.NewShell("localhost:5001") // Update with your IPFS node address if different

	// Fetch the file from IPFS
	reader, err := sh.Cat(cid)
	if err != nil {
		return "", fmt.Errorf("error while fetching file from IPFS: %w", err)
	}
	defer reader.Close()

	// Read the content into a buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", fmt.Errorf("error while reading file from IPFS: %w", err)
	}

	// Detect MIME type
	mime := mimetype.Detect(buf.Bytes())

	// Check if the MIME type indicates an image
	if strings.HasPrefix(mime.String(), "image/") {
		// If the file is an image, return an empty string
		return "", nil
	}

	// If it's not an image, try to decode the content as JSON
	var jsonData interface{}
	err = json.Unmarshal(buf.Bytes(), &jsonData)
	if err != nil {
		return "", fmt.Errorf("error while parsing JSON: %w", err)
	}

	// Convert JSON back to a formatted string (optional)
	jsonString, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error while formatting JSON: %w", err)
	}

	return string(jsonString), nil
}
