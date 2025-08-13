package sanitizer

type Sanitizer struct {
	Filters []Filter
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		Filters: []Filter{},
	}
}

func (s *Sanitizer) AddFilter(filter Filter) {
	s.Filters = append(s.Filters, filter)
}

func (s *Sanitizer) CleanString(connectionString string) any {
	for _, filter := range s.Filters {
		connectionString = filter.Apply(connectionString)
	}
	return connectionString
}
