package powermanager

import (
	"fmt"

	"github.com/godbus/dbus"
)

func systemdHandler(op string) func() (bool, func() error) {
	return func() (bool, func() error) {
		conn, err := dbus.SystemBus()
		if err != nil {
			return false, nil
		}

		loginManager := conn.Object("org.freedesktop.login1", "/org/freedesktop/login1")
		if loginManager == nil {
			return false, nil
		}

		var can string
		err = loginManager.Call(fmt.Sprintf("org.freedesktop.login1.Manager.Can%s", op), 0).Store(&can)
		if err != nil || can != "yes" {
			return false, nil
		}

		return true, func() error {
			return loginManager.Call(fmt.Sprintf("org.freedesktop.login1.Manager.%s", op), 0, false).Err
		}
	}
}
