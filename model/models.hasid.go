package model

func (u *User) GetID() string {
	return u.ID.String()
}

func (p *Project) GetID() string {
	return p.ID.String()
}

func (s *Server) GetID() string {
	return s.ID.String()
}

func (s *Session) GetID() string {
	return s.ID
}

func (e *Environment) GetID() string {
	return e.ID.String()
}
