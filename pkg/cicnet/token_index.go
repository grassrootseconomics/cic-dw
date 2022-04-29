package cicnet

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"math/big"
)

func (c *CicNet) EntryCount(ctx context.Context) (big.Int, error) {
	var tokenCount big.Int

	err := c.ethClient.CallCtx(
		ctx,
		eth.CallFunc(w3.MustNewFunc("entryCount()", "uint256"), c.tokenIndex).Returns(&tokenCount),
	)
	if err != nil {
		return big.Int{}, err
	}

	return tokenCount, nil
}

func (c *CicNet) AddressAtIndex(ctx context.Context, index *big.Int) (common.Address, error) {
	var address common.Address

	err := c.ethClient.CallCtx(
		ctx,
		eth.CallFunc(w3.MustNewFunc("entry(uint256 _idx)", "address"), c.tokenIndex, index).Returns(&address),
	)
	if err != nil {
		return [20]byte{}, err
	}

	return address, nil
}
