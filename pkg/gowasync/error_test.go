package gowasync

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestErrorSyncGroup(t *testing.T) {
	// try to add errors concurrently, then validate GetFirst
	errGr := ErrorSyncGroup{}

	rand.Seed(time.Now().Unix())

	n := 100
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			randSleepT := time.Duration(rand.Intn(5)) * time.Nanosecond
			time.Sleep(randSleepT)

			newErr := errors.New("error " + strconv.Itoa(i))
			errGr.Add(newErr)

			wg.Done()
		}(i)
	}
	wg.Wait()

	firstErr := errGr.GetFirst()
	assert.NotNil(t, firstErr)
	assert.Error(t, firstErr)

	t.Run("Succeeds when no errors in sync group", func(t *testing.T) {
		errGr := ErrorSyncGroup{}
		assert.Nil(t, errGr.GetFirst())
	})

}
