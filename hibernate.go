package powermanager

import (
	"errors"

	"github.com/godbus/dbus"
)

func CanHibernate() bool {
	for _, handler := range hibernateHandlers {
		ok, _ := handler()
		if ok {
			return true
		}
	}
	return false
}

func Hibernate() error {
	for _, handler := range hibernateHandlers {
		ok, fn := handler()
		if ok {
			return fn()
		}
	}

	return errors.New("No hibernate handler")
}

func hibernateUPower() (bool, func() error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return false, nil
	}

	upower := conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower")
	if upower == nil {
		return false, nil
	}

	canHibernate, err := upower.GetProperty("org.freedesktop.UPower.CanHibernate")
	if err != nil || canHibernate.Value() != true {
		return false, nil
	}

	var hibernateAllowed bool
	err = upower.Call("org.freedesktop.UPower.HibernateAllowed", 0).Store(&hibernateAllowed)
	if err != nil || !hibernateAllowed {
		return false, nil
	}

	return true, func() error {
		return upower.Call("org.freedesktop.UPower.Hibernate", 0).Err
	}
}

var hibernateHandlers = []func() (bool, func() error){
	systemdHandler("Hibernate"),
	hibernateUPower,
}
