package query

type Params map[string]string

func (p Params) Get(key string) string {
	return p[key]
}

func (p Params) Contains(key string) bool {
	_, ok := p[key]
	return ok
}

func FromURLQuery(q map[string][]string) Params {
	params := make(Params)
	for k, v := range q {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}
