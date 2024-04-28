package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tedmax100/CouponRushSystem/internal/api/dto"
	"github.com/tedmax100/CouponRushSystem/internal/coupon"
	"github.com/tedmax100/CouponRushSystem/internal/types"
	"github.com/tedmax100/CouponRushSystem/internal/user"
)

type Handler struct {
	couponService *coupon.CouponActiveService
	userService   *user.UserSertive
}

func NewHandler() *Handler {
	return &Handler{}
}

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

	h.couponService.ReserveCoupon(c, couponActive, reserveReq.UserID)
}

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
