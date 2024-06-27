package main

import (
    "log"
    "net"
	"os"
	"fmt"
	"strconv"
	"context"
    "bufio"
    "strings"
    "errors"
    "time"
    "sync"
    "encoding/gob"
    "os/signal"
    "io"
    "io/ioutil"

    "google.golang.org/grpc"
    pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-3/proto"
)

type FulcrumServer struct {
	pb.UnimplementedFulcrumServer
    id int
    state map[string]map[string]int
    vClocks map[string][]int
    otherServers []*grpc.ClientConn
    mu sync.Mutex
}

func NewFulcrumServer(id int) *FulcrumServer {
    s := &FulcrumServer{
        id:     id,
        state:  make(map[string]map[string]int),
        vClocks: make(map[string][]int),
    }

    // Initialize the otherServers slice
    for i := 0; i < 3; i++ {
        if i != id {
            conn, err := grpc.Dial(fmt.Sprintf("dist03%d:%d", 6+i, 50056+i), grpc.WithInsecure())
            if err != nil {
                log.Fatalf("Failed to connect to server %d: %v", i, err)
            }
            s.otherServers = append(s.otherServers, conn)
        }
    }

    return s
}

func (s *FulcrumServer) ProcessVanguardMessage(ctx context.Context, in *pb.Message) (*pb.Acknowledgement, error) {
    fmt.Println("Vanguard request received:", in.Sector, in.Base)
    // Get the stored vector clock for the sector
    storedClock, ok := s.vClocks[in.Sector]
    if !ok {
        storedClock = make([]int, len(s.vClocks))
    }

    // Compare the incoming vector clock with the stored vector clock
    for i := range in.VectorClock {
        if int(in.VectorClock[i]) < storedClock[i] {
            return nil, errors.New("stale read")
        }
    }

    // Open the .txt file
    filename := fmt.Sprintf("Sector%s.txt", strings.Title(in.Sector))
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Create a scanner to read the file
    scanner := bufio.NewScanner(file)

    // Loop over each line in the file
    for scanner.Scan() {
        // Split the line into sector, base, and soldiers
        parts := strings.Fields(scanner.Text())
        if len(parts) != 3 {
            continue
        }

        // Check if the sector and base match the input
        if parts[0] == in.Sector && parts[1] == in.Base {
            // Cast to []int32
            storedClock32 := make([]int32, len(storedClock))
            for i, v := range storedClock {
                storedClock32[i] = int32(v)
            }

            // Return the number of soldiers
            return &pb.Acknowledgement{
                Acknowledgement: parts[2],
                VectorClock:     storedClock32,
            }, nil
        }
    }

    // If no matching sector and base were found, return an error
    return nil, errors.New("sector and base not found")
}

func (s *FulcrumServer) AgregarBase(sector string, base string, quantity int) {
    // If the quantity is 0 and the sector does not exist, return an error
    if quantity == 0 {
        if _, ok := s.state[sector]; !ok {
            log.Println("Cannot create a new sector with a quantity of 0")
            return
        }
    }

    // Update the state
    if _, ok := s.state[sector]; !ok {
        s.state[sector] = make(map[string]int)
    }
    s.state[sector][base] = quantity

    // Update the vector clock
    if _, ok := s.vClocks[sector]; !ok {
        s.vClocks[sector] = make([]int, 3)
    }
    s.vClocks[sector][s.id-1]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("AgregarBase %s %s %d\n", sector, base, quantity))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector, "Agregar", strconv.Itoa(quantity))
}

func (s *FulcrumServer) RenombrarBase(sector string, base string, newBase string) {
    // If the sector or the base does not exist, return an error
    if _, ok := s.state[sector]; !ok {
        log.Println("Sector does not exist:", sector)
        return
    }
    if _, ok := s.state[sector][base]; !ok {
        log.Println("Base does not exist:", base)
        return
    }

    // Update the state
    s.state[sector][newBase] = s.state[sector][base]
    delete(s.state[sector], base)

    // Update the vector clock
    s.vClocks[sector][s.id-1]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("RenombrarBase %s %s %s\n", sector, base, newBase))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector, "Renombrar", newBase)
}

