package main

import (
	"log"
	"bufio"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"fmt"
	"strings"
	"strconv"
	"syscall"
	"github.com/illarion/gonotify"
	"github.com/spf13/viper"
)

//var read_log1 string = "/var/log/monit.log"
//var read_log2 string = "/var/log/virtualmin/remote-browser.eu_error_log"
var do_trace bool = false
var pidfile string
var ownlog string
var read_log1 string
var read_log2 string

type User struct {
    id string
    ip  string
    port uint
}

var users []User

func main() {
// Set location of config 
	viper.SetConfigName("proc_logs") // name of config file (without extension)
	viper.AddConfigPath("/etc/")   // path to look for the config file in

// Read config
	read_config()

// Write pidfile
        err := writePidFile(pidfile)
        if err != nil { 
                log.Fatalf("Pidfile could not be written: %v", err)
        }
	defer os.Remove(pidfile)

// Open log file
	f, err := os.OpenFile(ownlog, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
    		log.Fatalf("error opening ownlog: %v", err)
	}
	defer f.Close()
//	wrt := io.MultiWriter(os.Stdout, f)
	wrt := io.MultiWriter(f)
	log.SetOutput(wrt)

// Inform about trace
	log.Println("Trace set to: ", do_trace)

// Read users
	read_users()

// Catch signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP)
	go catch_signals(signals)

// Open read_log1
        r1, err := os.Open(read_log1)
        if err != nil {
                log.Fatalf("error opening read_log1: %v", err)
        }
        defer r1.Close()
        n1, err := r1.Seek(0, 2)

// Open read_log2
        r2, err := os.Open(read_log2)
        if err != nil {
                log.Fatalf("error opening read_log2: %v", err)
        }
        defer r2.Close()
        n2, err := r2.Seek(0, 2)

// Setup inotify watcher
        watcher, err := gonotify.NewFileWatcher(gonotify.IN_MODIFY | gonotify.IN_MOVED_FROM | gonotify.IN_MOVED_TO | gonotify.IN_CREATE, 
                        read_log1, read_log2)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

// Read watcher events
        for {
		e := <-watcher.C
		if do_trace {
			log.Println("event: ", e.Wd, e.Name, e.Mask)
		}
		switch e.Wd {
		case 1:
			switch e.Mask {
			case 2:
			n1 = proc_log1(r1, n1)
			case 64:
			r1.Close()
			case 256:
        		r1, err = os.Open(read_log1)
        		if err != nil {
                		log.Fatalf("error opening read_log1: %v", err)
        		}
		        n1, err = r1.Seek(0, 0)
			}
		case 2:
                        switch e.Mask {
                        case 2:
                        n2 = proc_log1(r2, n2)
                        case 64:
                        r2.Close()
                        case 256:
                        r2, err = os.Open(read_log2)
                        if err != nil {
                                log.Fatalf("error opening read_log2: %v", err)
                        }
                        n2, err = r2.Seek(0, 0)
                        }
		default:
		log.Println("Invalid event number")
		}
	}
}

func proc_log1(f *os.File, p int64) int64 {
	b1 := make([]byte, 500)
        _, _ = f.Seek(p, 0)
  	m, err := f.Read(b1)
	if m > 0 && err == nil {
		t := string(b1)[:m]
		mes := strings.Split(t, "\n")
    		for _, ames := range mes[:len(mes)-1] {
        		log.Print(ames)
    		}
		p = p + int64(m)
	} 
	return p
}

func proc_log2(f *os.File, p int64) int64 {
        b1 := make([]byte, 500)
        _, _ = f.Seek(p, 0)
        m, err := f.Read(b1)
        if m > 0 && err == nil {
                t := string(b1)[:m]
                mes := strings.Split(t, "\n")
                for _, ames := range mes[:len(mes)-1] {
                        log.Print(ames)
                }
                p = p + int64(m)
        }
        return p
}

// Write a pid file, but first make sure it doesn't exist with a running pid.
func writePidFile(pidFile string) error {
	// Read in the pid file as a slice of bytes.
	if piddata, err := ioutil.ReadFile(pidFile); err == nil {
		// Convert the file contents to an integer.
		if pid, err := strconv.Atoi(string(piddata)); err == nil {
			// Look for the pid in the process list.
			if process, err := os.FindProcess(pid); err == nil {
				// Send the process a signal zero kill.
				if err := process.Signal(syscall.Signal(0)); err == nil {
					// We only get an error if the pid isn't running, or it's not ours.
					return fmt.Errorf("pid already running: %d", pid)
				}
			}
		}
	}
	// If we get here, then the pidfile didn't exist,
	// or the pid in it doesn't belong to the user running this app.
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0664)
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

func catch_signals(c <-chan os.Signal){
	for {
		s := <-c
		log.Println("Got signal:", s)
		read_config()
		read_users()
	}
}

func read_config() {
        err := viper.ReadInConfig() // Find and read the config file
        if err != nil { // Handle errors reading the config file
                log.Fatalf("Config file not found: %v", err)
        }

        pidfile = viper.GetString("pid_file")
        if pidfile =="" { // Handle errors reading the config file
                log.Fatalf("Filename for pidfile unknown: %v", err)
        }
        ownlog = viper.GetString("own_log")
        if ownlog =="" { // Handle errors reading the config file
                log.Fatalf("Filename for ownlog unknown: %v", err)
        }
        read_log1 = viper.GetString("read_log1")
        if read_log1 =="" { // Handle errors reading the config file
                log.Fatalf("Filename for read_log1 unknown: %v", err)
        }
        read_log2 = viper.GetString("read_log2")
        if read_log2 =="" { // Handle errors reading the config file
                log.Fatalf("Filename for read_log2 unknown: %v", err)
        }
        do_trace = viper.GetBool("do_trace")

	if do_trace {
		log.Println("do_trace: ",do_trace)
		log.Println("own_log; ",ownlog)
		log.Println("pid_file: ",pidfile)
		log.Println("read_log1: ",read_log1)
		log.Println("read_log2: ",read_log2)
	}
}
