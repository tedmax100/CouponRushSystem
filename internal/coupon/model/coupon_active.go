package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

type CouponActive struct {
	ID              uint64            `json:"id" db:"id"`
	ActiveDate      time.Time         `json:"date" db:"date"`
	ActiveBeginTime time.Time         `json:"beginTime" db:"begin_time"`
	ActiveEndTime   time.Time         `json:"endTime" db:"end_time"`
	State           CouponActiveState `json:"state" db:"state"`
}

func (s CouponActive) IsValidToReserve(now time.Time) bool {
	return s.State == OPENING || s.ActiveBeginTime.Before(now) && s.ActiveEndTime.After(now)
}

func (s CouponActive) IsValidToPurchase(now time.Time) bool {
	return s.State == OPENING || now.After(s.ActiveEndTime)
}

type CouponActiveState int

const (
	NOT_OPEN CouponActiveState = iota
	OPENING
	CLOSED
)

func (s CouponActiveState) String() string {
	return [...]string{"NOT_OPEN", "OPENING", "CLOSED"}[s]
}

func (s *CouponActiveState) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}

	str := string(asBytes)
	switch str {
	case "NOT_OPEN":
		*s = NOT_OPEN
	case "OPENING":
		*s = OPENING
	case "CLOSED":
		*s = CLOSED
	default:
		return errors.New("Invalid CouponActiveState")
	}

	return nil
}

func (s CouponActiveState) Value() (driver.Value, error) {
	return s.String(), nil
}
