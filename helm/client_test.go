package helm

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	New("https://helm.astronomer.io")
}

//func TestDownloadChart(t *testing.T) {
//	chart, err := DownloadChart("https://helm.astronomer.io/index.yaml", "airflow", "")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	fmt.Println(chart)
//}
