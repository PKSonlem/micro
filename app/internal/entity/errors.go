package entity

import "errors"

var (
	ErrorCreatingHouse       = errors.New("error creating house")
	ErrorCreatingFlat        = errors.New("error creating flat")
	ErrorUpdateModeratorFlat = errors.New("error update flat")
)
