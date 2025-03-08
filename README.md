# calc

Package calc provides advanced parsers for floats, ints, ..., based on Go constants calculator.

# Introduction

When it comes to manually writing numbers, the usual literal format are not always the easiest.
How do you write the equivalent of 1 day but in seconds? "24*60*60" is probably easier than
figuring out it is "86400". Why, then asking your users to provide CLI arguments (or inputs in a textfield)
as a number, when you could easily ask them to enter it as a basic formula?

Go has figure that out, and has created a powerful [constants] systems that can be
used to higly improve parsing basic types.

[constants]: [https://go.dev/blog/constants](https://go.dev/blog/constants)

## Examples

### Int

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

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
