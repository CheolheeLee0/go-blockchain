package main

import (
	// 여러 가지 도구들을 가져와요
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 먼저 중요한 개념들을 알아볼까요?
// 1. 이더리움: 디지털 세계에서 사용하는 특별한 돈이에요. 마치 게임 속 금화 같은 거예요.
// 2. 스마트 컨트랙트: 컴퓨터가 자동으로 실행하는 약속이에요. 예를 들어, "사과 5개를 주면 바나나 3개를 줄게"라는 약속을 컴퓨터가 지켜주는 거예요.
// 3. 개인 키: 여러분만 아는 비밀 암호예요. 이걸로 디지털 세계에서 여러분임을 증명해요.
// 4. 주소: 디지털 세계에서 여러분의 집 주소 같은 거예요. 다른 사람들이 여기로 이더리움을 보낼 수 있어요.
// 5. 가스: 이더리움 세계에서 무언가를 하려면 내야 하는 요금이에요. 버스를 탈 때 요금을 내는 것처럼요.
// 이제 코드를 간단히 설명해 드릴게요:

// 우리가 만들 특별한 장난감(컨트랙트)의 설명서예요
const contractABI = `[{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
const contractBytecode = `608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea264697066735822122061dbfcf1f879e56edf094df20c09ce7cf9831da35005f49c9bc1d06b04e4798964736f6c63430008070033`

// SimpleStorage 컨트랙트를 위한 구조체를 정의해요
type SimpleStorage struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

func main() {
	// 이더리움이라는 특별한 놀이터에 연결해요
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
	if err != nil {
		log.Fatal(err) // 연결이 안 되면 슬퍼해요 :(
	}

	// 비밀 열쇠(개인 키)로 우리의 특별한 장난감을 만들 준비를 해요
	privateKey, err := crypto.HexToECDSA("your-private-key")
	if err != nil {
		log.Fatal(err) // 비밀 열쇠가 없으면 슬퍼해요 :(
	}

	// 비밀 열쇠로 우리의 주소를 만들어요
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("주소 만들기가 실패했어요 :(")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 우리 차례(논스)가 언제인지 알아봐요
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err) // 차례를 모르면 슬퍼해요 :(
	}

	// 장난감을 만드는 데 필요한 요금(가스 가격)을 알아봐요
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err) // 요금을 모르면 슬퍼해요 :(
	}

	// 장난감을 만들 때 필요한 정보들을 준비해요
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // 1은 메인 이터(메인넷)를 의미해요
	if err != nil {
		log.Fatal(err) // 준비가 안 되면 슬퍼해요 :(
	}
	auth.Nonce = big.NewInt(int64(nonce)) // 우리 차례를 기억해요
	auth.Value = big.NewInt(0)            // 장난감을 만드는 데 돈이 들지 않아요
	auth.GasLimit = uint64(300000)        // 장난감을 만드는 데 필요한 최대 요금이에요
	auth.GasPrice = gasPrice              // 요금을 설정해요

	// contractABI 문자열을 abi.ABI 타입으로 변환
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal("ABI 파싱 실패:", err)
	}

	// contractBytecode 문자열을 바이트 슬라이스로 변환
	bytecode := common.FromHex(contractBytecode)

	// 드디어 우리의 특별한 장난감(컨트랙트)을 만들어요!
	address, tx, _, err := bind.DeployContract(auth, parsedABI, bytecode, client)
	if err != nil {
		log.Fatal(err) // 장난감 만들기가 실패하면 슬퍼해요 :(
	}

	// 장난감이 잘 만들어졌다고 알려줘요
	fmt.Printf("새로운 장난감이 여기 생겼어요: %s\n", address.Hex())
	fmt.Printf("장난감을 만든 영수증 번호예요: %s\n", tx.Hash().Hex())

	// 컨트랙트 인스턴스 생성
	simpleStorage := &SimpleStorage{
		address: address,
		abi:     parsedABI,
		client:  client,
	}

	// 값 저장하기
	err = simpleStorage.Store(auth, big.NewInt(42))
	if err != nil {
		log.Fatal("값 저장하기 실패:", err)
	}
	fmt.Println("42를 저장했어요!")

	// 값 조회하기
	value, err := simpleStorage.Retrieve(nil)
	if err != nil {
		log.Fatal("값 조회하기 실패:", err)
	}
	fmt.Printf("저장된 값이에요: %s\n", value.String())
}

// Store 함수는 컨트랙트에 값을 저장해요
func (s *SimpleStorage) Store(auth *bind.TransactOpts, value *big.Int) error {
	input, err := s.abi.Pack("store", value)
	if err != nil {
		return err
	}

	tx := types.NewTransaction(
		auth.Nonce.Uint64(),
		s.address,
		big.NewInt(0),
		auth.GasLimit,
		auth.GasPrice,
		input,
	)

	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}

	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), s.client, signedTx)
	return err
}

// Retrieve 함수는 컨트랙트에서 값을 조회해요
func (s *SimpleStorage) Retrieve(opts *bind.CallOpts) (*big.Int, error) {
	input, err := s.abi.Pack("retrieve")
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &s.address,
		Data: input,
	}

	result, err := s.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var value *big.Int
	err = s.abi.UnpackIntoInterface(&value, "retrieve", result)
	return value, err
}
