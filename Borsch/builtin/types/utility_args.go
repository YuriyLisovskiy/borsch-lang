package types

import (
	"strings"
)

// checkNumberOfArgs checks the number of args passed in.
func checkNumberOfArgs(name string, nArgs, nResults, min, max int) error {
	if min == max {
		if nArgs != max {
			return ErrorNewf(TypeError, "%s() takes exactly %d arguments (%d given)", name, max, nArgs)
		}
	} else {
		if nArgs > max {
			return ErrorNewf(TypeError, "%s() takes at most %d arguments (%d given)", name, max, nArgs)
		}

		if nArgs < min {
			return ErrorNewf(TypeError, "%s() takes at least %d arguments (%d given)", name, min, nArgs)
		}
	}

	if nArgs > nResults {
		return ErrorNewf(TypeError, "Internal error: not enough arguments supplied to Unpack*/Parse*")
	}

	return nil
}

func ParseExactArgs(args Tuple, format string, results ...*Object) error {
	nResults := len(results)
	nArgs := len(args)
	min := nResults
	formatPair := strings.SplitN(format, ":", 2)
	name := formatPair[0]
	err := checkNumberOfArgs(name, nArgs, nResults, min, nResults)
	if err != nil {
		return err
	}

	formatTypes := formatPair[1]
	for i, c := range formatTypes {
		arg := args[i]
		result := results[i]
		switch c {
		// any, object
		case 'a', 'o', 'O':
			if (arg == Nil || arg == nil) && c == 'O' {
				break
			}

			*result = arg
		case 's', 'S':
			if _, ok := arg.(String); !ok {
				if c == 'S' {
					break
				}

				return ErrorNewf(TypeError, "%s() argument %d must be str, not %s", name, i+1, arg.Class().Name)
			}

			*result = arg
		case 'b', 'B':
			if _, ok := arg.(Bool); !ok {
				if c == 'B' {
					break
				}

				return ErrorNewf(TypeError, "%s() argument %d must be bool, not %s", name, i+1, arg.Class().Name)
			}

			*result = arg
		case 'i', 'I':
			if _, ok := arg.(Int); !ok {
				if c == 'I' {
					break
				}

				return ErrorNewf(TypeError, "%s() argument %d must be int, not %s", name, i+1, arg.Class().Name)
			}

			*result = arg
		case 'r', 'R':
			switch x := arg.(type) {
			case Int:
				*result = Real(x)
			case Real:
				*result = x
			default:
				if c == 'R' {
					break
				}

				return ErrorNewf(TypeError, "%s() argument %d must be real, not %s", name, i+1, arg.Class().Name)
			}
		case 'l', 'L':
			if _, ok := arg.(*List); !ok {
				if c == 'L' {
					break
				}

				return ErrorNewf(TypeError, "%s() argument %d must be list, not %s", name, i+1, arg.Class().Name)
			}

			*result = arg

		case 'd', 'D':
			if _, ok := arg.(Dict); !ok {
				if c == 'D' {
					break
				}

				return ErrorNewf(TypeError, "%s() argument %d must be Dict, not %s", name, i+1, arg.Class().Name)
			}

			*result = arg
		default:
			return ErrorNewf(
				TypeError,
				"Unknown/Unimplemented format character %q in ParseExactArgs called from %s",
				c,
				name,
			)
		}
	}

	return nil
}

// UnpackTuple unpacks the args tuple into the results.
//
// Up to the caller to set default values.
func UnpackTuple(args Tuple, name string, min int, max int, results ...*Object) error {
	err := checkNumberOfArgs(name, len(args), len(results), min, max)
	if err != nil {
		return err
	}

	for i := range args {
		*results[i] = args[i]
	}

	return nil
}
