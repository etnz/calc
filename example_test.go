package calc_test

import (
	"fmt"
	"strconv"

	"github.com/etnz/calc"
)

// Using [calc.Int] you can parse basic literals,
// exactly as [strconv.ParseInt] can do, and evaluate more advanced
// expressions.
//
// The example below showcases some more posibilities.
func ExampleInt() {

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

	// Output:
	// Literal: 2 = 2
	// Package strconv: 2 = 2
	// Algebra: 2+2 = 4
	// Safe overflow: 1<<100 + 2 - 1<<100 = 2
	// Octal: 0777 = 511
	// Hex: 0xFF = 255
	// Binary: 0b1010 ^ 0b0101 = 15
	// Mixed: 0xFF - 0b11111110 = 1
}

// When writing expressions it can be handy to use predefined constants.
// It is possible to prepare a [calc.Scope] with predefined variables.
func ExampleScope_Assign() {
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

	// Output:
	// Time: 2.5*d = 216000
	// Time: 2*d + 4*h = 187200
	// Float: 2.5*s = 2.5
}

// When writing expressions using variables, it is possible
// to "pack" them in their own namespace to avoid conflicts.
func ExampleScope_Import() {
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

	// Output:
	// Time: 2*time.D + 4*time.H = 187200
}
