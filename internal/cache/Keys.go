package cache

import "fmt"

func BookingKey(
	id uint,
) string {

	return fmt.Sprintf(
		"booking:%d",
		id,
	)
}