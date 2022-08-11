package oracle

import oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"

func (m *Module) RefreshDataSourceInfo(height int64, dataSource oracletypes.DataSource) error {
	return nil
}

//data source
//--------------------------------------------------------------
//id           INT NOT NULL PRIMARY KEY,						+
//create_block BIGINT NOT NULL,									?
//edit_block   BIGINT,											?
//name         TEXT NOT NULL,									+
//description  TEXT,											+
//executable   TEXT,											?
//fee          COIN[],											+
//owner        TEXT NOT NULL REFERENCES account (address),		+
//sender       TEXT,											?
//timestamp    TIMESTAMP WITHOUT TIME ZONE						?
