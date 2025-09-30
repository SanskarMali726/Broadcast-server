package main

import (
	"fmt"
	"os"
	"github.com/SanskarMali726/Broadcast-server/encryption"
	"github.com/joho/godotenv"
)

func sanskar(){
	err := godotenv.Load()
	if err != nil{
		panic(err)
	}
	key := os.Getenv("kEY")
	ciphertext,nonce,err :=encryption.Encrypt([]byte(key),"Hello sanskar",)
	if err != nil{
		panic(err)
	}

	msg,err := encryption.Decrypt([]byte(key),nonce,ciphertext)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(msg))
}