package scheduler

import (
	"binp/storage"
	"binp/util"
	"context"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() Scheduler {
	return Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Init(store *storage.Store) {
	logger := util.GetLogger()
	s.AddFunc("@hourly", func() {
		logger.Info().Msg("Checking for expired snippets...")
		count, err := store.DeleteExpiredSnippets()
		if err != nil {
			logger.Error().Err(err).Int("count", count).Msg("Failed to delete expired snippets")
		}
		logger.Info().Int("count", count).Msg("Expired snippets deleted")
	})
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}

func (s *Scheduler) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}
