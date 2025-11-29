package entity

const (
	CreatedStatus      = "created"
	ApprovedStatus     = "approved"
	DeclinedStatus     = "declined"
	OnModerationStatus = "on moderation"
)

var validStatuses = map[string]bool{
	CreatedStatus:      true,
	ApprovedStatus:     true,
	DeclinedStatus:     true,
	OnModerationStatus: true,
}

func IsValidStatus(status string) bool {
	return validStatuses[status]
}
