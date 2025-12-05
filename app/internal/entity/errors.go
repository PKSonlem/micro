package entity

import "errors"

var (
	ErrorCreatingHouse       = errors.New("error creating house")
	ErrorCreatingFlat        = errors.New("error creating flat")
	ErrorUpdateModeratorFlat = errors.New("error update flat")
	ErrorCreatingUser        = errors.New("error creating user")
	ErrorLoginUser           = errors.New("error login user")
	ErrorCreateSubs          = errors.New("error creating subscription")
	ErrorHouseNotFound       = errors.New("house not found")
)
