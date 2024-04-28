package dto

type CommonResponse struct {
	Code    int    `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ReserveCouponRequest struct {
	ActiveID uint64 `json:"active_id"`
	UserID   uint32 `json:"user_id"`
}
type ReserveCouponResponse struct {
	CommonResponse
}
