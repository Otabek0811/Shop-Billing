package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create staff godoc
// @ID create_staff
// @Router /staff [POST]
// @Summary Create Staff
// @Description Create Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param staff body models.CreateStaff true "CreateStaffRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateStaff(c *gin.Context) {

	var createStaff models.CreateStaff
	err := c.ShouldBindJSON(&createStaff)
	if err != nil {
		h.handlerResponse(c, "error staff should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Staff().Create(c.Request.Context(), &createStaff)
	if err != nil {
		h.handlerResponse(c, "storage.staff.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff resposne", http.StatusCreated, resp)
}

// GetByID staff godoc
// @ID get_by_id_staff
// @Router /staff/{id} [GET]
// @Summary Get By ID Staff
// @Description Get By ID Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdStaff(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Staff resposne", http.StatusOK, resp)
}

// GetList staff godoc
// @ID get_list_staff
// @Router /staff [GET]
// @Summary Get List Staff
// @Description Get List Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search_by_name query string false "search_by_name"
// @Param search_by_branch_id query string false "search_by_branch_id"
// @Param search_by_staff_type query string false "search_by_staff_type"
// @Param search_by_tarif_id query string false "search_by_tarif_id"
// @Param search_balance_from query string false "search_balance_from"
// @Param search_balance_to query string false "search_balance_to"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListStaff(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Staff offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Staff limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Staff().GetList(c.Request.Context(), &models.StaffGetListRequest{
		Offset: offset,
		Limit:  limit,
		SearchByName: c.Query("search_by_name"),
		SearchByBranchID: c.Query("search_by_branch_id"),
		SearchBYStaffType: c.Query("search_by_staff_type"),
		SearchByTarifID: c.Query("search_by_tarif_id"),
		SearchBalanceFrom: c.Query("search_balance_from"),
		SearchBalanceTo: c.Query("search_balance_to"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Staff resposne", http.StatusOK, resp)
}

// Update staff godoc
// @ID update_staff
// @Router /staff/{id} [PUT]
// @Summary Update Staff
// @Description Update Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param staff body models.UpdateStaff true "UpdateStaffRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateStaff(c *gin.Context) {

	var (
		id           string = c.Param("id")
		updateStaff models.UpdateStaff
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateStaff)
	if err != nil {
		h.handlerResponse(c, "error Staff should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateStaff.Id = id
	rowsAffected, err := h.strg.Staff().Update(c.Request.Context(), &updateStaff)
	if err != nil {
		h.handlerResponse(c, "storage.Staff.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Staff.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: updateStaff.Id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff resposne", http.StatusAccepted, resp)
}

// Patch staff godoc
// @ID patch_staff
// @Router /staff/{id} [PATCH]
// @Summary Patch Staff
// @Description Patch Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param staff body models.PatchRequest true "PatchStaffRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) PatchStaff(c *gin.Context) {

	var (
		id          string = c.Param("id")
		patchStaff models.PatchRequest
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&patchStaff)
	if err != nil {
		h.handlerResponse(c, "error Staff should bind json", http.StatusBadRequest, err.Error())
		return
	}

	patchStaff.ID = id
	rowsAffected, err := h.strg.Staff().Patch(c.Request.Context(), &patchStaff)
	if err != nil {
		h.handlerResponse(c, "storage.Staff.patch", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Staff.patch", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff resposne", http.StatusAccepted, resp)
}

// Delete staff godoc
// @ID delete_staff
// @Router /staff/{id} [DELETE]
// @Summary Delete Staff
// @Description Delete Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteStaff(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}


	err := h.strg.Staff().Delete(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "delete Staff resposne", http.StatusNoContent, nil)
}


