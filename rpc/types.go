package rpc

import (
	commonProtoc "github.com/irisnet/blockchain-rpc/codegen/server/model"
	"github.com/irisnet/irishub-server/models/document"
)

func BuildResponseAddress(address string) commonProtoc.Address {
	return commonProtoc.Address{
		Chain: "",
		App:   "",
		Addr:  address,
	}
}

func BuildResponseCoins(coins document.Coins) []*commonProtoc.Coin {
	var (
		modelCoins []*commonProtoc.Coin
	)

	if len(coins) > 0 {
		for _, v := range coins {
			modelCoin := commonProtoc.Coin{
				Denom:  v.Denom,
				Amount: float64(v.Amount),
			}
			modelCoins = append(modelCoins, &modelCoin)
		}
	}

	return modelCoins
}
