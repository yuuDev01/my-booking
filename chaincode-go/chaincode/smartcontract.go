package chaincode

import (
	"encoding/json"
	"fmt"
	// "time"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// 숙소 정보 stuct
type Accommodation struct{
	AID				string		`json:"aid"`    		 // 숙소번호
	Accmmodation	string 		`json:"accommodation"`	
	Room			string 		`json:"room"`
	Price			int			`json:"price"`			// 가격
	Locate			string 		`json:"locate"`			// 위치
	Available		bool		`json:"available"`  	// 예약가능여부
	BookingList		[]string	`json:"BookingList"`    // 예약 내역 
} 

// 예약 정보 struct
type Booking struct{
	BID 			string		`json:"bid"`			// 예약 번호
	UID				string		`json:"uid"`  			// 아이디
	Guest			string 		`json:"Guest"`    		// 방문자명
	Accommodation	string 		`json:"accommodation"`  // 숙소명
	Room			string 		`json:"room"`           // 호수	
	CheckIn			string 		`json:"checkIn"`    	//check-in
	CheckOut		string		`json:"checkOut"`   	//check-out
	Payment			string 		`json:"payment"`    	//결제방법
}

type User struct{
	UID				string		`json:"uid"`			// uid 
	BookingList		[]string	`json:"BookingList"`    // 예약 내역 
}

// 회원가입

// 숙소 조회 
func (s *SmartContract) GetAllAccommodations(ctx contractapi.TransactionContextInterface) ([]*AccommodationInfo, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var accommodations []*Accommodation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var accommodation Accommodation
		err = json.Unmarshal(queryResponse.Value, &accommodationInfo)
		if err != nil {
			return nil, err
		}
		accommodations = append(accommodations, &accommodation)
	}

	return accommodations, nil
}
// 관리자
// 숙소 등록 함수 Register - 관리자가 숙소를 등록(숙소이름, 호수, 가격, 위치 등) 
func (s *SmartContract) RegisterAccommodation(ctx contractapi.TransactionContextInterface, uid string, accmmodation string, room string, price int, locate string  ) error {
	//원장에 숙소 정보를 저장하는 것

	// 접근제어 (AccessControl)
	if uid != "admin"{
		return fmt. Errorf("Caller is not admin")
	}
	
	// 숙소 아이디 자동 생성	
	assetJSON, err := ctx.GetStub().GetState("ACCOMMODATION")
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if assetJSON == nil {
		ctx.GetStub().PutState("ACCOMMODATION", []byte("0"))
	}

	accommodationINT,_ := strconv.Atoi(string(assetJSON))
	accommodationINT += 1
	accommodationID := "accommodation"+strconv.Itoa(accommodationINT)
	
	// 현재 시간 생성
	// nowTime := time.Now()
	// unixTime := nowTime.Unix()
	
	// 숙소 정보 데이터 생성
	accommodation := Accommodation{
		AID             : accommodationID,
		Accmmodation	: accmmodation,
		Room            : room,
		Price           : price,
		Locate			: locate,
		Available		: false,
	}

	assetJSON, err = json.Marshal(accommodation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accommodationID, assetJSON)
}
// 숙소 수정 함수 update - 관리자가 숙소 정보를 수정(가격, 위치)
// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAccommodation(ctx contractapi.TransactionContextInterface, uid string, aid string, price int, locate string) error {
	// 접근제어 (AccessControl)
	if uid != "admin"{
		return fmt. Errorf("Caller is not admin")
	}
	
	exists, err := s.AssetExists(ctx, aid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	accommodation := Accommodation{
		AID             : aid,
		Accmmodation	: accmmodation,
		Room            : room,
		Price           : price,
		Locate			: locate,
		Available		: false,
	}

	assetJSON, err := json.Marshal(accommodation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// 숙소 삭제 함수 delete - 관리자가 숙소 정보를 삭제
func (s *SmartContract) DeleteAccommodationInfo(ctx contractapi.TransactionContextInterface, uid string, aid string) error {
	// 접근제어 (AccessControl)
	if uid != "admin"{
		return fmt. Errorf("Caller is not admin")
	}

	exists, err := s.AccommodationExists(ctx, aid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", aid)
	}

	return ctx.GetStub().DelState(aid)
}

// 숙소 존재 여부 확인 함수
func (s *SmartContract) AccommodationExists(ctx contractapi.TransactionContextInterface, aid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(aid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}


// 사용자 
// 예약 함수 booking - 사용자가 숙소를 예약(방문자명, 결제방식, 체크인/체크아웃 )
func (s *SmartContract) Booking(ctx contractapi.TransactionContextInterface, uid string, guest string, accmmodation string, checkin string,  checkout string, payment string ) error {
// 숙소 아이디 자동 생성	
	assetJSON, err := ctx.GetStub().GetState("BOOKING")
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if assetJSON == nil {
		ctx.GetStub().PutState("BOOKING", []byte("0"))
	}

	bookingINT,_ := strconv.Atoi(string(assetJSON))
	bookingINT += 1
	bookingID := "booking"+strconv.Itoa(bookingINT)
	

	// 예약 데이터 생성
	booking := Booking{
		BID				: bookingID,
		Guest			: guest,
		AID             : aid,
		Accmmodation	: accmmodation,
		Room            : room,
		CheckIn         : checkin,
		CheckOut		: checkout,
		Payment			: payment
	}

	assetJSON, err = json.Marshal(booking)
	if err != nil {
		return err
	}

	// 예약한 숙소의 BookingList에 bid 추가
	AccommodationJSON, _ := ctx.GetStub().GetState(aid)
	var accommodation Accommodation
	err = json.Unmarshal(AccommodationJSON, &accommodation)
	if err != nil{
		return err
	}
	accommodation.BookingList = append(accommodation.BookingList, bookingID)

	// 사용자 데이터에 있는 BookingList에 bookingID를 추가


	return ctx.GetStub().PutState(bookingID, assetJSON)


}
// 예약 정보 수정 함수 update - 예약 정보수정(방문자명, 결제방식) * 가능하다면 체크인/체크아웃날짜 변경
func (s *SmartContract) UpdateBooking(ctx contractapi.TransactionContextInterface, bid string, uid string, accmmodation string, room string, price int, locate string  ) error {
	exists, err := s.BookingExists(ctx, bid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", bid)
	}

	// overwriting original asset with new asset
	booking := Booking{
		BID				: bid,
		Guest			: guest,
		AID             : aid,
		Accmmodation	: accmmodation,
		Room            : room,
		CheckIn         : checkin,
		CheckOut		: checkout,
		Payment			: payment
	}

	assetJSON, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(bid, assetJSON)


}

// 예약 취소 함수 cancel/delete - 예약 취소
func (s *SmartContract) DeleteBooking(ctx contractapi.TransactionContextInterface, bid string) error {
// 예약 확인 
	exists, err := s.BookingExists(ctx, bid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", bid)
	}

	return ctx.GetStub().DelState(bid)


}
// 예약 존재 여부 확인 함수
func (s *SmartContract) BookingExists(ctx contractapi.TransactionContextInterface, bid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
// 사용자별 예약 리스트 조회
func (s *SmartContract) AllBooking(ctx contractapi.TransactionContextInterface, uid string )  ([]*Booking, error) error {
	// 사용자별 예약 리스트 존재 여부 확인
	assetJSON, err := ctx.GetStub().GetState(uid)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var bookings []*Booking
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var booking Booking
		err = json.Unmarshal(queryResponse.Value, &booking)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}

	return assets, nil

}


