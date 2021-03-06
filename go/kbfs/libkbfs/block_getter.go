// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package libkbfs

import (
	"github.com/keybase/client/go/kbfs/data"
	"github.com/keybase/client/go/kbfs/kbfscrypto"
	"github.com/keybase/client/go/kbfs/libkey"
	"golang.org/x/net/context"
)

// blockGetter provides the API for the block retrieval worker to obtain blocks.
type blockGetter interface {
	getBlock(
		context.Context, libkey.KeyMetadata, data.BlockPointer, data.Block,
		DiskBlockCacheType) error
	assembleBlock(
		context.Context, libkey.KeyMetadata, data.BlockPointer, data.Block,
		[]byte, kbfscrypto.BlockCryptKeyServerHalf) error
	assembleBlockLocal(
		context.Context, libkey.KeyMetadata, data.BlockPointer, data.Block,
		[]byte, kbfscrypto.BlockCryptKeyServerHalf) error
}

// realBlockGetter obtains real blocks using the APIs available in Config.
type realBlockGetter struct {
	config blockOpsConfig
}

// getBlock implements the interface for realBlockGetter.
func (bg *realBlockGetter) getBlock(
	ctx context.Context, kmd libkey.KeyMetadata, blockPtr data.BlockPointer,
	block data.Block, cacheType DiskBlockCacheType) error {
	bserv := bg.config.BlockServer()
	buf, blockServerHalf, err := bserv.Get(
		ctx, kmd.TlfID(), blockPtr.ID, blockPtr.Context, cacheType)
	if err != nil {
		return err
	}

	if _, isLocal := bserv.(blockServerLocal); isLocal {
		return assembleBlockLocal(
			ctx, bg.config.keyGetter(), bg.config.Codec(),
			bg.config.cryptoPure(), kmd, blockPtr, block, buf, blockServerHalf)
	}
	return assembleBlock(
		ctx, bg.config.keyGetter(), bg.config.Codec(), bg.config.cryptoPure(),
		kmd, blockPtr, block, buf, blockServerHalf)
}

func (bg *realBlockGetter) assembleBlock(ctx context.Context,
	kmd libkey.KeyMetadata, ptr data.BlockPointer, block data.Block, buf []byte,
	serverHalf kbfscrypto.BlockCryptKeyServerHalf) error {
	return assembleBlock(ctx, bg.config.keyGetter(), bg.config.Codec(),
		bg.config.cryptoPure(), kmd, ptr, block, buf, serverHalf)
}

func (bg *realBlockGetter) assembleBlockLocal(ctx context.Context,
	kmd libkey.KeyMetadata, ptr data.BlockPointer, block data.Block, buf []byte,
	serverHalf kbfscrypto.BlockCryptKeyServerHalf) error {
	return assembleBlockLocal(ctx, bg.config.keyGetter(), bg.config.Codec(),
		bg.config.cryptoPure(), kmd, ptr, block, buf, serverHalf)
}
