package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const updateDuration = 1 * time.Second

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {
	var counts map[string]int
	var countsLock sync.Mutex
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("Connecting to database...")
	db, err := mgo.Dial("mongodb")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("Closing database connection...")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	log.Println("Connecting to nsq...")
	consumer, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fmt.Println("error when we are creating a new consumer")
		fatal(err)
		return
	}

	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()

		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := consumer.ConnectToNSQLookupd("nsqlookupd:4161"); err != nil {
		fmt.Println("error when we connect to nsql lookupd")
		fatal(err)
		return
	}

	ticker := time.NewTicker(updateDuration)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-ticker.C:
			doCount(&countsLock, &counts, pollData)
		case <-termChan:
			consumer.Stop()
		case <-consumer.StopChan:
			// finished
			return
		}
	}
}

func doCount(countsLock *sync.Mutex, counts *map[string]int, pollData *mgo.Collection) {
	countsLock.Lock()
	defer countsLock.Unlock()
	if len(*counts) == 0 {
		log.Println("No new votes, skipping database update")
		return
	}

	ok := true

	for option, count := range *counts {
		query := bson.M{"options": bson.M{"$in": []string{option}}}
		up := bson.M{"$inc": bson.M{"results." + option: count}}
		if _, err := pollData.UpdateAll(query, up); err != nil {
			log.Println("failded to update:", err)
			ok = false
		}
	}

	if ok {
		*counts = nil
	}
}
