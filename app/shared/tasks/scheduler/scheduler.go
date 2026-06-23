package scheduler

import (
	"log"

	"github.com/imohamedsheta/xapp/app/x"
)

// RegisterSchedule defines all periodic tasks
func RegisterSchedule() {
	scheduler := x.Scheduler()
	if scheduler == nil {
		log.Println("[SCHEDULER] Scheduler not initialized, skipping schedule registration")
		return
	}
}
