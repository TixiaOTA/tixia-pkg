package train

import (
 "encoding/json"
 "testing"
)

func TestAssignSeatsRequestValidation(t *testing.T) {
 good:=AssignSeatsRequest{OrderNumber:"TN.123",ProviderBookings:[]AssignSeatProviderBooking{{Provider:ProviderKAI,Sequence:1,Journeys:[]AssignSeatJourney{{Sequence:1,Seats:[]AssignSeat{{CabinCode:"LUX-1",SeatRow:2,SeatColumn:"A",PassengerSequence:0},{CabinCode:"LUX-1",SeatRow:2,SeatColumn:"B",PassengerSequence:1}}}}}}}
 if err:=good.Validate();err!=nil{t.Fatalf("valid request rejected: %v",err)}
 duplicate:=good
 duplicate.ProviderBookings=append([]AssignSeatProviderBooking(nil),good.ProviderBookings...)
 duplicate.ProviderBookings[0].Journeys=append([]AssignSeatJourney(nil),good.ProviderBookings[0].Journeys...)
 duplicate.ProviderBookings[0].Journeys[0].Seats=append([]AssignSeat(nil),good.ProviderBookings[0].Journeys[0].Seats...)
 duplicate.ProviderBookings[0].Journeys[0].Seats[1].SeatColumn="A"
 if err:=duplicate.Validate();err==nil{t.Fatal("duplicate seat accepted")}
 negative:=good
 negative.ProviderBookings=append([]AssignSeatProviderBooking(nil),good.ProviderBookings...)
 negative.ProviderBookings[0].Journeys=append([]AssignSeatJourney(nil),good.ProviderBookings[0].Journeys...)
 negative.ProviderBookings[0].Journeys[0].Seats=append([]AssignSeat(nil),good.ProviderBookings[0].Journeys[0].Seats...)
 negative.ProviderBookings[0].Journeys[0].Seats[0].PassengerSequence=-1
 if err:=negative.Validate();err==nil{t.Fatal("negative passenger accepted")}
}

func TestAssignSeatsResponseUsesOneBasedSequences(t *testing.T){
 body:=[]byte(`{"result":{"error_code":0,"response":[{"sequence":1,"status":"SUCCESS","journeys":[{"sequence":1,"seats":[{"seat":"LUX-1,2A","passenger_sequence":1,"seat_code":13404194}]}]}]}}`)
 var response AssignSeatsResponse
 if err:=json.Unmarshal(body,&response);err!=nil{t.Fatal(err)}
 if response.Result.Response[0].Journeys[0].Seats[0].PassengerSequence!=1{t.Fatal("provider response passenger_sequence must remain one-based in SDK")}
 if response.Result.Response[0].Journeys[0].Seats[0].SeatCode.String()!="13404194"{t.Fatal("seat code not decoded")}
}
