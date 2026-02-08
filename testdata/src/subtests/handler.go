package subtests // want package:"testalign source order"

type Handler struct{}

func (h *Handler) Get() error { return nil }

func (h *Handler) Post() error { return nil }

func (h *Handler) Delete() error { return nil }
