package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"example.com/myproject/mypack"
)

func checkSite(url string) string {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {

		return fmt.Sprintf("%s: hata %v", url, err)
	}

	resp.Body.Close()
	elapsed := time.Since(start)
	return fmt.Sprintf("%s: %v", url, elapsed)
}

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		sonuc := checkSite(url)
		fmt.Printf("Worker %d -> %s\n", id, sonuc)
		results <- sonuc
	}
}

func main() {
	mypack.SayMe()
	fmt.Println("mypack paketinden sayMe fonksiyonu çağrıldı.")

	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.wikipedia.org",
		"https://www.sahibinden.com",
	}

	jobs := make(chan string, len(urls))
	results := make(chan string, len(urls))

	var wg sync.WaitGroup

	start := time.Now()

	for w := 1; w <= 6; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	wg.Wait()
	close(results)

	elapsed := time.Since(start) 

	fmt.Println("\nTüm işler tamamlandı, sonuçlar:")

	for res := range results {
		fmt.Println(res)
	}

	fmt.Printf("\nToplam geçen süre: %v\n", elapsed)

	// var myArray = []int{1, 2, 3, 4, 5, 10, 15, 20, 25, 30}
	// exercise(myArray)

	go fiberHandler() // Fiber uygulamasını başlat
	go Hi()     // Fiber uygulamasını başlat
	go helper() // net/http paketinden helper fonksiyonunu çağır
	select {}   // Sonsuz döngüde bekle, programın kapanmaması için

}
