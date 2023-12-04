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
	UID				string		`json:"uid"`  			// 아이디
	Accommodation	string 		`json:"accommodation"`	 // 숙소명
	Room			string 		`json:"room"`			// 방호수
	Price			int			`json:"price"`			// 가격
	Locate			string 		`json:"locate"`			// 위치
	Address			string 		`json:"address"`		// 주소
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

// 모든 숙소 구조체
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Accommodation
}

// 회원가입


// 관리자
// 숙소 등록 함수 Register - 관리자가 숙소를 등록(uid, auth, accommodation,room, price, locate 등) 
func (s *SmartContract) RegisterAccommodation(ctx contractapi.TransactionContextInterface, uid string, authorization string, accommodation string, room string, price int, locate string, address string  ) error {
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
		UID				: uid,
		Accommodation	: accommodation,
		Room            : room,
		Price           : price,
		Locate			: locate,
		Address 		: address,
		Available		: true,
		BookingList		: make([]string,0),
	}

	assetJSON, err = json.Marshal(accommodationInfo)
	if err != nil {
		return err
	}
	
	ctx.GetStub().PutState("ACMCOUNT", []byte(strconv.Itoa(acmcountINT))) // 숙소번호 원장에 저장, 다음 호출에 사용하기 위해
	
	return ctx.GetStub().PutState(acmID, assetJSON)
}
// 숙소 수정 함수 update - 관리자가 숙소 정보를 수정(가격, 이용여부)
// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAccommodation(ctx contractapi.TransactionContextInterface, uid string, aid string, price int, available bool) error {
	
	// 숙소 존재 여부 확인 
	accommodationInfo, err := s.GetAccommodation(ctx, aid)
	if err != nil {
		return err
	}
	// if !accommodationInfo {
	// 	return fmt.Errorf("the asset %s does not exist", aid)
	// }
	
	// 숙소 수정을 요청한 사람이 해당 숙소를 등록한 사람인지 확인
	if accommodationInfo.UID != uid {
		return fmt.Errorf("this accommdation does not available.")
	}

	// 요청한 수정 수행
	accommodationInfo.Price = price
	accommodationInfo.Available = available

	assetJSON, err := json.Marshal(accommodationInfo)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(aid, assetJSON)
}

