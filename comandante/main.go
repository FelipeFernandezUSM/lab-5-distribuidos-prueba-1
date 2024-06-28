package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/FelipeFernandezUSM/lab-5-distribuidos-prueba-1/proto"
	"google.golang.org/grpc"
)

type LogEntry struct {
	SectorInfo    string
	VectorClock   []int32
	FulcrumServer string
}

var logEntries []LogEntry

type server struct {
	pb.UnimplementedCommandServer
	brokerClient pb.BrokerClient
	clientClocks map[string][]int32
}

func (s *server) GetIngenieros(ctx context.Context, in *pb.Comando) (*pb.Response, error) {
	// Obtener el reloj vectorial más reciente del cliente
	clientClock, ok := s.clientClocks[in.ClientId]
	if !ok {
		clientClock = make([]int32, len(s.clientClocks))
	}

	// Enviar el comando al servidor Broker
	message := &pb.Message{
		Sector:      in.GetSector(),
		Base:        in.GetBase(),
		VectorClock: clientClock,
	}
	notif, err := s.brokerClient.Mediate(ctx, message)
	if err != nil {
		return nil, err
	}

	// Actualizar el reloj vectorial del cliente
	s.clientClocks[in.ClientId] = notif.GetVectorClock()

	// Registrar el comando y la respuesta
	logEntry := LogEntry{
		SectorInfo:    fmt.Sprintf("GetIngenieros %s %s", in.GetSector(), in.GetBase()),
		VectorClock:   notif.GetVectorClock(),
		FulcrumServer: notif.GetFulcrumServer(),
	}
	logEntries = append(logEntries, logEntry)
	fmt.Printf("Respuesta recibida. Registrado %v\n", logEntry)

	// Devolver la respuesta del servidor Broker al usuario
	return &pb.Response{
		Notificacion:  notif.GetNotificacion(),
		FulcrumServer: notif.GetFulcrumServer(),
		VectorClock:   notif.GetVectorClock(),
	}, nil
}

func main() {
	// Crear un listener en el puerto TCP
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	// Crear un objeto de servidor gRPC
	grpcServer := grpc.NewServer()

	// Conectar con el servidor Broker
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fallo al conectar con el servidor Broker: %v", err)
	}
	defer conn.Close()

	// Crear un cliente Broker
	brokerClient := pb.NewBrokerClient(conn)

	// Crear un nuevo servidor Command
	commandServer := &server{brokerClient: brokerClient, clientClocks: make(map[string][]int32)}
	fmt.Println("Servidor Command en ejecución...")

	// Adjuntar el servicio Command al servidor gRPC
	pb.RegisterCommandServer(grpcServer, commandServer)

	// Iniciar el servidor gRPC (bloqueante)
	go func() {
		for {
			fmt.Print("Enter sector (ingresar solo el nombre, primera letra mayuscula): ")
			var sector string
			fmt.Scanln(&sector)

			fmt.Print("Enter base (ingresar solo el nombre, primera letra mayuscula): ")
			var base string
			fmt.Scanln(&base)

			cmd := &pb.Comando{Sector: sector, Base: base}

			res, err := commandServer.GetIngenieros(context.Background(), cmd)
			if err != nil {
				log.Fatalf("Failed to execute command: %v", err)
			}

			// Print the response
			fmt.Println("Response:", res.GetNotificacion())
			fmt.Println("Vector Clock:", res.GetVectorClock())
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir el servidor gRPC: %v", err)
	}
}
