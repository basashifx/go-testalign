package noviolation // want package:"testalign source order"

type Service struct{}

func (s *Service) Create() error { return nil }

func (s *Service) Read() error { return nil }

func (s *Service) Delete() error { return nil }
