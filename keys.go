package main

import (
  "errors"
  "fmt"

  "encoding/pem"
  "encoding/asn1"

  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"

  "os"
  "io/ioutil"
)

func loadKeys() ([]byte, error) {
  pubkey, err := ioutil.ReadFile("priv/ussd_pub.pem")

  if err != nil {
    fmt.Printf("ERROR in loadKeys: %v\nCreating new key...\n", err)
    return createKeys()
  }

  return pubkey, nil
}

func createKeys() ([]byte, error) {
  reader   := rand.Reader
  key, err := rsa.GenerateKey(reader, 2048)
  if err != nil {
    fmt.Printf("ERROR in createKeys: %v\n", err)
    return make([]byte, 0), err
  }

  errPriv := savePrivatePem("priv/ussd_priv.pem", key)
  errPub  := savePublicPem("priv/ussd_pub.pem", key.PublicKey)

  if errPriv != nil || errPub != nil {
    return make([]byte, 0), errors.New("Error in key creation")
  }

  return ioutil.ReadFile("priv/ussd_pub.pem")
}

func savePrivatePem(filename string, key *rsa.PrivateKey) error {
  pemFile, err := os.Create(filename)
  defer pemFile.Close()

  if err != nil {
    fmt.Printf("ERROR in savePrivateKey: %v\n", pemFile)
    return err
  }

  privateKey := &pem.Block{
    Type: "PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(key),
  }

  err = pem.Encode(pemFile, privateKey)
  if err != nil {
    fmt.Printf("ERROR in savePrivateKey: %v\n", err)
    return err
  }

  return nil
}

func savePublicPem(filename string, pubkey rsa.PublicKey) error {
  asn1Bytes, err := asn1.Marshal(pubkey)
  if err != nil {
    fmt.Printf("ERROR in savePublicKey: ﬁv\n", pubkey)
    return err
  }

  pemPubFile, err := os.Create(filename)
  defer pemPubFile.Close()
  if err != nil {
    fmt.Printf("ERROR in savePublicKey: %v\n", pemPubFile)
    return err
  }

  publicKey := &pem.Block{
    Type: "PUBLIC KEY",
    Bytes: asn1Bytes,
  }

  err = pem.Encode(pemPubFile, publicKey)
  if err != nil {
    fmt.Printf("ERROR in savePublicKey: %v\n", err)
    return err
  }

  return nil
}
