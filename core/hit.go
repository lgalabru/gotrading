package core

type Hit struct {
	Endpoint       *Endpoint `json:"endpoint"`
	IsBaseToQuote  bool      `json:"isBaseToQuote"`
	SoldCurrency   Currency  `json:"soldCurrency"`
	BoughtCurrency Currency  `json:"boughtCurrency"`
}

func (h Hit) IsEqual(hit Hit) bool {
	return h.Endpoint.IsEqual(*hit.Endpoint)
}

func (h Hit) Description() string {
	var str string
	if h.IsBaseToQuote {
		str = "[" + string(h.Endpoint.From) + "-/+" + string(h.Endpoint.To) + "]@" + h.Endpoint.Exchange.Name
	} else {
		str = "[" + string(h.Endpoint.From) + "+/-" + string(h.Endpoint.To) + "]@" + h.Endpoint.Exchange.Name
	}
	return str
}

func (h Hit) ID() string {
	var str string
	if h.IsBaseToQuote {
		str = string(h.Endpoint.From) + "-" + string(h.Endpoint.To) + "@" + h.Endpoint.Exchange.Name
	} else {
		str = string(h.Endpoint.To) + "-" + string(h.Endpoint.From) + "@" + h.Endpoint.Exchange.Name
	}
	return str
}
