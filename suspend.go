package powermanager

import (
	"errors"

	"github.com/godbus/dbus"
)

func CanSuspend() bool {
	for _, handler := range suspendHandlers {
		ok, _ := handler()
		if ok {
			return true
		}
	}
	return false
}

func Suspend() error {
	for _, handler := range suspendHandlers {
		ok, fn := handler()
		if ok {
			return fn()
		}
	}

	return errors.New("No suspend handler")
}

func suspendUPower() (bool, func() error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return false, nil
	}

	upower := conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower")
	if upower == nil {
		return false, nil
	}

	canSuspend, err := upower.GetProperty("org.freedesktop.UPower.CanSuspend")
	if err != nil || canSuspend.Value() != true {
		return false, nil
	}

	var suspendAllowed bool
	err = upower.Call("org.freedesktop.UPower.SuspendAllowed", 0).Store(&suspendAllowed)
	if err != nil || !suspendAllowed {
		return false, nil
	}

	return true, func() error {
		return upower.Call("org.freedesktop.UPower.Suspend", 0).Err
	}
}

var suspendHandlers = []func() (bool, func() error){
	systemdHandler("Suspend"),
	suspendUPower,
}