func (s *FulcrumServer) ActualizarValor(sector string, base string, newValue int) {
    // If the sector or the base does not exist, return an error
    if _, ok := s.state[sector]; !ok {
        log.Println("Sector does not exist:", sector)
        return
    }
    if _, ok := s.state[sector][base]; !ok {
        log.Println("Base does not exist:", base)
        return
    }

    // Update the state
    s.state[sector][base] = newValue

    // Update the vector clock
    s.vClocks[sector][s.id-1]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("ActualizarValor %s %s %d\n", sector, base, newValue))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector, "Actualizar", strconv.Itoa(newValue))
}

func (s *FulcrumServer) BorrarBase(sector string, base string) {
    // If the sector or the base does not exist, return an error
    if _, ok := s.state[sector]; !ok {
        log.Println("Sector does not exist:", sector)
        return
    }
    if _, ok := s.state[sector][base]; !ok {
        log.Println("Base does not exist:", base)
        return
    }

    // Update the state
    delete(s.state[sector], base)

    // Update the vector clock
    s.vClocks[sector][s.id-1]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("BorrarBase %s %s\n", sector, base))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector, "Borrar", base)
}

func (s *FulcrumServer) updateSectorFile(sector string, command string, value string) {
    // Open the sector file
    f, err := os.OpenFile(fmt.Sprintf("Sector%s.txt", sector), os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    defer f.Close()

    // Handle the command
    switch command {
    case "Agregar":
        // Write to the sector file
        for base, quantity := range s.state[sector] {
            _, err = f.WriteString(fmt.Sprintf("%s %s %d\n", sector, base, quantity))
            if err != nil {
                log.Println(err)
                return
            }
        }
    case "Renombrar":
        for base := range s.state[sector] {
            // Read the file line by line
            lines, err := ioutil.ReadFile(fmt.Sprintf("Sector%s.txt", sector))
            if err != nil {
                log.Println(err)
                return
            }

            // Split the file into lines
            linesSlice := strings.Split(string(lines), "\n")

            // For each line, check if it contains the base that you want to rename
            for i, line := range linesSlice {
                if strings.Contains(line, base) {
                    // If it does, replace the base name with the new name
                    linesSlice[i] = strings.Replace(line, base, value, 1)
                }
            }

            // Write the lines back to the file
            err = ioutil.WriteFile(fmt.Sprintf("Sector%s.txt", sector), []byte(strings.Join(linesSlice, "\n")), 0644)
            if err != nil {
                log.Println(err)
                return
            }
        }
    case "Actualizar":
        for base := range s.state[sector] {
            // Read the file line by line
            lines, err := ioutil.ReadFile(fmt.Sprintf("Sector%s.txt", sector))
            if err != nil {
                log.Println(err)
                return
            }
    
            // Split the file into lines
            linesSlice := strings.Split(string(lines), "\n")
    
            // For each line, check if it contains the base that you want to update
            for i, line := range linesSlice {
                if strings.Contains(line, base) {
                    // If it does, split the line into words
                    words := strings.Fields(line)
                    // Replace the old value with the new one
                    words[2] = value
                    // Join the words back into a line
                    linesSlice[i] = strings.Join(words, " ")
                }
            }
    
            // Write the lines back to the file
            err = ioutil.WriteFile(fmt.Sprintf("Sector%s.txt", sector), []byte(strings.Join(linesSlice, "\n")), 0644)
            if err != nil {
                log.Println(err)
                return
            }
        }
    case "Borrar":
        // Read the file line by line
        lines, err := ioutil.ReadFile(fmt.Sprintf("Sector%s.txt", sector))
        if err != nil {
            log.Println(err)
            return
        }
    
        // Split the file into lines
        linesSlice := strings.Split(string(lines), "\n")
    
        // Create a new slice to hold the lines that don't contain the base
        newLines := make([]string, 0)
    
        // For each line, check if it contains the base that you want to delete
        for _, line := range linesSlice {
            if !strings.Contains(line, value) {
                // If it doesn't, add it to the newLines slice
                newLines = append(newLines, line)
            }
        }
    
        // Write the newLines back to the file
        err = ioutil.WriteFile(fmt.Sprintf("Sector%s.txt", sector), []byte(strings.Join(newLines, "\n")), 0644)
        if err != nil {
            log.Println(err)
            return
        }
    }

    // Write the vector clock
    _, err = f.WriteString(fmt.Sprintf("[%d,%d,%d]\n", s.vClocks[sector][0], s.vClocks[sector][1], s.vClocks[sector][2]))
    if err != nil {
        log.Println(err)
    }
}

func (s *FulcrumServer) ApplyCommand(ctx context.Context, command *pb.CommandRequest) (*pb.CommandResponse, error) {
    fmt.Println("Command received:", command.Action, command.Sector, command.Base)
    switch command.Action {
    case "AgregarBase":
        s.AgregarBase(command.Sector, command.Base, int(command.Value))
    case "ActualizarValor":
        s.ActualizarValor(command.Sector, command.Base, int(command.Value))
    case "RenombrarBase":
        s.RenombrarBase(command.Sector, command.Base, command.NewBase)
    case "BorrarBase":
        s.BorrarBase(command.Sector, command.Base)
    default:
        return nil, fmt.Errorf("unknown action: %s", command.Action)
    }

    vectorClock := make([]int32, len(s.vClocks[command.Sector]))
	for i, v := range s.vClocks[command.Sector] {
		vectorClock[i] = int32(v)
	}

	// Return the vector clock of the modified sector
	return &pb.CommandResponse{
		VectorClock: vectorClock,
	}, nil
}

func (s *FulcrumServer) ApplyPropagation(ctx context.Context, p *pb.Propagation) (*pb.PropagationResponse, error) {
    // Lock the server state for writing
    s.mu.Lock()
    defer s.mu.Unlock()

    // Get the current state and vector clock for the sector
    currentState, currentVC := s.state[p.Sector], s.vClocks[p.Sector]
    
    // If the sector doesn't exist in the state map, initialize it
    if currentState == nil {
        currentState = make(map[string]int)
        s.state[p.Sector] = currentState
    }

    // If the sector doesn't exist in the vector clocks map, initialize it
    if currentVC == nil {
        currentVC = make([]int, 3)
        s.vClocks[p.Sector] = currentVC
    }

    // Convert the incoming state to map[string]int
    incomingState := make(map[string]int)
    for k, v := range p.State {
        incomingState[k] = int(v)
    }

    // Compare the incoming vector clock with the current vector clock
    for i, incomingTime := range p.VectorClock {
        // Ensure currentVC is long enough
        for len(currentVC) <= i {
            currentVC = append(currentVC, 0)
        }

        incomingTimeInt := int(incomingTime)
        if incomingTimeInt > currentVC[i] {
            // The incoming state is more recent, so update the server state and vector clock
            currentState = incomingState
            currentVC[i] = incomingTimeInt
        } else if incomingTimeInt < currentVC[i] {
            // The server state is more recent, so ignore the incoming state
            continue
        } else {
            // The incoming state and server state are concurrent, so resolve the conflict
            for k, v := range incomingState {
                // If currentState is nil, initialize it
                if currentState == nil {
                    currentState = make(map[string]int)
                }

                if v2, ok := currentState[k]; !ok || v > v2 {
                    // If the key is not in the current state, or the incoming value is greater,
                    // update the current state with the incoming value
                    currentState[k] = v
                }
            }
        }
    }

    // If the sector doesn't exist in the state map, initialize it
    if s.state[p.Sector] == nil {
        s.state[p.Sector] = make(map[string]int)
    }

    // If the sector doesn't exist in the vector clocks map, initialize it
    if s.vClocks[p.Sector] == nil {
        s.vClocks[p.Sector] = make([]int, 3)
    }

    // Update the server state and vector clock
    s.state[p.Sector] = currentState
    s.vClocks[p.Sector] = currentVC

    // Check if the log file exists
    if _, err := os.Stat("log.txt"); err == nil {
        // If it exists, delete the log file
        err := os.Remove("log.txt")
        if err != nil {
            return nil, fmt.Errorf("failed to delete log file: %w", err)
        }
    }

    // Open the log file
    logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }
    defer logFile.Close()

    // Open the sector file
    sectorFile, err := os.OpenFile(fmt.Sprintf("Sector%s.txt", p.Sector), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open sector file: %w", err)
    }
    defer sectorFile.Close()

    // Write to the sector file
    for base, value := range s.state[p.Sector] {
        _, err = fmt.Fprintf(sectorFile, "%s %s %d\n", p.Sector, base, value)
        if err != nil {
            return nil, fmt.Errorf("failed to write to sector file: %w", err)
        }
    }

    return &pb.PropagationResponse{Success: true, Message: "Propagation applied successfully"}, nil
}

