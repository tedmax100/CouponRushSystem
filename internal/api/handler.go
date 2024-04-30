package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tedmax100/CouponRushSystem/internal/api/dto"
	"github.com/tedmax100/CouponRushSystem/internal/coupon"
	"github.com/tedmax100/CouponRushSystem/internal/types"
	"github.com/tedmax100/CouponRushSystem/internal/user"
	"go.uber.org/fx"
)

type HandlerParams struct {
	fx.In

	Ctx           context.Context
	CouponService *coupon.CouponActiveService
	UserService   *user.UserSertive
}

type Handler struct {
	couponService *coupon.CouponActiveService
	userService   *user.UserSertive
}

func NewHandler(p HandlerParams) *Handler {
	return &Handler{
		couponService: p.CouponService,
		userService:   p.UserService,
	}
}

// ReserveCoupon reserves a coupon
// @Summary Reserve a coupon
// @Description Reserve a coupon
// @Tags Coupons
// @Accept  json
// @Produce  json
// @Param reserveReq body dto.ReserveCouponRequest true "Reserve Coupon Request"
// @Success 200 {object} dto.ReserveCouponResponse
// @Failure 400 {object} dto.CommonErrorResponse
// @Failure 401 {object} dto.CommonErrorResponse
// @Failure 404 {object} dto.CommonErrorResponse
// @Failure 500 {object} dto.CommonErrorResponse
// @Router /coupons/reserve [post]
func (h *Handler) ReserveCoupon(c *gin.Context) {
	var reserveReq dto.ReserveCouponRequest

	now := time.Now().UTC()

	if err := c.ShouldBindJSON(&reserveReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonErrorResponse{
			Code:    http.StatusBadRequest,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	if _, err := h.userService.GetUser(reserveReq.UserID); err != nil {
		if errors.Is(err, types.ErrorUserNotFound) {
			c.JSON(http.StatusUnauthorized, dto.CommonErrorResponse{
				Code:    http.StatusUnauthorized,
				Path:    c.Request.URL.Path,
				Message: types.ErrorUserNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CommonErrorResponse{
			Code:    http.StatusInternalServerError,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	couponActive, err := h.couponService.GetActive(c, reserveReq.ActiveID)
	if err != nil {
		if errors.Is(err, types.ErrorCouponActiveNotFound) {
			c.JSON(http.StatusUnauthorized, dto.CommonErrorResponse{
				Code:    http.StatusUnauthorized,
				Path:    c.Request.URL.Path,
				Message: types.ErrorUserNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusNotFound, dto.CommonErrorResponse{
			Code:    http.StatusNotFound,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	if !couponActive.IsValidToReserve(now) {
		c.JSON(http.StatusForbidden, dto.CommonErrorResponse{
			Code:    http.StatusForbidden,
			Path:    c.Request.URL.Path,
			Message: "coupon active is not valid to reserve",
		})
		return
	}

	if err := h.couponService.ReserveCoupon(c, couponActive, reserveReq.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CommonErrorResponse{
			Code:    http.StatusInternalServerError,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.ReserveCouponResponse{
		UserID:       reserveReq.UserID,
		ActiveID:     reserveReq.ActiveID,
		ReservedTime: time.Now().UTC(),
	})
}

// PurchaseCoupon purchases a coupon
// @Summary Purchase a coupon
// @Description Purchase a coupon
// @Tags Coupons
// @Accept  json
// @Produce  json
// @Param purchaseReq body dto.PurchaseCouponRequest true "Purchase Coupon Request"
// @Success 200 {object} dto.PurchaseCouponResponse
// @Failure 400 {object} dto.CommonErrorResponse
// @Failure 401 {object} dto.CommonErrorResponse
// @Failure 404 {object} dto.CommonErrorResponse
// @Failure 500 {object} dto.CommonErrorResponse
// @Router /coupons/purchase [post]
func (h *Handler) PurchaseCoupon(c *gin.Context) {
	var reserveReq dto.ReserveCouponRequest

	now := time.Now().UTC()

	if err := c.ShouldBindJSON(&reserveReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonErrorResponse{
			Code:    http.StatusBadRequest,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	if _, err := h.userService.GetUser(reserveReq.UserID); err != nil {
		if errors.Is(err, types.ErrorUserNotFound) {
			c.JSON(http.StatusUnauthorized, dto.CommonErrorResponse{
				Code:    http.StatusUnauthorized,
				Path:    c.Request.URL.Path,
				Message: types.ErrorUserNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CommonErrorResponse{
			Code:    http.StatusInternalServerError,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	couponActive, err := h.couponService.GetActive(c, reserveReq.ActiveID)
	if err != nil {
		if errors.Is(err, types.ErrorCouponActiveNotFound) {
			c.JSON(http.StatusUnauthorized, dto.CommonErrorResponse{
				Code:    http.StatusUnauthorized,
				Path:    c.Request.URL.Path,
				Message: types.ErrorUserNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusNotFound, dto.CommonErrorResponse{
			Code:    http.StatusNotFound,
			Path:    c.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	if !couponActive.IsValidToPurchase(now) {
		c.JSON(http.StatusForbidden, dto.CommonErrorResponse{
			Code:    http.StatusForbidden,
			Path:    c.Request.URL.Path,
			Message: "coupon active is not valid to purchase",
		})
		return
	}

	purchasedCoupon, err :=
		h.couponService.PurchaseCoupon(c, couponActive, reserveReq.UserID)
	if err != nil {
		if errors.Is(err, types.ErrorNoCouponToPurchase) {
			c.JSON(http.StatusForbidden, dto.CommonErrorResponse{
				Code:    types.ErrorNoCouponToPurchase.Code,
				Path:    c.Request.URL.Path,
				Message: types.ErrorNoCouponToPurchase.Error(),
			})
			return
		}
		if errors.Is(err, types.ErrorUserAlreadyPurchasedCoupon) {
			c.JSON(http.StatusForbidden, dto.CommonErrorResponse{
				Code:    types.ErrorUserAlreadyPurchasedCoupon.Code,
				Path:    c.Request.URL.Path,
				Message: types.ErrorNoCouponToPurchase.Error(),
			})
			return
		}
		c.JSON(http.StatusForbidden, dto.CommonErrorResponse{
			Code:    http.StatusForbidden,
			Path:    c.Request.URL.Path,
			Message: "coupon active is not valid to purchase",
		})
		return
	}

	c.JSON(http.StatusOK, dto.PurchaseCouponResponse{
		UserID:     reserveReq.UserID,
		ActiveID:   reserveReq.ActiveID,
		CouponCode: purchasedCoupon.Code,
	})
}
