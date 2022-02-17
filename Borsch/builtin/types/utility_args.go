package types

func ParseTupleAndKeywords(args Tuple, kwargs StringDict, format string, keywords []string, results ...*Object) error {
	if keywords != nil && len(results) != len(keywords) {
		return ErrorNewf(TypeError, "Internal error: supply the same number of results and keywords")
	}

	var opsBuf [16]formatOp
	min, name, kwOnlyI, ops := parseFormat(format, opsBuf[:0])
	err := checkNumberOfArgs(name, len(args)+len(kwargs), len(results), min, len(ops))
	if err != nil {
		return err
	}

	// Check all the kwargs are in keywords
	// O(N^2) Slow but keywords is usually short
	for kwargName := range kwargs {
		for _, kw := range keywords {
			if kw == kwargName {
				goto found
			}
		}
		return ErrorNewf(TypeError, "%s() got an unexpected keyword argument '%s'", name, kwargName)
	found:
	}

	// Walk through all the results we want
	for i, op := range ops {

		var (
			arg Object
			kw  string
		)
		if i < len(keywords) {
			kw = keywords[i]
			arg = kwargs[kw]
		}

		// Consume ordered args first -- they should not require keyword only or also be specified via keyword
		if i < len(args) {
			if i >= kwOnlyI {
				return ErrorNewf(TypeError, "%s() specifies argument '%s' that is keyword only", name, kw)
			}

			if arg != nil {
				return ErrorNewf(TypeError, "%s() got multiple values for argument '%s'", name, kw)
			}

			arg = args[i]
		}

		// Unspecified args retain their default value
		if arg == nil {
			continue
		}

		result := results[i]
		switch op.code {
		case 'O':
			*result = arg
		case 'Z', 'z':
			if _, ok := arg.(NilType); ok {
				*result = arg
				break
			}
			fallthrough
		case 'U', 's':
			if _, ok := arg.(String); !ok {
				return ErrorNewf(TypeError, "%s() argument %d must be str, not %s", name, i+1, arg.Type().Name)
			}

			*result = arg
		case 'i':
			if _, ok := arg.(Int); !ok {
				return ErrorNewf(TypeError, "%s() argument %d must be int, not %s", name, i+1, arg.Type().Name)
			}

			*result = arg
		case 'p':
			if _, ok := arg.(Bool); !ok {
				return ErrorNewf(TypeError, "%s() argument %d must be bool, not %s", name, i+1, arg.Type().Name)
			}

			*result = arg
		case 'd':
			switch x := arg.(type) {
			case Int:
				*result = Real(x)
			case Real:
				*result = x
			default:
				return ErrorNewf(TypeError, "%s() argument %d must be float, not %s", name, i+1, arg.Type().Name)
			}

		default:
			return ErrorNewf(
				TypeError,
				"Unknown/Unimplemented format character %q in ParseTupleAndKeywords called from %s",
				op,
				name,
			)
		}
	}

	return nil
}

// ParseTuple parses tuple only.
func ParseTuple(args Tuple, format string, results ...*Object) error {
	return ParseTupleAndKeywords(args, nil, format, nil, results...)
}

type formatOp struct {
	code     byte
	modifier byte
}

func parseFormat(format string, in []formatOp) (min int, name string, kwOnlyI int, ops []formatOp) {
	name = "функція"
	min = -1
	kwOnlyI = 0xFFFF
	ops = in[:0]

	N := len(format)
	for i := 0; i < N; {
		op := formatOp{code: format[i]}
		i++
		if i < N {
			if mod := format[i]; mod == '*' || mod == '#' {
				op.modifier = mod
				i++
			}
		}

		switch op.code {
		case ':', ';':
			name = format[i:]
			i = N
		case '$':
			kwOnlyI = len(ops)
		case '|':
			min = len(ops)
		default:
			ops = append(ops, op)
		}
	}
	if min < 0 {
		min = len(ops)
	}

	return
}

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

// UnpackTuple unpacks the args tuple into the results.
//
// Up to the caller to set default values.
func UnpackTuple(args Tuple, kwargs StringDict, name string, min int, max int, results ...*Object) error {
	if len(kwargs) != 0 {
		return ErrorNewf(TypeError, "%s() не приймає аргументи у вигляді словника", name)
	}

	err := checkNumberOfArgs(name, len(args), len(results), min, max)
	if err != nil {
		return err
	}

	for i := range args {
		*results[i] = args[i]
	}

	return nil
}
