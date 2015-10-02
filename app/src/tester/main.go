package main

import "log"
import "fmt"
import "os"
import "time"

var MySQLContainers = []string{"mysql-alpha", "mysql-beta"}

func main() {
	logger := log.New(os.Stderr, "[tester] ", log.LstdFlags)
	session, e := NewSession(logger, MySQLContainers...)
	noError(e)

	run := func() {
		noError(session.Status())
		noError(session.Query())
		noError(session.Mutate())
		noError(session.Query())
	}

	reset := func() { noError(session.EnsureStarted("mysql-alpha", "mysql-beta")) }
	start := func(name string) { noError(session.Start(name)) }
	stop := func(name string) { noError(session.Stop(name)) }

	fillDisk := func() {
		noError(session.Create("filler"))
		time.Sleep(1 * time.Second)
		logger.Println("waiting for filler container.")

		noError(session.Start("filler"))
		logger.Println("waiting for filler to fill up.")
		time.Sleep(10 * time.Second) // wait for disk to fill up
	}

	undoFillDisk := func() {
		noError(session.Destroy("filler"))
	}

	reconfigure := func(instances ...string) {
		noError(session.Reconfigure(instances...))
		time.Sleep(1 * time.Second)
	}

	runError := func() {
		noError(session.Status("mysql-alpha", "mysql-beta"))

		// enough runs so the balancer visit as many nodes as it can.
		for i := 0; i < 5; i++ {
			if e := session.Query(); e != nil {
				logger.Println("got expected error:", e)
				return
			}
			if e := session.Mutate(); e != nil {
				logger.Println("got expected error:", e)
				return
			}
			if e := session.Query(); e != nil {
				logger.Println("got expected error:", e)
				return
			}
		}

		logger.Println("got no errors (unexpected).")
	}

	_, _, _ = run, start, stop
	_, _ = fillDisk, undoFillDisk

	for {
		reset()
		reconfigure("alpha", "beta")

		// run()
		// stop("mysql-alpha")
		// runError()

		// reconfigure("beta")
		// run()
		// stop("mysql-beta")
		// runError()

		// start("mysql-alpha")
		// runError()
		// reconfigure("alpha")
		// run()
		// start("mysql-beta")
		// reconfigure("alpha", "beta")
		// run()

		run()
		run()

		fillDisk()
		runError()
		runError()
		undoFillDisk()
	}
}

func noError(e error) {
	if e != nil {
		fmt.Println("UNEXPECTED ERROR:", e)
	}
}
