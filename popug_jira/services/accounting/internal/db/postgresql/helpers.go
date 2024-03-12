package postgresql

import (
	"fmt"
	"reflect"
	"time"
)

func retryFunc(maxAttempts int, backoffInt time.Duration, handleFunc func() error) error {
	attempt := 0
	backoff := backoffInt

	funcName := reflect.TypeOf(handleFunc).Name()

	for attempt < maxAttempts {
		if err := handleFunc(); err == nil {
			return nil
		}

		attempt++
		if attempt == maxAttempts {
			return fmt.Errorf("maximum number of attempts for success complete function %q reached", funcName)
		}

		time.Sleep(backoffInt)
		backoff += backoffInt
	}

	return fmt.Errorf("maximum number of attempts for success complete function %q reached", funcName)
}
