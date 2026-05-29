package worker

var EmailQueue chan EmailJob

func StartWorkers() {
	//Queue capacity = 100 jobs
	EmailQueue =
		make(
			chan EmailJob,
			100,
		)
	//5 workers
	for i := 1; i <= 5; i++ {

		go StartEmailWorker(
			i,
			EmailQueue,
		)
	}
}