package powermanager

import (
	"log"
	"testing"
)

func TestSuspend(t *testing.T) {
	log.Print(Suspend())
}
