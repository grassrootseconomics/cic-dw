package cicnet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
)

type CicNet struct {
	ethClient  *w3.Client
	tokenIndex common.Address
}

func NewCicNet(rpcEndpoint string, tokenIndex common.Address) *CicNet {
	ethClient := w3.MustDial(rpcEndpoint)

	return &CicNet{
		ethClient:  ethClient,
		tokenIndex: tokenIndex,
	}
}

func (c *CicNet) Close() error {
	err := c.ethClient.Close()
	if err != nil {
		return err
	}

	return nil
}
