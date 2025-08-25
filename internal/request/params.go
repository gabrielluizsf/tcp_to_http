package request

import "github.com/i9si-sistemas/stringx"

type Params map[string]string

func (p Params) Set(lineTarget, endpoint string) {
	separator := "/"
	getParts := func(s string) []string {
		return stringx.New(s).Trim(separator).Split(separator)
	}
	targetParts := getParts(lineTarget)
	endpointParts := getParts(endpoint)

	for i := range endpointParts {
		ep := stringx.New(endpointParts[i])
		if ep.HasPrefix("{") && ep.HasSuffix("}") {
			key := ep[1 : len(ep)-1]
			p[key.String()] = targetParts[i]
		} else if ep.HasPrefix(":") {
			key := ep[1:]
			p[key.String()] = targetParts[i]
		}
	}
}

func (p Params) Get(paramName string) string {
	return p[paramName]
}

func (p Params) Reset() Params {
	return make(Params)
}
