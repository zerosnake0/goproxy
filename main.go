package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	port        int
	directRepos []string
	goProxy     = strings.TrimRight(os.Getenv("GOPROXY"), "/")
)

func init() {
	if goProxy == "" {
		goProxy = "http:/"
	}
	fmt.Println("GOPROXY", goProxy)
}

func getRedirect(repo string) string {
	for _, s := range directRepos {
		if strings.HasPrefix(repo, "/"+s+"/") {
			return "http:/" + repo
		}
	}
	return goProxy + repo
}

func parse() error {
	var direct string
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.StringVar(&direct, "direct", "", "direct site file")
	flag.Parse()

	if direct != "" {
		fp, err := os.Open(direct)
		if err != nil {
			return err
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			directRepos = append(directRepos, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := parse(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		repo := request.URL.String()
		redirect := getRedirect(repo)
		fmt.Println(repo, "=>", redirect)
		http.Redirect(writer, request, redirect, http.StatusMovedPermanently)
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
