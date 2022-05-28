package model

import (
	"errors"

	"github.com/hb-go/gorm"
	"github.com/siprtcio/balanceservice/model/orm"
)

type BalanceService struct {
	orm.Model `gorm:"-" json:"-" swaggerignore:"true"`
	AuthId    string  `gorm:"column:auth_i_d" json:"auth_id" validate:"required" binding:"required"`      // tenant authid.
	Balance   float64 `gorm:"column:balance_money" json:"balance" validate:"required" binding:"required"` // balance in USD($)
}

func GetBalanceByParentAuthIdAuthId(parentAuthId, authId string) (*BalanceService, error) {
	var balService BalanceService
	var err error
	if err = DB().Where(map[string]interface{}{"auth_i_d": authId, "parent_auth_id": parentAuthId}).
		First(&balService).Error; err != nil {
		//there might be database down we need to see on this one
		return nil, err
	}
	return &balService, err
}

func GetBalanceByAuthId(authId string) (*BalanceService, error) {
	var balService BalanceService
	var err error
	if err = DB().Where(map[string]interface{}{"auth_i_d": authId}).
		First(&balService).Error; err != nil {
		//there might be database down we need to see on this one
		return nil, err
	}
	return &balService, err
}

func UpdateBalanceByParentAuthIdAuthId(parentAuthId string, bal *BalanceService) error {
	if count := DB().Model(bal).Where(map[string]interface{}{"auth_i_d": bal.AuthId, "parent_auth_id": parentAuthId}).
		UpdateColumn("balance_money", gorm.Expr("balance_money + ?", bal.Balance)).RowsAffected; count == 0 {
		return errors.New("Not Found")
	}
	return nil
}

func UpdateBalanceByAuthId(bal *BalanceService) error {
	if count := DB().Model(&bal).Where(map[string]interface{}{"auth_i_d": bal.AuthId}).
		UpdateColumn("balance_money", gorm.Expr("balance_money + ?", bal.Balance)).RowsAffected; count == 0 {
		return errors.New("Not Found")
	}
	return nil
}

func IsVendor(authId string) bool {
	var count = 0
	DB().Model(BalanceService{}).Where(map[string]interface{}{"auth_i_d": authId, "type": "vendor"}).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

func (s BalanceService) TableName() string {
	return "account"
}
