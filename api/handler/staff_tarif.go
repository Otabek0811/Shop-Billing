package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create staff_tarif godoc
// @ID create_staff_tarif
// @Router /staff_tarif [POST]
// @Summary Create Staff_Tarif
// @Description Create Staff_Tarif
// @Tags Staff_Tarif
// @Accept json
// @Procedure json
// @Param staff_tarif body models.CreateStaff_Tarif true "CreateStaff_TarifRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateStaff_Tarif(c *gin.Context) {

	var createStaff_Tarif models.CreateStaff_Tarif
	err := c.ShouldBindJSON(&createStaff_Tarif)
	if err != nil {
		h.handlerResponse(c, "error staff_tarif should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Staff_Tarif().Create(c.Request.Context(), &createStaff_Tarif)
	if err != nil {
		h.handlerResponse(c, "storage.staff_tarif.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Staff_Tarif().GetByID(c.Request.Context(), &models.Staff_TarifPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff_Tarif resposne", http.StatusCreated, resp)
}

// GetByID staff_tarif godoc
// @ID get_by_id_staff_tarif
// @Router /staff_tarif/{id} [GET]
// @Summary Get By ID Staff_Tarif
// @Description Get By ID Staff_Tarif
// @Tags Staff_Tarif
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdStaff_Tarif(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Staff_Tarif().GetByID(c.Request.Context(), &models.Staff_TarifPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Staff_Tarif resposne", http.StatusOK, resp)
}

// GetList staff_tarif godoc
// @ID get_list_staff_tarif
// @Router /staff_tarif [GET]
// @Summary Get List Staff_Tarif
// @Description Get List Staff_Tarif
// @Tags Staff_Tarif
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListStaff_Tarif(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Staff_Tarif offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Staff_Tarif limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Staff_Tarif().GetList(c.Request.Context(), &models.Staff_TarifGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Staff_Tarif resposne", http.StatusOK, resp)
}

// Update staff_tarif godoc
// @ID update_staff_tarif
// @Router /staff_tarif/{id} [PUT]
// @Summary Update Staff_Tarif
// @Description Update Staff_Tarif
// @Tags Staff_Tarif
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param staff_tarif body models.UpdateStaff_Tarif true "UpdateStaff_TarifRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateStaff_Tarif(c *gin.Context) {

	var (
		id           string = c.Param("id")
		updateStaff_Tarif models.UpdateStaff_Tarif
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateStaff_Tarif)
	if err != nil {
		h.handlerResponse(c, "error Staff_Tarif should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateStaff_Tarif.Id = id
	rowsAffected, err := h.strg.Staff_Tarif().Update(c.Request.Context(), &updateStaff_Tarif)
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Staff_Tarif.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Staff_Tarif().GetByID(c.Request.Context(), &models.Staff_TarifPrimaryKey{Id: updateStaff_Tarif.Id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff_Tarif resposne", http.StatusAccepted, resp)
}

// Patch staff_tarif godoc
// @ID patch_staff_tarif
// @Router /staff_tarif/{id} [PATCH]
// @Summary Patch Staff_Tarif
// @Description Patch Staff_Tarif
// @Tags Staff_Tarif
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param staff_tarif body models.PatchRequest true "PatchStaff_TarifRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) PatchStaff_Tarif(c *gin.Context) {

	var (
		id          string = c.Param("id")
		patchStaff_Tarif models.PatchRequest
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&patchStaff_Tarif)
	if err != nil {
		h.handlerResponse(c, "error Staff_Tarif should bind json", http.StatusBadRequest, err.Error())
		return
	}

	patchStaff_Tarif.ID = id
	rowsAffected, err := h.strg.Staff_Tarif().Patch(c.Request.Context(), &patchStaff_Tarif)
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.patch", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Staff_Tarif.patch", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Staff_Tarif().GetByID(c.Request.Context(), &models.Staff_TarifPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff_Tarif resposne", http.StatusAccepted, resp)
}

// Delete staff_tarif godoc
// @ID delete_staff_tarif
// @Router /staff_tarif/{id} [DELETE]
// @Summary Delete Staff_Tarif
// @Description Delete Staff_Tarif
// @Tags Staff_Tarif
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteStaff_Tarif(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Staff_Tarif().Delete(c.Request.Context(), &models.Staff_TarifPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff_Tarif.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "delete Staff_Tarif resposne", http.StatusNoContent, nil)
}
