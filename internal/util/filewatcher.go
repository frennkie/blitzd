package util

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func FileWatcher(module, title string, absFilePathIn string, updateFunc func(module, title string, absFilePath string)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := watcher.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					updateFunc(module, title, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(absFilePathIn)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
