package at_rest


//I copy/pasted alot from https://github.com/gtank/cryptopasta
//in an attempt to get the encryption right.
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"errors"
	"bytes"
	"fmt"
)
var key_name string = "EARKEY"
//This function makes a handy, cross-platform way to 
//automaticallly encrypt and read smallish configuration files 
//that hold service login credentials, encrypton keys and so forth.
//You give EncryptedAtRest the name of the original un-encrypted 
//file and it returns a single byte array of un-encrypted bytes 
//that are in the file.
//The magic takes place when you first call the function
//on an un-encrypted file.  
//If an unencrypted file is there, a new encrypted file is created
//with an .aes extension and the original unencrypted file is 
//overwritten with random bytes and then deleted.
//If the original file is not there, it looks for the .aes encrypted 
//one, decrypts it and returns the decrypted bytes.
//The function uses the contents of the envirnment variable "EARKEY" 
//to encrypt and decrypt files.  You can ether set the envirnment
//variable EARKEY to your own highly random 32 byte key encoded 
//in base64, or you can just leave it blank and the function will
//create EARKEY for you and set it to a highly randomized and  
//encoded key.
func EncryptedAtRest( file_name string ) ([]byte, error) {
	ear_key_str :=  os.Getenv(key_name)
	if len(ear_key_str) == 0 {
		newkey := NewEncryptionKeyAsBase64()
		os.Setenv(key_name, newkey)
		ear_key_str = os.Getenv(key_name)
		if newkey != ear_key_str {
			return nil, errors.New("Unable to set encryption key.")
		}
	}
	ear_key := DecodeBase64Key(ear_key_str)
	enc_file_name := fmt.Sprintf("%s.aes",file_name)
	if _, err := os.Stat(file_name); !os.IsNotExist(err) {
		//file needs to be encrypted
		if _, err := os.Stat(enc_file_name); !os.IsNotExist(err) {
			//remove the old one if it is there
			os.Remove(enc_file_name)
		}
		//read the source file into a byte array
		dat, dat_err := ioutil.ReadFile(file_name)
		if dat_err != nil {
			return nil, dat_err
		}
		//create an encrypted byte array
		enc_dat, enc_dat_err := Encrypt(dat, ear_key)
		if enc_dat_err != nil {
			return nil, enc_dat_err
		}		
		//create the .aes file and write the encrypted data to it
		write_err := ioutil.WriteFile(enc_file_name, enc_dat, 0644)
		if write_err != nil {
			return nil, write_err
		}
		_, err := io.ReadFull(rand.Reader, dat[:])
		if err != nil {
			return nil, err
		}
		//write the source file with random bytes
		write_err = ioutil.WriteFile(file_name, dat[:], 0644)
		if write_err != nil {
			return nil, write_err
		}		
		//delete the file
		os.Remove(file_name)
	}
	//open the encrypted file and read to byte array
	enc_dat, enc_dat_err := ioutil.ReadFile(enc_file_name)
	if enc_dat_err != nil {
		return nil, enc_dat_err
	}		
	//decrypt the byte array
	dat, dat_err := Decrypt(enc_dat, ear_key)
	if dat_err != nil {
		return nil, dat_err
	}
	//return a pointer to the decrypted byte array
	return dat[:], nil
}
func NewEncryptionKeyAsBase64() string {
	return base64.StdEncoding.EncodeToString(NewEncryptionKey()[:])
}
func DecodeBase64Key(key string) *[32]byte {
	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(err)
	}
	if len(data) != 32 {	
		panic( fmt.Sprintf("Key len is %s, must be 32", len(data)))
	}
	rval := [32]byte{}
	_, copy_err := io.ReadFull(bytes.NewReader(data),rval[:])
	if copy_err != nil {
		panic (err) 
	}
	return &rval
}

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}