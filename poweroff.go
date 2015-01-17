package powermanager

import "errors"

func CanPowerOff() bool {
	for _, handler := range powerOffHandlers {
		ok, _ := handler()
		if ok {
			return true
		}
	}
	return false
}

func PowerOff() error {
	for _, handler := range powerOffHandlers {
		ok, fn := handler()
		if ok {
			return fn()
		}
	}

	return errors.New("No powerOff handler")
}

var powerOffHandlers = []func() (bool, func() error){
	systemdHandler("PowerOff"),
}
