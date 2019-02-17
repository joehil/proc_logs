package main

import (
//	"strings"
//	"os"
//	"os/exec"
//	"io"
//	"encoding/csv"
//	"bufio"
//	"log"
//	"strconv"
)

func suppress_field(nr int, word string, do_log bool, fields []string) bool {
        if len(fields) > nr {
                if fields[nr] == word {
                        return false
                }
        }
	return do_log
}

