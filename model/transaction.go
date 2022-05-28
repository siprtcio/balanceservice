package model

import (
	"time"

	"github.com/siprtcio/balanceservice/model/orm"
)

type Transaction struct {
	orm.Model        `gorm:"-"`
	TxId             string    `gorm:"column:tx_id" json:"tx_id"`
	StripeId         string    `gorm:"column:stripe_tx_id" json:"stripe_tx_id"`
	Created          time.Time `gorm:"column:created" json:"created"`
	Updated          time.Time `gorm:"column:updated" json:"updated"`
	AccId            string    `gorm:"column:acc_i_d" json:"acc_i_d"`
	Amount           float64   `gorm:"column:amount" json:"amount"`
	Currency         string    `gorm:"column:currency" json:"currency"`
	RechargeAmount   float64   `gorm:"column:recharge_amount" json:"recharge_amount"`
	RechargeCurrency string    `gorm:"column:recharge_currency" json:"recharge_currency"`
	TxType           string    `gorm:"column:tx_type" json:"tx_type"`
	ParentAccId      string    `gorm:"column:parent_acc_id" json:"parent_acc_id"`
	Status           string    `gorm:"column:pay_pal_status" json:"status"`
	ConversionRate   float64   `gorm:"column:conversion_rate" json:"conversion_rate"`
}

func (trans *Transaction) CreateTransaction() error {
	trans.Created = time.Now()
	trans.Updated = time.Now()
	trans.RechargeAmount = trans.Amount
	trans.RechargeCurrency = "USD"
	trans.StripeId = ""
	trans.Currency = "USD"
	trans.TxType = "MANUAL_RECHARGE"
	trans.ConversionRate = 0
	trans.Status = "COMPLETED"
	if err := DB().Create(trans).Error; err != nil {
		return err
	}
	return nil
}

//vendor tenant signup
//vendor change username/oassword
func (trans Transaction) TableName() string {
	return "transaction"
}
