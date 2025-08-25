package request

import "github.com/i9si-sistemas/stringx"

type QueryParams map[string]string

func (q QueryParams) Parse(target string) {
	qStr := ""
	if idx := stringx.New(target).IndexOf("?"); idx != -1 {
		qStr = target[idx+1:]
	} else {
		return
	}

	pairs := stringx.New(qStr).Split("&")
	for _, pair := range pairs {
		kv := stringx.New(pair).Split("=")
		if len(kv) == 2 {
			q[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			q[kv[0]] = stringx.Empty.String()
		}
	}
}

func (q QueryParams) Get(key string) string {
	return q[key]
}

func (q QueryParams) Reset() QueryParams {
	return make(QueryParams)
}
