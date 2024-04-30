package api

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"
	//docs "github.com/tedmax100/CouponRushSystem/api/docs"
	//swaggerfiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(engine *gin.Engine, p HandlerParams) {
	h := NewHandler(p)

	docs.SwaggerInfo.BasePath = "/api"
	v1 := engine.Group("/api/v1")
	v1.POST("/coupons/reservations", h.ReserveCoupon)
	v1.POST("/coupons/purchases", h.PurchaseCoupon)
}
