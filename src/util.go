package main

import (
	"fmt"
	"os"
	"time"
)

func useQuickstart(s *server) bool {
	if len(os.Args) <= 1 {
		return false
	}
	if os.Args[1] == "s" {
		s.runIperfServer(-1)
		time.Sleep(300 * time.Second)
	} else if os.Args[1] == "c" {
		port := "5001"
		if len(os.Args) > 3 {
			port = os.Args[3]
		}
		s.runIperfClient(-1, os.Args[2], port, func(results string, failed bool) error {
			if !failed {
				fmt.Println("Passed with results")
				fmt.Println(results)
			} else {
				fmt.Fprintf(os.Stderr, "Failed with results")
				fmt.Fprintf(os.Stderr, results)
			}
			return nil
		})
		time.Sleep(60 * time.Second)
	} else {
		fmt.Fprintf(os.Stderr, "Recieved args %s, pos 1 needs to be 's' or 'c'", os.Args)
	}

	return true
}
