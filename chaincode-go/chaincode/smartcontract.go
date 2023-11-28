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
	Accommodation	string 		`json:"accommodation"`	 // 숙소명
	Room			string 		`json:"room"`			// 방호수
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
	AID				string		`json:"aid"`    		 // 숙소번호
	Accommodation	string 		`json:"accommodation"`  // 숙소명
	Room			string 		`json:"room"`           // 호수	
	CheckIn			string 		`json:"checkIn"`    	//check-in
	CheckOut		string		`json:"checkOut"`   	//check-out
	Price			int			`json:"price"`			// 가격
	Payment			string 		`json:"payment"`    	//결제방법
}

type User struct{
	UID				string		`json:"uid"`			// uid 
	Authorization   string		`json:"authorization"`  // 권한 ( 관리자 : admin, 사용자 : user) 
	BookingList		[]string	`json:"BookingList"`    // 예약 내역 
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Accommodation
}

// 회원가입

// 모든 숙소 조회 
func (s *SmartContract) GetAllAccommodations(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	startKey := "accommodation1"
    endKey := "accommodation999"
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	accommodations := []QueryResult{}
	// var accommodations []*Accommodation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		
		if err != nil {
			return nil, err
		}
		accommodation := new(Accommodation)
		_ = json.Unmarshal(queryResponse.Value, accommodation)
		queryResult := QueryResult{Key: queryResponse.Key, Record: accommodation}
		accommodations = append(accommodations, queryResult)
		// var accommodation Accommodation
		// err = json.Unmarshal(queryResponse.Value, &accommodation)
		// if err != nil {
		// 	return nil, err
		// }
		// accommodations = append(accommodations, &accommodation)
	}

	return accommodations, nil
}
// 관리자
// 숙소 등록 함수 Register - 관리자가 숙소를 등록(uid, auth, accommodation,room, price, locate 등) 
func (s *SmartContract) RegisterAccommodation(ctx contractapi.TransactionContextInterface, uid string, authorization string, accommodation string, room string, price int, locate string  ) error {
	//원장에 숙소 정보를 저장하는 것

	// 접근제어 (AccessControl)
	if authorization != "admin"{
		return fmt. Errorf("Caller is not admin")
	}
	
	// 숙소 아이디 자동 생성	
	assetJSON, err := ctx.GetStub().GetState("ACMCOUNT")
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if assetJSON == nil {
		ctx.GetStub().PutState("ACMCOUNT", []byte("0"))
	}

	acmcountINT,_ := strconv.Atoi(string(assetJSON))
	acmcountINT += 1
	acmID := "accommodation"+strconv.Itoa(acmcountINT)
	// 현재 시간 생성
	// nowTime := time.Now()
	// unixTime := nowTime.Unix()
	
	// 숙소 정보 데이터 생성
	accommodationInfo := Accommodation{
		AID             : acmID,
		Accommodation	: accommodation,
		Room            : room,
		Price           : price,
		Locate			: locate,
		Available		: false,
	}

	assetJSON, err = json.Marshal(accommodationInfo)
	if err != nil {
		return err
	}
	
	ctx.GetStub().PutState("ACMCOUNT", []byte(strconv.Itoa(acmcountINT))) // 숙소번호 원장에 저장, 다음 호출에 사용하기 위해
	
	return ctx.GetStub().PutState(acmID, assetJSON)
}
// 숙소 수정 함수 update - 관리자가 숙소 정보를 수정(가격, 위치)
// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAccommodation(ctx contractapi.TransactionContextInterface, uid string, authorization string, aid string,accommodation string, room string,  price int, locate string) error {
	// 접근제어 (AccessControl) : 관리자만 접근 가능
	if authorization != "admin"{
		return fmt. Errorf("Caller is not admin")
	}
	
	// 숙소 존재 여부 확인 
	exists, err := s.AccommodationExists(ctx, aid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", aid)
	}

	// overwriting original asset with new asset
	accommodationInfo := Accommodation{
		AID             : aid,
		Accommodation	: accommodation,
		Room            : room,
		Price           : price,
		Locate			: locate,
		Available		: false,
	}

	assetJSON, err := json.Marshal(accommodationInfo)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(aid, assetJSON)
}

