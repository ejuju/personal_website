package app

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPRequest struct {
	ID            string
	CreatedAt     time.Time
	VisitorID     string
	URL           string
	IPAddress     string
	ContentLength int64
	TimeToHandle  time.Duration
}

type Visitor struct {
	ID               string
	FirstVisitedAt   time.Time
	FirstVisitedPage string
}

type Report struct {
	From time.Time
	To   time.Time

	NumVisitors         int
	NumRequests         int
	MostRequestedURLs   map[string]int
	AverageTimeToHandle time.Duration
}

func (r *Report) String() string {
	out := fmt.Sprintf("Analytics report (%s to %s)\n", r.From.Format(time.RFC3339), r.To.Format(time.RFC3339))
	out += "---"
	out += fmt.Sprintf("Number of visitors: %v\n", r.NumVisitors)
	out += fmt.Sprintf("Number of requests: %v\n", r.NumRequests)
	out += fmt.Sprintf("Average time to handle a request: %v\n", r.AverageTimeToHandle)
	out += "Most requested URLs:\n"
	for url, numRequests := range r.MostRequestedURLs {
		out += fmt.Sprintf("\t%v requests for %s\n", numRequests, url)
	}
	return out
}

const visitorIDCookieName = "visitor_id"

func newAnalyticsMiddleware(db DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var visitor *Visitor

			// Check if tracking cookie is present
			cookie, err := r.Cookie(visitorIDCookieName)
			if err != nil {
				// Here, an error means the cookie is not present (= new visitor),
				// create a new visitor and set cookie
				visitor, err = createVisitorAndSetCookie(db, w, r)
				if err != nil {
					log.Println(err)
					respondErrorPage(w, http.StatusInternalServerError, "failed to store visitor")
					return
				}
			} else {
				// Get visitor from cookie
				visitorID := cookie.Value
				visitor, err = db.GetVisitor(visitorID)
				if err != nil {
					log.Println(err)
					respondErrorPage(w, http.StatusInternalServerError, "failed to get visitor")
					return
				}
				if visitor == nil {
					// Here, no visitor was found for the provided cookie,
					// create a new visitor and set cookie
					visitor, err = createVisitorAndSetCookie(db, w, r)
					if err != nil {
						log.Println(err)
						respondErrorPage(w, http.StatusInternalServerError, "failed to store visitor")
						return
					}
				}
			}

			before := time.Now()
			next.ServeHTTP(w, r) // serve request
			after := time.Now()

			// Store request in DB
			req := &HTTPRequest{
				ID:            newID(32),
				CreatedAt:     time.Now(),
				VisitorID:     visitor.ID,
				URL:           r.URL.String(),
				IPAddress:     r.RemoteAddr,
				ContentLength: r.ContentLength,
				TimeToHandle:  after.Sub(before),
			}
			err = db.NewHTTPRequest(req)
			if err != nil {
				log.Println(err)
				return
			}
		})
	}
}

func newAnalyticsCookie(visitorID string) *http.Cookie {
	return &http.Cookie{
		Name:     visitorIDCookieName,
		Value:    visitorID,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
}

func createVisitorAndSetCookie(db DB, w http.ResponseWriter, r *http.Request) (*Visitor, error) {
	visitor := &Visitor{
		ID:               newID(32),
		FirstVisitedAt:   time.Now(),
		FirstVisitedPage: r.URL.String(),
	}
	err := db.NewVisitor(visitor)
	if err != nil {
		return nil, err
	}
	http.SetCookie(w, newAnalyticsCookie(visitor.ID))
	return visitor, nil
}

func sendAnalyticsReport(e Emailer, r *Report) error {
	return e.Send(&Email{
		Sender:        "bot@juliensellier.com",
		Recipient:     "admin@juliensellier.com",
		Subject:       "Analytics report for juliensellier.com",
		PlainTextBody: fmt.Sprintf("New analytics report:\n%s", r.String()),
	})
}

func handleAnalyticsReportNotification(emailer Emailer) {
	for t := range time.Tick(time.Minute) {
		var from time.Time
		switch {
		// daily report every day at midnight
		case t.Hour() == 0 && t.Minute() == 0:
			from = t.Add(-24 * time.Hour)

		// weekly report every monday at 7 AM
		case t.Hour() == 7 && t.Minute() == 0 && t.Weekday() == time.Monday:
			from = t.Add(-24 * time.Hour * 7)

		// monthly report every 1st of the month at 7 AM
		case t.Day() == 0 && t.Hour() == 7 && t.Minute() == 0:
			from = t.Add(-24 * time.Hour * 30)
		}
		err := sendAnalyticsReport(emailer, generateReport(from, t))
		if err != nil {
			log.Println(err)
		}
	}
}

func generateReport(from, to time.Time) *Report {
	return nil
}
