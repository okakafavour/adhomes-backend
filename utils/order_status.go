package utils

var AllowedStatuses = []string{
	"Processing",
	"Preparing",
	"Out for Delivery",
	"Delivered",
}

func IsValidStatus(status string) bool {
	for _, s := range AllowedStatuses {
		if s == status {
			return true
		}
	}
	return false
}
