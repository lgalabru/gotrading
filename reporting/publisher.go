package reporting

type Publisher struct {
}

func (pub *Publisher) Init(params string) {
}

func (pub *Publisher) Send(report Report) {
  report.Encoded()
}
