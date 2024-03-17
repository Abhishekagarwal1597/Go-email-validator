package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Print("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarcRecord\n")

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "@"){
				splitString := strings.Split(scanner.Text(),"@")
				if len(splitString) != 2{
					fmt.Sprintf("Kindly enter a valid email address: %s", scanner.Text())
				}
				checkDomain(splitString[1])
		}else{
			errMsg := fmt.Sprintf("Kindly include the @ in the email address: %s", scanner.Text())
			log.Fatal(errMsg)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}
}

//After executing this code, hasSPF and hasDMARC will indicate whether SPF and DMARC records exist
//for the given domain, and spfRecord and dmarcRecord will contain the actual record contents if they were found.
func checkDomain(domain string) {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	//This code checks if the domain of the email address has any MX (Mail Exchange) records.
	//If there are MX records, it indicates that the domain is set up to handle email,
	//which suggests that the email address might be valid.
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	//This function performs a DNS TXT lookup for the specified domain,
	//returning a slice of strings containing the TXT records associated with that domain.
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	for _, record := range txtRecords {
		//If a record starts with "v=spf1" (indicating an SPF record)
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	//This function performs a DNS TXT lookup for the DMARC record associated with the domain.
	//DMARC records are conventionally stored under the subdomain _dmarc.
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("ErrorL%v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v, %v, %v, %v, %v, %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
