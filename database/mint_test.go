package database_test

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveInflation() {

	// Save the data
	err := suite.database.SaveInflation(sdk.NewDecWithPrec(10050, 2), 100)
	suite.Require().NoError(err)

	// Verify the data
	var rows []dbtypes.InflationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected := dbtypes.NewInflationRow(100.50, 100)
	suite.Require().True(expected.Equal(rows[0]))

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with lower height
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(20000, 2), 90)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	// Verify the data
	rows = []dbtypes.InflationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected = dbtypes.NewInflationRow(100.50, 100)
	suite.Require().True(expected.Equal(rows[0]), "data should not change with lower height")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(30000, 2), 100)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	// Verify the data
	rows = []dbtypes.InflationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected = dbtypes.NewInflationRow(300.00, 100)
	suite.Require().True(expected.Equal(rows[0]), "data should change with same height")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(40000, 2), 110)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	// Verify the data
	rows = []dbtypes.InflationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected = dbtypes.NewInflationRow(400.00, 110)
	suite.Require().True(expected.Equal(rows[0]), "data should change with higher height")
}

func (suite *DbTestSuite) TestBigDipperDb_SaveMintParams() {
	privateKey, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	privateKeyBytes := crypto.FromECDSA(privateKey)

	mintParams := minttypes.NewParams(
		"loki",
		sdk.NewDecWithPrec(4, 1),
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(4, 1),
		sdk.NewDecWithPrec(8, 1),
		sdk.Coins{},
		uint64(60*60*8766/5),
		true,
		map[string]string{"eth": hexutil.Encode(privateKeyBytes)},
		[]string{"odin1pl07tk6hcpp2an3rug75as4dfgd743qp80g63g"},
		sdk.NewCoins(),
		[]string{"minigeo"},
		[]string{"odin1pl07tk6hcpp2an3rug75as4dfgd743qp80g63g"},
	)
	err = suite.database.SaveMintParams(types.NewMintParams(mintParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.MintParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams minttypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(mintParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
