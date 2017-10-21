package protocol

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestCreatePointToPoint(t *testing.T) {
	// Setup
	testPointToPoint := createTestPointToPoint()

	// Manually create PointToPoint
	actualName := "Name"
	actualIP := "192.168.0.0"
	actualPort := 8080
	actualAddress := strings.Join([]string{actualIP, strconv.Itoa(actualPort)}, ":")
	actualPointToPoint := PointToPoint{actualName, actualAddress}

	pointToPointsEqual := reflect.DeepEqual(testPointToPoint, actualPointToPoint)

	if !pointToPointsEqual {
		t.Errorf("CreatePointToPoint() gave %+v, should be %+v", testPointToPoint, actualPointToPoint)
	}
}

func TestToString(t *testing.T) {
	// Setup
	testPointToPoint := createTestPointToPoint()

	actualString := "Name: Name, address: 192.168.0.0:8080"

	testString := testPointToPoint.ToString()

	stringsEqual := reflect.DeepEqual(testString, actualString)

	if !stringsEqual {
		t.Errorf("ToString() on pointToPoint: %+v == %v, should be %v", testPointToPoint, testString, actualString)
	}
}

func TestGetAddress(t *testing.T) {
	// Setup
	testPointToPoint := createTestPointToPoint()

	actualAddress := "192.168.0.0:8080"

	testAddress := testPointToPoint.GetAddress()

	addressesEqual := reflect.DeepEqual(testAddress, actualAddress)

	if !addressesEqual {
		t.Errorf("GetAddress() on pointToPoint: %+v == %v, should be %v", testPointToPoint, testAddress, actualAddress)
	}
}

func TestGetName(t *testing.T) {
	// Setup
	testPointToPoint := createTestPointToPoint()

	actualName := "Name"

	testName := testPointToPoint.GetName()

	namesEqual := reflect.DeepEqual(testName, actualName)

	if !namesEqual {
		t.Errorf("GetName() on pointToPoint: %+v == %v, should be %v", testPointToPoint, testName, actualName)
	}
}

func createTestPointToPoint() PointToPoint {
	testName := "Name"
	testIP := "192.168.0.0"
	testPort := "8080"

	return CreatePointToPoint(testName, testIP, testPort)
}
