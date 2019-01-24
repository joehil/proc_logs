package main

import (
	"strings"
	"os"
	"io"
	"encoding/csv"
	"bufio"
	"log"
	"strconv"
)

type User struct {
    id string
    ip  string
    port uint
}

var users []User


// Do customized initialisation
func proc_init() {
	read_users()
}

func process_rules(message string, rlog string) bool {
//	fields := strings.Fields(message)

	return !strings.Contains(message, "AH01276")
	//return true
}

func read_users() {
    csvFile, _ := os.Open("/root/ibrowser-helper/envs.dat")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    reader.Comma = rune(59)
    users = nil
    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }
        i, err := strconv.Atoi(line[2])
        if err != nil {
                // handle error
                log.Fatal(err)
                os.Exit(2)
        }
        users = append(users, User{
            id: line[0],
            ip:  line[1],
            port: uint(i),
        })
    }
    log.Println("User data read")
    csvFile.Close()
    if do_trace {
        for _, us := range users{
                log.Println(us)
        }
    }
}
