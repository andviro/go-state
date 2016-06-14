package state

// Middleware is a simple modifier for state function. It receives a Func and returns modified Func.
type Middleware func(Func) Func

// Chain is a series of state modifiers that apply sequentially to the initial state Func
func Chain(mws ...Middleware) Middleware {
	return func(final Func) Func {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

// Use appends some more modifiers to the middleware Chain and produces new Chain
func (mw Middleware) Use(others ...Middleware) Middleware {
	return func(final Func) Func {
		return mw(Chain(others...)(final))
	}
}
