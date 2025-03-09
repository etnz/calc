// Package calc provides advanced parsers for floats, ints, ..., based on Go constants calculator.
//
// # Introduction
//
// When it comes to manually write numbers, using literals is not always the easiest way.
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
	"unicode"
	"unicode/utf8"
)

// Float32 computes the float expression.
func Float32(expr string) (float32, error) { return Scope{}.Float32(expr) }

// Float64 computes the float expression.
func Float64(expr string) (float64, error) { return Scope{}.Float64(expr) }

// Complex64 computes the complex expression.
func Complex64(expr string) (c complex64, err error) { return Scope{}.Complex64(expr) }

// Complex128 computes the complex expression.
func Complex128(expr string) (c complex128, err error) { return Scope{}.Complex128(expr) }

// Int computes the int expression.
func Int(expr string) (int64, error) { return Scope{}.Int(expr) }

// Uint computes the int expression.
func Uint(expr string) (uint64, error) { return Scope{}.Uint(expr) }

// Bool computes the bool expression.
func Bool(expr string) (bool, error) { return Scope{}.Bool(expr) }

// String computes the string expression.
func String(expr string) (string, error) { return Scope{}.String(expr) }

// Scope contains a set of [constant.Value] that can be referenced by their name.
//
// zero type is valid.
type Scope struct {
	p *types.Package
}

// eval expr in this Scope. nil value for 'p' is ok.
func (s Scope) eval(expr string) (constant.Value, error) {
	// c.main can be nil, and that is ok.
	tv, err := types.Eval(token.NewFileSet(), s.p, token.NoPos, expr)
	if err != nil {
		return nil, err
	}
	return tv.Value, nil
}

// return a non nil package.
func (s *Scope) pack() *types.Package {
	if s.p == nil {
		s.p = types.NewPackage("main", "main")
	}
	return s.p
}

// assign a value to the variable 'name' if not already defined.
func (s *Scope) assign(name string, tv types.TypeAndValue) {
	s.pack().Scope().Insert(types.NewConst(token.NoPos, s.pack(), name, tv.Type, tv.Value))
}

// Float64 evaluates 'expr' as a float64.
func (s Scope) Float64(expr string) (float64, error) {
	val, err := s.eval(expr)
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

// Float32  evaluates 'expr' as a float32.
func (s Scope) Float32(expr string) (float32, error) {
	val, err := s.eval(expr)
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

// Complex128  evaluates 'expr' as a complex128.
func (s Scope) Complex128(expr string) (cplx complex128, err error) {
	val, err := s.eval(expr)
	if err != nil {
		return
	}
	// Force conversion to a constant.Float type (or Unknown)
	fval := constant.ToComplex(val)
	if fval.Kind() == constant.Unknown {
		return cplx, fmt.Errorf("not representable as a complex (%v): %q", val.Kind(), expr)
	}

	r, _ := constant.Float64Val(constant.Real(fval)) // ignoring the bool about rounding
	i, _ := constant.Float64Val(constant.Imag(fval)) // ignoring the bool about rounding

	return complex(r, i), nil
}

// Complex64 evaluates 'expr' as a complex64.
func (s Scope) Complex64(expr string) (cplx complex64, err error) {
	val, err := s.eval(expr)
	if err != nil {
		return
	}
	// Force conversion to a constant.Float type (or Unknown)
	fval := constant.ToComplex(val)
	if fval.Kind() == constant.Unknown {
		return cplx, fmt.Errorf("not representable as a complex (%v): %q", val.Kind(), expr)
	}

	r, _ := constant.Float32Val(constant.Real(fval)) // ignoring the bool about rounding
	i, _ := constant.Float32Val(constant.Imag(fval)) // ignoring the bool about rounding

	return complex(r, i), nil
}

// Int evaluates 'expr' as a int64.
func (s Scope) Int(expr string) (int64, error) {
	val, err := s.eval(expr)
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

// Uint evaluates 'expr' as an int64.
func (s Scope) Uint(expr string) (uint64, error) {
	val, err := s.eval(expr)
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

// Bool evaluates 'expr' as a bool.
func (s Scope) Bool(expr string) (bool, error) {
	val, err := s.eval(expr)
	if err != nil {
		return false, err
	}
	if val.Kind() != constant.Bool {
		return false, fmt.Errorf("not representable as a bool (%v): %q", val.Kind(), expr)
	}
	return constant.BoolVal(val), nil
}

// String evaluates 'expr' as a string.
func (s Scope) String(expr string) (string, error) {
	val, err := s.eval(expr)
	if err != nil {
		return "", err
	}
	if val.Kind() != constant.String {
		return "", fmt.Errorf("not representable as a string (%v): %q", val.Kind(), expr)
	}
	return constant.StringVal(val), nil
}

// Assign evaluates 'expr' and assign its value to the variable 'name'.
//
// If the variable 'name' already exists, its value is not changed.
func (s *Scope) Assign(name, expr string) error {
	tv, err := types.Eval(token.NewFileSet(), s.p, token.NoPos, expr)
	if err != nil {
		return err
	}
	s.assign(name, tv)
	return nil
}

// AssignValue directly assign the runtime value 'v' to the variable 'name'.
// 'v' must be one of:
//
//	float64
//	float32
//	complex128
//	complex64
//	int64
//	int32
//	int16
//	int8
//	int
//	uint64
//	uint32
//	uint16
//	uint8
//	uint
//	bool
//	string
//
// If the variable 'name' already exists, its value is not changed.
func (s *Scope) AssignValue(name string, v any) {
	switch o := v.(type) {
	case float64:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedFloat],
			Value: constant.MakeFloat64(o),
		})
	case float32:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedFloat],
			Value: constant.MakeFloat64(float64(o)),
		})
	case complex128:
		x := constant.MakeFloat64(real(o))
		y := constant.MakeFloat64(imag(o))

		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedComplex],
			Value: constant.BinaryOp(x, token.ADD, constant.MakeImag(y)),
		})
	case complex64:
		x := constant.MakeFloat64(float64(real(o)))
		y := constant.MakeFloat64(float64(imag(o)))

		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedComplex],
			Value: constant.BinaryOp(x, token.ADD, constant.MakeImag(y)),
		})
	case int64:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeInt64(o),
		})
	case int32:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeInt64(int64(o)),
		})
	case int16:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeInt64(int64(o)),
		})
	case int8:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeInt64(int64(o)),
		})
	case int:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeInt64(int64(o)),
		})
	case uint64:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeUint64(o),
		})
	case uint32:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeUint64(uint64(o)),
		})
	case uint16:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeUint64(uint64(o)),
		})
	case uint8:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeUint64(uint64(o)),
		})
	case uint:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedInt],
			Value: constant.MakeUint64(uint64(o)),
		})
	case bool:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedBool],
			Value: constant.MakeBool(o),
		})
	case string:
		s.assign(name, types.TypeAndValue{
			Type:  types.Typ[types.UntypedString],
			Value: constant.MakeString(o),
		})
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

// Import another [Scope] inside this one.
//
// Exposed variables in 'lib' can be referenced as `<name>.<var>`.
//
// Following the rules of Go, only Capitalized variables are exposed.
//
// An error is returned if 'name' is exported.
func (s *Scope) Import(name string, lib *Scope) error {
	ch, _ := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(ch) {
		return fmt.Errorf("package names cannot be exported: %v", name)
	}
	pkgName := types.NewPkgName(token.NoPos, s.pack(), name, lib.pack())
	s.pack().Scope().Insert(pkgName)
	return nil
}
