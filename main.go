package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Kozical/samuel/app"
)

func main() {
	path := flag.String("config", "config.yaml", "Specify the path to the configuration yaml file. (ie: C:\\my\\path\\config.yaml) (default: config.yaml)")
	flag.Parse()

	config, err := app.Init(*path)

	if err != nil {
		panic(err)
	}

	client, err := app.New(config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("[%d]: Run() called\n", i)
			if err := client.Run(); err != nil {
				fmt.Printf("Error[%d]: %s\n", i, err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("work complete..")
}
