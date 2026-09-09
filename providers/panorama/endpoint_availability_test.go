package panorama

import (
	"encoding/json"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ffAvailabilityPayload = `{
	"result": {
		"error_code": 0,
		"error_msg": "",
		"response": {
			"ff_availability_provider": [
				{
					"provider": "sia",
					"program_name": "Singapore Airlines",
					"ff_availability": [
						{ "name": "KrisFlyer", "ff_code": "krisflyer" }
					],
					"status": "available",
					"schedule_id": 0,
					"sequence": 1,
					"carrier_codes": ["SQ"]
				}
			]
		}
	}
}`

func TestGetFFAvailability_ParseResponse(t *testing.T) {
	var resp GetFFAvailabilityResponse
	err := parseResponseBody([]byte(ffAvailabilityPayload), &resp)
	require.NoError(t, err)
	require.Equal(t, Success, resp.Result.ErrorCode)
	require.Len(t, resp.Result.Response.FFAvailabilityProvider, 1)

	provider := resp.Result.Response.FFAvailabilityProvider[0]
	assert.Equal(t, "sia", provider.Provider)
	assert.Equal(t, "Singapore Airlines", provider.ProgramName)
	assert.Equal(t, "available", provider.Status)
	assert.Equal(t, []string{"SQ"}, provider.CarrierCodes)
	require.Len(t, provider.FFAvailability, 1)
	assert.Equal(t, "KrisFlyer", provider.FFAvailability[0].Name)
	assert.Equal(t, "krisflyer", provider.FFAvailability[0].FFCode)
}

func TestGetFFAvailabilityRequest_Marshal(t *testing.T) {
	req := GetFFAvailabilityRequest{Signature: "sig-123"}
	body, err := json.Marshal(req)
	require.NoError(t, err)
	// The signature is the only body field (everything else travels in headers).
	assert.JSONEq(t, `{"signature":"sig-123"}`, string(body))
}

func TestGetFFAvailabilityRequest_Validate(t *testing.T) {
	require.NoError(t, (&GetFFAvailabilityRequest{Signature: "sig-123"}).Validate())
	require.Error(t, (&GetFFAvailabilityRequest{}).Validate())
}

const ssrAvailabilityPayload = `{
	"result": {
		"error_code": 0,
		"error_msg": "",
		"response": {
			"ssr_availability_provider": [
				{
					"provider": "sia",
					"ssr_availability": {
						"baggage": [
							{
								"journey_code": "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia",
								"segments": [
									{
										"carrier_code": "SQ",
										"carrier_number": "965",
										"origin": "CGK",
										"departure_date": "2021-08-15 19:00:00",
										"destination": "SIN",
										"arrival_date": "2021-08-15 21:50:00",
										"provider": "sia"
									}
								],
								"availability_type": "baggage",
								"ssrs": [
									{
										"ssr_code": "BAG05",
										"ssr_type": "BG",
										"image_url": "",
										"available_count": 1,
										"name": "5 Kg Baggage",
										"description": [],
										"ssr_ref": [],
										"service_charges": [
											{
												"charge_code": "fare",
												"charge_type": "SSR",
												"currency": "IDR",
												"pax_type": "ADT",
												"pax_count": 1,
												"amount": 29300,
												"foreign_currency": "IDR",
												"foreign_amount": 29300,
												"total": 29300,
												"sequence": 1
											}
										],
										"base_fare": 29300,
										"base_tax": 0,
										"base_nta": 29300,
										"base_commission": 0,
										"availability_type": "baggage",
										"journey_code": "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia"
									}
								]
							}
						],
						"meal": [
							{
								"journey_code": "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia",
								"availability_type": "meal",
								"ssrs": []
							}
						]
					},
					"status": "available",
					"schedule_id": 0,
					"sequence": 1
				}
			]
		}
	}
}`

func TestGetSSRAvailability_ParseResponse(t *testing.T) {
	var resp GetSSRAvailabilityResponse
	err := parseResponseBody([]byte(ssrAvailabilityPayload), &resp)
	require.NoError(t, err)
	require.Equal(t, Success, resp.Result.ErrorCode)
	require.Len(t, resp.Result.Response.SSRAvailabilityProvider, 1)

	provider := resp.Result.Response.SSRAvailabilityProvider[0]
	assert.Equal(t, "sia", provider.Provider)
	assert.Equal(t, "available", provider.Status)
	assert.Equal(t, 1, provider.Sequence)

	baggage := provider.SSRAvailability.Baggage
	require.Len(t, baggage, 1)
	require.Equal(t, "baggage", baggage[0].AvailabilityType)
	assert.Contains(t, baggage[0].JourneyCode, "SQ,965")
	require.Len(t, baggage[0].Segments, 1)
	assert.Equal(t, "SQ", baggage[0].Segments[0].CarrierCode)
	assert.Equal(t, "965", baggage[0].Segments[0].CarrierNumber)
	assert.Equal(t, "CGK", baggage[0].Segments[0].Origin)
	assert.Equal(t, "SIN", baggage[0].Segments[0].Destination)

	ssrs := baggage[0].SSRs
	require.Len(t, ssrs, 1)
	assert.Equal(t, "BAG05", ssrs[0].SSRCode)
	assert.Equal(t, "BG", ssrs[0].SSRType)
	assert.Equal(t, "5 Kg Baggage", ssrs[0].Name)
	assert.Equal(t, 1, ssrs[0].AvailableCount)
	assert.Equal(t, float64(29300), ssrs[0].BaseFare)
	assert.Equal(t, float64(29300), ssrs[0].BaseNta)
	require.Len(t, ssrs[0].ServiceCharges, 1)
	charge := ssrs[0].ServiceCharges[0]
	assert.Equal(t, "fare", charge.ChargeCode)
	assert.Equal(t, "SSR", charge.ChargeType)
	assert.Equal(t, "IDR", charge.Currency)
	assert.Equal(t, "ADT", charge.PaxType)
	assert.Equal(t, float64(29300), charge.Amount)
	assert.Equal(t, float64(29300), charge.Total)

	require.Len(t, provider.SSRAvailability.Meal, 1)
	assert.Equal(t, "meal", provider.SSRAvailability.Meal[0].AvailabilityType)
}

func TestGetSSRAvailabilityRequest_Validate(t *testing.T) {
	require.NoError(t, (&GetSSRAvailabilityRequest{Signature: "sig-123"}).Validate())
	require.Error(t, (&GetSSRAvailabilityRequest{}).Validate())
}

const seatAvailabilityPayload = `{
	"result": {
		"error_code": 0,
		"error_msg": "",
		"response": {
			"seat_availability_provider": [
				{
					"provider": "sia",
					"status": "available",
					"schedule_id": 0,
					"sequence": 1,
					"segments": [
						{
							"origin": "CGK",
							"origin_terminal": "3",
							"origin_utc": 7,
							"departure_date": "2021-08-15 19:00:00",
							"destination": "SIN",
							"destination_terminal": "3",
							"destination_utc": 8,
							"arrival_date": "2021-08-15 21:50:00",
							"segment_code": "SIA-SQ965-1",
							"segment_code2": "SQ-965-CGK-SIN",
							"route_code": "CGKSIN",
							"elapsed_time": "PT2H50M",
							"transit_count": 0,
							"transit_duration": "",
							"provider": "sia",
							"sequence": 1,
							"carrier_name": "Singapore Airlines",
							"carrier_code": "SQ",
							"carrier_type_code": "J",
							"carrier_number": "965",
							"operating_airline_code": "SQ",
							"seat_cabins": [
								{
									"cabin_name": "Economy",
									"cabin_code": "Y",
									"cabin_type_code": "P",
									"first_row_number": "31",
									"last_row_number": "47",
									"seat_rows": [
										{
											"row_number": "31",
											"seats": [
												{
													"seat_name": "31A",
													"seat_code": "A31",
													"column": "A",
													"availability": 0,
													"seat_type": "W",
													"seat_ref": "REF-31A",
													"service_charges": [
														{
															"charge_code": "fare",
															"charge_type": "SEAT",
															"currency": "IDR",
															"pax_type": "ADT",
															"pax_count": 1,
															"amount": 150000,
															"foreign_currency": "IDR",
															"foreign_amount": 150000,
															"total": 150000,
															"sequence": 1
														}
													],
													"base_fare": 150000,
													"base_tax": 0,
													"base_nta": 150000,
													"base_commission": 0,
													"description": [],
													"sequence": 1
												}
											],
											"sequence": 1
										}
									],
									"sequence": 1
								}
							]
						}
					]
				}
			]
		}
	}
}`

func TestGetSeatAvailability_ParseResponse(t *testing.T) {
	var resp GetSeatAvailabilityResponse
	err := parseResponseBody([]byte(seatAvailabilityPayload), &resp)
	require.NoError(t, err)
	require.Equal(t, Success, resp.Result.ErrorCode)
	require.Len(t, resp.Result.Response.SeatAvailabilityProvider, 1)

	provider := resp.Result.Response.SeatAvailabilityProvider[0]
	assert.Equal(t, "sia", provider.Provider)
	require.Len(t, provider.Segments, 1)

	segment := provider.Segments[0]
	assert.Equal(t, "CGK", segment.Origin)
	assert.Equal(t, "SIN", segment.Destination)
	assert.Equal(t, "SIA-SQ965-1", segment.SegmentCode)
	assert.Equal(t, "SQ", segment.CarrierCode)
	assert.Equal(t, "965", segment.CarrierNumber)
	require.Len(t, segment.SeatCabins, 1)

	cabin := segment.SeatCabins[0]
	assert.Equal(t, "Economy", cabin.CabinName)
	assert.Equal(t, "Y", cabin.CabinCode)
	assert.Equal(t, FlexibleString("31"), cabin.FirstRowNumber)
	assert.Equal(t, FlexibleString("47"), cabin.LastRowNumber)
	require.Len(t, cabin.SeatRows, 1)

	row := cabin.SeatRows[0]
	assert.Equal(t, FlexibleString("31"), row.RowNumber)
	require.Len(t, row.Seats, 1)

	seat := row.Seats[0]
	assert.Equal(t, "31A", seat.SeatName)
	assert.Equal(t, "A31", seat.SeatCode)
	assert.Equal(t, "A", seat.Column)
	assert.Equal(t, 0, seat.Availability)
	assert.Equal(t, "W", seat.SeatType)
	assert.Equal(t, "REF-31A", seat.SeatRef)
	assert.Equal(t, float64(150000), seat.BaseFare)
	require.Len(t, seat.ServiceCharges, 1)
	assert.Equal(t, "SEAT", seat.ServiceCharges[0].ChargeType)
	assert.Equal(t, float64(150000), seat.ServiceCharges[0].Amount)
}

func TestGetSeatAvailability_ParseNumericRowNumbers(t *testing.T) {
	// Some providers return row numbers as JSON numbers instead of strings.
	const payload = `{
		"result": {
			"error_code": 0,
			"error_msg": "",
			"response": {
				"seat_availability_provider": [
					{
						"provider": "sia",
						"status": "available",
						"schedule_id": 0,
						"sequence": 1,
						"segments": [
							{
								"origin": "CGK",
								"departure_date": "2021-08-15 19:00:00",
								"destination": "SIN",
								"arrival_date": "2021-08-15 21:50:00",
								"segment_code": "SIA-SQ965-1",
								"carrier_code": "SQ",
								"carrier_number": "965",
								"seat_cabins": [
									{
										"cabin_name": "Economy",
										"cabin_code": "Y",
										"first_row_number": 31,
										"last_row_number": 47,
										"seat_rows": [
											{
												"row_number": 31,
												"seats": [
													{
														"seat_name": "31A",
														"seat_code": "A31",
														"column": "A",
														"availability": 1,
														"seat_type": "N",
														"seat_ref": "REF-31A",
														"description": [],
														"sequence": 1
													}
												],
												"sequence": 1
											}
										],
										"sequence": 1
									}
								]
							}
						]
					}
				]
			}
		}
	}`

	var resp GetSeatAvailabilityResponse
	err := parseResponseBody([]byte(payload), &resp)
	require.NoError(t, err)
	require.Len(t, resp.Result.Response.SeatAvailabilityProvider, 1)

	cabin := resp.Result.Response.SeatAvailabilityProvider[0].Segments[0].SeatCabins[0]
	assert.Equal(t, FlexibleString("31"), cabin.FirstRowNumber)
	assert.Equal(t, FlexibleString("47"), cabin.LastRowNumber)
	assert.Equal(t, FlexibleString("31"), cabin.SeatRows[0].RowNumber)
	assert.Equal(t, 1, cabin.SeatRows[0].Seats[0].Availability)
}

func TestGetSeatAvailabilityRequest_Validate(t *testing.T) {
	require.NoError(t, (&GetSeatAvailabilityRequest{Signature: "sig-123"}).Validate())
	require.Error(t, (&GetSeatAvailabilityRequest{}).Validate())
}

const sellSSRsPayload = `{
	"result": {
		"error_code": 0,
		"error_msg": "",
		"response": {
			"total_price": 29300,
			"total_nta": 29300,
			"sell_ssr_provider": [
				{
					"provider": "sia",
					"journeys": [
						{
							"journey_code": "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia",
							"segments": [
								{
									"segment_code": "SIA-SQ965-1",
									"fare_code": "FARE-1",
									"carrier_code": "SQ",
									"carrier_number": "965",
									"origin": "CGK",
									"destination": "SIN",
									"departure_date": "2021-08-15 19:00:00",
									"arrival_date": "2021-08-15 21:50:00"
								}
							]
						}
					],
					"passengers": [
						{
							"passenger_number": 0,
							"first_name": "Adi",
							"last_name": "Laksono",
							"pax_type": "ADT"
						}
					]
				}
			]
		}
	}
}`

func TestSellSSRs_ParseResponse(t *testing.T) {
	var resp SellSSRsResponse
	err := parseResponseBody([]byte(sellSSRsPayload), &resp)
	require.NoError(t, err)
	require.Equal(t, Success, resp.Result.ErrorCode)
	assert.Equal(t, float64(29300), resp.Result.Response.TotalPrice)
	assert.Equal(t, float64(29300), resp.Result.Response.TotalNta)
	require.Len(t, resp.Result.Response.SellSSRProvider, 1)

	provider := resp.Result.Response.SellSSRProvider[0]
	assert.Equal(t, "sia", provider.Provider)
	require.Len(t, provider.Journeys, 1)
	assert.Contains(t, provider.Journeys[0].JourneyCode, "SQ,965")
	require.Len(t, provider.Journeys[0].Segments, 1)
	assert.Equal(t, "SIA-SQ965-1", provider.Journeys[0].Segments[0].SegmentCode)
	require.Len(t, provider.Passengers, 1)
	assert.Equal(t, 0, provider.Passengers[0].PassengerNumber)
	assert.Equal(t, "ADT", provider.Passengers[0].PaxType)
}

func TestSSRRequests_RedactSignature(t *testing.T) {
	const signature = "sentinel-signature-do-not-log"

	tests := map[string]any{
		"get SSR availability": GetSSRAvailabilityRequest{Signature: signature},
		"sell SSRs":            SellSSRsRequest{Signature: signature},
	}

	for name, request := range tests {
		t.Run(name, func(t *testing.T) {
			sanitized, err := logger.Sanitize(request)
			require.NoError(t, err)
			assert.Equal(t, "[REDACTED]", sanitized["signature"])

			logged, err := json.Marshal(sanitized)
			require.NoError(t, err)
			assert.NotContains(t, string(logged), signature)
		})
	}
}

func TestSellSSRsRequest_Marshal(t *testing.T) {
	req := SellSSRsRequest{
		Signature: "sig-123",
		SellSSRsRequest: []SellSSRProvider{
			{
				SellSSRs: []SellSSR{
					{
						JourneyCode:      "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia",
						AvailabilityType: "baggage",
						Passengers: []SellSSRPassenger{
							{PassengerNumber: 0, SSRCode: "BAG05"},
						},
					},
				},
				Provider: "sia",
			},
		},
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"signature": "sig-123",
		"sell_ssrs_request": [
			{
				"sell_ssrs": [
					{
						"journey_code": "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia",
						"availability_type": "baggage",
						"passengers": [
							{ "passenger_number": 0, "ssr_code": "BAG05" }
						]
					}
				],
				"provider": "sia"
			}
		]
	}`, string(body))
}

func TestSellSSRsRequest_Validate(t *testing.T) {
	valid := SellSSRsRequest{
		Signature: "sig-123",
		SellSSRsRequest: []SellSSRProvider{
			{
				SellSSRs: []SellSSR{
					{
						JourneyCode:      "J1",
						AvailabilityType: "baggage",
						Passengers:       []SellSSRPassenger{{PassengerNumber: 0, SSRCode: "BAG05"}},
					},
				},
				Provider: "sia",
			},
		},
	}
	require.NoError(t, valid.Validate())

	t.Run("missing signature", func(t *testing.T) {
		invalid := valid
		invalid.Signature = ""
		require.Error(t, invalid.Validate())
	})

	t.Run("empty sell_ssrs_request", func(t *testing.T) {
		invalid := valid
		invalid.SellSSRsRequest = nil
		require.Error(t, invalid.Validate())
	})

	t.Run("missing provider", func(t *testing.T) {
		invalid := valid
		invalid.SellSSRsRequest[0].Provider = ""
		require.Error(t, invalid.Validate())
	})

	t.Run("missing journey_code", func(t *testing.T) {
		invalid := valid
		invalid.SellSSRsRequest[0].SellSSRs[0].JourneyCode = ""
		require.Error(t, invalid.Validate())
	})

	t.Run("missing ssr_code", func(t *testing.T) {
		invalid := valid
		invalid.SellSSRsRequest[0].SellSSRs[0].Passengers[0].SSRCode = ""
		require.Error(t, invalid.Validate())
	})
}

const assignSeatsPayload = `{
	"result": {
		"error_code": 0,
		"error_msg": "",
		"response": {
			"total_price": 150000,
			"total_nta": 150000,
			"order_number": "ORD-123",
			"assign_seat_provider": [
				{
					"provider": "sia",
					"journeys": [
						{
							"journey_code": "SQ,965,CGK,2021-08-15 19:00:00,SIN,2021-08-15 21:50:00,sia",
							"segments": [
								{
									"segment_code": "SIA-SQ965-1",
									"carrier_code": "SQ",
									"carrier_number": "965",
									"origin": "CGK",
									"destination": "SIN",
									"departure_date": "2021-08-15 19:00:00",
									"arrival_date": "2021-08-15 21:50:00"
								}
							]
						}
					],
					"passengers": [
						{
							"passenger_number": 0,
							"first_name": "Adi",
							"last_name": "Laksono",
							"pax_type": "ADT",
							"seats": [
								{ "seat_code": "A31", "segment_code": "SIA-SQ965-1" }
							]
						}
					]
				}
			]
		}
	}
}`

func TestAssignSeats_ParseResponse(t *testing.T) {
	var resp AssignSeatsResponse
	err := parseResponseBody([]byte(assignSeatsPayload), &resp)
	require.NoError(t, err)
	require.Equal(t, Success, resp.Result.ErrorCode)
	assert.Equal(t, float64(150000), resp.Result.Response.TotalPrice)
	assert.Equal(t, float64(150000), resp.Result.Response.TotalNta)
	assert.Equal(t, "ORD-123", resp.Result.Response.OrderNumber)
	require.Len(t, resp.Result.Response.AssignSeatProvider, 1)

	provider := resp.Result.Response.AssignSeatProvider[0]
	assert.Equal(t, "sia", provider.Provider)
	require.Len(t, provider.Journeys, 1)
	require.Len(t, provider.Journeys[0].Segments, 1)
	assert.Equal(t, "SIA-SQ965-1", provider.Journeys[0].Segments[0].SegmentCode)
	require.Len(t, provider.Passengers, 1)
	assert.Equal(t, 0, provider.Passengers[0].PassengerNumber)
	assert.Len(t, provider.Passengers[0].Seats, 1)
}

func TestAssignSeatsRequest_Marshal(t *testing.T) {
	req := AssignSeatsRequest{
		Signature: "sig-123",
		AssignSeatsRequest: []AssignSeatProvider{
			{
				AssignSeats: []AssignSeat{
					{
						SegmentCode: "SIA-SQ965-1",
						Passengers:  []AssignSeatPassenger{{PassengerNumber: 0, SeatCode: "A31"}},
					},
				},
				Provider: "sia",
			},
		},
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"signature": "sig-123",
		"assign_seats_request": [
			{
				"assign_seats": [
					{
						"segment_code": "SIA-SQ965-1",
						"passengers": [
							{ "passenger_number": 0, "seat_code": "A31" }
						]
					}
				],
				"provider": "sia"
			}
		]
	}`, string(body))
}

