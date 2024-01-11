// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE
package encoding

import (
	"cosmossdk.io/x/tx/signing"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"

	enccodec "github.com/evmos/ethermint/encoding/codec"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// MakeConfig creates an EncodingConfig for testing
// and registers the interfaces
func MakeConfig(mb module.BasicManager) sdktestutil.TestEncodingConfig {
	ec := encodingConfig()
	enccodec.RegisterLegacyAminoCodec(ec.Amino)
	mb.RegisterLegacyAminoCodec(ec.Amino)
	enccodec.RegisterInterfaces(ec.InterfaceRegistry)
	mb.RegisterInterfaces(ec.InterfaceRegistry)
	return ec
}

// encodingConfig creates a new EncodingConfig and returns it
func encodingConfig() sdktestutil.TestEncodingConfig {
	cdc := amino.NewLegacyAmino()
	signingOptions := signing.Options{
		AddressCodec: address.Bech32Codec{
			Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
		},
		ValidatorAddressCodec: address.Bech32Codec{
			Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
		},
	}
	signingOptions.DefineCustomGetSigners(
		evmtypes.MsgEthereumTxCustomGetSigner.MsgType,
		evmtypes.MsgEthereumTxCustomGetSigner.Fn,
	)

	interfaceRegistry, _ := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles:     proto.HybridResolver,
		SigningOptions: signingOptions,
	})
	codec := amino.NewProtoCodec(interfaceRegistry)

	return sdktestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}
}
