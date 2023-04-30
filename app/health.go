package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type report struct {
	From time.Time
	To   time.Time

	// Website traffic
	NumVisitors         int
	NumRequests         int
	MostRequestedURLs   map[string]int
	AverageTimeToHandle time.Duration
}

func (r *report) String() string {
	out := fmt.Sprintf("# Health report (%s to %s)\n\n", r.From.Format(time.RFC3339), r.To.Format(time.RFC3339))
	out += "## Traffic\n\n"
	out += fmt.Sprintf("%-25s %s\n", "Number of visitors:", strconv.Itoa(r.NumVisitors))
	out += fmt.Sprintf("%-25s %s\n", "Number of requests:", strconv.Itoa(r.NumRequests))
	out += fmt.Sprintf("%-25s %s\n", "Avg. time to handle:", r.AverageTimeToHandle)
	out += "Most requested URLs:\n"
	for url, numRequests := range r.MostRequestedURLs {
		out += fmt.Sprintf("\t* %-10s %q\n", strconv.Itoa(numRequests), url)
	}
	return out
}

type httpRequest struct {
	ID            string
	CreatedAt     time.Time
	VisitorID     string
	URL           string
	UserAgent     string
	IPAddress     string
	ContentLength int64
	TimeToHandle  time.Duration
}

type visitor struct {
	ID               string
	FirstVisitedAt   time.Time
	FirstVisitedPage string
}

func generateReport(db DB, from, to time.Time) (*report, error) {
	numVisitors, err := db.CountVisitors(from, to)
	if err != nil {
		return nil, err
	}
	numHTTPRequests, err := db.CountHTTPRequests(from, to)
	if err != nil {
		return nil, err
	}
	averageTimeToHandle, err := db.GetAverageTimeToHandleHTTPRequest(from, to)
	if err != nil {
		return nil, err
	}
	mostRequestedURLs, err := db.GetMostRequestedURLs(from, to)
	if err != nil {
		return nil, err
	}
	return &report{
		From:                from,
		To:                  to,
		NumVisitors:         numVisitors,
		NumRequests:         numHTTPRequests,
		AverageTimeToHandle: averageTimeToHandle,
		MostRequestedURLs:   mostRequestedURLs,
	}, nil
}

const visitorIDCookieName = "visitor_id"

func newRequestTrackingMiddleware(db DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			before := time.Now()

			// Get visitor from cookie,
			// If not visitor found, create a new visitor and set cookie
			visitor, err := getVisitorFromRequest(db, r)
			if err != nil && !errors.Is(err, errVisitorNotFound) {
				log.Println(err)
				respondErrorPage(w, http.StatusInternalServerError, "failed to store visitor")
				return
			}
			if errors.Is(err, errVisitorNotFound) {
				visitor, err = createVisitorAndSetCookie(db, w, r)
				if err != nil {
					log.Println(err)
					respondErrorPage(w, http.StatusInternalServerError, "failed to store visitor")
					return
				}
			}

			next.ServeHTTP(w, r) // serve request
			after := time.Now()

			// Store request in DB
			req := &httpRequest{
				ID:            newID(32),
				CreatedAt:     time.Now(),
				VisitorID:     visitor.ID,
				URL:           r.URL.String(),
				IPAddress:     r.RemoteAddr,
				ContentLength: r.ContentLength,
				TimeToHandle:  after.Sub(before),
				UserAgent:     r.UserAgent(),
			}
			err = db.StoreHTTPRequest(req)
			if err != nil {
				log.Println(err)
				return
			}
		})
	}
}

func createVisitorAndSetCookie(db DB, w http.ResponseWriter, r *http.Request) (*visitor, error) {
	visitor := &visitor{
		ID:               newID(32),
		FirstVisitedAt:   time.Now(),
		FirstVisitedPage: r.URL.String(),
	}
	err := db.StoreVisitor(visitor)
	if err != nil {
		return nil, err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     visitorIDCookieName,
		Value:    visitor.ID,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})
	return visitor, nil
}

func getVisitorFromRequest(db DB, r *http.Request) (*visitor, error) {
	// Check if tracking cookie is present
	// Here, an error means the cookie is not present (= new visitor),
	cookie, err := r.Cookie(visitorIDCookieName)
	if err != nil {
		return nil, errVisitorNotFound // no visitor found
	}

	// Get visitor from cookie
	return db.GetVisitor(cookie.Value)
}

func doPeriodicHealthReport(config *Config, emailer Emailer, db DB) {
	// Generate and send report on startup
	report, err := generateReport(db, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		log.Println(err)
	}
	err = sendEmailToAdmin(config, emailer, "Health report for juliensellier.com", report.String())
	if err != nil {
		log.Println(err)
	}

	// Every minute, check if a cron job needs to be executed
	for t := range time.Tick(time.Minute) {
		// Get report "from" timestamp
		var from time.Time
		switch {
		default:
			continue
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
		// Generate and send report
		report, err := generateReport(db, from, t)
		if err != nil {
			log.Println(err)
			continue
		}
		err = sendEmailToAdmin(config, emailer, "Health report for juliensellier.com", report.String())
		if err != nil {
			log.Println(err)
		}
	}
}
