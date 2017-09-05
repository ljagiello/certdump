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
	if err != nil {
		return nil, err
	}
	output := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	return output, nil
}

func lookupGroup(groupname string) (int, error) {
	g, err := user.LookupGroup(groupname)
	if err != nil {
		return -1, err
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return -1, err
	}
	return gid, nil
}

func lookupUser(username string) (int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return -1, err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

func realMain() error {
	if len(os.Args) < 6 || len(os.Args) > 7 {
		return nil
	}

	var err error
	var data []byte

	// certdump <filepath> <user> <group> <mode> <type> <data>
	path := os.Args[1]
	user := os.Args[2]
	group := os.Args[3]
	mode := os.Args[4]

	if os.Args[5] == "pkcs8" {
		data, err = pkcs8Cert([]byte(os.Args[6]))
		if err != nil {
			return err
		}
	} else {
		data = []byte(os.Args[5])
	}

	if _, err = os.Stat(path); !os.IsNotExist(err) {
		err = os.Chmod(path, 0600)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(path, data, 0400)
	if err != nil {
		return err
	}

	// uid, gid, chown
	uid, err := lookupUser(user)
	if err != nil {
		return err
	}
	gid, err := lookupGroup(group)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}

	filemode, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		// if error return filemode 0400
		filemode = 0400
	}
	err = os.Chmod(path, os.FileMode(filemode))
	if err != nil {
		return err
	}

	return nil
}
