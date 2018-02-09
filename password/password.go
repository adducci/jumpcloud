package password

import (
    "crypto/sha512"
    "encoding/base64"
)


//Takes in a string, str, and returns its equivalent byte array
func stringToBytes(str string) []byte {
    return []byte(str)
}


//Takes in a byte array, b, and returns its Base64 string equivalent
func bytesToBase64String(b []byte) string {
    str := base64.StdEncoding.EncodeToString(b)
    return str
}

/*
Takes in a password as a string and returns its Base64 encoded
string of the password hashed with SHA512
*/
func Encrypt(password string) string {
    encrypted := sha512.Sum512(stringToBytes(password))
    return bytesToBase64String(encrypted[:])   
}