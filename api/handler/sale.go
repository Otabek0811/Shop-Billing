package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create sale godoc
// @ID create_sale
// @Router /sale [POST]
// @Summary Create Sale
// @Description Create Sale
// @Tags Sale
// @Accept json
// @Procedure json
// @Param sale body models.CreateSale true "CreateSaleRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateSale(c *gin.Context) {

	var (
		createSale models.CreateSale
	)
	err := c.ShouldBindJSON(&createSale)
	if err != nil {
		h.handlerResponse(c, "error sale should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Sale().Create(c.Request.Context(), &createSale)
	if err != nil {
		h.handlerResponse(c, "storage.sale.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Sale().GetByID(c.Request.Context(), &models.SalePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sale.getById", http.StatusInternalServerError, err.Error())
		return
	}
	err = h.strg.Biznes_Logic().Do_Staff_Transaction(c.Request.Context(), resp)
	if err != nil {
		h.handlerResponse(c, "storage.Sale.staff_transaction", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Sale response", http.StatusCreated, resp)
}

// GetByID sale godoc
// @ID get_by_id_sale
// @Router /sale/{id} [GET]
// @Summary Get By ID Sale
// @Description Get By ID Sale
// @Tags Sale
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdSale(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Sale().GetByID(c.Request.Context(), &models.SalePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sale.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Sale resposne", http.StatusOK, resp)
}

// GetList sale godoc
// @ID get_list_sale
// @Router /sale [GET]
// @Summary Get List Sale
// @Description Get List Sale
// @Tags Sale
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search_branch_id query string false "search_branch_id"
// @Param search_client_name query string false "search_client_name"
// @Param search_payment_type query string false "search_payment_type"
// @Param search_assistant_id query string false "search_assistant_id"
// @Param search_status query string false "search_status"
// @Param search_created_at_from query string false "search_created_at_from"
// @Param search_created_at_to query string false "search_created_at_to"
// @Param sort_price_type query string false "sort_price_type"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListSale(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Sale offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Sale limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Sale().GetList(c.Request.Context(), &models.SaleGetListRequest{
		Offset:              offset,
		Limit:               limit,
		SearchBranchID:      c.Query("search_branch_id"),
		SearchClientName:    c.Query("search_client_name"),
		SearchPaymentType:   c.Query("search_payment_type"),
		SearchAssistantID:   c.Query("search_assistant_id"),
		SearchCashierID:     c.Query("search_cashier_id"),
		SearchStatus:        c.Query("search_status"),
		SearchCreatedAtFrom: c.Query("search_created_at_from"),
		SearchCreatedAtTo:   c.Query("search_created_at_to"),
		SortPriceType:       c.Query("sort_price_type"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Sale.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Sale resposne", http.StatusOK, resp)
}

// Update sale godoc
// @ID update_sale
// @Router /sale/{id} [PUT]
// @Summary Update Sale
// @Description Update Sale
// @Tags Sale
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param sale body models.UpdateSale true "UpdateSaleRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateSale(c *gin.Context) {

	var (
		id         string = c.Param("id")
		updateSale models.UpdateSale
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateSale)
	if err != nil {
		h.handlerResponse(c, "error Sale should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateSale.Id = id
	rowsAffected, err := h.strg.Sale().Update(c.Request.Context(), &updateSale)
	if err != nil {
		h.handlerResponse(c, "storage.Sale.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Sale.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Sale().GetByID(c.Request.Context(), &models.SalePrimaryKey{Id: updateSale.Id})
	if err != nil {
		h.handlerResponse(c, "storage.Sale.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Sale resposne", http.StatusAccepted, resp)
}

// Patch sale godoc
// @ID patch_sale
// @Router /sale/{id} [PATCH]
// @Summary Patch Sale
// @Description Patch Sale
// @Tags Sale
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param sale body models.PatchRequest true "PatchSaleRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) PatchSale(c *gin.Context) {

	var (
		id        string = c.Param("id")
		patchSale models.PatchRequest
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&patchSale)
	if err != nil {
		h.handlerResponse(c, "error Sale should bind json", http.StatusBadRequest, err.Error())
		return
	}

	patchSale.ID = id
	rowsAffected, err := h.strg.Sale().Patch(c.Request.Context(), &patchSale)
	if err != nil {
		h.handlerResponse(c, "storage.Sale.patch", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Sale.patch", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Sale().GetByID(c.Request.Context(), &models.SalePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sale.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Sale resposne", http.StatusAccepted, resp)
}

// Delete sale godoc
// @ID delete_sale
// @Router /sale/{id} [DELETE]
// @Summary Delete Sale
// @Description Delete Sale
// @Tags Sale
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteSale(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Sale().Delete(c.Request.Context(), &models.SalePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sale.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "delete Sale resposne", http.StatusNoContent, nil)
}



