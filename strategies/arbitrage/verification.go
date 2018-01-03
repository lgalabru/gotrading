package arbitrage

type Verification struct {
	Report Report
}

func (ver *Verification) Init(exec Execution) {
}

func (ver *Verification) Run() {
}

func (ver *Verification) IsSuccessful() bool {
	return ver.Report.IsVerificationSuccessful
}
