package web

type FilterBuilder func(next Filter) Filter

type Filter func(c *Context)
