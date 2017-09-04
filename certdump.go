package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/docker/notary/tuf/utils"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func pkcs8Cert(rsaRaw []byte) ([]byte, error) {

	block, _ := pem.Decode(rsaRaw)
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privKey, err := utils.RSAToPrivateKey(rsaKey)
	if err != nil {
		return nil, err
	}
	der, err := utils.ConvertTUFKeyToPKCS8(privKey, nil)
	output := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	return output, nil
}

func realMain() error {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		return nil
	}
	var data []byte
	var err error
	// certdump <filepath> <owner> <type> <data>
	path := os.Args[1]
	owner := os.Args[2]
	if os.Args[3] == "pkcs8" {
		data, err = pkcs8Cert([]byte(os.Args[4]))
		if err != nil {
			return err
		}
	} else {
		data = []byte(os.Args[3])
	}

	err = ioutil.WriteFile(path, data, 0700)
	if err != nil {
		return err
	}
	u, err := user.Lookup(owner)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}
	gid := os.Getgid()
	err = os.Chmod(path, 0400)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}

	return nil
}
