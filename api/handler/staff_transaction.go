package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create staff_transaction godoc
// @ID create_staff_transaction
// @Router /staff_transaction [POST]
// @Summary Create StaffTransaction
// @Description Create StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param staff_transaction body models.CreateStaffTransaction true "CreateStaffTransactionRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateStaffTransaction(c *gin.Context) {

	var createStaffTransaction models.CreateStaffTransaction
	err := c.ShouldBindJSON(&createStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "error staff_transaction should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.StaffTransaction().Create(c.Request.Context(), &createStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "storage.staff_transaction.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create StaffTransaction resposne", http.StatusCreated, resp)
}

// GetByID staff_transaction godoc
// @ID get_by_id_staff_transaction
// @Router /staff_transaction/{id} [GET]
// @Summary Get By ID StaffTransaction
// @Description Get By ID StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdStaffTransaction(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id StaffTransaction resposne", http.StatusOK, resp)
}

// GetList staff_transaction godoc
// @ID get_list_staff_transaction
// @Router /staff_transaction [GET]
// @Summary Get List StaffTransaction
// @Description Get List StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search_Name query string false "search_by_name"
// @Param search_Address query string false "search_by_address"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListStaffTransaction(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list StaffTransaction offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list StaffTransaction limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.StaffTransaction().GetList(c.Request.Context(), &models.StaffTransactionGetListRequest{
		Offset: offset,
		Limit:  limit,
		SearchSaleID:      c.Query("search_sale_id"),
		SearchStaffID:    c.Query("search_staff_id"),
		SearchSourceType:   c.Query("search_source_type"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list StaffTransaction resposne", http.StatusOK, resp)
}

// Update staff_transaction godoc
// @ID update_staff_transaction
// @Router /staff_transaction/{id} [PUT]
// @Summary Update StaffTransaction
// @Description Update StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param staff_transaction body models.UpdateStaffTransaction true "UpdateStaffTransactionRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateStaffTransaction(c *gin.Context) {

	var (
		id           string = c.Param("id")
		updateStaffTransaction models.UpdateStaffTransaction
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "error StaffTransaction should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateStaffTransaction.Id = id
	rowsAffected, err := h.strg.StaffTransaction().Update(c.Request.Context(), &updateStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.StaffTransaction.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: updateStaffTransaction.Id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create StaffTransaction resposne", http.StatusAccepted, resp)
}

// Patch staff_transaction godoc
// @ID patch_staff_transaction
// @Router /staff_transaction/{id} [PATCH]
// @Summary Patch StaffTransaction
// @Description Patch StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param staff_transaction body models.PatchRequest true "PatchStaffTransactionRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) PatchStaffTransaction(c *gin.Context) {

	var (
		id          string = c.Param("id")
		patchStaffTransaction models.PatchRequest
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&patchStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "error StaffTransaction should bind json", http.StatusBadRequest, err.Error())
		return
	}

	patchStaffTransaction.ID = id
	rowsAffected, err := h.strg.StaffTransaction().Patch(c.Request.Context(), &patchStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.patch", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.StaffTransaction.patch", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create StaffTransaction resposne", http.StatusAccepted, resp)
}

// Delete staff_transaction godoc
// @ID delete_staff_transaction
// @Router /staff_transaction/{id} [DELETE]
// @Summary Delete StaffTransaction
// @Description Delete StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteStaffTransaction(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.StaffTransaction().Delete(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "delete StaffTransaction resposne", http.StatusNoContent, nil)
}
