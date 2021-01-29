package main

import (
"fmt"
"time"
"sync"
"regexp"
"os"
"bufio"
"net/http"
"net/url"
"io/ioutil"
"flag"
)

func main () {

	// The regex's map
	regs:=map[string]string {
		"links"    : `(https?|ftp|file)://[-A-Za-z0-9\+&@#/%?=~_|!:,.;]*[-A-Za-z0-9\+&@#/%=~_|]`,
		"awskeys"  : `([^A-Z0-9]|^)(AKIA|A3T|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{12,}`,
		"regex"    :`(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`,
		"domxss"   : `/((src|href|data|location|code|value|action)\s*["'\]]*\s*\+?\s*=)|((replace|assign|navigate|getResponseHeader|open(Dialog)?|showModalDialog|eval|evaluate|execCommand|execScript|setTimeout|setInterval)\s*["'\]]*\s*\()/`,
	}


	// Variables
	var concurrency int
	var mode string
		
	flag.IntVar(&concurrency, "c", 30, "Set concurrency for greater speed")
	flag.StringVar(&mode, "m", "links", "Set the regex to use (e.g. links,awskeys,domxss,endpoints)")
	flag.Parse()

	var wg sync.WaitGroup
	for i:=0; i<=concurrency; i++ {
		wg.Add(1)
		go func () {
			search_with_regex(mode, regs)
			wg.Done()	
		}()
		wg.Wait()
	}
}


// Search through all javascript with regex's
func search_with_regex(mode string, regs map[string]string) {

	time.Sleep(time.Millisecond * 10)
	client:= &http.Client{}
	

	scanner:=bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		jsLink:=scanner.Text()
		
		req,err:=http.NewRequest("GET", jsLink,nil)
		if err != nil {
			return
		}

		resp,err:=client.Do(req)
		if err != nil {
			return
		}

		bodyBuffer,err:=ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		// The body to grep for
		bodyString:=string(bodyBuffer)

		// Check to see if we are searching for links
		if mode == "links" {
			// Search for all links
			re:=regexp.MustCompile(regs["links"])
			match:=re.FindStringSubmatch(bodyString)
			if match != nil {
				fmt.Printf("%s\n", match[0])
			}
		}

		// Check to see if we are searching for apis
                if mode == "regex" {
                        // Search for all links
                        re:=regexp.MustCompile(regs["regex"])
                        match:=re.FindStringSubmatch(bodyString)
                        if match != nil {
                                fmt.Printf("%s\n", match[0])
                        }
                }

		// Check to see if we are searching for endpoints
                if mode == "endpoints" {
                        // Search for all links
                        re:=regexp.MustCompile(regs["links"])
                        match:=re.FindStringSubmatch(bodyString)
                        if match != nil {
				endpoint:=match[0]
				u,err:=url.Parse(endpoint)
				if err != nil { return }
                                fmt.Printf("%s\n", u.Path)
                        }
                }

	        // Check to see if we are searching for domxss
                if mode == "domxss" {
                        // Search for all links
                        re:=regexp.MustCompile(regs["domxss"])
                        match:=re.FindStringSubmatch(bodyString)
                        if match != nil {
                                fmt.Printf("%s : \t\t%q\n", match[0], jsLink)
                        }
                }

	}
}






















