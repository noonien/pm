package powermanager

import "errors"

func CanReboot() bool {
	for _, handler := range rebootHandlers {
		ok, _ := handler()
		if ok {
			return true
		}
	}
	return false
}

func Reboot() error {
	for _, handler := range rebootHandlers {
		ok, fn := handler()
		if ok {
			return fn()
		}
	}

	return errors.New("No reboot handler")
}

var rebootHandlers = []func() (bool, func() error){
	systemdHandler("Reboot"),
}
