package api

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/siprtcio/balanceservice/model"
)

// @Description This API allows you to retrieve your account balance or tanent account balance.
// @Summary Get Account balance by AuthID.
// @Tags Balance
// @Produce  json
// @Param   auth_id     	path    string     true        "Tiniyo Account Auth ID"
// @Success 200 {object} model.BalanceService
// @Failure 400 {string} Message "Bad Request"
// @Failure 404 {string} Message "Account balance not found"
// @security BasicAuth
// @Router /Accounts/{auth_id}/Balance [get]
func GetBalanceByAuthId(c *Context) error {
	authId := c.Param("auth_id")
	if authId == "" {
		return c.JSON(http.StatusBadRequest, "AuthId not Valid")
	}
	balanceDetails, err := model.GetBalanceByAuthId(authId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	if balanceDetails == nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, balanceDetails)
}

// @Description This API allows you to retrieve your tenant account balance using your vendor auth id.
// @Summary Get Account balance by vendor using AuthID.
// @Tags Balance
// @Produce  json
// @Param   vendor_auth_id     	path    string     true        "Tiniyo Vendor Account Auth ID"
// @Param   tenant_auth_id     	path    string     true        "Tiniyo Tenant Account Auth ID"
// @Success 200 {object} model.BalanceService
// @Failure 400 {string} Message "Bad Request"
// @Failure 404 {string} Message "Account balance not found"
// @Failure 412 {string} Message "StatusPreconditionFailed : Vendor authid is not vendor"
// @security BasicAuth
// @Router /Accounts/{vendor_auth_id}/Balances/{tenant_auth_id} [get]
func GetBalanceByVendorAuthIdAuthId(c *Context) error {
	vendorAuthId := c.Param("vendor_auth_id")
	tenantAuthId := c.Param("tenant_auth_id")
	if vendorAuthId == "" || tenantAuthId == "" {
		return c.JSON(http.StatusBadRequest, "AuthId not Valid")
	}
	if !model.IsVendor(vendorAuthId) {
		return c.JSON(http.StatusPreconditionFailed, "Vendor AuthID is not valid vendor")
	}
	balanceDetails, err := model.GetBalanceByParentAuthIdAuthId(vendorAuthId, tenantAuthId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	if balanceDetails == nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, balanceDetails)
}

// @Description This API allows you to recharge your tenant account balance using your vendor auth id. Balance can be credit or debit based on +/- ve value. Balance value is in USD($).
// @Summary Recharge Account balance by vendor using AuthID.
// @Tags Balance
// @Accept  json
// @Produce  json
// @Param   vendor_auth_id     	path    string     true        "Tiniyo Vendor Account Auth ID"
// @Param   tenant_auth_id     	path    string     true        "Tiniyo Tenant Account Auth ID"
// @Param message body model.BalanceService true "User Data"
// @Success 200 {string} Message "Balance updated"
// @Failure 400 {string} Message "Bad Request"
// @Failure 412 {string} Message "StatusPreconditionFailed : Vendor authid is not vendor"
// @Failure 406 {string} Message "StatusNotAcceptable: Json validation failed"
// @Failure 422 {string} Message "StatusUnprocessableEntity: update account balance failed"
// @security BasicAuth
// @Router /Accounts/{vendor_auth_id}/Balances/{tenant_auth_id} [patch]
func UpdateBalanceByVendorAuthIdAuthId(c *Context) error {
	vendorAuthId := c.Param("vendor_auth_id")
	tenantAuthId := c.Param("tenant_auth_id")

	if vendorAuthId == "" || tenantAuthId == "" {
		return c.JSON(http.StatusBadRequest, "AuthId not Valid")
	}
	if !model.IsVendor(vendorAuthId) {
		return c.JSON(http.StatusPreconditionFailed, "Vendor AuthID is not valid vendor")
	}
	balanceModel := model.BalanceService{}
	balanceModel.AuthId = tenantAuthId

	if err := c.Bind(&balanceModel); err != nil {
		return c.JSON(http.StatusBadRequest, "Binding Error")
	}

	if err := c.Validate(balanceModel); err != nil {
		return c.JSON(http.StatusNotAcceptable, balanceModel.AuthId)
	}

	err := model.UpdateBalanceByParentAuthIdAuthId(vendorAuthId, &balanceModel)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	TransactionModel := new(model.Transaction)

	transId, _ := uuid.NewV4()
	uuid := fmt.Sprintf("%s", transId)
	TransactionModel.TxId = uuid
	TransactionModel.AccId = tenantAuthId
	TransactionModel.ParentAccId = vendorAuthId
	TransactionModel.Amount = balanceModel.Balance
	err = TransactionModel.CreateTransaction()
	if err != nil {
		//We need to handle it here
		//we need to create both these request as database transaction
	}

	return c.JSON(http.StatusOK, "Balance updated")
}

func UpdateBalanceByAuthId(c *Context) error {
	authId := c.Param("auth_id")
	if authId == "" {
		return c.JSON(http.StatusBadRequest, "AuthId not Valid")
	}

	balanceModel := new(model.BalanceService)
	balanceModel.AuthId = authId
	if err := c.Bind(balanceModel); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}
	if err := c.Validate(balanceModel); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	err := model.UpdateBalanceByAuthId(balanceModel)
	if err != nil {
		return c.JSON(http.StatusPreconditionFailed, err)
	}
	return c.JSON(http.StatusOK, "Balance updated")
}

func UpdateGetBalanceByAuthId(c *Context) error {
	authId := c.Param("auth_id")
	if authId == "" {
		return c.JSON(http.StatusBadRequest, "AuthId not Valid")
	}

	balanceModel := new(model.BalanceService)
	balanceModel.AuthId = authId
	if err := c.Bind(balanceModel); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	if err := c.Validate(balanceModel); err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	err := model.UpdateBalanceByAuthId(balanceModel)
	if err != nil {
		return c.JSON(http.StatusPreconditionFailed, err)
	}

	balanceDetails, err := model.GetBalanceByAuthId(authId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	if balanceDetails == nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, balanceDetails)
}
