package powermanager

import "errors"

func CanHybridSleep() bool {
	for _, handler := range hybridSleepHandlers {
		ok, _ := handler()
		if ok {
			return true
		}
	}
	return false
}

func HybridSleep() error {
	for _, handler := range hybridSleepHandlers {
		ok, fn := handler()
		if ok {
			return fn()
		}
	}

	return errors.New("No hybridSleep handler")
}

var hybridSleepHandlers = []func() (bool, func() error){
	systemdHandler("HybridSleep"),
}
