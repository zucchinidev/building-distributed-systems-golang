# Twittervotes
Building Distributed Systems and Workign with Flexible Data.


The system I will build in this repository will prepare us for a future where all democratic elections happen online on Twitter, of course. Our solution will collect and count votes by querying Twitter's streaming API for mentions of specific hash tags, and each component will be capable of horizontally scaling to meet demand.

The ideas discussed here are directly applicable to any system that needs true-scale capabilities.

![Basic overview of the system we are going to build](./docs/distributed-system.png)

# Run app
```sh
$ cp setup_dist.sh setup.sh
```

Add your twitter application keys in the setup.sh file and grant execution privileges

```sh
$ chmod +u setup.sh
```

Start the MongoDB server and nslookup and nsqd daemons.

```sh
$ docker-compose up
```

Navigate to the counter folder and build and run it:

```sh
$ cd counter
$ go build -o counter
$ ./counter
```

Navigate to the twittervotes folder and build and run it. Ensure that you have the appropriate environment variables set; otherwise, you will see errors when you run the program:

```sh
$ cd ../twittervotes
$ go build -o twittervotes
$ ./twittervotes
```

Navigate to the api folder and build and run it:

```sh
cd ../api
go build -o api
./api
```


Navigate to the web folder and build and run it:

```sh
cd ../web
go build -o web
./web
```

Now that everything is running, open a browser and head to http://localhost:8081/.
Using the user interface, create a poll called Moods and input the options as happy,sad,fail,success.
These are common enough words that we are likely to see some relevant activity on Twitter.


Once you have created your poll, you will be taken to the view page where you will start to see the results coming in.
Wait for a few seconds and enjoy the UI updates in real time, showing real-time results:

![Basic overview of the system we are going to build](./docs/poll-example.png)


# TODO
Improve frontend system.
Add Vuejs dashboard
Create SPA.