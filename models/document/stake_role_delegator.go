package document

import (
	"time"

	"github.com/irisnet/irishub-server/modules/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/irisnet/irishub-server/models"
)

const (
	CollectionNmStakeRoleDelegator = "stake_role_delegator"
)

type Delegator struct {
	Address             string  `json:"address" bson:"address"`
	ValidatorAddr       string  `json:"pub_key" bson:"validator_addr"` // validator ValidatorAddress
	Shares              float64 `json:"shares" bson:"shares"`
	BondedTokens        float64
	UnbondingDelegation UnbondingDelegation `bson:"unbonding_delegation"`
}

type UnbondingDelegation struct {
	Balance Coins     `bson:"balance"`
	MinTime time.Time `bson:"min_time"`
}

type DelegatorShares struct {
	ValidatorAddr        string  `bson:"_id"`
	TotalShares          float64 `json:"total_shares" bson:"total_shares"`
	TotalUnbondingTokens float64 `bson:"unbonding_tokens"`
}

func (d Delegator) Name() string {
	return CollectionNmStakeRoleDelegator
}

func (d Delegator) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address, "validator_addr": d.ValidatorAddr}
}

func (d Delegator) Query(
	query bson.M, skip int, limit int, sorts ...string,
) (results []Delegator, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sorts...).
			Skip(skip).Limit(limit).All(&results)
	}
	return results, models.ExecCollection(d.Name(), exop)
}

func (d Delegator) GetDelegatorListByAddressAndValidatorAddrs(address string, valAddrs []string,
) ([]Delegator, error) {

	query := bson.M{
		"address": address,
		"validator_addr": &bson.M{
			"$in": valAddrs,
		},
	}
	sorts := make([]string, 0)

	delegator, err := d.Query(query, 0, len(valAddrs), sorts...)

	if err != nil {
		logger.Error.Println(err)
	}

	return delegator, err
}

func (d Delegator) GetDelegatorListByAddress(address string, skip int,
	limit int, sorts []string) ([]Delegator, error) {

	query := bson.M{
		"address": address,
		"shares": &bson.M{
			"$gt": 0,
		},
	}

	delegator, err := d.Query(query, skip, limit, sorts...)

	if err != nil {
		logger.Error.Println(err)
	}

	return delegator, err
}

func (d Delegator) GetTotalSharesByAddress(address string) ([]DelegatorShares, error) {
	var value []DelegatorShares

	q := func(c *mgo.Collection) error {
		m := []bson.M{
			{
				"$match": bson.M{
					"address": address,
				},
			},
			{"$unwind": "$unbonding_delegation.balance"},
			{
				"$group": bson.M{
					"_id":              "$validator_addr",
					"total_shares":     bson.M{"$sum": "$shares"},
					"unbonding_tokens": bson.M{"$sum": "$unbonding_delegation.balance.amount"},
				},
			},
		}
		return c.Pipe(m).All(&value)
	}

	err := models.ExecCollection(d.Name(), q)

	if err != nil {
		return nil, err
	}
	return value, nil

}
