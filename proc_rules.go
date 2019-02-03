package main

import (
	"strings"
	"os"
	"os/exec"
	"io"
	"encoding/csv"
	"bufio"
	"log"
	"strconv"
)

type User struct {
    id string
    ip string
    port uint
}

//var users []User
var users map[string]User


// Do customized initialisation
func proc_init() {
	read_users()
}

func process_rules(message string, rlog string) bool {
	fields := strings.Fields(message)

	if false {
		for i, field := range fields {
			log.Println(i,": ",field)
		}
	}

	if len(fields) > 16 {
		if fields[10] == "AH00957:" {
			f := fields[16]
			s := strings.Split(f, ":")
			f = s[1]
			log.Println(users[f].id, "wird gestartet!")
			cmd := exec.Command("/usr/bin/monit", "start", users[f].id)
			_ = cmd.Start()
		}
	}

        if len(fields) > 8 {
                if fields[4] == "error" && fields[8] == "cpu" {
                        f := fields[6]
                        f = strings.Replace(f, "'", "", 2)
                        log.Println(f, "wird gestoppt!")
                        cmd := exec.Command("/usr/bin/monit", "stop", f)
                        _ = cmd.Start()
                }
        }

        if len(fields) > 12 {
                if fields[4] == "error" && fields[12] == "125)" {
                        f := fields[6]
                        f = strings.Replace(f, "'", "", 2)
                        log.Println(f, " Container wird entfernt!")
                        cmd := exec.Command("/usr/bin/docker", "rm", f)
                        _ = cmd.Start()
                }
        }

        if len(fields) > 10 {
                if fields[10] == "AH00126:" {
                        return false
                }
        }

        if len(fields) > 10 {
                if fields[10] == "AH01114:" {
                        return false
                }
        }

        if len(fields) > 10 {
                if fields[10] == "AH01276:" {
                        return false
                }
        }

	return true
}

func read_users() {
    csvFile, _ := os.Open("/root/ibrowser-helper/envs.dat")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    reader.Comma = rune(59)
    users = make(map[string]User)
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
        users[line[2]] = User{
            id: line[0],
            ip:  line[1],
            port: uint(i),
        }
    }
    log.Println("User data read")
    csvFile.Close()
    if do_trace {
            log.Println(users)
    }
}

