### Minor changes to the library {#minor_library_changes}

#### golang/types

The `Var.Kind` method returns an enumeration of type `VarKind` that
classifies the variable (package-level, local, receiver, parameter,
result, or struct field). See issue #70250.

Callers of `NewVar` or `NewParam` are encouraged to call `Var.SetKind`
to ensure that this attribute is set correctly in all cases.
