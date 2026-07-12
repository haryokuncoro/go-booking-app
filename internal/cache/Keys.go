package cache

import (
	"fmt"

	"github.com/google/uuid"
)

func BookingKey(
	id uuid.UUID,
) string {

	return fmt.Sprintf(
		"booking:%s",
		id,
	)
}