package stdversion

import "reflect"

var _ = reflect.TypeFor[int]() // ERROR "reflect.TypeFor requires golang1.22 or later \(module is golang1.21\)"
