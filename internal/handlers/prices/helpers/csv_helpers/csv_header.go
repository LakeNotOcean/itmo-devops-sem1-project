package csvhelpers

var CSVHeader = []string{"id", "name", "category", "price", "create_date"}

func CheckHeader(header []string) bool {
	if len(header) != len(CSVHeader) {
		return false
	}

	for i := range header {
		if header[i] != CSVHeader[i] {
			return false
		}
	}
	return true
}
