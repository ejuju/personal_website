package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

type DB interface {
	StoreContactFormSubmission(*ContactFormSubmission) error

	StoreHTTPRequest(*httpRequest) error
	CountHTTPRequests(from, to time.Time) (int, error)
	GetAverageTimeToHandleHTTPRequest(from, to time.Time) (time.Duration, error)
	GetMostRequestedURLs(from, to time.Time) (map[string]int, error)

	StoreVisitor(*visitor) error
	GetVisitor(id string) (*visitor, error)
	CountVisitors(from, to time.Time) (int, error)
}

var (
	errVisitorNotFound = errors.New("visitor not found")
)

type boltDB struct {
	f *bbolt.DB
}

var (
	boltContactFormBucket  = []byte("contact_form_submissions")
	boltHTTPRequestsBucket = []byte("http_requests")
	boltVisitorBucket      = []byte("visitors")
)

func newBoltDB() *boltDB {
	// Open file
	db, err := bbolt.Open(".tmp/main.boltdb", 0666, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}

	// Ensure buckets are created
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucketID := range [][]byte{
			boltContactFormBucket,
			boltHTTPRequestsBucket,
			boltVisitorBucket,
		} {
			_, err := tx.CreateBucketIfNotExists(bucketID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &boltDB{f: db}
}

func (db *boltDB) close() error { return db.f.Close() }

func (db *boltDB) StoreContactFormSubmission(s *ContactFormSubmission) error {
	return db.f.Update(func(tx *bbolt.Tx) error {
		key := []byte(s.CreatedAt.Format(time.RFC3339) + s.ID)
		return tx.Bucket(boltContactFormBucket).Put(key, mustMarshalJSON(s))
	})
}

func (db *boltDB) StoreHTTPRequest(req *httpRequest) error {
	return db.f.Update(func(tx *bbolt.Tx) error {
		key := []byte(req.CreatedAt.Format(time.RFC3339) + req.ID)
		return tx.Bucket(boltHTTPRequestsBucket).Put(key, mustMarshalJSON(req))
	})
}

func (db *boltDB) CountHTTPRequests(from, to time.Time) (int, error) {
	count := 0
	return count, db.readTimeRange(boltHTTPRequestsBucket, from, to, func(k, v []byte) error {
		count++
		return nil
	})
}

func (db *boltDB) GetAverageTimeToHandleHTTPRequest(from, to time.Time) (time.Duration, error) {
	out := time.Duration(0)
	count := 0
	err := db.readTimeRange(boltHTTPRequestsBucket, from, to, func(k, v []byte) error {
		req := &httpRequest{}
		mustUnmarshalJSON(v, req)
		out += req.TimeToHandle
		count++
		return nil
	})
	if count != 0 {
		return out / time.Duration(count), err
	}
	return 0, err
}

func (db *boltDB) GetMostRequestedURLs(from, to time.Time) (map[string]int, error) {
	out := map[string]int{}
	return out, db.readTimeRange(boltHTTPRequestsBucket, from, to, func(k, v []byte) error {
		req := &httpRequest{}
		mustUnmarshalJSON(v, req)
		out[req.URL]++
		return nil
	})
}

func (db *boltDB) StoreVisitor(v *visitor) error {
	return db.f.Update(func(tx *bbolt.Tx) error {
		key := []byte(v.ID)
		return tx.Bucket(boltVisitorBucket).Put(key, mustMarshalJSON(v))
	})
}

func (db *boltDB) GetVisitor(id string) (*visitor, error) {
	v := &visitor{}
	return v, db.f.View(func(tx *bbolt.Tx) error {
		if raw := tx.Bucket(boltVisitorBucket).Get([]byte(id)); raw != nil {
			mustUnmarshalJSON(raw, v)
		} else {
			return errVisitorNotFound
		}
		return nil
	})
}

func (db *boltDB) CountVisitors(from, to time.Time) (int, error) {
	visitorIDs := map[string]struct{}{}
	err := db.readTimeRange(boltHTTPRequestsBucket, from, to, func(k, v []byte) error {
		req := &httpRequest{}
		mustUnmarshalJSON(v, req)
		visitorIDs[req.VisitorID] = struct{}{}
		return nil
	})
	return len(visitorIDs), err
}

func (db *boltDB) readTimeRange(bucket []byte, from, to time.Time, cb func(k, v []byte) error) error {
	return db.f.View(func(tx *bbolt.Tx) error {
		min := []byte(from.Format(time.RFC3339))
		max := []byte(to.Format(time.RFC3339))
		c := tx.Bucket(bucket).Cursor()
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			if err := cb(k, v); err != nil {
				return err
			}
		}
		return nil
	})
}

// periodic boltdb database file backup to an another file on disk.
func (db *boltDB) doPeriodicDBFileBackup(config *Config, emailer Emailer) {
	if err := db.backupDB(db.newBackupFname(time.Now())); err != nil {
		if sendEmailToAdmin(config, emailer, "DB backup failed", err.Error()) != nil {
			log.Println(err)
		}
	}
	for tick := range time.Tick(24 * time.Hour) {
		if err := db.backupDB(db.newBackupFname(tick)); err != nil {
			if err := sendEmailToAdmin(config, emailer, "DB backup failed", err.Error()); err != nil {
				log.Println(err)
			}
			continue
		}
	}
}

func (db *boltDB) newBackupFname(t time.Time) string {
	return t.Format("2006_01_02_15_04_05") + ".backup.boltdb"
}

func (db *boltDB) backupDB(fileName string) error {
	f, err := os.Create(".tmp/" + fileName)
	if err != nil {
		return err
	}
	return db.f.View(func(tx *bbolt.Tx) error { _, err := tx.WriteTo(f); return err })
}

func mustMarshalJSON(v any) []byte {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return out
}

func mustUnmarshalJSON(raw []byte, into any) {
	err := json.Unmarshal(raw, into)
	if err != nil {
		panic(err)
	}
}
