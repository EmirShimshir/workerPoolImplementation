package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// add constants
const (
	userCount   int = 100
	workerCount int = 12 // runtime.NumCPU()
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	jobs := make(chan User, userCount)    // make buffered User channel
	results := make(chan bool, userCount) // make buffered bool channel

	// use worker pool
	for i := 0; i < workerCount; i++ {
		go worker(jobs, results)
	}

	users := generateUsers(userCount)

	// write users to buffered channel
	for _, user := range users {
		jobs <- user
	}
	close(jobs)

	// print is job done
	for i := 0; i < userCount; i++ {
		fmt.Printf("job #%d done: %t\n", i+1, <-results)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

// workers read users form channel, save info and write true to results channel
func worker(jobs <-chan User, results chan<- bool) {
	for job := range jobs {
		saveUserInfo(job)
		results <- true
	}
}

func saveUserInfo(user User) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func generateUsers(count int) []User {
	users := make([]User, count)

	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}

	return users
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
