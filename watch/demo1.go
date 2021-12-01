package watch

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"strconv"
)

func Watch1() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	increment := 0
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("watch event not ok : increment info: ", increment)
					return
				}
				switch event.Op {
				case fsnotify.Create:
					log.Printf("%s add watch event\n", event.Name)
					err = watcher.Add(event.Name)
					if err != nil {
						log.Printf("%s add watch error %v increment %d \n", event.Name, err, increment)
						continue
					}
					increment++
				case fsnotify.Remove:
					log.Printf("%s remove watch event", event.Name)
					err = watcher.Remove(event.Name)
					if err != nil {
						log.Printf("%s remove watch error %v increment %d \n", event.Name, err, increment)
						continue
					}
					increment--
				case fsnotify.Rename:
					err = watcher.Remove(event.Name)
					if err != nil {
						log.Printf("%s remove watch error %v, increate %d \n", event.Name, err, increment)
						continue
					}
					increment--
					// rename： 原文件触发rename,重命名的文件触发create
					log.Printf("%s rename watch event", event.Name)
				}

				log.Println("event: ", event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("increment: %d error : %v ", increment, err)
			}
		}
	}()

	for idx := 0; idx < 100_000_000; idx++ {
		err = watcher.Add("D:\\my\\test_agent\\file\\" + strconv.Itoa(idx))
		if err != nil {
			log.Panicln(err)
		}
	}

	<-done
	log.Println("finished.........")
}
