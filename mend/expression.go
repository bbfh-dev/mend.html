package mend

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func ComputeExpression(params *Parameters, in string) string {
	if strings.HasPrefix(in, ".") {
		return computeParam(params, in)
	}

	parts := strings.SplitN(in, " ", 2)
	if len(parts) != 2 {
		goto return_invalid
	}
	switch parts[0] {

	case "quote":
		return fmt.Sprintf("%q", computeParam(params, parts[1]))

	case "length", "len", "size":
		result, ok := accessParam(params, parts[1])
		if !ok {
			return fmt.Sprintf("{{? NOT_FOUND %q ?}}", parts[1])
		}
		return fmt.Sprintf("%d", len(result.Array()))

	}

return_invalid:
	return fmt.Sprintf("{{? INVALID %q ?}}", in)
}

func accessParam(params *Parameters, in string) (gjson.Result, bool) {
	if strings.HasPrefix(in, ".") {
		if in == "." {
			return params.Get("@this"), true
		} else {
			result := params.Get(in[1:])
			if !result.Exists() {
				return gjson.Result{}, false
			}
			return result, true
		}
	}

	return gjson.Result{}, false
}

func computeParam(params *Parameters, in string) string {
	result, ok := accessParam(params, in)
	if !ok {
		return fmt.Sprintf("{{? NOT_FOUND %q ?}}", in)
	}
	return result.String()
}