func TestAssignSeatsRequest_Validate(t *testing.T) {
	valid := AssignSeatsRequest{
		Signature: "sig-123",
		AssignSeatsRequest: []AssignSeatProvider{
			{
				AssignSeats: []AssignSeat{
					{
						SegmentCode: "SIA-SQ965-1",
						Passengers:  []AssignSeatPassenger{{PassengerNumber: 0, SeatCode: "A31"}},
					},
				},
				Provider: "sia",
			},
		},
	}
	require.NoError(t, valid.Validate())

	t.Run("missing signature", func(t *testing.T) {
		invalid := valid
		invalid.Signature = ""
		require.Error(t, invalid.Validate())
	})

	t.Run("empty assign_seats_request", func(t *testing.T) {
		invalid := valid
		invalid.AssignSeatsRequest = nil
		require.Error(t, invalid.Validate())
	})

	t.Run("missing segment_code", func(t *testing.T) {
		invalid := valid
		invalid.AssignSeatsRequest[0].AssignSeats[0].SegmentCode = ""
		require.Error(t, invalid.Validate())
	})

	t.Run("missing seat_code", func(t *testing.T) {
		invalid := valid
		invalid.AssignSeatsRequest[0].AssignSeats[0].Passengers[0].SeatCode = ""
		require.Error(t, invalid.Validate())
	})
}
