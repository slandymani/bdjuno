package utils

import (
	"sync"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var once sync.Once
var cdc *codec.ProtoCodec

func GetCodec() codec.Codec {
	once.Do(func() {
		marshaler, _ := odinapp.MakeCodecs()
		cdc = codec.NewProtoCodec(marshaler.InterfaceRegistry())
	})
	return cdc
}

// UnpackMessage unpacks a message from a byte slice
func UnpackMessage[T proto.Message](cdc codec.Codec, bz []byte, _ T) T {
	var anyT codectypes.Any
	cdc.MustUnmarshalJSON(bz, &anyT)
	var cosmosMsg sdk.Msg
	if err := cdc.UnpackAny(&anyT, &cosmosMsg); err != nil {
		panic(err)
	}
	return cosmosMsg.(T)
}
