package archivehelpers

import (
	"fmt"
	"net/http"
	"strconv"
)

// отправка готового архива с csv-файлом
func SendArchiveToClient(w http.ResponseWriter, archiveBytes []byte, dataFileName string) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", dataFileName+".zip"))
	w.Header().Set("Content-Length", strconv.Itoa(len(archiveBytes)))

	if _, err := w.Write(archiveBytes); err != nil {
		fmt.Printf("Failed to send ZIP file: %v\n", err)
	}
}
