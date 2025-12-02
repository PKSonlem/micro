package entity

const (
	CreatedStatus      = "created"
	ApprovedStatus     = "approved"
	DeclinedStatus     = "declined"
	OnModerationStatus = "on moderation"
)

var validStatus = map[string]bool{
	CreatedStatus:      true,
	ApprovedStatus:     true,
	DeclinedStatus:     true,
	OnModerationStatus: true,
}

func IsValidStatus(status string) bool {
	return validStatus[status]
}
