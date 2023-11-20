package main


// import (
// 	"log"
// 	"os"
// )

// func logger(filePath string, prefix string, message string) error {
// 	// Open or create a log file
// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Set log output to the file
// 	log.SetOutput(file)
// 	log.SetPrefix(prefix + " ")
// 	// Write the message to the log file
// 	log.Println(message)

// 	return nil
// }