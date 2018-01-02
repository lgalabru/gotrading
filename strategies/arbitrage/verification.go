package arbitrage

type Verification struct {
	IsSuccessful bool
}

func (ver *Verification) Init(exec Execution) {
}

func (ver *Verification) Run() {
}

func (ver *Verification) BuildReport() Report {
	return Report{}
}
