package app

import (
	"reflect"
)

func validateKeeper(keepers ...any) {
	for _, k := range keepers {
		if k == nil || reflect.ValueOf(k).IsZero() {
			panic("Keeper is nil. You should initialize the keepers properly before using")
		}
	}
}
