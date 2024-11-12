package crypt

import (
	"ComicCollector/main/backend/utils/crypt/auth"
	"ComicCollector/main/backend/utils/env"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

var rsaKey *rsa.PrivateKey = nil

func InitRSAKey() bool {
	rsaKeyPath := env.GetRSAFilename()

	// check if the RSA key already exists
	// if not, generate a new one
	if _, err := os.Stat(rsaKeyPath); os.IsNotExist(err) {
		err = generateRSAKey(rsaKeyPath)
		if err != nil {
			log.Fatal("Error when generating a new RSA key")
		}
	}

	// load and return the existing rsa key
	key, err := loadRSAKey(rsaKeyPath)
	if err != nil {
		log.Fatal("Error when loading the RSA key")
	}

	rsaKey = key
	auth.KeySetup(rsaKey)

	return true
}

func generateRSAKey(filepath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Couldn't generate a new RSA private key")
	}

	// create the file on the disk
	pemPrivateFile, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		return err
	}

	// encode the private key to the PEM format
	var pemPrivateBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// write the private key to the file
	err = pem.Encode(pemPrivateFile, pemPrivateBlock)
	if err != nil {
		log.Println(err)
		return err
	}

	// close the file
	err = pemPrivateFile.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("A new RSA key has been generated successfully")

	return nil
}

func loadRSAKey(filepath string) (*rsa.PrivateKey, error) {
	// read the content of the local file
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// decode the data
	decodedData, _ := pem.Decode([]byte(data))
	privateKey, err := x509.ParsePKCS1PrivateKey(decodedData.Bytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return privateKey, nil
}
