package common

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sea-project/go-logger"
	"log"
)

func TestProcess() {
	privateKey, err := crypto.HexToECDSA("c0dc14ef79b0088e0c33f8d2f7b8a512ecf096c4639a9756e7e8e7b191e08f0c")
	if err != nil {
		logger.Error(err)
	}
	// 私钥生成公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	//
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("publicKeyBytes:", publicKeyBytes)

	data := []byte("1000")
	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0xd63087bea9f1800eed943829fc1d61e7869764805baa3259078c1caf3d4f5a48
	// 私钥签名
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		logger.Error(err)
	}

	fmt.Println(hexutil.Encode(signature)) // 0xa4fd0da0c27edf07c85af6ec3c0a9f7a987e86f08371b7d842a4c75a549694db0d2351a5bafceaa917486909938c95472997b4c5c4ba0ef6b783540c964629c701

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		logger.Error(err)
	}

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches) // true
	ecdsaPub, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		logger.Error(err)
	}
	// 公钥地址
	ethAddress := crypto.PubkeyToAddress(*ecdsaPub).String()
	fmt.Println("ethAddress:", ethAddress)

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		logger.Error(err)
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches) // true

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified) // true
}
/*
func ServiceVerify(wallet_addr, msg, sign string) bool {
	data := []byte(msg)
	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0xd63087bea9f1800eed943829fc1d61e7869764805baa3259078c1caf3d4f5a48
	// signature := []byte(sign)
	signature, err := hexutil.Decode(sign)
	if err != nil {
		logger.Error(err)
	}

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		logger.Error(err)
	}

	ecdsaPub, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		logger.Error(err)
	}

	ethAddress := crypto.PubkeyToAddress(*ecdsaPub)
	fmt.Println("ethAddress:", ethAddress.String())

	rawAddress, err := hexutil.Decode(wallet_addr)
	if err != nil {
		logger.Error(err)
	}
	matches := bytes.Equal(ethAddress.Bytes(), rawAddress)
	return matches
}

 */

func ServiceVerify(address string, nonce string, signature string) (bool, error) {
	// personal_sign format with data
	messagePrefix := "\x19Ethereum Signed Message:\n"
	lenMsg := fmt.Sprint(len(nonce))
	data := []byte(messagePrefix + lenMsg + nonce)
	// message hash
	fmt.Println(string(data))
	hash := crypto.Keccak256Hash(data)
	fmt.Println("data Hash hex:", hash.Hex())

	// signature format
	sign, err := hexutil.Decode(signature)
	if err != nil {
		logger.Error(err)
		return false, err 
	}
	// format where V is 0 or 1
	if sign[64] == 27 || sign[64] == 28 {
		sign[64] -= 27
	}

	// public key recv from signature
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sign)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	ecdsaPub, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	// address from public key which recv from signature
	signAddress := crypto.PubkeyToAddress(*ecdsaPub)
	fmt.Println("address from signature:", signAddress.String())

	// sender address match address recv from signature
	rawAddress, err := hexutil.Decode(address)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	matches := bytes.Equal(signAddress.Bytes(), rawAddress)
	return matches, nil
}

func Test() {
	a := crypto.VerifySignature([]byte("459A0136E53B122e902f8bC9f13154c53C43aBF5"),
		[]byte("d63087bea9f1800eed943829fc1d61e7869764805baa3259078c1caf3d4f5a48"),
		[]byte("1000"))
	fmt.Println(a)
}

func PrivateKeyToAddr() {
	// 私钥转钱包地址
	privateKey, err := crypto.HexToECDSA("私钥")
	if err != nil {
		logger.Error(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priv := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(priv, address)
	/*
	 * 公钥：在secp256k1规范下，由私钥和规范中指定的生成点计算出的坐标(x, y)
	 *      非压缩格式公钥： [前缀0x04] + x + y (65字节)
	 *      压缩格式公钥：[前缀0x02或0x03] + x ，其中前缀取决于 y 的符号
	 */
}


func Tests() {
	privateKey, err := crypto.HexToECDSA("4e2aef520ef5a5c0737676c2555f1b4bc3ecb015a2e976898c71da3776e06e27")
	if err != nil {
		logger.Error(err)
	}
	hash := crypto.Keccak256Hash([]byte("Example `personal_sign` message"))

	a, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil  {
		fmt.Println(err)
	}
	fmt.Println(string(a))

}
