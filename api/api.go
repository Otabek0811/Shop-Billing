package api

import (
	_ "app/api/docs"
	"app/api/handler"
	"app/config"
	"app/pkg/logger"
	"app/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewApi(r *gin.Engine, cfg *config.Config, storage storage.StorageI, logger logger.LoggerI) {

	handler := handler.NewHandler(cfg, storage, logger)

	r.POST("/branch", handler.CreateBranch)
	r.GET("/branch/:id", handler.GetByIdBranch)
	r.GET("/branch", handler.GetListBranch)
	r.PUT("/branch/:id", handler.UpdateBranch)
	r.PATCH("/branch/:id", handler.PatchBranch)
	r.DELETE("/branch/:id", handler.DeleteBranch)

	r.POST("/staff_tarif", handler.CreateStaff_Tarif)
	r.GET("/staff_tarif/:id", handler.GetByIdStaff_Tarif)
	r.GET("/staff_tarif", handler.GetListStaff_Tarif)
	r.PUT("/staff_tarif/:id", handler.UpdateStaff_Tarif)
	r.PATCH("/staff_tarif/:id", handler.PatchStaff_Tarif)
	r.DELETE("/staff_tarif/:id", handler.DeleteStaff_Tarif)

	r.POST("/staff", handler.CreateStaff)
	r.GET("/staff/:id", handler.GetByIdStaff)
	r.GET("/staff", handler.GetListStaff)
	r.PUT("/staff/:id", handler.UpdateStaff)
	r.PATCH("/staff/:id", handler.PatchStaff)
	r.DELETE("/staff/:id", handler.DeleteStaff)

	r.POST("/sale", handler.CreateSale)
	r.GET("/sale/:id", handler.GetByIdSale)
	r.GET("/sale", handler.GetListSale)
	r.PUT("/sale/:id", handler.UpdateSale)
	r.PATCH("/sale/:id", handler.PatchSale)
	r.DELETE("/sale/:id", handler.DeleteSale)

	r.GET("/staff_transaction/:id", handler.GetByIdStaffTransaction)
	r.GET("/staff_transaction", handler.GetListStaffTransaction)

	r.GET("/get_top_staff", handler.GetTopStaff)
	r.GET("/get_top_branch", handler.GetTopBranchs)



	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
