package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var reserved = map[string]int{}

func loadReserved(filename string) map[string]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	var reservations = make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		login := parts[0]
		number, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println(err)
		} else {
			reservations[login] = number
		}
	}

	return reservations
}

func checkIfReserved(value int, reservations map[string]int) bool {
	for _, v := range reservations {
		if v == value {
			return true
		}
	}
	return false
}

func main() {
	var counter int = 0

	portPtr := flag.Int("port", 19229, "Bind to specified UDP port.")
	addrPtr := flag.String("bind", "127.0.0.1", "Address to bind on.")

	flag.Parse()

	if _, err := os.Stat("reserved.lst"); err == nil {
		fmt.Println("Loading reserved numbers from 'reserved.lst'...")
		reserved = loadReserved("reserved.lst")
	}

	if len(reserved) > 0 {
		for login, number := range reserved {
			fmt.Printf("Reserved %d for '%s'.\n", number, login)
		}
	} else {
		fmt.Println("No numbers reserved for anyone.")
	}

	addr := net.UDPAddr{
		Port: *portPtr,
		IP:   net.ParseIP(*addrPtr),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	var buf []byte = make([]byte, 256)
	fmt.Println("Listening for requests...")
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		received := string(buf[:n-1])

		is_reserved := "no"
		assigned := -1
		if n-1 > 0 {
			if val, ok := reserved[received]; ok {
				assigned = val
				is_reserved = "yes"
			} else {
				counter += 1
				for checkIfReserved(counter, reserved) == true {
					counter += 1
				}
				assigned = counter
			}
		}

		out := fmt.Sprintf("%d:%s\n", assigned, received)
		conn.WriteToUDP([]byte(out), addr)
		fmt.Printf(
			"%s (%d bytes) '%s' -> %d (reserved: %s)\n",
			addr, n-1, received, assigned, is_reserved,
		)
	}
}
