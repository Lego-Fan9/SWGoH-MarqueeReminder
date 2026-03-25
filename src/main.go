package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Lego-Fan9/MarqueeReminder/comlinkevent"
	"github.com/Lego-Fan9/MarqueeReminder/env"
	log "github.com/sirupsen/logrus"
)

func main() {
	/* env.LoadTemplate()
	MainTask()

	return*/
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	env.LoadTemplate()

	wg.Add(1)

	go MainLoop(ctx, &wg)

	<-sigCh
	log.Info("Shutting down...")

	cancel()

	wg.Wait()
	log.Info("All workers done, exiting")
}

func MainLoop(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Panicf("failed to load timezone: %v", err)
	}

	for {
		now := time.Now().In(loc)

		next := time.Date(
			now.Year(), now.Month(), now.Day(),
			15, 0, 0, 0,
			loc,
		)

		if !now.Before(next) {
			next = next.Add(24 * time.Hour)
		}

		wait := time.Until(next)
		log.Infof("Next run at %s (in %s)", next, wait)

		timer := time.NewTimer(wait)

		select {
		case <-ctx.Done():
			log.Info("shutting down MainLoop")
			timer.Stop()

			return

		case <-timer.C:
			MainTask()
		}
	}
}

func MainTask() {
	marquees, err := comlinkevent.GetActiveMarquees()
	if err != nil {
		log.Errorf("Failed to get marquees: %v", err)

		return
	}

	localization, err := comlinkevent.GetLocalization()
	if err != nil {
		log.Errorf("Failed to get localization: %v", err)

		return
	}

	units, err := comlinkevent.GetUnits()
	if err != nil {
		log.Errorf("Failed to get units: %v", err)

		return
	}

	for _, marquee := range marquees {
		err = PostMarqueeDiscord(marquee, localization, units)
		if err != nil {
			log.Errorf("Failed to post marquee to discord: %v", err)
		}
	}
}
