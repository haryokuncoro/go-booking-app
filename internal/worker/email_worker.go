package worker

import (
	"fmt"
	"time"
)

func StartEmailWorker(
	id int,
	jobs <-chan EmailJob,
) {

	for job := range jobs {

		fmt.Printf(
			"[Worker %d] Sending email to %s\n",
			id,
			job.UserEmail,
		)

		time.Sleep(
			3 * time.Second,
		)

		fmt.Printf(
			"[Worker %d] Email sent for room %s\n",
			id,
			job.RoomName,
		)
	}
}