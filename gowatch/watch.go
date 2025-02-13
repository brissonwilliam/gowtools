package watch

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"
)

func Start(withInitRun bool, serviceFilters []string) error {
	log := logger.NewLogger()
	db, err := db.Connect(config.GetDatabaseConfiguration())
	if err != nil {
		log.WithError(err).Error("Could not connect to database")
		return err
	}

	c := cron.New()

	watches := []service.WatchService{
		anomaly.NewWatch(db),
		inactivity.NewWatch(db),
	}

	// validate names
	for _, filter := range serviceFilters {
		valid := false
		for _, s := range watches {
			if filter == s.Name() {
				valid = true
				break
			}
		}
		if !valid {
			newErr := errors.New("invalid --services provided: " + filter)
			log.Error(newErr)
			return newErr
		}
	}

	// add services to cron runtime
	started := 0
	for i := range watches {
		s := watches[i]
		if !shouldRun(s.Name(), serviceFilters) {
			continue
		}

		// parse / validate cron
		sched := s.GetCronSchedule()
		cronSched, err := cron.ParseStandard(sched)
		if err != nil {
			s.Log().WithError(err).Error("could not read cron tab schedule of ")
			return err
		}

		// add to cron
		s.Log().Info("adding service with cron schedule " + sched)

		f := wrapServiceRun(s)
		_, err = c.AddFunc(sched, f)
		if err != nil {
			s.Log().WithError(err).Error("could not set cron func")
			return err
		}

		// init run
		if withInitRun {
			s.Log().Info("doing initial run")
			go f()
		}

		// print next run
		nextRun := cronSched.Next(time.Now())
		s.Log().Info("next schedule is " + nextRun.String())

		started++
	}
	c.Start()

	if started < 1 {
		log.Warn("no services configured to start. Iddling")
	}

	webSrv := echo.New()
	webSrv.HideBanner = true
	webSrv.GET("/ping", Ping)

	log.Info("starting http server on :8080")
	return webSrv.Start(":8080") // blocking
}

func wrapServiceRun(s service.WatchService) func() {
	return func() {
		now := time.Now()
		s.Log().Info("starting watch service run")

		s.Run()

		duration := time.Since(now)
		s.Log().WithField("duration", duration.String()).Info("watch service run DONE")
	}
}

func shouldRun(serviceName string, serviceFilters []string) bool {
	if len(serviceFilters) < 1 {
		// default empty = all = true
		return true
	}
	for _, filter := range serviceFilters {
		if serviceName == filter {
			return true
		}
	}
	return false
}

func Ping(ctx echo.Context) error {
	return ctx.String(200, "pong")
}
