package cleanup

import (
	"context"
	"log"
	"time"

	"github.com/apialerts/hooks/internal/db"
)

type Job struct {
	db *db.Queries
}

func New(db *db.Queries) *Job {
	return &Job{db: db}
}

func (j *Job) Start(ctx context.Context) {
	// Run immediately on startup
	go func() {
		j.run(ctx)

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.run(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (j *Job) run(ctx context.Context) {
	deleted, err := j.db.DeleteExpiredEndpoints(ctx)
	if err != nil {
		log.Printf("cleanup error: %v", err)
		return
	}
	if deleted > 0 {
		log.Printf("cleanup: deleted %d expired endpoints", deleted)
	}
}
