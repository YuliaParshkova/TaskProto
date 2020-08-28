package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"

	pb "github.com/YuliaParshkova/TaskProto/Client/proto/consignment"
)

const (
	address         = "localhost:50051"
	defaultFilename = "command.json"
)

func parseJSON(file string) (*pb.Coefficients, error) {
	var coeff *pb.Coefficients
	fileBody, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(fileBody, &coeff)
	return coeff, err
}

func main() {

	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect to port: %v", err)
	}
	defer connection.Close()

	client := pb.NewSolverClient(connection)

	coeff, err := parseJSON(defaultFilename)
	if err != nil {
		log.Fatalf("can not parse .json file: %v", err)
	}
	resp, err := client.Solve(context.Background(), coeff) //CreateCommand(context.Background(), command)
	if err != nil {
		log.Fatalf("can not get response: %v", err)
	}
	log.Printf("NRoots: %t", resp.NRoots)

	getAll, err := client.GetAll(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("can not get response: %v", err)
	}
	for _, v := range getAll.Solutions {
		log.Println(v)
	}

}
