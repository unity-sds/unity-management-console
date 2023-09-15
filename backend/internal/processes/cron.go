package processes

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"time"
)

func RunSync() {

	store, err := database.NewGormDatastore()
	if err != nil {
		log.WithError(err).Error("Error creating datastore")
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Start a goroutine that runs the function every time the ticker ticks
	go func() {
		for {
			select {
			case <-ticker.C:
				fetchAllApplications(store)
			}
		}
	}()

	// Keep the main goroutine running
	select {}

}
