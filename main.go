// Package calc provides advanced parsers for floats, ints, ..., based on Go constants calculator.
//
// # Introduction
//
// When it comes to manually writing numbers, the usual literal format are not always the easiest.
// How do you write the equivalent of 1 day but in seconds? `24*60*60` is probably easier than
// figuring out it is `86400`. Why, then asking your users to provide CLI arguments (or inputs in a textfield)
// as a number, when you could easily ask them to enter it as a basic formula?
//
// Go has figure that out, and has created a powerful [constants] systems that can be
// used to higly improve parsing basic types.
//
// [constants]: https://go.dev/blog/constants
package calc

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"math"
)

// evalConst reads any valid Go constant expression, and returns its value.
//
// It supports string, bool float and int computation, all that is represented by
// [constant.Value] interface.
// The result is always a [constant.Value].
//
// Most of the time you might be interested in converting it to a float64, or
// an int64, or an uint64.
// See constant.Float64Val() ... functions to convert the result down to your need.
func evalConst(expr string) (constant.Value, error) {
	tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, expr)
	if err != nil {
		return nil, err
	}
	return tv.Value, nil
}

// Float64 computes the float expression.
func Float64(expr string) (float64, error) {
	val, err := evalConst(expr)
	if err != nil {
		return math.NaN(), err
	}
	// Force conversion to a constant.Float type (or Unknown)
	fval := constant.ToFloat(val)
	if fval.Kind() == constant.Unknown {
		return math.NaN(), fmt.Errorf("not representable as a float (%v): %q", val.Kind(), expr)
	}
	f, _ := constant.Float64Val(fval) // ignoring the bool about rounding
	return f, nil
}

// Float32 computes the float expression.
func Float32(expr string) (float32, error) {
	val, err := evalConst(expr)
	if err != nil {
		return float32(math.NaN()), err
	}
	// Force conversion to a constant.Float type (or Unknown)
	fval := constant.ToFloat(val)
	if fval.Kind() == constant.Unknown {
		return float32(math.NaN()), fmt.Errorf("not representable as a float (%v): %q", val.Kind(), expr)
	}
	f, _ := constant.Float32Val(fval) // ignoring the bool about rounding
	return f, nil
}

// Complex128 computes the complex expression.
func Complex128(expr string) (c complex128, err error) {
	val, err := evalConst(expr)
	if err != nil {
		return
	}
	// Force conversion to a constant.Float type (or Unknown)
	fval := constant.ToComplex(val)
	if fval.Kind() == constant.Unknown {
		return c, fmt.Errorf("not representable as a complex (%v): %q", val.Kind(), expr)
	}

	r, _ := constant.Float64Val(constant.Real(fval)) // ignoring the bool about rounding
	i, _ := constant.Float64Val(constant.Imag(fval)) // ignoring the bool about rounding

	return complex(r, i), nil
}

// Complex64 computes the complex expression.
func Complex64(expr string) (c complex64, err error) {
	val, err := evalConst(expr)
	if err != nil {
		return
	}
	// Force conversion to a constant.Float type (or Unknown)
	fval := constant.ToComplex(val)
	if fval.Kind() == constant.Unknown {
		return c, fmt.Errorf("not representable as a complex (%v): %q", val.Kind(), expr)
	}

	r, _ := constant.Float32Val(constant.Real(fval)) // ignoring the bool about rounding
	i, _ := constant.Float32Val(constant.Imag(fval)) // ignoring the bool about rounding

	return complex(r, i), nil
}

// Int computes the int expression.
func Int(expr string) (int64, error) {
	val, err := evalConst(expr)
	if err != nil {
		return 0, err
	}
	// Force conversion to a constant.Float type (or Unknown)
	ival := constant.ToInt(val)
	if ival.Kind() == constant.Unknown {
		return 0, fmt.Errorf("not representable as an int (%v): %q", val.Kind(), expr)
	}
	i, ok := constant.Int64Val(ival)
	if !ok {
		return 0, fmt.Errorf("not exactly representable as an int64: %q", expr)
	}
	return i, nil
}

// Uint computes the int expression.
func Uint(expr string) (uint64, error) {
	val, err := evalConst(expr)
	if err != nil {
		return 0, err
	}
	// Force conversion to a constant.Float type (or Unknown)
	ival := constant.ToInt(val)
	if ival.Kind() == constant.Unknown {
		return 0, fmt.Errorf("not representable as an int (%v): %q", val.Kind(), expr)
	}
	i, ok := constant.Uint64Val(ival)
	if !ok {
		return 0, fmt.Errorf("not exactly representable as an uint64: %q", expr)
	}
	return i, nil
}

// Bool computes the bool expression.
func Bool(expr string) (bool, error) {
	val, err := evalConst(expr)
	if err != nil {
		return false, err
	}
	if val.Kind() != constant.Bool {
		return false, fmt.Errorf("not representable as a bool (%v): %q", val.Kind(), expr)
	}
	return constant.BoolVal(val), nil
}

// String computes the string expression.
func String(expr string) (string, error) {
	val, err := evalConst(expr)
	if err != nil {
		return "", err
	}
	if val.Kind() != constant.String {
		return "", fmt.Errorf("not representable as a string (%v): %q", val.Kind(), expr)
	}
	return constant.StringVal(val), nil
}
