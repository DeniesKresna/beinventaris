package Helpers

import (
	"os"
	"strings"
)

func GetAppDomain() string {
	return os.Getenv("DOMAIN")
}

func DeleteFile(filepath string) error {
	err := os.Remove("./" + filepath)
	if err != nil {
		return err
	}
	return nil
}

// return academy path. Asset/Academies/{additional_string}
func AcademyPath(additional string) string {
	return strings.Join([]string{"Assets/Academies/", additional}, "")
}

// return inventory path. Asset/Inventories/{additional_string}
func InventoryPath(additional string) string {
	return strings.Join([]string{"Assets/Inventories/", additional}, "")
}

// return history path. Asset/Histories/{additional_string}
func HistoryPath(additional string) string {
	return strings.Join([]string{"Assets/Histories/", additional}, "")
}

// return inventory documents path. Asset/Inventories/Documents/{additional_string}
func InventoryDocumentsPath(additional string) string {
	return strings.Join([]string{"Assets/Inventories/Documents/", additional}, "")
}
