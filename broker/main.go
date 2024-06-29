package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/FelipeFernandezUSM/lab-5-distribuidos-prueba-1/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBrokerServer
}

// Initialize the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

var fulcrumServers = []string{
	"fulcrum1:50056",
	"fulcrum2:50057",
	"fulcrum3:50058",
}

// Función que redirige a un ingeniero a un servidor Fulcrum
func (s *server) RedirigirIngeniero(ctx context.Context, in *pb.RequestIngeniero) (*pb.DireccionFullcrum, error) {
	// Choose a random Fulcrum server
	address := fulcrumServers[rand.Intn(len(fulcrumServers))]
	fmt.Printf("Redirigir Ingeniero a %v\n", address) // Print the address of the Fulcrum server

	return &pb.DireccionFullcrum{Address: address}, nil // Return the address of the Fulcrum server
}

// Función que medía entre un sector y una base
func (s *server) Mediate(ctx context.Context, in *pb.Mensaje) (*pb.Notificacion, error) {
	address := fulcrumServers[rand.Intn(len(fulcrumServers))]

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewFulcrumClient(conn)

	// Send the message to the Fulcrum server
	message := &pb.Mensaje{
		Sector:      in.GetSector(),
		Base:        in.GetBase(),
		VectorClock: in.GetVectorClock(),
	}
	not, err := client.ProcessCommandMessage(ctx, message) // Send the message to the Fulcrum server
	fmt.Printf("Mensaje de Mediación %v, %v\n", in.GetSector(), in.GetBase())
	if err != nil {
		return nil, err
	}

	return not, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &server{})

	fmt.Printf("Broker listening at %v\n", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
