package app

import (
	"crypto/sha1"
	"encoding/base32"
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
	NumRequestPerURL    map[string]int
	AverageTimeToHandle time.Duration
}

func (r *report) String() string {
	out := fmt.Sprintf("# Health report (%s to %s)\n\n", r.From.Format(time.RFC3339), r.To.Format(time.RFC3339))
	out += "## Traffic\n\n"
	out += fmt.Sprintf("%-25s %s\n", "Number of visitors:", strconv.Itoa(r.NumVisitors))
	out += fmt.Sprintf("%-25s %s\n", "Number of requests:", strconv.Itoa(r.NumRequests))
	out += fmt.Sprintf("%-25s %s\n", "Avg. time to handle:", r.AverageTimeToHandle)
	out += "Most requested URLs:\n"
	for url, numRequests := range r.NumRequestPerURL {
		out += fmt.Sprintf("\t* %-10s %q\n", strconv.Itoa(numRequests), url)
	}
	return out
}

type httpRequest struct {
	ID            string
	CreatedAt     time.Time
	URL           string
	VisitorHash   string
	UserAgent     string
	ContentLength int64
	TimeToHandle  time.Duration
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
	numRequestPerURL, err := db.GetNumRequestPerURL(from, to)
	if err != nil {
		return nil, err
	}
	return &report{
		From:                from,
		To:                  to,
		NumVisitors:         numVisitors,
		NumRequests:         numHTTPRequests,
		AverageTimeToHandle: averageTimeToHandle,
		NumRequestPerURL:    numRequestPerURL,
	}, nil
}

func newRequestTrackingMiddleware(db DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			before := time.Now()
			next.ServeHTTP(w, r) // serve request
			after := time.Now()

			visitorHash, err := newVisitorHash(r)
			if err != nil {
				panic(err)
			}

			// Store request in DB
			req := &httpRequest{
				ID:            newID(32),
				CreatedAt:     time.Now(),
				VisitorHash:   visitorHash,
				URL:           r.URL.String(),
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

func getIPAddr(r *http.Request) (string, error) {
	// Check for reverse proxy header
	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}
	return remoteAddr, nil
}

func newVisitorHash(r *http.Request) (string, error) {
	// Hash IP addr and user-agent
	hash := sha1.New()
	ip, err := getIPAddr(r)
	if err != nil {
		return "", err
	}
	_, err = hash.Write([]byte(ip + r.UserAgent()))
	if err != nil {
		panic(err)
	}
	// Return base32 hex encoded hash
	return base32.HexEncoding.EncodeToString(hash.Sum(nil)), nil
}

func doPeriodicHealthReport(config *Config, emailer Emailer, db DB) {
	// Generate and send report on startup
	report, err := generateReport(db, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		log.Println(err)
	}
	err = sendEmailToAdmin(config, emailer, "New startup report for juliensellier.com", report.String())
	if err != nil {
		log.Println(err)
	}

	// Every minute, check if a cron job needs to be executed
	for t := range time.Tick(time.Minute) {
		// Get report "from" timestamp
		var from time.Time
		subjectPrefix := ""
		switch {
		default:
			continue
		case t.Hour() == 7 && t.Minute() == 0 && t.Weekday() == time.Monday:
			// weekly report every monday at 7 AM
			subjectPrefix = "Last week"
			from = t.Add(-24 * time.Hour * 7)
		case t.Day() == 0 && t.Hour() == 7 && t.Minute() == 0:
			// monthly report every 1st of the month at 7 AM
			subjectPrefix = "Last month"
			from = t.Add(-24 * time.Hour * 30)
		}
		// Generate and send report
		report, err := generateReport(db, from, t)
		if err != nil {
			log.Println(err)
			continue
		}
		err = sendEmailToAdmin(config, emailer, subjectPrefix+" on juliensellier.com", report.String())
		if err != nil {
			log.Println(err)
		}
	}
}
