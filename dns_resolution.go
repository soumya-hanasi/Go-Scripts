// dns_resolver.go
// ------------------------------------------------------------
// Output format:
// timestamp,status,dns,resolvedip,timetaken(ms),error
//
// Example:
// 2026-04-22 13:00:01,SUCCESS,google.com,[142.250.183.14,142.250.183.15],10,No Error
// ------------------------------------------------------------

package main

import (
	"bufio"   // Buffered I/O for efficient file reading/writing
	"fmt"     // Formatted I/O (printing and string formatting)
	"net"     // Networking (DNS resolution)
	"os"      // OS operations (files, arguments)
	"strings" // String manipulation
	"time"    // Time measurement
)

func main() {

	// Validate arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: dns_resolver <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open input file
	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	// ✅ Open output file in APPEND mode
	out, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	scanner := bufio.NewScanner(in)

	// Write header only if file is empty
	fileInfo, _ := out.Stat()
	if fileInfo.Size() == 0 {
		writer.WriteString("timestamp,status,dns,resolvedip,timetaken(ms),error\n")
	}

	// Process each DNS
	for scanner.Scan() {

		dns := strings.TrimSpace(scanner.Text())

		// Skip empty/comment lines
		if dns == "" || strings.HasPrefix(dns, "#") {
			continue
		}

		// Current timestamp
		timestamp := time.Now().Format("2026-01-02 15:04:05")

		// Start timer
		start := time.Now()

		// DNS resolution
		ips, err := net.LookupHost(dns)

		// Time taken in milliseconds
		duration := time.Since(start).Milliseconds()

		if err != nil || len(ips) == 0 {
			// ❌ Failure case
			writer.WriteString(fmt.Sprintf(
				"%s,FAIL,%s,[],%d,%v\n",
				timestamp, dns, duration, err,
			))
		} else {
			// ✅ Success case
			ipList := "[" + strings.Join(ips, ",") + "]"

			writer.WriteString(fmt.Sprintf(
				"%s,SUCCESS,%s,%s,%d,No Error\n",
				timestamp, dns, ipList, duration,
			))
		}
	}

	// Handle scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	// Flush buffer
	writer.Flush()

	fmt.Printf("Done. Data appended to %s\n", outputFile)
}