// 숙소 삭제 함수 delete - 관리자가 숙소 정보를 삭제
func (s *SmartContract) DeleteAccommodation(ctx contractapi.TransactionContextInterface, uid string, aid string) error {
	// 숙소 존재 여부 확인 
	accommodationInfo, err := s.GetAccommodation(ctx, aid)
	if err != nil {
		return err
	}

	// 숙소 수정을 요청한 사람이 해당 숙소를 등록한 사람인지 확인
	if accommodationInfo.UID != uid {
		return fmt.Errorf("this accommdation does not available.")
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

// 관리자마다 각자 본인이 등록한 숙소 조회
func (s *SmartContract) GetMyAccommodations(ctx contractapi.TransactionContextInterface, uid string) ([]QueryResult, error) {
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
		accommodation.BookingList  =make([]string,0)
		// 조회한 uid가 존재할 경우
		if accommodation.UID == uid{
			queryResult := QueryResult{Key: queryResponse.Key, Record: accommodation}
			accommodations = append(accommodations, queryResult)
		}
		
	}
	
	return accommodations, nil
}

// 특정 숙소 조회  - 수정에 필요. 인자(숙소번호)
func (s *SmartContract) GetAccommodation(ctx contractapi.TransactionContextInterface, aid string) (*Accommodation, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	acmBytes, err := ctx.GetStub().GetState(aid)
	if err != nil {
		return nil, err
	}
	

	if acmBytes == nil {
		return nil, fmt.Errorf("%s does not exist", aid)
	}

	accommodation := new(Accommodation)
	_ = json.Unmarshal(acmBytes, accommodation)
	
	return accommodation, nil
}






// user 기능
// 모든 숙소 조회 .
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
		accommodation.BookingList  =make([]string,0)
		// 숙소 이용가능할 경우에만 조회
		if accommodation.Available != false{
			queryResult := QueryResult{Key: queryResponse.Key, Record: accommodation}
			accommodations = append(accommodations, queryResult)
		}
	}
	return accommodations, nil
}
// 사용자 
// 예약 함수 booking - 사용자가 숙소를 예약(uid, aid,  방문자명, 체크인/체크아웃 , 결제방식 )
func (s *SmartContract) Booking(ctx contractapi.TransactionContextInterface, uid string,  aid string, guest string,  checkin string,  checkout string, payment string ) error {
	// 현재 예약 가능한 숙소인지 확인
	accommodationInfo, err := s.GetAccommodation(ctx, aid)
	if err != nil {
		return err
	}
	if accommodationInfo.Available != true {
		return fmt.Errorf("this accommdation does not available.")
	}

	// 예약 아이디 자동 생성	
	bookingJSON, err := ctx.GetStub().GetState("BOOKCOUNT")
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if bookingJSON == nil {
		ctx.GetStub().PutState("BOOKCOUNT", []byte("0"))
	}

	bookingINT,_ := strconv.Atoi(string(bookingJSON))
	bookingINT += 1
	bookingID := "booking"+strconv.Itoa(bookingINT)
	
	// 예약 데이터 생성
	booking := Booking{
		BID				: bookingID,
		UID				: uid,
		Guest			: guest,
		AID             : aid, 
		Accommodation	: accommodationInfo.Accommodation,
		Room            : accommodationInfo.Room,
		CheckIn         : checkin,
		CheckOut		: checkout,
		Price			: accommodationInfo.Price,
		Payment			: payment,
	}
	
	bookingJSON, err = json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("error1: %v", err)
	}
	ctx.GetStub().PutState("BOOKCOUNT", []byte(strconv.Itoa(bookingINT)))

	// 예약한 숙소의 BookingList에 bid 추가
	// accommodationJSON, _ := ctx.GetStub().GetState(aid)
	
	// var accommodation Accommodation
	// err = json.Unmarshal(accommodationJSON, &accommodation)
	// if err != nil{
	// 	return fmt.Errorf("error2: %v", err)
	// }
	accommodationInfo.BookingList = append(accommodationInfo.BookingList, bookingID)
	// *** 예약시 다른사람이 예약 못하게 available 를 false로 바꿈 > 추후 여러 사람이 다른날짜에 예약가능하도록 바꿈
	accommodationInfo.Available = false

	accommodationJSON, err := json.Marshal(accommodationInfo)
	ctx.GetStub().PutState(aid, accommodationJSON)

	// 사용자 데이터에 있는 BookingList에 bookingID를 추가 
	userJSON, _ := ctx.GetStub().GetState(uid)
	var user User
	if userJSON != nil { // 사용자 정보 존재시
		err = json.Unmarshal(userJSON, &user) // json을 구조체로 변경하여 user에 넣기 
		if err != nil{
			return fmt.Errorf("error3: %v", err)
		}		
	}	
	user.UID = uid // User 객체 uid값 넣기
	user.BookingList = append(user.BookingList, bookingID)  // 예약 번호 넣기
	userJSON, err = json.Marshal(user) 
	ctx.GetStub().PutState(uid, userJSON)

	// 예약 정보 넣기
	return ctx.GetStub().PutState(bookingID, bookingJSON)
}


// 예약 정보 수정 함수 update - 예약 정보수정(방문자 이름만 가능) * 가능하다면 체크인/체크아웃날짜 변경
func (s *SmartContract) UpdateBooking(ctx contractapi.TransactionContextInterface, bid string, uid string, guest string  ) error {
	
	// booking정보 받아오기 
	assetJSON, err := ctx.GetStub().GetState(bid)
	var booking Booking
	if assetJSON != nil { // 사용자 정보 존재시
		err = json.Unmarshal(assetJSON, &booking) // json을 구조체로 변경하여 user에 넣기 
	if err != nil{
		return fmt.Errorf("error3: %v", err)
		}		
	}
	
	// 예약자인지 확인
	if booking.UID != uid{
		return fmt.Errorf("예약 ID %s 의 예약자가 아닙니다.", bid)
		
	}
	booking.Guest = guest // User 객체 uid값 넣기
	assetJSON, err = json.Marshal(booking) 
	return ctx.GetStub().PutState(bid, assetJSON)
}