func (s *FulcrumServer) PropagateChanges() {
    fmt.Println("Propagating changes...")
    // Iterate over all other servers
    for _, otherServer := range s.otherServers {
        // Create a Fulcrum client
        fulcrumClient := pb.NewFulcrumClient(otherServer)

        // Create a context with a timeout
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // Iterate over all sectors
        for sector, state := range s.state {
            // Convert the state map to map[string]int32
            stateInt32 := make(map[string]int32)
            for k, v := range state {
                stateInt32[k] = int32(v)
            }

            // Convert the vector clock to []int32
            vClockInt32 := make([]int32, len(s.vClocks[sector]))
            for i, v := range s.vClocks[sector] {
                vClockInt32[i] = int32(v)
            }

            // Prepare the message with the current state and vector clock for the sector
            message := &pb.Propagation{
                Sector:      sector,
                State:       stateInt32,
                VectorClock: vClockInt32,
            }

            // Send the message to the other server
            _, err := fulcrumClient.ApplyPropagation(ctx, message)
            if err != nil {
                log.Println("Failed to propagate changes to server:", err)
            }
        }
    }
}

func (s *FulcrumServer) saveVectorClocks() error {
    file, err := os.Create("vectorClocks.gob")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := gob.NewEncoder(file)
    err = encoder.Encode(s.vClocks)
    if err != nil {
        return err
    }

    return nil
}

