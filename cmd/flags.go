package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	email "github.com/MarlonCorreia/ephemail/internal/email"
)

func CheckValidDomain(m *email.EmailModel, d string) (bool, error) {
	availableDomains, err := m.ListEmailDomains()
	if err != nil {
		return false, err
	}
	for _, domain := range availableDomains {
		if d == domain {
			return true, nil
		}
	}
	return false, nil
}

func ParseFlags(m *email.EmailModel) error {
	listDomains := flag.Bool("list-domains", false, "List available domains")
	user := flag.String("user", "", "Username for the e-mail")
	domain := flag.String("domain", "", "Domain for the e-mail")
	flag.Parse()

	if *listDomains {
		domainsList, err := m.ListEmailDomains()
		if err != nil {
			return errors.New(fetchDomainError)
		}

		domains := strings.Join(domainsList, " - ")
		fmt.Println(domains)
		os.Exit(0)
	} else if *user != "" {
		isValid, err := CheckValidDomain(m, *domain)
		if err != nil {
			return errors.New(fetchDomainError)
		} else if !isValid {
			return errors.New(domainUnavailableErr)
		}
		m.User = *user
		m.Domain = *domain
	}

	return nil
}
