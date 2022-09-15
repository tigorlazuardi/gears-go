package prelude

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	once         = &sync.Once{}
	ctrlcContext = context.Background()
)

/*
First call to this function will make a signal capturer that will catch SIGINT (CTRL+C)
and a new context that will receive Done signal when SIGINT signal is received,
then that context will be returned

Subsequent calls will not make a new signal capturer and
only returns the same context as the first created context.
*/
func GetCtrlCContext() context.Context {
	once.Do(func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		ctx, exit := context.WithCancel(ctrlcContext)
		ctrlcContext = ctx

		go func() {
			for sig := range c {
				fmt.Printf("\n[Prelude] Received Signal: %s\n", sig)
				exit()
				break
			}
		}()
	})

	return ctrlcContext
}
