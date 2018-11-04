package main

import (
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/joeshaw/envdecode"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var db *mgo.Session

var config struct {
	MongoService        string `env:"BDSG_MONGO_SERVICE,required"`
	MongoDBName         string `env:"BDSG_MONGO_DB_NAME,required"`
	MongoCollectionName string `env:"BDSG_MONGO_COLLECTION_NAME,required"`
	NSQDService         string `env:"BDSG_NSQD_SERVICE,required"`
	TwitterApiUrl       string `env:"BDSG_TWITTER_API_URL,required"`
}

var twitterConfiguration struct {
	ConsumerKey    string `env:"BDSG_TWITTER_KEY,required"`
	ConsumerSecret string `env:"BDSG_TWITTER_SECRET,required"`
	AccessToken    string `env:"BDSG_TWITTER_ACCESSTOKEN,required"`
	AccessSecret   string `env:"BDSG_TWITTER_ACCESSSECRET,required"`
}

func init() {
	if err := envdecode.Decode(&config); err != nil {
		log.Fatal(err)
	}

	if err := envdecode.Decode(&twitterConfiguration); err != nil {
		log.Fatal(err)
	}
}

type poll struct {
	Options []string
}

func main() {
	var stopLock sync.Mutex // protects stop
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stopLock.Lock()
		stop = true
		stopLock.Unlock()
		log.Println("Stopping....")
		stopChan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	if err := dialdb(); err != nil {
		log.Fatalln("failed to dial MongoDB:", err)
	}
	defer closedb()

	votes := make(chan string) // chan for votes
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stopLock.Lock()
			if stop {
				stopLock.Unlock()
				return
			}
			stopLock.Unlock()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	<-publisherStoppedChan
}

func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial(config.MongoService)
	return err
}

func closedb() {
	db.Close()
	log.Println("closed database connection")
}

func loadOptions() ([]string, error) {
	var options []string
	iterator := db.DB(config.MongoDBName).C(config.MongoCollectionName).Find(nil).Iter()
	var p poll
	for iterator.Next(&p) {
		options = append(options, p.Options...)
	}
	iterator.Close()
	return options, iterator.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopChan := make(chan struct{}, 1)
	pub, err := nsq.NewProducer(config.NSQDService, nsq.NewConfig())
	if err != nil {
		fmt.Println("Error when create a new producer: ", err)
		panic("Error when create a producer")
	}
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // publish vote
		}
		log.Println("Publisher: stopping")
		pub.Stop()
		log.Println("Publisher: stopped")
		stopChan <- struct{}{}
	}()
	return stopChan
}
