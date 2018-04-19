package main

import (
  "errors"
  "fmt"

  "encoding/pem"
  "encoding/asn1"
  "encoding/base64"

  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "crypto/sha256"

  "os"
  "io/ioutil"
)

const (
  NILBLOCK    = "decoded pemblock is nil"
  INVALID_PEM = "invalid pem file"

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
    return make([]byte, 0), errors.New("Error in keys creation")
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

func encrypt(pubkeyString []byte, raw string) (string, error) {
  nilStr   := ""
  block, _ := pem.Decode(pubkeyString)
  if block == nil {
    fmt.Printf("ERROR in encrypt data: block is nil\n")
    return nilStr, errors.New(NILBLOCK)
  }

  if got, want := block.Type, "PUBLIC KEY"; got != want {
    fmt.Printf("ERROR in encrypt data: invalid pem file\n")
    return nilStr, errors.New(INVALID_PEM)
  }

  pubkey, err := x509.ParsePKCS1PublicKey(block.Bytes)
  if err != nil {
    fmt.Printf("ERROR parsing pubkey PEM block to RSA Public Key format\n")
    return nilStr, err
  }

  cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey,
      []byte(raw), nil)
  if err != nil {
    fmt.Printf("ERROR in data encryption\n")
    return nilStr, err
  }

  cipherString := base64.StdEncoding.EncodeToString(cipher)
  return cipherString, nil
}

func decrypt(privkeyString []byte, cipherString string) (string, error) {
  nilStr := ""
  block, _ := pem.Decode(privkeyString)
  if block == nil {
    fmt.Printf("ERROR in decrypt data: block is nil\n")
    return nilStr, errors.New(NILBLOCK)
  }

  if got, want := block.Type, "PRIVATE KEY"; got != want {
    fmt.Printf("ERROR in decrypt data: invalid pem file\n")
    return nilStr, errors.New(INVALID_PEM)
  }

  privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
  if err != nil {
    fmt.Printf("ERROR parsing privkey PEM block to RSA Private Key format\n")
    return nilStr, err
  }

  cipher, _ := base64.StdEncoding.DecodeString(cipherString)
  rawString, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privkey, cipher, nil)
  if err != nil {
    fmt.Printf("ERROR in data decryption")
    return nilStr, err
  }

  return string(rawString), nil
}
