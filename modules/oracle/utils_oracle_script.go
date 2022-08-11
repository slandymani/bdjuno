package oracle

//TODO: 2 READERS; 2 PARSERS; 2 COBRA COMMANDS, FIND VALUES, SHOW DUPLICATED TABLES

//oracle script
//-----------------------------------------------------------
//id              INT NOT NULL PRIMARY KEY,						+
//create_block    BIGINT NOT NULL,								?
//edit_block      BIGINT,										?
//name            TEXT NOT NULL,								+
//description     TEXT,											+
//schema          TEXT,											+
//source_code_url TEXT,											+
//owner           TEXT NOT NULL REFERENCES account (address),	+
//sender          TEXT,											?
//timestamp       TIMESTAMP WITHOUT TIME ZONE					?
//-----------------------------------------------------------

//ID            OracleScriptID `protobuf:"varint,1,opt,name=id,proto3,casttype=OracleScriptID" json:"id,omitempty"`
//Owner         string         `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
//Name          string         `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
//Description   string         `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
//Filename      string         `protobuf:"bytes,5,opt,name=filename,proto3" json:"filename,omitempty"`
//Schema        string         `protobuf:"bytes,6,opt,name=schema,proto3" json:"schema,omitempty"`
//SourceCodeURL string
