//read command line arguments with flag
//read file
	//parse by comma
	//print first part as question
	//wait for input
	//compare user entry with 2nd part
	//keep score
	//keep iterating till all is done

package main  

import ( 
	"fmt"
	"flag"
	"bufio"
	"os"
	"log"
	"strings"
	"encoding/csv"
	"io"
	"time"
)
//fmethod that compares the user input answer to the answer in the csv
func CompareAnswers(ip string, correct *int, done chan <-bool) {
	
	ansscan := bufio.NewScanner(os.Stdin)
	for ansscan.Scan() {
		if( strings.Compare(ansscan.Text(),ip) == 0 ) {
			fmt.Printf("answer was %v \n", ansscan.Text())
			*correct++
		}
		done <- true
		break
	}
}

func main () {
	filename := flag.String("fname", "quiz.csv", "string")
	timer := flag.Int("timer", 30, "timer value")
	flag.Parse()
	//fmt.Printf("name of file is %v \n", *filename)

	f, err := os.Open(*filename)
	if( err != nil ) {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	//reading each line of file and parsing with comma based separator
	/*
	contents := bufio.NewScanner(f)

	ctr:= 0
	correct := 0
	for(contents.Scan()) {
		ctr++
		line := contents.Text()
		tokens := strings.Split(line, ",")
		ansscan := bufio.NewScanner(os.Stdin)
		for i:= range tokens {
			fmt.Printf("Question [%v] -> %v \n", ctr, tokens[i] )
			fmt.Printf("Answer [%v] -> ", ctr)
			for ansscan.Scan() {
				if( strings.Compare(ansscan.Text(),tokens[i+1]) == 0 ) {
					fmt.Printf("answer was %v \n", ansscan.Text())
					correct++
				}
				
				break
			}
			break
		}
		err = contents.Err()
		if( err != nil ) {
			log.Fatal(err)
		}
	}
	*/

	//using CSV parser
	t := time.NewTimer(time.Duration(*timer)*time.Second)
	
	ctr:= 0
	correct := 0
	r := csv.NewReader(bufio.NewReader(f))
	loop:
	for {
		line, err := r.Read()
		if( err == io. EOF ) {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		ctr++
	
		fmt.Printf("Question [%v] -> %v \n", ctr, line[0] )
		fmt.Printf("Answer [%v] -> ", ctr)
		done := make(chan bool)
		
		go CompareAnswers(line[1], &correct, done)
		
		select {
		case <-done:
		case <-t.C:
			fmt.Println("Timer expired")
			break loop
		}

	}
	t.Stop()
	fmt.Printf("Total Questions %v \n", ctr)
	fmt.Printf("Number of correct answers by user is %v \n", correct)
	fmt.Printf("Number of wrong answers by user is %v \n", ctr-correct)
	
}

