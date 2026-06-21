package main

import (
	"fmt"
	"log"
	"net"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	log.Println("A client has connected", c.RemoteAddr())
	defer c.Close()

	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	// Read data from the client
	n, err := c.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data from client:", err)
		return
	}

	// Print the received data
	client_value := string(buffer[:n])
	fmt.Printf("Received data from client: %s\n", client_value)
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"), // Replace with your AWS region
		Credentials: credentials.NewStaticCredentials("AKIAWHRR72TTEAHQ4LY2", "LYUHavuXaOjZntdFDhXCc9ZkwsdZo5yoZlKPEI6R", ""),
	}))

	// Create a DynamoDB client.
	dynamodb_client := dynamodb.New(sess)

	table_name := "clab-dynamodb-table"
	key_name := "UserID"
	key_value := client_value

	input := &dynamodb.GetItemInput{
		TableName: aws.String(table_name),
		Key: map[string]*dynamodb.AttributeValue{
			key_name: {
				S: aws.String(key_value),
			},
		},
	}

	result, err := dynamodb_client.GetItem(input)
	if err != nil {
		fmt.Println("Error getting item from DynamoDB", err)
		return
	}

	// Check if the item was found.
	if result.Item == nil {
		fmt.Println("Item not found")
		return
	}

	// Process the retrieved item.
	var tabledata []string
	for key, value := range result.Item {
		// fmt.Printf("%s: %v\n", key, value)
		formattedString := fmt.Sprintf("%s: %v\n", key, value)
		tabledata = append(tabledata, formattedString)
	}
	messageforclient := ""
	for _, str := range tabledata {
		messageforclient += str
	}

	c.Write([]byte(messageforclient))
}
