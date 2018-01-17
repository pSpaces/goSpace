package protocol

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestCreatePointToPoint(t *testing.T) {
	config := createTestConfig()

	// Setup
	testPointToPoint := createTestPointToPoint(config)

	// Manually create PointToPoint
	actualName := "Name"
	actualIP := "192.168.0.0"
	actualPort := 8080
	actualAddress := strings.Join([]string{actualIP, strconv.Itoa(actualPort)}, ":")
	actualConfig := config
	actualPointToPoint := PointToPoint{actualName, actualAddress, actualConfig}

	pointToPointsEqual := reflect.DeepEqual(*testPointToPoint, actualPointToPoint)

	if !pointToPointsEqual {
		t.Errorf("CreatePointToPoint() gave %+v, should be %+v", *testPointToPoint, actualPointToPoint)
	}
}

func TestToString(t *testing.T) {
	// Setup
	config := createTestConfig()
	testPointToPoint := createTestPointToPoint(config)

	actualString := "Name: Name, address: 192.168.0.0:8080"

	testString := testPointToPoint.ToString()

	stringsEqual := reflect.DeepEqual(testString, actualString)

	if !stringsEqual {
		t.Errorf("ToString() on pointToPoint: %+v == %v, should be %v", testPointToPoint, testString, actualString)
	}
}

func TestGetAddress(t *testing.T) {
	// Setup
	config := createTestConfig()
	testPointToPoint := createTestPointToPoint(config)

	actualAddress := "192.168.0.0:8080"

	testAddress := testPointToPoint.GetAddress()

	addressesEqual := reflect.DeepEqual(testAddress, actualAddress)

	if !addressesEqual {
		t.Errorf("GetAddress() on pointToPoint: %+v == %v, should be %v", testPointToPoint, testAddress, actualAddress)
	}
}

func TestGetName(t *testing.T) {
	// Setup
	config := createTestConfig()
	testPointToPoint := createTestPointToPoint(config)

	actualName := "Name"

	testName := testPointToPoint.GetName()

	namesEqual := reflect.DeepEqual(testName, actualName)

	if !namesEqual {
		t.Errorf("GetName() on pointToPoint: %+v == %v, should be %v", testPointToPoint, testName, actualName)
	}
}

func createTestPointToPoint(config *tls.Config) *PointToPoint {
	testName := "Name"
	testIP := "192.168.0.0"
	testPort := "8080"

	return CreatePointToPoint(testName, testIP, testPort, config)
}

func createTestConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("../test/resources/certificates/server.pem", "../test/resources/certificates/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	// create a pool of trusted certs
	certPool := x509.NewCertPool()
	pemFile, _ := ioutil.ReadFile("../test/resources/certificates/client.pem")
	certPool.AppendCertsFromPEM(pemFile)

	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: certPool}
	config.Rand = rand.Reader

	return &config
}
