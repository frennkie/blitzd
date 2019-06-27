package metric

type SetStatic interface {
	Set()
}

func Set(a SetStatic) { a.Set() }

type StartTimeBased interface {
	Start()
}

func Start(a StartTimeBased) { a.Start() }

type UpdateEventBased interface {
	Update()
}

func Update(a UpdateEventBased) { a.Update() }
