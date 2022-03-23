package helper

func MapRDGetKeybyValue(request_data map[string][]string, key string) (string, error){
	for k, v := range request_data {
		if(k == key) {
			return v[0], nil
		}
    }
	return "", ErrNotFound
}