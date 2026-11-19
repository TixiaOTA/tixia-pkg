package train

import (
 "context"
 "strings"
 "time"
 "github.com/TixiaOTA/tixia-pkg/http/client"
 otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)
const assignSeatsTimeout = 30 * time.Second
const getBookingTimeout = 30 * time.Second

// AssignSeats signs in independently of the search/commit session.
// Seat changes must not be retried automatically after ambiguous failures.
func (c *clientImpl) AssignSeats(ctx context.Context, req AssignSeatsRequest) (*AssignSeatsResponse, *ProviderError) {
 ctx, span, _ := otelpkg.StartSpanWithID(ctx,"provider","orbis/train","AssignSeats")
 defer span.End()
 if err:=req.Validate(); err!=nil {return nil,newProviderError(ErrorInvalidInput,err.Error())}
 signin,err:=c.signIn(ctx);if err!=nil{return nil,err}
 headers:=map[string]string{"Action":"assign_seats","Signature":signin.Result.Response.Signature}
 timeout:=assignSeatsTimeout
 noRetry:=client.RetryConfig{MaxRetries:0}
 httpResp,sendErr:=c.sendRequest(ctx,req,bookingTrainPath,&headers,client.RequestOptions{Timeout:&timeout,Retry:&noRetry})
 if sendErr!=nil{return nil,newProviderError(ErrorInternal,sendErr.Error())}
 var resp AssignSeatsResponse
 if err:=decodeResponse(httpResp,&resp);err!=nil{return nil,err}
 if len(resp.Result.Response)!=len(req.ProviderBookings){return nil,newProviderError(ErrorInternal,"Orbis assign seats result count mismatch")}
 requested:=map[int]bool{}
 for _,booking:=range req.ProviderBookings {requested[booking.Sequence]=true}
 for _,booking:=range resp.Result.Response {
  if !requested[booking.Sequence] || strings.TrimSpace(booking.Status)=="" {return nil,newProviderError(ErrorInternal,"Orbis assign seats returned unexpected booking")}
  delete(requested,booking.Sequence)
 }
 return &resp,nil
}

// GetBooking uses a fresh sign-in and is safe for reconciliation after a timeout.
func (c *clientImpl) GetBooking(ctx context.Context, req GetBookingRequest) (*GetBookingResponse, *ProviderError) {
 ctx,span,_:=otelpkg.StartSpanWithID(ctx,"provider","orbis/train","GetBooking")
 defer span.End()
 if err:=req.Validate();err!=nil{return nil,newProviderError(ErrorInvalidInput,err.Error())}
 signin,err:=c.signIn(ctx);if err!=nil{return nil,err}
 headers:=map[string]string{"Action":"get_booking","Signature":signin.Result.Response.Signature}
 timeout:=getBookingTimeout
 httpResp,sendErr:=c.sendRequest(ctx,req,bookingTrainPath,&headers,client.RequestOptions{Timeout:&timeout})
 if sendErr!=nil{return nil,newProviderError(ErrorInternal,sendErr.Error())}
 var resp GetBookingResponse
 if err:=decodeResponse(httpResp,&resp);err!=nil{return nil,err}
 if strings.TrimSpace(resp.Result.Response.OrderNumber)=="" {return nil,newProviderError(ErrorInternal,"Orbis get booking returned empty order number")}
 if resp.Result.Response.OrderNumber!=req.OrderNumber {return nil,newProviderError(ErrorInternal,"Orbis get booking returned different order number")}
 return &resp,nil
}
