package panorama

// FilterFeesByType mengembalikan fees yang sesuai dengan feeType
func FilterFeesByType(fees []BookingFee, feeType string) []BookingFee {
	var result []BookingFee
	for _, fee := range fees {
		if fee.FeeType == feeType {
			result = append(result, fee)
		}
	}
	return result
}