// 예약 취소 함수 cancel/delete - 예약 취소
func (s *SmartContract) DeleteBooking(ctx contractapi.TransactionContextInterface, bid string, uid string, aid string) error {

	// 예약정보 해당 예약 삭제
	// 사용자 BookingList에 해당 예약 삭제
	userJSON, err := ctx.GetStub().GetState(uid)
	var user User
	if userJSON != nil { // 사용자 정보 존재시
		err = json.Unmarshal(userJSON, &user) // json을 구조체로 변경하여 user에 넣기 
		if err != nil{
			return fmt.Errorf("error3: %v", err)
		}		
	}	
	// user.UID = uid // User 객체 uid값 넣기
	//user의 bookinglist에서 해당 예약번호 삭제
	user.BookingList = DeleteBid(user.BookingList, bid)
	// user.BookingList = user.BookingList.remove(bid)
	userJSON, err = json.Marshal(user)
	ctx.GetStub().PutState(uid, userJSON)

	// 숙소 정보 BookingList에 해당 예약 삭제
	accommodationJSON, _ := ctx.GetStub().GetState(aid)
	var accommodation Accommodation
	err = json.Unmarshal(accommodationJSON, &accommodation)
	if err != nil{
		return err
	}
	accommodation.Available = true
	accommodation.BookingList = DeleteBid(accommodation.BookingList, bid)
	// accommodation.BookingList = accommodation.BookingList.remove(bid)
	accommodationJSON, err = json.Marshal(accommodation)
	ctx.GetStub().PutState(aid, accommodationJSON)

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

// 사용자별 예약 리스트 조회 ([]Booking, error)
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
		// return nil,fmt.Errorf("the asset %s does not exist", uid)
		return nil,fmt.Errorf("error1: %v", err)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		// return nil,err
		return nil,fmt.Errorf("error2: %v", err)
	}
	// bookinglist := make([]string, len(user.BookingList))
	// for i,v := range user.BookingList{
	// 	bookinglist[i] = v
	// }


	// 2-1. booking에 받아온 예약 정보 저장 포인터 배열넣기, null값일 경우 json input 오류남. 이부분 처리하기
	bookings := []Booking{}

	for _, bookingid := range user.BookingList {
		bookingJSON, err := ctx.GetStub().GetState(bookingid)
		

		mybooking := new(Booking)
		err = json.Unmarshal(bookingJSON, mybooking)

		if err != nil {
			return nil,fmt.Errorf("error3: %v", err)
			// return nil,err
		}

		bookings = append(bookings, *mybooking)
		}
		
	
	// 2. booking에 받아온 예약 정보 저장 포인터 배열넣기, null값일 경우 json input 오류남. 이부분 처리하기
	// var bookings []Booking
	// for _, bookingid := range user.BookingList{
	// 	bookingJSON, err := ctx.GetStub().GetState(bookingid)

	// 	mybooking := new(Booking)
	// 	err = json.Unmarshal(bookingJSON, &mybooking)
	// 	if err != nil {
	// 		return nil,fmt.Errorf("error3: %v", err)
	// 		// return nil,err
	// 	}

	// 	bookings = append(bookings, *mybooking)
	// }	

	return bookings, nil

}


// 취소 후 다시 예약리스트 만들기
func DeleteBid(bookinglist []string ,bid string) []string {
	var newbookinglist []string
	for idx, id := range bookinglist{
		if id == bid {
			newbookinglist = append(bookinglist[:idx], bookinglist[idx+1:]...)
		}
	}
    return newbookinglist
}


// docker rm -f $(docker ps -aq) 