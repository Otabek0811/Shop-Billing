package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)


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
// @Param search_sale_id query string false "search_sale_id"
// @Param search_source_type query string false "search_source_type"
// @Param search_price_type query string false "search_price_type"
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
		SortPriceType:   c.Query("search_price_type"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list StaffTransaction resposne", http.StatusOK, resp)
}
