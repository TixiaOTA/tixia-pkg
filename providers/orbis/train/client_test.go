package train

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/providers/panorama"
)

func TestSearchOWAndRT(t *testing.T) {
	for _, tc := range []struct {
		name string
		direction string
		journeys []SearchJourney
		scheduleCount int
	}{
		{"one-way", DirectionOW, []SearchJourney{{Origin:"GMR",Destination:"BD",DepartureDate:"2027-01-15"}},1},
		{"round-trip", DirectionRT, []SearchJourney{{Origin:"GMR",Destination:"BD",DepartureDate:"2027-01-15"},{Origin:"BD",Destination:"GMR",DepartureDate:"2027-01-17"}},2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/booking/train" {t.Errorf("wrong method/path: %s %s",r.Method,r.URL.Path)}
				if r.Header.Get("Action") != "search" || r.Header.Get("Signature") != "secret" {t.Error("wrong provider headers")}
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {t.Fatal(err)}
				if _,found := body["signature"];found {t.Error("signature leaked into body")}
				var journeys []SearchJourney
				if err:=json.Unmarshal(body["journey_list"],&journeys);err!=nil || len(journeys)!=len(tc.journeys) {t.Errorf("journey count invalid: %v",err)}
				var direction string
				_ = json.Unmarshal(body["direction"],&direction)
				if direction != tc.direction {t.Errorf("wrong direction: %s",direction)}
				w.Header().Set("Content-Type","application/json")
				_,_ = w.Write([]byte(`{"result":{"error_code":0,"error_msg":"","response":{"schedules":[{"journey_list":{"origin":"GMR","destination":"BD","departure_date":"2027-01-15"},"journeys":[{"journey_code":"J1","fares":[{"fare_code":"F1","service_charge_summary":[{"pax_type":"ADT","total_price":176500,"total_commission":-7350,"service_charges":[{"charge_type":"FARE","amount":84500}]}]}]}]}` + func()string {if tc.scheduleCount==2{return `,{"journey_list":{"origin":"BD","destination":"GMR","departure_date":"2027-01-17"},"journeys":[]}`};return ""}() + `]}}}`))
			}))
			defer server.Close()
			c,err:=New(Config{BaseURL:server.URL});if err!=nil{t.Fatal(err)}
			resp,providerErr:=c.Search(context.Background(),SearchRequest{Signature:"secret",Provider:"kai",Direction:tc.direction,Adult:2,JourneyList:tc.journeys})
			if providerErr!=nil {t.Fatal(providerErr)}
			if len(resp.Result.Response.Schedules)!=tc.scheduleCount {t.Fatalf("schedules: %d",len(resp.Result.Response.Schedules))}
			fare:=resp.Result.Response.Schedules[0].Journeys[0].Fares[0]
			if fare.FareCode!="F1" || fare.ServiceChargeSummary[0].TotalCommission!=-7350 {t.Errorf("fare decoding failed: %+v",fare)}
		})
	}
}

func TestSearchProviderError(t *testing.T) {
	server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,_ *http.Request) {
		_,_=w.Write([]byte(`{"result":{"error_code":4020,"error_msg":"no schedules","response":""}}`))
	}))
	defer server.Close()
	c,err:=New(Config{BaseURL:server.URL});if err!=nil{t.Fatal(err)}
	resp,providerErr:=c.Search(context.Background(),SearchRequest{Signature:"secret",Provider:"kai",Direction:"OW",Adult:1,JourneyList:[]SearchJourney{{Origin:"GMR",Destination:"BD",DepartureDate:"2027-01-15"}}})
	if resp!=nil || providerErr==nil || providerErr.Result.ErrorCode!=panorama.NoSchedules {t.Fatalf("expected provider error: %#v, %#v",resp,providerErr)}
}

func TestSearchValidation(t *testing.T) {
	valid:=SearchRequest{Signature:"secret",Provider:"kai",Direction:"OW",Adult:1,JourneyList:[]SearchJourney{{Origin:"GMR",Destination:"BD",DepartureDate:"2027-01-15"}}}
	if err:=valid.Validate();err!=nil {t.Fatal(err)}
	valid.Signature=""
	if err:=valid.Validate();err==nil {t.Fatal("missing signature should fail")}
	valid.Signature="secret"
	valid.Direction="RT"
	if err:=valid.Validate();err==nil {t.Fatal("RT requires two legs")}
}
