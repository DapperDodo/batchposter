# batchposter

Buffer POST requests in batches. Flush when buffer is full or when N seconds have elapsed. Threadsafe.
The server/endpoint must support multiline POST requests.

### use case

Google Analytics, but useful for any endpoint that supports multiline POST

### installation

  go get -u github.com/DapperDodo/batchposter

### example

    func main() {
    
    	errorlogger := log.New(os.Stderr, "[GA BATCH ERROR] ", log.Llongfile)
    	
    	// "http://www.google-analytics.com/batch" 	: the url to POST to
    	// 20										: the maximum batch size (`flush immediately when 20 POSTS are buffered`)
    	// time.Second*5 							: the maximum time between flushes (`at least flush every N seconds`)
    	// errorlogger 								: logger where errors are reported to
  
    	ga := batchposter.New("http://www.google-analytics.com/batch", 20, time.Second*5, errorlogger)
    
    	for {
    		uuid := strconv.Itoa(1000000 + rand.Intn(8999999))
    
    		payload := "v=1&tid=UA-XXXXXXX-1&cid=" + uuid + "&t=event&ec=accounts&ea=create&aip=1&ds=api&uid=" + uuid
    		fmt.Println(payload)
    		err := ga.Post(payload)
    		if err != nil {
    			fmt.Println(err)
    		}
    		time.Sleep(time.Second)
    
    	}
    }
