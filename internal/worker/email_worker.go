package worker

import (
	"booking-app/internal/logger"
	"go.uber.org/zap"
	"time"
)

func StartEmailWorker(
	id int,
	jobs <-chan EmailJob,
) {

	for job := range jobs {

		logger.Log.Info(
			"sending email",

			zap.Int(
				"worker_id",
				id,
			),

			zap.String(
				"email",
				job.UserEmail,
			),
		)

		time.Sleep(
			3 * time.Second,
		)

		logger.Log.Info(
			"email sent",

			zap.Int(
				"worker_id",
				id,
			),

			zap.String(
				"room_name",
				job.RoomName,
			),
		)
	}
}
