package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/mingi3442/mingicoin/utils"
)

var w *wallet

const (
	fileName string = "mingicoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

func hasWalletFile() bool {
	_, err := os.Stat(fileName) //mingicoin.wallet file이 있는지 확인 후 err를 return
	return !os.IsNotExist(err)  // return 받은 error가 file이 없을 경우 생기는 err일 경우 true return 그러나 not(!)이니 파일이 없을 경우 true return
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key) // privateKey를 받아서 byte slice를 return
	utils.HandleErr(err)
	os.WriteFile(fileName, bytes, 0644) // return받은 값을 파일로 작성 (0644는 read&write)
}
func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func encodeBigInt(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInt(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, Wallet().privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInt(r.Bytes(), s.Bytes())
}

func resoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	r, s, err := resoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := resoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(), //privateKey를 생성할 때 사용한 curve를 사용
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() { // 만약 wallet이 있다면 true를 반환
			w.privateKey = restoreKey() //wallet이 있다면 키를 복원
		} else {
			key := createPrivKey() //walletdl 없다면 privateKey생성
			persistKey(key)        //key를 byte로 변환해서 파일에 저장
			w.privateKey = key     //
		}
		w.Address = aFromK(w.privateKey) //if문이 끝난 후 wallet의 address값을 넣는다
	}
	return w
}
