package app

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

type DB interface {
	StoreHTTPRequest(*httpRequest) error
	CountHTTPRequests(from, to time.Time) (int, error)
	GetAverageTimeToHandleHTTPRequest(from, to time.Time) (time.Duration, error)
	GetNumRequestPerURL(from, to time.Time) (map[string]int, error)
	CountVisitors(from, to time.Time) (int, error)
}

type boltDB struct {
	dbDirPath string // with trailing slash
	f         *bbolt.DB
}

var (
	boltHTTPRequestsBucket = []byte("http_requests")
)

func newBoltDB() *boltDB {
	// Open file
	dbDirPath := os.Getenv("DB_DIR_PATH")
	if dbDirPath == "" {
		dbDirPath = ".tmp/"
	}
	db, err := bbolt.Open(dbDirPath+"main.boltdb", 0666, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}

	// Ensure buckets are created
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucketID := range [][]byte{
			boltHTTPRequestsBucket,
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

	return &boltDB{f: db, dbDirPath: dbDirPath}
}

func (db *boltDB) close() error { return db.f.Close() }

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

func (db *boltDB) GetNumRequestPerURL(from, to time.Time) (map[string]int, error) {
	out := map[string]int{}
	return out, db.readTimeRange(boltHTTPRequestsBucket, from, to, func(k, v []byte) error {
		req := &httpRequest{}
		mustUnmarshalJSON(v, req)
		out[req.URL]++
		return nil
	})
}

func (db *boltDB) CountVisitors(from, to time.Time) (int, error) {
	visitorHashes := map[string]struct{}{}
	err := db.readTimeRange(boltHTTPRequestsBucket, from, to, func(k, v []byte) error {
		req := &httpRequest{}
		mustUnmarshalJSON(v, req)
		visitorHashes[req.VisitorHash] = struct{}{}
		return nil
	})
	return len(visitorHashes), err
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
	ticker := time.NewTicker(24 * time.Hour)
	for tick := range ticker.C {
		if err := db.backupDB(db.newBackupFname(tick)); err != nil {
			if err := sendEmailToAdmin(config, emailer, "DB backup failed", err.Error()); err != nil {
				log.Println(err)
			}
			continue
		}
	}
}

func (db *boltDB) newBackupFname(t time.Time) string {
	// 2006: Year, 01: Month, 02: Day
	// 15: Hour, 04: Minute, 05: second
	return t.Format("20060102_150405") + "_back.boltdb"
}

func (db *boltDB) backupDB(fileName string) error {
	f, err := os.Create(db.dbDirPath + fileName)
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
