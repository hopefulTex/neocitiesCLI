package util

import (
	"fmt"
	"neocitiesCli/api"
	"os"
	"strings"
)

func Login() (api.Config, error) {
	config := api.Config{}

	domain, isSubdomain, err := GetDomainPrompt()
	if err != nil {
		return config, err
	}
	config.Domain = domain
	config.IsSubdomain = isSubdomain

	api_key, err := GetAPIkeyPrompt(domain)
	if err != nil {
		return config, err
	}
	config.APIKey = api_key

	return config, nil
}

func GetDomainPrompt() (string, bool, error) {
	var isSubdomain bool = false
	var err error
	isSubdomainString := ""
	domain := os.Getenv("NEOCITIES_SITENAME")
	for domain == "" {
		prompt := "Enter your Neocities Domain: "
		fmt.Print(prompt)
		_, err = fmt.Scanln(&domain)
		if err != nil && err.Error() != "unexpected newline" {
			fmt.Print("error: \n", err)
			return "", isSubdomain, err
		}
	}
	for isSubdomainString != "y" && isSubdomainString != "n" {
		fmt.Printf("Is this a subdomain (i.e. neocities.org/%s) or a custom domain?", domain)
		fmt.Print(" [y/n]\ny:  subdomain\nn:  custom domain\n")
		_, err = fmt.Scanln(&isSubdomainString)
		if err != nil && err.Error() != "unexpected newline" {
			return "", isSubdomain, err
		}
		if isSubdomainString == "Y" || isSubdomainString == "s" {
			isSubdomainString = "y"
		}
		if isSubdomainString == "N" || isSubdomainString == "c" {
			isSubdomainString = "n"
		}
		if isSubdomainString == "y" {
			// if !strings.HasPrefix(domain, "neocities.org/") {
			// 	domain = "neocities.org/" + domain
			// }
			isSubdomain = true
		}
	}
	return domain, isSubdomain, nil
}

func GetAPIkeyPrompt(domain string) (string, error) {
	var err error

	pass := os.Getenv("NEOCITIES_PASSWORD")
	for pass == "" {
		prompt := "Enter Password for " + domain + ": "
		fmt.Print(prompt)
		_, err = fmt.Scanln(&pass)
		if err != nil && err.Error() != "unexpected newline" {
			return "", err
		}
	}

	api_key, err := api.GetAPIkey(domain, pass)
	if err != nil {
		return "", err
	}
	fmt.Printf("Your API Key is: %s\n", api_key)

	if api_key == "" {
		err = fmt.Errorf("could not find API Key")
		return api_key, err
	}
	return api_key, nil
}

func ResetPrompt() bool {
	var input string
	fmt.Print("Are you sure you want to reset the configuration file? [y/n]\n")
	_, err := fmt.Scanln(&input)
	if err != nil && err.Error() != "unexpected newline" {
		return false
	}
	if input == "y" || input == "Y" {
		return true
	}
	return false
}

func RemoveGitIgnore(files []string) ([]string, error) {
	var new_files []string
	var isIgnored bool = false
	ignore, err := os.ReadFile(".gitignore")
	if err != nil {
		return []string{}, err
	}
	ignore_lines := strings.Split(string(ignore), "\n")
	for _, file := range files {
		isIgnored = false
		for _, ignore := range ignore_lines {

			if strings.Contains(ignore, "*") {
				// start
				if strings.HasPrefix(ignore, "*") {
					ignore = strings.TrimPrefix(ignore, "*")
					if strings.HasSuffix(file, ignore) {
						isIgnored = true
						break
					}
				}
				// end
				if strings.HasSuffix(ignore, "*") {
					ignore = strings.TrimSuffix(ignore, "*")
					if strings.HasPrefix(file, ignore) {
						isIgnored = true
						break
					}
				}
				// middle
				first, last, _ := strings.Cut(ignore, "*")
				if strings.HasPrefix(file, first) {
					if strings.HasSuffix(file, last) {
						isIgnored = true
						break
					}
				}
				continue
			}
			if ignore == file {
				isIgnored = true
				break
			}
		}
		if !isIgnored {
			new_files = append(new_files, file)
		}
	}
	return new_files, nil
}

func Upload(conn *api.Connection, paths []string) error {
	var errs []error
	var files []api.UploadFile

	for _, path := range paths {
		text, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		file := api.UploadFile{
			Name: path,
			File: text,
		}
		files = append(files, file)
	}
	conn.Upload(files)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
