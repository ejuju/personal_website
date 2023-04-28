package app

type DB interface {
	NewContactFormSubmission(*ContactFormSubmission) error
}

type inMemoryDB struct {
	contactFormSubmissions map[string]*ContactFormSubmission
}

func (db *inMemoryDB) NewContactFormSubmission(s *ContactFormSubmission) error {
	db.contactFormSubmissions[s.ID] = s
	return nil
}