func (s *FulcrumServer) loadVectorClocks() error {
    file, err := os.Open("vectorClocks.gob")
    if err != nil {
        if os.IsNotExist(err) {
            // If the file doesn't exist, that's okay; we'll just start with empty vector clocks
            return nil
        }
        return err
    }
    defer file.Close()

    decoder := gob.NewDecoder(file)
    err = decoder.Decode(&s.vClocks)
    if err != nil {
        if err == io.EOF {
            // If the file is empty, that's okay; we'll just start with empty vector clocks
            return nil
        }
        return err
    }

    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run main.go <server_id>")
        os.Exit(1)
    }

    id, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Invalid server ID:", os.Args[1])
        os.Exit(1)
    }

    // Initialize the server
    s := NewFulcrumServer(id)

    // Load the vector clocks
    if err := s.loadVectorClocks(); err != nil {
        log.Fatalf("Failed to load vector clocks: %v", err)
    }

    // Create a channel to receive OS signals
    sigs := make(chan os.Signal, 1)
    // Register the channel to receive interrupt signals
    signal.Notify(sigs, os.Interrupt)

    // Start a goroutine to propagate changes every 30 seconds
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            // Propagate changes to all other servers
            s.PropagateChanges()
        }
    }()

    // Start a gRPC server
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50055+id))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterFulcrumServer(grpcServer, s)

    // Start a goroutine that will stop the server when an interrupt signal is received
    go func() {
        <-sigs
        log.Println("Stopping server...")
        grpcServer.Stop()
        if err := s.saveVectorClocks(); err != nil {
            log.Fatalf("Failed to save vector clocks: %v", err)
        }
        os.Exit(0)
    }()

    log.Printf("Fulcrum Server %v is running...", id)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}