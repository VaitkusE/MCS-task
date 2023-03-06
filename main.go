// Edgaras Vaitkus task for KNA MCS software dev position.
// Program accepts multiple line inputs for one interval - user can press enter as many times as needed.
// If user types something down during an interval and does not press enter before its end, the input will be sent via channel in next interval.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func readInput(input chan<- string) { //input reader function, with a channel sending input to main routine
	for {
		var u string
		_, err := fmt.Scan(&u)
		if err != nil {
			panic(err)
		}
		input <- u
	}
}
func findMostFrequentByte(slice []byte) (byte, int) { // itterate through slice of bytes, put keys/values into map, which serves as duplicate remover and duplicate counter.
	m := map[byte]int{}
	var highestCount int
	var mostPopularChar byte

	for _, a := range slice { //pereinam visa lista, mape susiskaiciuoja kiekvieno simbolio pasikartojimas
		m[a]++
		if m[a] > highestCount {
			highestCount = m[a]
			mostPopularChar = a
		}
	}
	return mostPopularChar, highestCount
}
func removeByte(slice []byte, mostFrequentByte byte) []byte { // remove highest byte, to make finding top 3 characters more efficient and easy to read
	newList := []byte{}
	for _, a := range slice {
		if a != mostFrequentByte {
			newList = append(newList, a)
		}
	}
	return newList
}
func printIntervalStatistics(currentIntervalData []byte, sessionData []byte, sessionTimeInS float32) { //display interval data, specified by task description

	fmt.Println("\n=============================================")
	fmt.Printf("Typing speed of the session is: %.3f characters/sec\n", float32(len(sessionData))/float32(sessionTimeInS))
	fmt.Println("Interval statistics:")
	fmt.Println("\nCharacter count: ", len(currentIntervalData))
	fmt.Println("\nTop 3 most frequent characters ranked in descending order: ")
	for i := 0; i < 3; i++ {
		var mostFreqChar, count = findMostFrequentByte(currentIntervalData)
		fmt.Println("Top ", i+1, ":", string([]byte{mostFreqChar}), "| counted :", count, " times")
		currentIntervalData = removeByte(currentIntervalData, mostFreqChar)
	}
	fmt.Println("=============================================")
}
func printSessionStatistics(intervalData []byte, sessionTimeInS float32) { //display session data, specified by task description
	fmt.Println("\n=============================================")
	fmt.Println("Session statistics:")
	fmt.Println("\nCharacter count: ", len(intervalData))
	fmt.Printf("\nTyping speed of the session is: %.3f characters/sec\n", float32(len(intervalData))/float32(sessionTimeInS))
	fmt.Println("\nTop 3 most frequent characters ranked in descending order: ")
	for i := 0; i < 3; i++ {
		var mostFreqChar, count = findMostFrequentByte(intervalData)
		fmt.Println("Top ", i+1, ":", string([]byte{mostFreqChar}), "| counted :", count, " times")
		intervalData = removeByte(intervalData, mostFreqChar)
	}
	fmt.Println("=============================================")
}

func main() {
	userInput := make(chan string)
	sessionData := []string{}
	var iTime, sTime time.Duration

	fmt.Println("Input time of interval in seconds:")
	_, err := fmt.Scanln(&iTime)
	if err != nil {
		panic(err)
	}
	fmt.Println("Input time of session in seconds:")
	_, err = fmt.Scanln(&sTime)
	if err != nil {
		panic(err)
	}

	var intervalRepetitions int = int(sTime / iTime) // Looping through intervals will be managed by for loop, and this variable will control how many iterations to perform.
	if sTime%iTime != 0 {
		intervalRepetitions++
	}

	for { //this for loop blocks code from further execution until 'start' is entered.
		fmt.Println("Input 'start' to start session.")
		reader := bufio.NewReader(os.Stdin)
		waitingOnStart, _ := reader.ReadString('\n')

		if strings.Compare(waitingOnStart, "start\n") == 1 {
			go readInput(userInput)
			break
		} else {
			fmt.Println("Input invalid: expected 'start'.")
		}
	}

	sessionStartTime := time.Now()
	sessionTimer := time.NewTimer(sTime * time.Second)
	for i := 0; i < intervalRepetitions; i++ { //for loop makes one iteration for each interval, requested by user.
		intervalData := []string{}
		intervalTimer := time.NewTimer(iTime * time.Second)
		log.Println("Interval start")

		/* case 1 - if we receive from channel, add data to variables in main routine
		case 2 - if session timer fires, terminate last interval and call method for printing session statistics.
		case 3 - if interval timer fires, terminate the interval and proceed to next iteration. */

	timerNotExpired:
		select {
		case intervalInput := <-userInput:
			intervalData = append(intervalData, intervalInput)
			sessionData = append(sessionData, intervalInput)

			if strings.Contains(intervalInput, "stop") {
				fmt.Println("Terminating session preemptively.")
				sessionTimeInS := time.Now().Sub(sessionStartTime)
				printIntervalStatistics([]byte(strings.Join(intervalData, "")), []byte(strings.Join(sessionData, "")), float32(sessionTimeInS.Seconds()))
				printSessionStatistics([]byte(strings.Join(sessionData, "")), float32(sessionTimeInS.Seconds()))
				goto sessionEnd
			}
			goto timerNotExpired
		case <-sessionTimer.C:
			log.Println("Interval end")
			sessionTimeInS := time.Now().Sub(sessionStartTime)
			printIntervalStatistics([]byte(strings.Join(intervalData, "")), []byte(strings.Join(sessionData, "")), float32(sessionTimeInS.Seconds()))

		case <-intervalTimer.C:
			log.Println("Interval end")
			sessionTimeInS := time.Now().Sub(sessionStartTime)
			printIntervalStatistics([]byte(strings.Join(intervalData, "")), []byte(strings.Join(sessionData, "")), float32(sessionTimeInS.Seconds()))
		}

		if i == intervalRepetitions-1 { //on the very last repetition, print session statistics, if it hasn't been done already.
			sessionTimeInS := time.Now().Sub(sessionStartTime)
			printSessionStatistics([]byte(strings.Join(sessionData, "")), float32(sessionTimeInS.Seconds()))
		}
	}
sessionEnd:
}
