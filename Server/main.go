package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	pb "github.com/YuliaParshkova/TaskProto/proto/consignment"
)

const (
	port = ":50051"
)

type repository interface {
	Solve(client *pb.Coefficients) (*pb.Solution, error)
	GetAll() []*pb.Solution
}

//Repository ... Наша база данных
type Repository struct {
	solutions []*pb.Solution
}

//Solve
func (r *Repository) Solve(coeff *pb.Coefficients) (*pb.Solution, error) {
	t := CalcResult(coeff)
	m := append(r.solutions, t)
	r.solutions = m
	return t, nil
}

//GetAll ...
func (r *Repository) GetAll() []*pb.Solution {
	return r.solutions
}

type service struct {
	repo repository
}

func (s *service) GetAll(ctx context.Context, request *pb.GetRequest) (*pb.Solutions, error) {
	solutions := s.repo.GetAll()

	return &pb.Solutions{
		Solutions: solutions,
	}, nil
}

func (s *service) Solve(ctx context.Context, coefficients *pb.Coefficients) (*pb.Solution, error) {
	solve, err := s.repo.Solve(coefficients)
	if err != nil {
		return nil, err
	}
	return &pb.Solution{
		Coefs:  solve.Coefs,
		NRoots: solve.NRoots,
	}, nil
}

func main() {
	repo := &Repository{}

	//Настройка gRPC сервера
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}

	server := grpc.NewServer()

	//Регистрируем наш сервис для сервера
	ourService := &service{repo}
	pb.RegisterSolverServer(server, ourService)
	//Чтобы выходные параметры сервера сохранялись в go-runtime
	reflection.Register(server)

	log.Println("gRPC server running on port:", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to server from port: %v", err)
	}
}

func CalcResult(coeff *pb.Coefficients) *pb.Solution {
	a := coeff.A
	b := coeff.B
	c := coeff.C

	var Nroots int32

	if (a == 0 && b != 0) || (a != 0 && c == 0 && b == 0) || (a == b && c == 0) {
		Nroots = 1
	} else if a == 0 && b == 0 {
		Nroots = 0
	} else {
		D := b*b - 4*a*c
		if D < 0 {
			Nroots = 0
		} else if D > 0 {
			Nroots = 2
		} else {
			Nroots = 1
		}
	}

	return &pb.Solution{
		Coefs: &pb.Coefficients{
			A: a,
			B: b,
			C: c,
		},
		NRoots: Nroots,
	}
}