// 숙소 삭제 함수 delete - 관리자가 숙소 정보를 삭제
func (s *SmartContract) DeleteAccommodation(ctx contractapi.TransactionContextInterface, uid string, authorization string, aid string) error {
	// 접근제어 (AccessControl) : 관리자만 접근 가능
	if authorization != "admin"{
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
// 예약 함수 booking - 사용자가 숙소를 예약(uid, 방문자명,aid, 숙소, 호수, 체크인/체크아웃 , 결제방식 )
func (s *SmartContract) Booking(ctx contractapi.TransactionContextInterface, uid string, authorization string, guest string, aid string, myaccommodation string, room string, checkin string,  checkout string, price int, payment string ) error {
	// 권한 확인 : user만 예약 가능
	if authorization != "user"{
		return fmt. Errorf("Caller is not user")
	}
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
	ctx.GetStub().PutState("BOOKCOUNT", []byte(strconv.Itoa(bookingINT)))

	// 예약 데이터 생성
	booking := Booking{
		BID				: bookingID,
		UID				: uid,
		Guest			: guest,
		AID             : aid, 
		Accommodation	: myaccommodation,
		Room            : room,
		CheckIn         : checkin,
		CheckOut		: checkout,
		Price			: price,
		Payment			: payment,
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
	UserJSON, _ := ctx.GetStub().GetState(uid)
	var user User
	err = json.Unmarshal(UserJSON, &user)
	if err != nil{
		return err
	}
	user.BookingList = append(user.BookingList, bookingID)

	return ctx.GetStub().PutState(bookingID, assetJSON)
}


// 예약 정보 수정 함수 update - 예약 정보수정(방문자명, 결제방식) * 가능하다면 체크인/체크아웃날짜 변경
func (s *SmartContract) UpdateBooking(ctx contractapi.TransactionContextInterface, bid string, guest string, aid string, accommodation string, room string, checkin string, checkout string, payment string  ) error {
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
		Accommodation	: accommodation,
		Room            : room,
		CheckIn         : checkin,
		CheckOut		: checkout,
		Payment			: payment,
	}

	assetJSON, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(bid, assetJSON)


}

// 예약 취소 함수 cancel/delete - 예약 취소
func (s *SmartContract) DeleteBooking(ctx contractapi.TransactionContextInterface, bid string, uid string) error {

	// 예약정보 해당 예약 삭제
	// 사용자 BookingList에 해당 예약 삭제
	// 숙소 정보 BookingList에 해당 예약 삭제
	
	// 예약 확인
	exists, err := s.BookingExists(ctx, bid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", bid)
	}

    // 예약리스트에서 삭제 / 사용자 
	return ctx.GetStub().DelState(bid)


}
// 예약 존재 여부 확인 함수
func (s *SmartContract) BookingExists(ctx contractapi.TransactionContextInterface, bid string) (bool, error) {
	bookingJSON, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return bookingJSON != nil, nil
}

// 사용자별 예약 리스트 조회
func (s *SmartContract) AllBooking(ctx contractapi.TransactionContextInterface, uid string )  ([]Booking, error) {
	// user의 bookinglist 가져옴
	// booking에서 bookinglist에 있는 bid를 하나씩 넣어서 booking정보를 받아옴
	// 해당 사용자의 모든 예약 정보가 있는 booking포인터 배열 리턴
	
	
	// 1. 사용자정보에서 bookinglist가져오기
	userJSON, err := ctx.GetStub().GetState(uid)
	if err != nil {
		return nil,err
	}
	if userJSON ==nil{
		return nil,fmt.Errorf("the asset %s does not exist", uid)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil,err
	}
	var bookinglist = user.BookingList
	


	// 받아온 예약 정보 저장 포인터 배열넣기
	var bookings []Booking
	for _, bookingid := range bookinglist{
		bookingJSON, err := ctx.GetStub().GetState(bookingid)

		mybooking := new(Booking)
		err = json.Unmarshal(bookingJSON, &mybooking)
		if err != nil {
			return nil,err
		}

		bookings = append(bookings, *mybooking)
	}	

	return bookings, nil

}


