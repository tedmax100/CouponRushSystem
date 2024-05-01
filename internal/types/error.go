package types

import "fmt"

type SystenError struct {
	Code    int
	Message string
}

func (e *SystenError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

var (
	ErrorUserNotFound                    = &SystenError{Code: 1001, Message: "User not found"}
	ErrorCouponActiveNotFound            = &SystenError{Code: 1002, Message: "Coupon active not found"}
	ErrorCouponNotFound                  = &SystenError{Code: 1003, Message: "Coupon not found"}
	ErrorUserNotReserveCouponActive      = &SystenError{Code: 1004, Message: "User not reserve coupon active"}
	ErrorUserAlreadyPurchasedCoupon      = &SystenError{Code: 1005, Message: "User already purchased coupon"}
	ErrorNoCouponToPurchase              = &SystenError{Code: 1006, Message: "No coupon to purchase"}
	ErrorUserAlreadyReservedCouponActive = &SystenError{Code: 1007, Message: "User already reserved coupon active"}
	ErrorCouponActiveNotValidToPurchase  = &SystenError{Code: 1008, Message: "Coupon active is not valid to purchase"}
	ErrorCouponActiveNotValidToReserve   = &SystenError{Code: 1009, Message: "Coupon active is not valid to reserve"}
)
