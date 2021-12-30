package action

import "sort"

func newContextSet() *contextSet {
	return &contextSet{
		m: map[string]*Context{},
	}
}

type contextSet struct {
	m map[string]*Context
}

func (c *contextSet) Add(context *Context) {
	c.m[context.Function] = context
}

func (c *contextSet) List() (contexts []*Context) {
	for _, context := range c.m {
		contexts = append(contexts, context)
	}
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Function < contexts[j].Function
	})
	return contexts
}
