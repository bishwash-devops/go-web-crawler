package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

/* Interfaces */
/* A Type implements an interface by defining the required methods
Example : fmt has Stringer Interface
		: io Writer Interface

*/

/* Reflection */
/* Type information, basic operations available at run-time. A little goes a long way.
Print and Printf are not built into the language,
Functionality that printf needs is packaged into package called reflect
Example :

*/

/*
	Parallelism: Running multiple things simultaneously
	Concurrency: Deal with multiple things simultaneously, co-ordination of parallel computations

	Goroutines let you run multiple computations simultaneously
	Channels let you coordinate the computations, by explicit communication

*/

func main() {
	// Stringer Interface
	fmt.Printf("Testing This %s!", new(World))
	fmt.Println()

	Pointer()

	officePlace[0] = "Cambridge"
	officePlace[1] = "Queens"

	fmt.Printf("Hello, %s\n", NewYork)

	day := time.Now().Weekday()
	fmt.Printf("Hellow, %s (%d)\n", day, day)

	//Writer Interface, to any IO Writer
	fmt.Fprintf(os.Stdout, "Hello!\n")

	//MultiWriter
	h := crc32.NewIEEE()
	w := io.MultiWriter(h, os.Stdout)
	fmt.Fprintf(w, "Hello, Multiwriter\n")
	fmt.Printf("hash=%#x\n", h.Sum32())

	//Struct
	lang := Lang{"Go", 2009, "http://golang.org/"}
	fmt.Printf("%v\n", lang)
	fmt.Printf("%+v\n", lang)
	fmt.Printf("%#v\n", lang)
	fmt.Println(lang.Name)

	// Struct to XML
	data, err := xml.MarshalIndent(lang, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

	// Struct to Json
	data1, err1 := json.Marshal(lang)
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Printf("%s\n", data1)

	// JSON to Struct
	input, err := os.Open("./lang.json")
	if err != nil {
		log.Fatal(err)
	}
	//Copy input file to StdOut - memcopy. Similar to cat
	//io.Copy(os.Stdout, input)

	dec := json.NewDecoder(input)
	for {
		var lang Lang
		err := dec.Decode(&lang)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("%v\n", lang)
	}

	// JSON to Struct abstracted
	do(func(lang Lang) {
		fmt.Printf("%v\n", lang)
	})

	// JSON to XML abstracted
	do(func(lang Lang) {
		data, err := xml.MarshalIndent(lang, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
	})

	//// Toy Web Crawler
	//start := time.Now()
	//do(func(lang Lang){
	//	count(lang.Name, lang.URL)
	//})
	//fmt.Printf("%.2fs total\n", time.Since(start).Seconds())

	//// Toy Web Crawler - In Parallel
	//start := time.Now()
	//do(func(lang Lang){
	//	go count(lang.Name, lang.URL)
	//})
	//time.Sleep(10 * time.Second )
	//fmt.Printf("%.2fs total\n", time.Since(start).Seconds())

	// Toy Web Crawler - Concurrently
	start := time.Now()
	c := make(chan string)
	n := 0
	do(func(lang Lang) {
		n++
		go count_concurrent(lang.Name, lang.URL, c)
	})

	timeout := time.After(1 * time.Second)
	for i := 0; i < n; i++ {
		select {
		case result := <-c:
			fmt.Print(result)
		case <-timeout:
			fmt.Print("Timed out\n")
			return
		}
		fmt.Print(<-c)
	}
	fmt.Printf("%.2fs total\n", time.Since(start).Seconds())

}

type World struct{}

func (w *World) String() string {
	return "MyWorld"
}

/* Pointers */
func Pointer() {
	x := 15
	a := &x // memory address
	fmt.Println(a)
	fmt.Println(*a) // *value

	*a = 5

	fmt.Println(x)

}

/* Printing Stringers */
type Office int

var officePlace [2]string

const (
	Boston Office = iota
	NewYork
)

func (o Office) String() string {
	return "Google, " + officePlace[o]
}

type Lang struct {
	Name string
	Year int
	URL  string
}

func do(f func(Lang)) {
	input, err := os.Open("./lang.json")
	if err != nil {
		log.Fatal(err)
	}
	//Copy input file to StdOut - memcopy. Similar to cat
	//io.Copy(os.Stdout, input)

	dec := json.NewDecoder(input)
	for {
		var lang Lang
		err := dec.Decode(&lang)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		f(lang)
	}
}

func count(name, url string) {
	start := time.Now()
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s: %s", name, err)
		return
	}
	n, _ := io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	fmt.Printf("%s %d [%.2fs]\n", name, n, time.Since(start).Seconds())
}

func count_concurrent(name, url string, c chan<- string) {
	start := time.Now()
	r, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprintf("%s: %s", name, err)
		return
	}
	n, _ := io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	c <- fmt.Sprintf("%s %d [%.2fs]\n", name, n, time.Since(start).Seconds())
}
