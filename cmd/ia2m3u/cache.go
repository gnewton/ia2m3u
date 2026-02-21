package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io"
	"io/ioutil"
	"log"
)

var db *bolt.DB

var DBBucketName = "ia"

func initializeCache(dbFileName string) error {
	if db != nil {
		log.Fatal("DB is not nil; only run initializeCache() once!")
	}

	var err error

	db, err = bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DBBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}

func getKey(url string) (io.Reader, bool) {
	var v []byte

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DBBucketName))
		v = b.Get([]byte(url))
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if v != nil {
		log.Println("**************************** Cache hit", url)
		var buf2 bytes.Buffer
		err := gunzipper(&buf2, v)
		if err != nil {
			log.Fatal(err)
		}
		//return bytes.NewReader(buf2), true
		return &buf2, true
	}
	return nil, false
}

func addToCache(url string, body []byte) error {
	log.Println("    Cache add:", url)
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DBBucketName))

		var gzbuf bytes.Buffer
		gzipper(&gzbuf, []byte(body))

		return b.Put([]byte(url), gzbuf.Bytes())
	})
}

func gzipper(w io.Writer, data []byte) {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	_, err := gw.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func gunzipper(w io.Writer, data []byte) error {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()

	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}
