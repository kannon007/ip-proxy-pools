package support

import (
	"fmt"
	"gopkg.in/resty.v1"
	"strconv"
	"testing"
)

func TestGetIPs(test *testing.T) {
	for _, v := range GetIPs() {
		fmt.Println(v.Ip + ":" + v.Port + "  " + strconv.FormatInt(v.Ms, 10) + "ms")
	}

}

func TestServer(test *testing.T) {
	//Server()
}

func TestClient(test *testing.T) {

	client := resty.SetProxy("http://125.108.68.49:9000")
	client.SetHTTPMode()
	resp, err := client.R().Get("http://2019.ip138.com/ic.asp")

	// explore response object
	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())
}
