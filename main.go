package main

import (
	"flag"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	PathToKey = flag.String("keypath", "/etc/ssh-flux-key", "unix path to the TOTP URI key used")
)

func main() {
	flag.Parse()

	Token := ReadToken()
	otpkey, err := otp.NewKeyFromURL(Token)
	if err != nil {
		log.Fatalf("Unable to read OTP URI %v", err)
	}

	for {
		otp, err := totp.GenerateCode(otpkey.Secret(), time.Now())
		if err != nil {
			log.Printf("Unable to Generate code ??? %v", err)
			time.Sleep(time.Second * 30)
			continue
		}

		port, err := strconv.ParseInt(otp[:5], 10, 64)
		if err != nil {
			log.Printf("Unable to Parse code ??? %v", err)
			time.Sleep(time.Second * 30)
			continue
		}

		if port > 65536 {
			port, err = strconv.ParseInt(otp[:4], 10, 64)
			if err != nil {
				log.Printf("Unable to Parse code ??? %v", err)
				time.Sleep(time.Second * 30)
				continue
			}
		}
		go ReroutePort(int(port), 22)
		time.Sleep(time.Second * 30)
	}
}

func ReroutePort(port int, sshport int) {
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-p", "tcp", "--dport", fmt.Sprint(port), "-j", "REDIRECT", "--to-port", fmt.Sprint(sshport))
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	time.Sleep(time.Second * 30)
	cmd = exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-p", "tcp", "--dport", fmt.Sprint(port), "-j", "REDIRECT", "--to-port", fmt.Sprint(sshport))
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
}

func ReadToken() string {
	_, err := os.Stat(*PathToKey)
	if err != nil {
		Hostname, _ := os.Hostname()
		opts := totp.GenerateOpts{
			Issuer:      "sshflux",
			AccountName: Hostname,
		}
		newkey, err := totp.Generate(opts)
		err = ioutil.WriteFile(*PathToKey, []byte(newkey.String()), 0600)
		if err != nil {
			log.Fatalf("Unable to read/write to %s | %v", *PathToKey, err)
		}
		return newkey.String()
	}

	bytes, err := ioutil.ReadFile(*PathToKey)
	line := strings.Split(string(bytes), "\n")[0]
	return strings.Trim(line, "\r\n\t ")
}
