package discovery

func getEtcdURL() ([]string, error) {
	s := []string{"http://192.168.254.249:2379", "http://192.168.254.248:2379", "http://192.168.254.247:2379"}
	return s, nil
}
