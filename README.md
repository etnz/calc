# calc

Package calc provides advanced parsers for floats, ints, ..., based on Go constants calculator.

# Introduction

When it comes to manually writing numbers, the usual literal format are not always the easiest.
How do you write the equivalent of 1 day but in seconds? `24*60*60` is probably easier than
figuring out it is `86400`. Why, then asking your users to provide CLI arguments (or inputs in a textfield)
as a number, when you could easily ask them to enter it as a basic formula?

Go has figure that out, and has created a powerful [constants] systems that can be
used to higly improve parsing basic types.

[constants]: [https://go.dev/blog/constants](https://go.dev/blog/constants)

## Functions

### func [Bool](/main.go#L43)

`func Bool(expr string) (bool, error)`

Bool computes the bool expression.

### func [Complex128](/main.go#L34)

`func Complex128(expr string) (c complex128, err error)`

Complex128 computes the complex expression.

### func [Complex64](/main.go#L31)

`func Complex64(expr string) (c complex64, err error)`

Complex64 computes the complex expression.

### func [Float32](/main.go#L25)

`func Float32(expr string) (float32, error)`

Float32 computes the float expression.

### func [Float64](/main.go#L28)

`func Float64(expr string) (float64, error)`

Float64 computes the float expression.

### func [Int](/main.go#L37)

`func Int(expr string) (int64, error)`

Int computes the int expression.

Using [calc.Int] you can parse basic literals,
exactly as [strconv.ParseInt] can do, and evaluate more advanced
expressions.

The example below showcases some more posibilities.

```golang
package main

import (
	"fmt"
	"strconv"

	"github.com/etnz/calc"
)

func main() {

	exp := "2"
	v, _ := calc.Int(exp)
	fmt.Println("Literal:", exp, "=", v)

	// For comparison:
	v, _ = strconv.ParseInt(exp, 10, 64)
	fmt.Println("Package strconv:", exp, "=", v)

	exp = "2+2"
	v, _ = calc.Int(exp)
	fmt.Println("Algebra:", exp, "=", v)

	exp = "1<<100 + 2 - 1<<100"
	v, _ = calc.Int(exp)
	fmt.Println("Safe overflow:", exp, "=", v)

	exp = "0777"
	v, _ = calc.Int(exp)
	fmt.Println("Octal:", exp, "=", v)

	exp = "0xFF"
	v, _ = calc.Int(exp)
	fmt.Println("Hex:", exp, "=", v)

	exp = "0b1010 ^ 0b0101"
	v, _ = calc.Int(exp)
	fmt.Println("Binary:", exp, "=", v)

	exp = "0xFF - 0b11111110"
	v, _ = calc.Int(exp)
	fmt.Println("Mixed:", exp, "=", v)

}

```

 Output:

```
Literal: 2 = 2
Package strconv: 2 = 2
Algebra: 2+2 = 4
Safe overflow: 1<<100 + 2 - 1<<100 = 2
Octal: 0777 = 511
Hex: 0xFF = 255
Binary: 0b1010 ^ 0b0101 = 15
Mixed: 0xFF - 0b11111110 = 1
```

### func [String](/main.go#L46)

`func String(expr string) (string, error)`

String computes the string expression.

### func [Uint](/main.go#L40)

`func Uint(expr string) (uint64, error)`

Uint computes the int expression.

## Types

### type [Scope](/main.go#L51)

`type Scope struct { ... }`

Scope contains a set of [constant.Value] that can be referenced by their name.

zero type is valid.

#### func (*Scope) [Assign](/main.go#L207)

`func (s *Scope) Assign(name, expr string) error`

Assign evaluates 'expr' and assign its value to the variable 'name'.

If the variable 'name' already exists, its value is not changed.

When writing expressions it can be handy to use predefined constants.
It is possible to prepare a [calc.Scope] with predefined variables.

```golang
package main

import (
	"fmt"

	"github.com/etnz/calc"
)

func main() {
	var (
		exp string
		v   int64
		c   calc.Scope
	)

	// prepare the Scope with useful constants.
	c.Assign("s", "1") // 1 second
	c.Assign("m", "60*s")
	c.Assign("h", "60*m")
	c.Assign("d", "24*h")
	c.Assign("w", "7*d")

	// Constants in Go are powerful:
	// expression can use floating point and still lead
	// to an exact integer.
	exp = "2.5*d"
	v, _ = c.Int(exp)
	fmt.Println("Time:", exp, "=", v)

	exp = "2*d + 4*h"
	v, _ = c.Int(exp)
	fmt.Println("Time:", exp, "=", v)

	// Constants in Go are powerful 2
	// It doesn't matter how you have defined them.
	exp = "2.5*s"
	f, _ := c.Float64(exp)
	fmt.Println("Float:", exp, "=", f)

}

```

 Output:

```
Time: 2.5*d = 216000
Time: 2*d + 4*h = 187200
Float: 2.5*s = 2.5
```

#### func (*Scope) [AssignValue](/main.go#L237)

`func (s *Scope) AssignValue(name string, v any)`

AssignValue directly assign the runtime value 'v' to the variable 'name'.
'v' must be one of:

```go
float64
float32
complex128
complex64
int64
int32
int16
int8
int
uint64
uint32
uint16
uint8
uint
bool
string
```

If the variable 'name' already exists, its value is not changed.

#### func (Scope) [Bool](/main.go#L181)

`func (s Scope) Bool(expr string) (bool, error)`

Bool evaluates 'expr' as a bool.

#### func (Scope) [Complex128](/main.go#L109)

`func (s Scope) Complex128(expr string) (cplx complex128, err error)`

Complex128  evaluates 'expr' as a complex128.

#### func (Scope) [Complex64](/main.go#L127)

`func (s Scope) Complex64(expr string) (cplx complex64, err error)`

Complex64 evaluates 'expr' as a complex64.

#### func (Scope) [Float32](/main.go#L94)

`func (s Scope) Float32(expr string) (float32, error)`

Float32  evaluates 'expr' as a float32.

#### func (Scope) [Float64](/main.go#L79)

`func (s Scope) Float64(expr string) (float64, error)`

Float64 evaluates 'expr' as a float64.

#### func (*Scope) [Import](/main.go#L337)

`func (s *Scope) Import(name string, lib *Scope)`

Import another [Scope] inside this one.

Exposed variables in 'lib' can be referenced as `<name>.<var>`.

Following the rules of Go, only Capitalized variables are exposed.

Nesting Scopes is not supported (by Go).

When writing expressions using variables, it is possible
to "pack" them in their own namespace to avoid conflicts.

```golang
package main

import (
	"fmt"

	"github.com/etnz/calc"
)

func main() {
	var (
		exp    string
		v      int64
		c, lib calc.Scope
	)

	// prepare a library of variables.
	lib.Assign("S", "1") // 1 second
	lib.Assign("M", "60*S")
	lib.Assign("H", "60*M")
	lib.Assign("D", "24*H")

	// Import that lib as 'time' to make it usable.
	c.Import("time", &lib)

	exp = "2*time.D + 4*time.H"
	v, _ = c.Int(exp)
	fmt.Println("Time:", exp, "=", v)

}

```

 Output:

```
Time: 2*time.D + 4*time.H = 187200
```

#### func (Scope) [Int](/main.go#L145)

`func (s Scope) Int(expr string) (int64, error)`

Int evaluates 'expr' as a int64.

#### func (Scope) [String](/main.go#L193)

`func (s Scope) String(expr string) (string, error)`

String evaluates 'expr' as a string.

#### func (Scope) [Uint](/main.go#L163)

`func (s Scope) Uint(expr string) (uint64, error)`

Uint evaluates 'expr' as an int64.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
