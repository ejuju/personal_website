package app

type DB interface {
	// Contact messages
	NewContactFormSubmission(*ContactFormSubmission) error
	// Website analytics
	NewHTTPRequest(*HTTPRequest) error
	NewVisitor(*Visitor) error
	GetVisitor(id string) (*Visitor, error)
}

type inMemoryDB struct {
	contactFormSubmissions map[string]*ContactFormSubmission
	httpRequests           map[string]*HTTPRequest
	visitors               map[string]*Visitor
}

func newInMemoryDB() *inMemoryDB {
	return &inMemoryDB{
		contactFormSubmissions: map[string]*ContactFormSubmission{},
		httpRequests:           map[string]*HTTPRequest{},
		visitors:               map[string]*Visitor{},
	}
}

func (db *inMemoryDB) NewContactFormSubmission(s *ContactFormSubmission) error {
	db.contactFormSubmissions[s.ID] = s
	return nil
}

func (db *inMemoryDB) NewHTTPRequest(pv *HTTPRequest) error   { db.httpRequests[pv.ID] = pv; return nil }
func (db *inMemoryDB) NewVisitor(v *Visitor) error            { db.visitors[v.ID] = v; return nil }
func (db *inMemoryDB) GetVisitor(id string) (*Visitor, error) { v, _ := db.visitors[id]; return v, nil }
