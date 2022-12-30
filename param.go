package goblin

// Param is a parameter.
type Param struct {
	key   string
	value string
}

// getParams gets parameters.
func (t *tree) getParams() *[]Param {
	ps, _ := t.paramsPool.Get().(*[]Param)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

// putParams puts parameters.
func (t *tree) putParams(ps *[]Param) {
	if ps != nil {
		t.paramsPool.Put(ps)
	}
}
