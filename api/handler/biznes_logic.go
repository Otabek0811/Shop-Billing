package handler

import (
	"app/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get Top Staff godoc
// @ID get_top_staff
// @Router /get_top_staff [GET]
// @Summary Get Top Staff
// @Description Get Top Staff
// @Tags Top Staff
// @Accept json
// @Procedure json
// @Param from_date query string false "from_date"
// @Param to_date query string false "to_date"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetTopStaff(c *gin.Context) {

	var (
		fromDate string = c.Query("from_date")
		toDate   string = c.Query("to_date")
	)

	resp, err := h.strg.Biznes_Logic().GetTopStaff(c.Request.Context(), &models.TopStaffRequest{
		FromDate: fromDate,
		ToDate:   toDate,
	})
	if err != nil {
		h.handlerResponse(c, "storage.Get.Top.Staff", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get Top Staff response", http.StatusOK, resp)
}


// Get Top Branchs godoc
// @ID get_top_branch
// @Router /get_top_branch [GET]
// @Summary Get Top Branchs
// @Description Get Top Branchs
// @Tags Top Branchs
// @Accept json
// @Procedure json
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetTopBranchs(c *gin.Context) {

	resp, err := h.strg.Biznes_Logic().GetTopBranchByDate(c.Request.Context())
	if err != nil {
		h.handlerResponse(c, "storage.Get.Top.Branch", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get Top Branch response", http.StatusOK, resp)
}

