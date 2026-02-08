package multifile // want package:"testalign source order"

type Order struct{}

func (o *Order) Place() error { return nil }

func (o *Order) Cancel() error { return nil }
