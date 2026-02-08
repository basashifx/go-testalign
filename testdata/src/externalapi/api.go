package externalapi // want package:"testalign source order"

type API struct{}

func (a *API) Get() error { return nil }

func (a *API) Post() error { return nil }

func (a *API) Delete() error { return nil }
