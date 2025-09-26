package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func Encrypt(key []byte,message string) ([]byte,[]byte,error){
	//key := make([]byte,32)
	//_,err := io.ReadFull(rand.Reader,key)
	//if err != nil{
	//	return nil,err
	//}

	block,err := aes.NewCipher(key)
	if err != nil{
		return nil,nil,err
	}

	aead,err := cipher.NewGCM(block)
	if err != nil{
		return nil,nil,err
	}

	nonce := make([]byte,aead.NonceSize())
	_,err = io.ReadFull(rand.Reader,nonce)
	if err != nil{
		return nil,nil,err
	}

	cipherText := aead.Seal(nil,nonce,[]byte(message),nil)

	return cipherText,nonce,nil

}

func Decrypt(key,nonce,cipherText []byte) ([]byte,error){

	block,err :=aes.NewCipher(key)
	if err != nil{
		return nil,err
	}

	aead,err := cipher.NewGCM(block)
	if err != nil{
		return nil,err
	}

	msg,err := aead.Open(nil,nonce,cipherText,nil)
	if err != nil{
		return nil,err
	}

	return msg,nil

}
