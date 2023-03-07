//Auto generated with MetaApi https://github.com/exyzzy/metaapi
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
    "strconv"
)

import	"time"


var testDb *sql.DB
var configdb map[string]interface{}
const testDbName = "testtodo"

// ======= helpers

//assumes a configlocaldb.json file as:
//{
//    "Host": "localhost",
//    "Port": "5432",
//    "User": "dbname",
//    "Pass": "dbname",
//    "Name": "dbname",
//    "SSLMode": "disable"
//}
func loadConfig() {
	fmt.Println("  loadConfig")
	file, err := os.Open("configlocaldb.json")
	if err != nil {
		log.Panicln("Cannot open configlocaldb file", err.Error())
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configdb)
	if err != nil {
		log.Panicln("Cannot get local configurationdb from file", err.Error())
	}
}

func createDb(db *sql.DB, dbName string, owner string) (err error) {
	ss := fmt.Sprintf("CREATE DATABASE %s OWNER %s", dbName, owner)
	fmt.Println("  " + ss)
	_, err = db.Exec(ss)
	return
}

func setTzDb(db *sql.DB) (err error) {
	ss := fmt.Sprintf("SET TIME ZONE UTC")
	fmt.Println("  " + ss)
	_, err = db.Exec(ss)
	return
}

func dropDb(db *sql.DB, dbName string) (err error) {
	ss := fmt.Sprintf("DROP DATABASE %s", dbName)
	fmt.Println("  " + ss)
	_, err = db.Exec(ss)
	return
}

func rowExists(db *sql.DB, query string, args ...interface{}) (exists bool, err error) {
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	fmt.Println("  " + query)
	err = db.QueryRow(query, args...).Scan(&exists)
	return
}

func tableExists(db *sql.DB, table string) (valid bool, err error) {

	valid, err = rowExists(db, "SELECT 1 FROM pg_tables WHERE tablename = $1", table)
	return
}

func initTestDb() (err error) {
	loadConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"sslmode=%s", configdb["Host"], configdb["Port"], configdb["User"], configdb["Pass"], configdb["SSLMode"])
	testDb, err = sql.Open("postgres", psqlInfo)
	return
}

func TestMain(m *testing.M) {
	//test setup
	err := initTestDb()
	if err != nil {
		log.Panicln("cannot initTestDb ", err.Error())
	}

	err = createDb(testDb, testDbName, configdb["User"].(string))
	if err != nil {
		log.Panicln("cannot CreateDb ", err.Error())
	}

	err = setTzDb(testDb)
	if err != nil {
		log.Panicln("cannot setTzDb ", err.Error())
	}

	//run tests
	exitVal := m.Run()

	//test teardown
	err = dropDb(testDb, testDbName)
	if err != nil {
		log.Panicln("cannot DropDb ", err.Error())
	}
	os.Exit(exitVal)
}

type compareType func(interface{}, interface{}) bool

func noCompare(result, expect interface{}) bool {
	fmt.Printf("  noCompare: %v, %v -  %T, %T \n", result, expect, result, expect)
	return (true)
}

func defaultCompare(result, expect interface{}) bool {
	fmt.Printf("  defaultCompare: %v, %v -  %T, %T \n", result, expect, result, expect)
	return (result == expect)
}

func jsonCompare(result, expect interface{}) bool {
	fmt.Printf("  jsonCompare: %v, %v -  %T, %T \n", result, expect, result, expect)

	//json fields can be any order after db return, so read into map[string]interface and look up
	resultMap := make(map[string]interface{})
	expectMap := make(map[string]interface{})

	if reflect.TypeOf(result).String() == "sql.NullString" {
		err := json.Unmarshal([]byte(result.(sql.NullString).String), &resultMap)
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal([]byte(expect.(sql.NullString).String), &expectMap)
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := json.Unmarshal([]byte(result.(string)), &resultMap)
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal([]byte(expect.(string)), &expectMap)
		if err != nil {
			log.Panic(err)
		}
	}

	for k, v := range expectMap {
		if v != resultMap[k] {
			fmt.Printf("Key: %v, Result: %v, Expect: %v", k, resultMap[k], v)
			return false
		}
	}
	return true


	for k, v := range expectMap {
		if v != resultMap[k] {
			fmt.Printf("Key: %v, Result: %v, Expect: %v", k, resultMap[k], v)
			return false
		}
	}
	return true
}

func stringCompare(result, expect interface{}) bool {

	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Panic(err)
	}
	expectJson, err := json.Marshal(expect)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("  stringCompare: %v, %v -  %T, %T \n", string(resultJson), string(expectJson), result, expect)
	return (strings.TrimSpace(string(resultJson)) == strings.TrimSpace(string(expectJson)))
}

//psgl truncs reals at 6 digits
func realCompare(result, expect interface{}) bool {

	fmt.Printf("  realCompare: %v, %v -  %T, %T \n", result, expect, result, expect)

	var resultStr string
	var expectStr string
	if reflect.TypeOf(result).String() == "sql.NullFloat64" {
		resultStr = strconv.FormatFloat(result.(sql.NullFloat64).Float64, 'f', 6, 32)
		expectStr = strconv.FormatFloat(expect.(sql.NullFloat64).Float64, 'f', 6, 32)
	} else {
		resultStr = strconv.FormatFloat(float64(result.(float32)), 'f', 6, 32)
		expectStr = strconv.FormatFloat(float64(expect.(float32)), 'f', 6, 32)
	}
	return (resultStr == expectStr)
}

//iterate through each field of struct and apply the compare function to each field based on compareType map
func equalField(result, expect interface{}, compMap map[string]compareType) error {

	u := reflect.ValueOf(expect)
	v := reflect.ValueOf(result)
	typeOfS := u.Type()

	for i := 0; i < u.NumField(); i++ {

		if !(compMap[typeOfS.Field(i).Name])(v.Field(i).Interface(), u.Field(i).Interface()) {
			return fmt.Errorf("Field: %s, Result: %v, Expect: %v", typeOfS.Field(i).Name, v.Field(i).Interface(), u.Field(i).Interface())
		}
	}
	return nil
}


//table specific 


const community_resultstableName = "community_results"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testCommunity_result = [2]Community_result{  Id: 1, Uuid: "MA3Kp161j7zoqNTO", ProductUuid: "upvsWyorW2bKO8wi", MetricGraphs: "6vghIXrBRd8MXnTX",  Id: 2, Uuid: "jcWikvLxBHUBBYcF", ProductUuid: "v2d9sVX5qYRFCWyp", MetricGraphs: "ofyCoTJnbEDrJ3o2" }

var updateCommunity_result = Community_result Id: 1, Uuid: "gyM56NYjNua5f1Nx", ProductUuid: "QQlGmvETA2Se4Ebc", MetricGraphs: "3PlQJ4NN3ATy7e7p"

//compare functions
var compareCommunity_results = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"ProductUuid": defaultCompare,
	"MetricGraphs": defaultCompare,

}

func reverseCommunity_results(community_results []Community_result) (result []Community_result) {

	for i := len(community_results) - 1; i >= 0; i-- {
		result = append(result, community_results[i])
	}
	return
}

// ======= tests: Community_result =======

func TestCreateTableCommunity_results(t *testing.T) {
	fmt.Println("==CreateTableCommunity_results")
	err := CreateTableCommunity_results(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableCommunity_results " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableCommunity_results")
	}
	exists, err := tableExists(testDb, "community_results")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(community_results) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateCommunity_result(t *testing.T) {
	fmt.Println("==CreateCommunity_result")
	result, err := testCommunity_result[0].CreateCommunity_result(testDb)
	if err != nil {
		t.Errorf("cannot CreateCommunity_result " + err.Error())
	} else {
		fmt.Println("  Done: CreateCommunity_result")
	}
	err = equalField(result, testCommunity_result[0], compareCommunity_results)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveCommunity_result(t *testing.T) {
	fmt.Println("==RetrieveCommunity_result")
	result, err := testCommunity_result[0].RetrieveCommunity_result(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveCommunity_result " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveCommunity_result")
	}
	err = equalField(result, testCommunity_result[0], compareCommunity_results)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllCommunity_results(t *testing.T) {
	fmt.Println("==RetrieveAllCommunity_results")
	_, err := testCommunity_result[1].CreateCommunity_result(testDb)
	if err != nil {
		t.Errorf("cannot CreateCommunity_result " + err.Error())
	} else {
		fmt.Println("  Done: CreateCommunity_result")
	}
	result, err := RetrieveAllCommunity_results(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllCommunity_results " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllCommunity_results")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseCommunity_results(testCommunity_result[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareCommunity_results)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateCommunity_result(t *testing.T) {
	fmt.Println("==UpdateCommunity_result")
	result, err := updateCommunity_result.UpdateCommunity_result(testDb)
	if err != nil {
		t.Errorf("cannot UpdateCommunity_result " + err.Error())
	} else {
		fmt.Println("  Done: UpdateCommunity_result")
	}
	err = equalField(result, updateCommunity_result, compareCommunity_results)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const product_reviewstableName = "product_reviews"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testProduct_review = [2]Product_review{  Id: 1, Uuid: "TbxB9vqEGvmCgwgD", Title: "md2eLOalIHSieGc6", Author: "NO8fvoYb3Cgyc66B", CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), HeroImage: "GIUtqRnhk9UbCfyN", Summary: "mQshZgCTJZxcnGHO", Rating: "NwJ9WOcNqwjzgJ2U", Score: -1.8100513540586696, BodyHtml: "PZ3AfsDlAtaEbdPF", UserId: 1, ProductUuid: "J515PeMQrqqGO5JC", ExpectedTst: sql.NullInt32{207197061, true}, ExpectedDeep: sql.NullInt32{139552820, true}, ExpectedRem: sql.NullInt32{1877740286, true}, ExpectedOnset: sql.NullInt32{688553233, true}, ExpectedWakefulness: sql.NullInt32{1958669997, true}, ExpectedEfficiency: sql.NullInt32{1990196944, true}, ExpectedAccuracy: sql.NullInt32{321494717, true}, ReviewedResults: sql.NullString{"iCvWtamJy1uQd7jz", true}, Validated: true, StartDate: time.Now().UTC().Truncate(time.Microsecond), EndDate: time.Now().UTC().Truncate(time.Microsecond), IsPreview: true,  Id: 2, Uuid: "u2v4JtvmcIGQY6xg", Title: "rBrcxISPeV0CVIlu", Author: "HNeRsczVvtjNLx3Q", CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), HeroImage: "NPYAf8dgHjyL4Klg", Summary: "nJPzqlgzhaZDZujT", Rating: "NdNEFWPWFCkJdKWB", Score: -1.0110692387315474, BodyHtml: "OEsX4jbo24MUcDbo", UserId: 2, ProductUuid: "9dhOXls2HyGced7A", ExpectedTst: sql.NullInt32{304103460, true}, ExpectedDeep: sql.NullInt32{1819012522, true}, ExpectedRem: sql.NullInt32{1793670605, true}, ExpectedOnset: sql.NullInt32{1130764826, true}, ExpectedWakefulness: sql.NullInt32{172123921, true}, ExpectedEfficiency: sql.NullInt32{1892531180, true}, ExpectedAccuracy: sql.NullInt32{1888952258, true}, ReviewedResults: sql.NullString{"1OPpPcBT52sYTHJ2", true}, Validated: true, StartDate: time.Now().UTC().Truncate(time.Microsecond), EndDate: time.Now().UTC().Truncate(time.Microsecond), IsPreview: false }

var updateProduct_review = Product_review Id: 1, Uuid: "Rvy5KUXJ3HAGZ0OV", Title: "yDU0YvFWFevyoP65", Author: "L1jANx3gipXuXKOq", CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), HeroImage: "j7H0vfU15Jp8LOhz", Summary: "vEyIKoOrtBDcZK4Y", Rating: "c4eMoSYLbaNP2RXu", Score: 0.7239705943478476, BodyHtml: "MbMwZJ5eJTuWVJcI", UserId: 1, ProductUuid: "rMAEc5SklSNQmJnT", ExpectedTst: sql.NullInt32{401195116, true}, ExpectedDeep: sql.NullInt32{167281835, true}, ExpectedRem: sql.NullInt32{1740527693, true}, ExpectedOnset: sql.NullInt32{685876249, true}, ExpectedWakefulness: sql.NullInt32{1172040498, true}, ExpectedEfficiency: sql.NullInt32{2082901529, true}, ExpectedAccuracy: sql.NullInt32{193111067, true}, ReviewedResults: sql.NullString{"6i8m9lXUVHsOHmlt", true}, Validated: false, StartDate: time.Now().UTC().Truncate(time.Microsecond), EndDate: time.Now().UTC().Truncate(time.Microsecond), IsPreview: false

//compare functions
var compareProduct_reviews = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Title": defaultCompare,
	"Author": defaultCompare,
	"CreationDate": stringCompare,
	"LastUpdated": stringCompare,
	"HeroImage": defaultCompare,
	"Summary": defaultCompare,
	"Rating": defaultCompare,
	"Score": defaultCompare,
	"BodyHtml": defaultCompare,
	"UserId": defaultCompare,
	"ProductUuid": defaultCompare,
	"ExpectedTst": defaultCompare,
	"ExpectedDeep": defaultCompare,
	"ExpectedRem": defaultCompare,
	"ExpectedOnset": defaultCompare,
	"ExpectedWakefulness": defaultCompare,
	"ExpectedEfficiency": defaultCompare,
	"ExpectedAccuracy": defaultCompare,
	"ReviewedResults": defaultCompare,
	"Validated": defaultCompare,
	"StartDate": stringCompare,
	"EndDate": stringCompare,
	"IsPreview": defaultCompare,

}

func reverseProduct_reviews(product_reviews []Product_review) (result []Product_review) {

	for i := len(product_reviews) - 1; i >= 0; i-- {
		result = append(result, product_reviews[i])
	}
	return
}

// ======= tests: Product_review =======

func TestCreateTableProduct_reviews(t *testing.T) {
	fmt.Println("==CreateTableProduct_reviews")
	err := CreateTableProduct_reviews(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableProduct_reviews " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableProduct_reviews")
	}
	exists, err := tableExists(testDb, "product_reviews")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(product_reviews) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateProduct_review(t *testing.T) {
	fmt.Println("==CreateProduct_review")
	result, err := testProduct_review[0].CreateProduct_review(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_review " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_review")
	}
	err = equalField(result, testProduct_review[0], compareProduct_reviews)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveProduct_review(t *testing.T) {
	fmt.Println("==RetrieveProduct_review")
	result, err := testProduct_review[0].RetrieveProduct_review(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveProduct_review " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveProduct_review")
	}
	err = equalField(result, testProduct_review[0], compareProduct_reviews)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllProduct_reviews(t *testing.T) {
	fmt.Println("==RetrieveAllProduct_reviews")
	_, err := testProduct_review[1].CreateProduct_review(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_review " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_review")
	}
	result, err := RetrieveAllProduct_reviews(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_reviews " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllProduct_reviews")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseProduct_reviews(testProduct_review[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareProduct_reviews)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateProduct_review(t *testing.T) {
	fmt.Println("==UpdateProduct_review")
	result, err := updateProduct_review.UpdateProduct_review(testDb)
	if err != nil {
		t.Errorf("cannot UpdateProduct_review " + err.Error())
	} else {
		fmt.Println("  Done: UpdateProduct_review")
	}
	err = equalField(result, updateProduct_review, compareProduct_reviews)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const product_listingstableName = "product_listings"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testProduct_listing = [2]Product_listing{  Id: 1, Uuid: "sj8tFaclwuipdlrn", VendorUuid: "354L70wRlcU1Rr78", CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), PurchaseLink: "JRTR8lCLjpXUqw3i", Price: 2.122724586423068, ProductUuid: "w7SovTINCUuk6Xeq", IsManufacterer: true,  Id: 2, Uuid: "FtaEvqYxupK8FiJ4", VendorUuid: "BvZRRvAR1APHyEXh", CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), PurchaseLink: "zr9UWymSccF7n1vL", Price: 0.683031955115693, ProductUuid: "yTkts4mvVEXDN5A5", IsManufacterer: true }

var updateProduct_listing = Product_listing Id: 1, Uuid: "KKUBuSq2fjNNJtKF", VendorUuid: "3vIOlfZXp9FrhPbH", CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), PurchaseLink: "4A5V5iBOlR84ZiTv", Price: -0.716582511915298, ProductUuid: "vAnn58TUYkoNKRnR", IsManufacterer: false

//compare functions
var compareProduct_listings = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"VendorUuid": defaultCompare,
	"CreationDate": stringCompare,
	"LastUpdated": stringCompare,
	"PurchaseLink": defaultCompare,
	"Price": defaultCompare,
	"ProductUuid": defaultCompare,
	"IsManufacterer": defaultCompare,

}

func reverseProduct_listings(product_listings []Product_listing) (result []Product_listing) {

	for i := len(product_listings) - 1; i >= 0; i-- {
		result = append(result, product_listings[i])
	}
	return
}

// ======= tests: Product_listing =======

func TestCreateTableProduct_listings(t *testing.T) {
	fmt.Println("==CreateTableProduct_listings")
	err := CreateTableProduct_listings(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableProduct_listings " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableProduct_listings")
	}
	exists, err := tableExists(testDb, "product_listings")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(product_listings) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateProduct_listing(t *testing.T) {
	fmt.Println("==CreateProduct_listing")
	result, err := testProduct_listing[0].CreateProduct_listing(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_listing " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_listing")
	}
	err = equalField(result, testProduct_listing[0], compareProduct_listings)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveProduct_listing(t *testing.T) {
	fmt.Println("==RetrieveProduct_listing")
	result, err := testProduct_listing[0].RetrieveProduct_listing(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveProduct_listing " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveProduct_listing")
	}
	err = equalField(result, testProduct_listing[0], compareProduct_listings)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllProduct_listings(t *testing.T) {
	fmt.Println("==RetrieveAllProduct_listings")
	_, err := testProduct_listing[1].CreateProduct_listing(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_listing " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_listing")
	}
	result, err := RetrieveAllProduct_listings(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_listings " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllProduct_listings")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseProduct_listings(testProduct_listing[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareProduct_listings)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateProduct_listing(t *testing.T) {
	fmt.Println("==UpdateProduct_listing")
	result, err := updateProduct_listing.UpdateProduct_listing(testDb)
	if err != nil {
		t.Errorf("cannot UpdateProduct_listing " + err.Error())
	} else {
		fmt.Println("  Done: UpdateProduct_listing")
	}
	err = equalField(result, updateProduct_listing, compareProduct_listings)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const photo_urlstableName = "photo_urls"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testPhoto_url = [2]Photo_url{  Id: 1, Uuid: "IKgNy2ixeMcfLHl3", Url: "vzFTUl0k2rOQIi12", Title: "Yppiujc6vmQitbon", AltText: "WIYe4N5XJkVGb9ks", Keywords: "iqwWxLpPAP6OQ9SD", UserId: 1, CreationDate: time.Now().UTC().Truncate(time.Microsecond),  Id: 2, Uuid: "0axp5zyr1RyPFFNC", Url: "Yy3xOdYvMVNA8Jlf", Title: "9HrftxvdcbQrZLkb", AltText: "emzF9TySSV8KlBJp", Keywords: "3NGUZmD3KsIyBSUM", UserId: 2, CreationDate: time.Now().UTC().Truncate(time.Microsecond) }

var updatePhoto_url = Photo_url Id: 1, Uuid: "MeVlDX9Qn6RYSUMk", Url: "OnZoAYRNGufGRsP1", Title: "gd694JkKtMob4tGZ", AltText: "fjZA6oHekDI0cfGE", Keywords: "hEnFnCSZkss3jlNB", UserId: 1, CreationDate: time.Now().UTC().Truncate(time.Microsecond)

//compare functions
var comparePhoto_urls = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Url": defaultCompare,
	"Title": defaultCompare,
	"AltText": defaultCompare,
	"Keywords": defaultCompare,
	"UserId": defaultCompare,
	"CreationDate": stringCompare,

}

func reversePhoto_urls(photo_urls []Photo_url) (result []Photo_url) {

	for i := len(photo_urls) - 1; i >= 0; i-- {
		result = append(result, photo_urls[i])
	}
	return
}

// ======= tests: Photo_url =======

func TestCreateTablePhoto_urls(t *testing.T) {
	fmt.Println("==CreateTablePhoto_urls")
	err := CreateTablePhoto_urls(testDb)
	if err != nil {
		t.Errorf("cannot CreateTablePhoto_urls " + err.Error())
	} else {
		fmt.Println("  Done: CreateTablePhoto_urls")
	}
	exists, err := tableExists(testDb, "photo_urls")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(photo_urls) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreatePhoto_url(t *testing.T) {
	fmt.Println("==CreatePhoto_url")
	result, err := testPhoto_url[0].CreatePhoto_url(testDb)
	if err != nil {
		t.Errorf("cannot CreatePhoto_url " + err.Error())
	} else {
		fmt.Println("  Done: CreatePhoto_url")
	}
	err = equalField(result, testPhoto_url[0], comparePhoto_urls)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrievePhoto_url(t *testing.T) {
	fmt.Println("==RetrievePhoto_url")
	result, err := testPhoto_url[0].RetrievePhoto_url(testDb)
	if err != nil {
		t.Errorf("cannot RetrievePhoto_url " + err.Error())
	} else {
		fmt.Println("  Done: RetrievePhoto_url")
	}
	err = equalField(result, testPhoto_url[0], comparePhoto_urls)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllPhoto_urls(t *testing.T) {
	fmt.Println("==RetrieveAllPhoto_urls")
	_, err := testPhoto_url[1].CreatePhoto_url(testDb)
	if err != nil {
		t.Errorf("cannot CreatePhoto_url " + err.Error())
	} else {
		fmt.Println("  Done: CreatePhoto_url")
	}
	result, err := RetrieveAllPhoto_urls(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllPhoto_urls " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllPhoto_urls")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reversePhoto_urls(testPhoto_url[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], comparePhoto_urls)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdatePhoto_url(t *testing.T) {
	fmt.Println("==UpdatePhoto_url")
	result, err := updatePhoto_url.UpdatePhoto_url(testDb)
	if err != nil {
		t.Errorf("cannot UpdatePhoto_url " + err.Error())
	} else {
		fmt.Println("  Done: UpdatePhoto_url")
	}
	err = equalField(result, updatePhoto_url, comparePhoto_urls)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const vendorstableName = "vendors"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testVendor = [2]Vendor{  Id: 1, Uuid: "N1JjkciyiW30xQb6", Name: "DrX5xc2IuxuqsLQf", Link: sql.NullString{"FCPa7XXstRzbuDE9", true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), AffiliateId: sql.NullString{"Vj6y3H2XUpAt7hn2", true},  Id: 2, Uuid: "bMkp38uONswtB9PC", Name: "7Ft2v9mM8MAdeBVr", Link: sql.NullString{"a6qIJstcwctnWUQg", true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), AffiliateId: sql.NullString{"5xV6kGaEmNu0uLMa", true} }

var updateVendor = Vendor Id: 1, Uuid: "Bcjmk201U0Alis5Z", Name: "q0U4oOrrMweqnAdv", Link: sql.NullString{"HNxO34F1rnoA38AL", true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), AffiliateId: sql.NullString{"UOXJE0eN9vqu3VsM", true}

//compare functions
var compareVendors = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Name": defaultCompare,
	"Link": defaultCompare,
	"CreationDate": stringCompare,
	"LastUpdated": stringCompare,
	"AffiliateId": defaultCompare,

}

func reverseVendors(vendors []Vendor) (result []Vendor) {

	for i := len(vendors) - 1; i >= 0; i-- {
		result = append(result, vendors[i])
	}
	return
}

// ======= tests: Vendor =======

func TestCreateTableVendors(t *testing.T) {
	fmt.Println("==CreateTableVendors")
	err := CreateTableVendors(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableVendors " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableVendors")
	}
	exists, err := tableExists(testDb, "vendors")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(vendors) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateVendor(t *testing.T) {
	fmt.Println("==CreateVendor")
	result, err := testVendor[0].CreateVendor(testDb)
	if err != nil {
		t.Errorf("cannot CreateVendor " + err.Error())
	} else {
		fmt.Println("  Done: CreateVendor")
	}
	err = equalField(result, testVendor[0], compareVendors)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveVendor(t *testing.T) {
	fmt.Println("==RetrieveVendor")
	result, err := testVendor[0].RetrieveVendor(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveVendor " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveVendor")
	}
	err = equalField(result, testVendor[0], compareVendors)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllVendors(t *testing.T) {
	fmt.Println("==RetrieveAllVendors")
	_, err := testVendor[1].CreateVendor(testDb)
	if err != nil {
		t.Errorf("cannot CreateVendor " + err.Error())
	} else {
		fmt.Println("  Done: CreateVendor")
	}
	result, err := RetrieveAllVendors(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllVendors " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllVendors")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseVendors(testVendor[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareVendors)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateVendor(t *testing.T) {
	fmt.Println("==UpdateVendor")
	result, err := updateVendor.UpdateVendor(testDb)
	if err != nil {
		t.Errorf("cannot UpdateVendor " + err.Error())
	} else {
		fmt.Println("  Done: UpdateVendor")
	}
	err = equalField(result, updateVendor, compareVendors)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const solutionstableName = "solutions"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testSolution = [2]Solution{  Id: 1, Uuid: "NHZXdGnt4fzdpSYV", Name: "hbANM7trWkVMStKX", ProductUuid: sql.NullString{"zuPRFdscq30mkqMd", true}, TechniqueUuid: "WmxVcow6YAYdH8hT", UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), IsPublic: true, Validated: false, LastUpdated: time.Now().UTC().Truncate(time.Microsecond), Archived: false, Description: sql.NullString{"kyTjRouPgvRJ32Wz", true},  Id: 2, Uuid: "cfJUdTAhQXsEwDvD", Name: "ZDYqx9EUzBhTtWeu", ProductUuid: sql.NullString{"7WcARlTMEAkkeW6v", true}, TechniqueUuid: "LDKPdK6VvbtUmcRb", UserId: 2, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), IsPublic: true, Validated: false, LastUpdated: time.Now().UTC().Truncate(time.Microsecond), Archived: false, Description: sql.NullString{"2gBhqS4UqRl1kObH", true} }

var updateSolution = Solution Id: 1, Uuid: "PGDkCK87QT3XJoRm", Name: "6uCvZDXNUMp6gRGy", ProductUuid: sql.NullString{"barObwv9lVDUNjYa", true}, TechniqueUuid: "RNc7JUteam7ilAmf", UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), IsPublic: true, Validated: false, LastUpdated: time.Now().UTC().Truncate(time.Microsecond), Archived: false, Description: sql.NullString{"0jhKYLBUG7qrKbcl", true}

//compare functions
var compareSolutions = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Name": defaultCompare,
	"ProductUuid": defaultCompare,
	"TechniqueUuid": defaultCompare,
	"UserId": defaultCompare,
	"CreatedAt": stringCompare,
	"IsPublic": defaultCompare,
	"Validated": defaultCompare,
	"LastUpdated": stringCompare,
	"Archived": defaultCompare,
	"Description": defaultCompare,

}

func reverseSolutions(solutions []Solution) (result []Solution) {

	for i := len(solutions) - 1; i >= 0; i-- {
		result = append(result, solutions[i])
	}
	return
}

// ======= tests: Solution =======

func TestCreateTableSolutions(t *testing.T) {
	fmt.Println("==CreateTableSolutions")
	err := CreateTableSolutions(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableSolutions " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableSolutions")
	}
	exists, err := tableExists(testDb, "solutions")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(solutions) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateSolution(t *testing.T) {
	fmt.Println("==CreateSolution")
	result, err := testSolution[0].CreateSolution(testDb)
	if err != nil {
		t.Errorf("cannot CreateSolution " + err.Error())
	} else {
		fmt.Println("  Done: CreateSolution")
	}
	err = equalField(result, testSolution[0], compareSolutions)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveSolution(t *testing.T) {
	fmt.Println("==RetrieveSolution")
	result, err := testSolution[0].RetrieveSolution(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveSolution " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveSolution")
	}
	err = equalField(result, testSolution[0], compareSolutions)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllSolutions(t *testing.T) {
	fmt.Println("==RetrieveAllSolutions")
	_, err := testSolution[1].CreateSolution(testDb)
	if err != nil {
		t.Errorf("cannot CreateSolution " + err.Error())
	} else {
		fmt.Println("  Done: CreateSolution")
	}
	result, err := RetrieveAllSolutions(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSolutions " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllSolutions")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseSolutions(testSolution[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareSolutions)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateSolution(t *testing.T) {
	fmt.Println("==UpdateSolution")
	result, err := updateSolution.UpdateSolution(testDb)
	if err != nil {
		t.Errorf("cannot UpdateSolution " + err.Error())
	} else {
		fmt.Println("  Done: UpdateSolution")
	}
	err = equalField(result, updateSolution, compareSolutions)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const productstableName = "products"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testProduct = [2]Product{  Id: 1, Uuid: "xloSa7tRk11ZaFr2", Name: "8KSPaK0eFLPsC1I2", Model: "ddlAqbr71RsSvCkq", Category: 1, Manufacturer: "ggwCexKwduTakFqI", ImageLink: "4rIn3GlblvCbZP5V", PurchaseLink: "1PGqwO4n3lNP6Tfr", AffiliatedId: "RvbtfRThXejyCvNa", Description: "ztnj2J7ZSK1vyDPN", UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), IsPublic: false, Validated: true, Archived: false, ManufacturerUuid: sql.NullString{"eQAmKTQJOZHWLQKX", true}, ManufacturerProductId: sql.NullString{"lhpkWBX94EcRGDLG", true}, ManufacturerProductLink: sql.NullString{"mJSgJTBOv32SYsJv", true}, Msrp: sql.NullFloat64{-0.3473862637047197, true},  Id: 2, Uuid: "1WvTmpukmQgxfoje", Name: "LUX3e4H3sF5qv3Hd", Model: "Y88sRMN98vx9jcTr", Category: 2, Manufacturer: "WHOd80BIs9YW8SQg", ImageLink: "C7179Kzaqe7DrjAX", PurchaseLink: "9YAe35nGABbjfrYe", AffiliatedId: "t7yCIPGrCkrHy3fb", Description: "u3jgWcvyhpQYqNAv", UserId: 2, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), IsPublic: true, Validated: false, Archived: true, ManufacturerUuid: sql.NullString{"rP9ELyBaOUBGJpt3", true}, ManufacturerProductId: sql.NullString{"nmtDSv0pvUmWkkp9", true}, ManufacturerProductLink: sql.NullString{"8hB2eo4LSMbEgwil", true}, Msrp: sql.NullFloat64{-0.23733370167418588, true} }

var updateProduct = Product Id: 1, Uuid: "LEULJ0AWrcPrZjJL", Name: "OvTuaTmwIsVliXs0", Model: "IHUv9sJtx7QOXADd", Category: 1, Manufacturer: "EvTChYlUrAe3vWts", ImageLink: "ebavH81bdVE2IBHC", PurchaseLink: "5K4G6Y1Ip8d7kFSX", AffiliatedId: "0mteCTzWfvyXFAAV", Description: "dzv9XC93TClOKQiQ", UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), IsPublic: false, Validated: true, Archived: false, ManufacturerUuid: sql.NullString{"vF1oqhxOQ7kEekgC", true}, ManufacturerProductId: sql.NullString{"6MuYwfH6JQtC7ORp", true}, ManufacturerProductLink: sql.NullString{"vBr5WVM6NjDlXumV", true}, Msrp: sql.NullFloat64{-0.18294254228544482, true}

//compare functions
var compareProducts = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Name": defaultCompare,
	"Model": defaultCompare,
	"Category": defaultCompare,
	"Manufacturer": defaultCompare,
	"ImageLink": defaultCompare,
	"PurchaseLink": defaultCompare,
	"AffiliatedId": defaultCompare,
	"Description": defaultCompare,
	"UserId": defaultCompare,
	"CreatedAt": stringCompare,
	"IsPublic": defaultCompare,
	"Validated": defaultCompare,
	"Archived": defaultCompare,
	"ManufacturerUuid": defaultCompare,
	"ManufacturerProductId": defaultCompare,
	"ManufacturerProductLink": defaultCompare,
	"Msrp": defaultCompare,

}

func reverseProducts(products []Product) (result []Product) {

	for i := len(products) - 1; i >= 0; i-- {
		result = append(result, products[i])
	}
	return
}

// ======= tests: Product =======

func TestCreateTableProducts(t *testing.T) {
	fmt.Println("==CreateTableProducts")
	err := CreateTableProducts(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableProducts " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableProducts")
	}
	exists, err := tableExists(testDb, "products")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(products) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateProduct(t *testing.T) {
	fmt.Println("==CreateProduct")
	result, err := testProduct[0].CreateProduct(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct")
	}
	err = equalField(result, testProduct[0], compareProducts)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveProduct(t *testing.T) {
	fmt.Println("==RetrieveProduct")
	result, err := testProduct[0].RetrieveProduct(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveProduct " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveProduct")
	}
	err = equalField(result, testProduct[0], compareProducts)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllProducts(t *testing.T) {
	fmt.Println("==RetrieveAllProducts")
	_, err := testProduct[1].CreateProduct(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct")
	}
	result, err := RetrieveAllProducts(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProducts " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllProducts")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseProducts(testProduct[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareProducts)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateProduct(t *testing.T) {
	fmt.Println("==UpdateProduct")
	result, err := updateProduct.UpdateProduct(testDb)
	if err != nil {
		t.Errorf("cannot UpdateProduct " + err.Error())
	} else {
		fmt.Println("  Done: UpdateProduct")
	}
	err = equalField(result, updateProduct, compareProducts)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const articlestableName = "articles"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testArticle = [2]Article{  Id: 1, Uuid: "nQnrytN7X7Scj3B1", Title: "KOqcsdqsW7EVQdAH", Summary: "Mg9cfl8loqwW9phL", Urls: "EGZQYYw7xahTCOHl", Author: "Z7RIYpo9tmdD1FPZ", Category: "hJlCmn66ivmW1OVy", ImageLink: "ao7EOwr9DEpScJlD", CreationDate: time.Now().UTC().Truncate(time.Microsecond), UserId: 1, ArticleLink: "nFroVwIckkdVtWZZ", BodyHtml: "orbMO6TWNCNuoGgG", Validated: true,  Id: 2, Uuid: "fWmEyGA4hk5pLg4y", Title: "iFQUTa27t4dZ2v9s", Summary: "qtav8nPHxpoI7f1b", Urls: "RC8HI1ZrEuC56Ap0", Author: "VYzijBFfOx2jfccO", Category: "rSM8oaHQ4V1dvnWu", ImageLink: "EUT1Pyjg5LCO4pLv", CreationDate: time.Now().UTC().Truncate(time.Microsecond), UserId: 2, ArticleLink: "Umaky8MkT0Bdfmsi", BodyHtml: "J1W7UXJCcF0tjcWB", Validated: true }

var updateArticle = Article Id: 1, Uuid: "UnvcWndVDBQSVQO7", Title: "d1lVUasfRUlk4yRO", Summary: "AovvuHl4UKoTDah8", Urls: "vHzImAGeZTkanDR9", Author: "5TUmmImfGpwnui8n", Category: "bkkTAc5DHukW8xcE", ImageLink: "iXfRURiGyeDxn6Lc", CreationDate: time.Now().UTC().Truncate(time.Microsecond), UserId: 1, ArticleLink: "WN5t6lLMiN81PzMi", BodyHtml: "hRbKQCxbny9EiYYV", Validated: true

//compare functions
var compareArticles = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Title": defaultCompare,
	"Summary": defaultCompare,
	"Urls": defaultCompare,
	"Author": defaultCompare,
	"Category": defaultCompare,
	"ImageLink": defaultCompare,
	"CreationDate": stringCompare,
	"UserId": defaultCompare,
	"ArticleLink": defaultCompare,
	"BodyHtml": defaultCompare,
	"Validated": defaultCompare,

}

func reverseArticles(articles []Article) (result []Article) {

	for i := len(articles) - 1; i >= 0; i-- {
		result = append(result, articles[i])
	}
	return
}

// ======= tests: Article =======

func TestCreateTableArticles(t *testing.T) {
	fmt.Println("==CreateTableArticles")
	err := CreateTableArticles(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableArticles " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableArticles")
	}
	exists, err := tableExists(testDb, "articles")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(articles) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateArticle(t *testing.T) {
	fmt.Println("==CreateArticle")
	result, err := testArticle[0].CreateArticle(testDb)
	if err != nil {
		t.Errorf("cannot CreateArticle " + err.Error())
	} else {
		fmt.Println("  Done: CreateArticle")
	}
	err = equalField(result, testArticle[0], compareArticles)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveArticle(t *testing.T) {
	fmt.Println("==RetrieveArticle")
	result, err := testArticle[0].RetrieveArticle(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveArticle " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveArticle")
	}
	err = equalField(result, testArticle[0], compareArticles)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllArticles(t *testing.T) {
	fmt.Println("==RetrieveAllArticles")
	_, err := testArticle[1].CreateArticle(testDb)
	if err != nil {
		t.Errorf("cannot CreateArticle " + err.Error())
	} else {
		fmt.Println("  Done: CreateArticle")
	}
	result, err := RetrieveAllArticles(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllArticles " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllArticles")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseArticles(testArticle[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareArticles)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateArticle(t *testing.T) {
	fmt.Println("==UpdateArticle")
	result, err := updateArticle.UpdateArticle(testDb)
	if err != nil {
		t.Errorf("cannot UpdateArticle " + err.Error())
	} else {
		fmt.Println("  Done: UpdateArticle")
	}
	err = equalField(result, updateArticle, compareArticles)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const product_trialstableName = "product_trials"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testProduct_trial = [2]Product_trial{  Id: 1, Uuid: "IzFPoYOSlcMIgtSL", UserId: 1, ProductUuid: "VaqaUCmGfevb3fxQ", StartDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, EndDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), Archived: false,  Id: 2, Uuid: "JLbajHfxIKhj81Nb", UserId: 2, ProductUuid: "t9PQD1pJUQ3RpPJ6", StartDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, EndDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), Archived: true }

var updateProduct_trial = Product_trial Id: 1, Uuid: "2wHVYnGvm1rvZePI", UserId: 1, ProductUuid: "34AcxDuqswElJBjK", StartDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, EndDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), Archived: true

//compare functions
var compareProduct_trials = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"UserId": defaultCompare,
	"ProductUuid": defaultCompare,
	"StartDate": stringCompare,
	"EndDate": stringCompare,
	"CreationDate": stringCompare,
	"Archived": defaultCompare,

}

func reverseProduct_trials(product_trials []Product_trial) (result []Product_trial) {

	for i := len(product_trials) - 1; i >= 0; i-- {
		result = append(result, product_trials[i])
	}
	return
}

// ======= tests: Product_trial =======

func TestCreateTableProduct_trials(t *testing.T) {
	fmt.Println("==CreateTableProduct_trials")
	err := CreateTableProduct_trials(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableProduct_trials " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableProduct_trials")
	}
	exists, err := tableExists(testDb, "product_trials")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(product_trials) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateProduct_trial(t *testing.T) {
	fmt.Println("==CreateProduct_trial")
	result, err := testProduct_trial[0].CreateProduct_trial(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_trial " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_trial")
	}
	err = equalField(result, testProduct_trial[0], compareProduct_trials)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveProduct_trial(t *testing.T) {
	fmt.Println("==RetrieveProduct_trial")
	result, err := testProduct_trial[0].RetrieveProduct_trial(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveProduct_trial " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveProduct_trial")
	}
	err = equalField(result, testProduct_trial[0], compareProduct_trials)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllProduct_trials(t *testing.T) {
	fmt.Println("==RetrieveAllProduct_trials")
	_, err := testProduct_trial[1].CreateProduct_trial(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_trial " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_trial")
	}
	result, err := RetrieveAllProduct_trials(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_trials " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllProduct_trials")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseProduct_trials(testProduct_trial[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareProduct_trials)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateProduct_trial(t *testing.T) {
	fmt.Println("==UpdateProduct_trial")
	result, err := updateProduct_trial.UpdateProduct_trial(testDb)
	if err != nil {
		t.Errorf("cannot UpdateProduct_trial " + err.Error())
	} else {
		fmt.Println("  Done: UpdateProduct_trial")
	}
	err = equalField(result, updateProduct_trial, compareProduct_trials)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const hypnostatstableName = "hypnostats"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testHypnostat = [2]Hypnostat{  Id: 1, UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), Source: "yoifLE11tck3cF1l", Use: true, Hypno: "{\"name\": \"b0LHOXgj5PookEWg\", \"age\": 353294265, \"city\": \"CH21eHlZ0EcnU1aVlhxN\"}", HypnoModel: 1889044415, MotionLen: 604237443, EnergyLen: 2100245084, HeartrateLen: 95875133, StandLen: 755490870, SleepsleepLen: 91179353, SleepinbedLen: 195476648, BeginBedRel: 0.8921058694980185, BeginBed: -0.8946272093272022, EndBed: -0.9755562425511233, HrMin: -0.2707283955828474, HrMax: -1.2631590594960072, HrAvg: 1.2130421584402868, Tst: -0.703701710713742, TimeAwake: -0.7933276631730521, TimeRem: 1.4247185506189952, TimeLight: -0.3449231470229793, TimeDeep: -0.7029106969619003, NumAwakes: 0.7080815980066608, Score: -1.5937590169935723, ScoreDuration: -0.6615799175783774, ScoreEfficiency: -0.8175374214661579, ScoreContinuity: -2.036210180596706, UtcOffset: 1746165052, SleepOnset: 1.485964302470144, SleepEfficiency: 0.7770039791706532,  Id: 2, UserId: 2, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), Source: "UYU6GNJBT1paXUM3", Use: true, Hypno: "{\"name\": \"JXbCXGppiqPoaPLv\", \"age\": 1419419863, \"city\": \"ucvvlWp6HhE2EOJ9iVvS\"}", HypnoModel: 86221317, MotionLen: 1139929738, EnergyLen: 980422869, HeartrateLen: 1316237552, StandLen: 1236633126, SleepsleepLen: 1724996731, SleepinbedLen: 875431319, BeginBedRel: -0.3098268031014817, BeginBed: 0.3116907691364911, EndBed: 0.018324615233548947, HrMin: 1.8897203166194505, HrMax: -0.37770771311743256, HrAvg: -0.194526037744297, Tst: 0.046232157764804216, TimeAwake: -1.1462875445057696, TimeRem: -1.53910969138551, TimeLight: 0.8786646422761399, TimeDeep: 1.8063418428573175, NumAwakes: 0.8419430270719513, Score: 1.201608038578108, ScoreDuration: 0.5661263182608003, ScoreEfficiency: -0.9969416183323425, ScoreContinuity: 1.6857210664921822, UtcOffset: 2078050787, SleepOnset: 0.5955705246712403, SleepEfficiency: -1.0230477147682437 }

var updateHypnostat = Hypnostat Id: 1, UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), Source: "JNN7cwzNv2GTgX6H", Use: true, Hypno: "{\"name\": \"pvvYreSTxoTx67LA\", \"age\": 1743010716, \"city\": \"ENUJJOVVngJ4EPewcd2p\"}", HypnoModel: 1599424472, MotionLen: 752526487, EnergyLen: 281441495, HeartrateLen: 1121072461, StandLen: 201148367, SleepsleepLen: 408002103, SleepinbedLen: 1956836132, BeginBedRel: 1.4598217810853225, BeginBed: -1.0127747416031438, EndBed: -1.4254993788607107, HrMin: -0.9204775170323252, HrMax: 1.9974131858858404, HrAvg: 1.9525320414855099, Tst: -1.0697743116914997, TimeAwake: -1.0755300878917289, TimeRem: 2.267608509177378, TimeLight: 2.1628992262838853, TimeDeep: 0.8831927252159804, NumAwakes: -0.8516208811024304, Score: -1.1433009234145968, ScoreDuration: -0.8326705933680514, ScoreEfficiency: 0.1706764727452954, ScoreContinuity: -0.5804045370164347, UtcOffset: 33172211, SleepOnset: 0.5914598912776481, SleepEfficiency: -1.0310254050687497

//compare functions
var compareHypnostats = map[string]compareType{
	"Id": defaultCompare,
	"UserId": defaultCompare,
	"CreatedAt": stringCompare,
	"Source": defaultCompare,
	"Use": defaultCompare,
	"Hypno": jsonCompare,
	"HypnoModel": defaultCompare,
	"MotionLen": defaultCompare,
	"EnergyLen": defaultCompare,
	"HeartrateLen": defaultCompare,
	"StandLen": defaultCompare,
	"SleepsleepLen": defaultCompare,
	"SleepinbedLen": defaultCompare,
	"BeginBedRel": defaultCompare,
	"BeginBed": defaultCompare,
	"EndBed": defaultCompare,
	"HrMin": defaultCompare,
	"HrMax": defaultCompare,
	"HrAvg": defaultCompare,
	"Tst": defaultCompare,
	"TimeAwake": defaultCompare,
	"TimeRem": defaultCompare,
	"TimeLight": defaultCompare,
	"TimeDeep": defaultCompare,
	"NumAwakes": defaultCompare,
	"Score": defaultCompare,
	"ScoreDuration": defaultCompare,
	"ScoreEfficiency": defaultCompare,
	"ScoreContinuity": defaultCompare,
	"UtcOffset": defaultCompare,
	"SleepOnset": defaultCompare,
	"SleepEfficiency": defaultCompare,

}

func reverseHypnostats(hypnostats []Hypnostat) (result []Hypnostat) {

	for i := len(hypnostats) - 1; i >= 0; i-- {
		result = append(result, hypnostats[i])
	}
	return
}

// ======= tests: Hypnostat =======

func TestCreateTableHypnostats(t *testing.T) {
	fmt.Println("==CreateTableHypnostats")
	err := CreateTableHypnostats(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableHypnostats " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableHypnostats")
	}
	exists, err := tableExists(testDb, "hypnostats")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(hypnostats) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateHypnostat(t *testing.T) {
	fmt.Println("==CreateHypnostat")
	result, err := testHypnostat[0].CreateHypnostat(testDb)
	if err != nil {
		t.Errorf("cannot CreateHypnostat " + err.Error())
	} else {
		fmt.Println("  Done: CreateHypnostat")
	}
	err = equalField(result, testHypnostat[0], compareHypnostats)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveHypnostat(t *testing.T) {
	fmt.Println("==RetrieveHypnostat")
	result, err := testHypnostat[0].RetrieveHypnostat(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveHypnostat " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveHypnostat")
	}
	err = equalField(result, testHypnostat[0], compareHypnostats)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllHypnostats(t *testing.T) {
	fmt.Println("==RetrieveAllHypnostats")
	_, err := testHypnostat[1].CreateHypnostat(testDb)
	if err != nil {
		t.Errorf("cannot CreateHypnostat " + err.Error())
	} else {
		fmt.Println("  Done: CreateHypnostat")
	}
	result, err := RetrieveAllHypnostats(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllHypnostats " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllHypnostats")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseHypnostats(testHypnostat[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareHypnostats)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateHypnostat(t *testing.T) {
	fmt.Println("==UpdateHypnostat")
	result, err := updateHypnostat.UpdateHypnostat(testDb)
	if err != nil {
		t.Errorf("cannot UpdateHypnostat " + err.Error())
	} else {
		fmt.Println("  Done: UpdateHypnostat")
	}
	err = equalField(result, updateHypnostat, compareHypnostats)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const product_commentstableName = "product_comments"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testProduct_comment = [2]Product_comment{  Id: 1, Uuid: "9d6hCgRHGm0t1NZT", UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), Comments: sql.NullString{"pWrDXu0O31G1ggZl", true}, Rating: 367997908, ProductUuid: "eQl7WG1FrvGqEnmR",  Id: 2, Uuid: "j7XFe1ATGE7e17ge", UserId: 2, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), Comments: sql.NullString{"i9gubJXhDWmuCURE", true}, Rating: 1427299315, ProductUuid: "CDHMQQ7A9xHj78CA" }

var updateProduct_comment = Product_comment Id: 1, Uuid: "XCxqc7OTuT4Tbmxb", UserId: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), LastUpdated: time.Now().UTC().Truncate(time.Microsecond), Comments: sql.NullString{"G3td0tNhdANVOwWU", true}, Rating: 892145459, ProductUuid: "Lfuo7NLaxUGA3ztv"

//compare functions
var compareProduct_comments = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"UserId": defaultCompare,
	"CreatedAt": stringCompare,
	"LastUpdated": stringCompare,
	"Comments": defaultCompare,
	"Rating": defaultCompare,
	"ProductUuid": defaultCompare,

}

func reverseProduct_comments(product_comments []Product_comment) (result []Product_comment) {

	for i := len(product_comments) - 1; i >= 0; i-- {
		result = append(result, product_comments[i])
	}
	return
}

// ======= tests: Product_comment =======

func TestCreateTableProduct_comments(t *testing.T) {
	fmt.Println("==CreateTableProduct_comments")
	err := CreateTableProduct_comments(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableProduct_comments " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableProduct_comments")
	}
	exists, err := tableExists(testDb, "product_comments")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(product_comments) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateProduct_comment(t *testing.T) {
	fmt.Println("==CreateProduct_comment")
	result, err := testProduct_comment[0].CreateProduct_comment(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_comment " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_comment")
	}
	err = equalField(result, testProduct_comment[0], compareProduct_comments)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveProduct_comment(t *testing.T) {
	fmt.Println("==RetrieveProduct_comment")
	result, err := testProduct_comment[0].RetrieveProduct_comment(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveProduct_comment " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveProduct_comment")
	}
	err = equalField(result, testProduct_comment[0], compareProduct_comments)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllProduct_comments(t *testing.T) {
	fmt.Println("==RetrieveAllProduct_comments")
	_, err := testProduct_comment[1].CreateProduct_comment(testDb)
	if err != nil {
		t.Errorf("cannot CreateProduct_comment " + err.Error())
	} else {
		fmt.Println("  Done: CreateProduct_comment")
	}
	result, err := RetrieveAllProduct_comments(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_comments " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllProduct_comments")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseProduct_comments(testProduct_comment[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareProduct_comments)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateProduct_comment(t *testing.T) {
	fmt.Println("==UpdateProduct_comment")
	result, err := updateProduct_comment.UpdateProduct_comment(testDb)
	if err != nil {
		t.Errorf("cannot UpdateProduct_comment " + err.Error())
	} else {
		fmt.Println("  Done: UpdateProduct_comment")
	}
	err = equalField(result, updateProduct_comment, compareProduct_comments)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const community_product_trial_statstableName = "community_product_trial_stats"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testCommunity_product_trial_stat = [2]Community_product_trial_stat{  Id: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), ProductUuid: "eCiKmVSelYE94dPu", Source: "470gmXV5srNaWLXL", Days: 1147277284, UserCount: 1920997527, Tst: sql.NullFloat64{-0.28605823971346195, true}, TimeAwake: sql.NullFloat64{0.07871284265989487, true}, TimeRem: sql.NullFloat64{-0.7522555533819817, true}, TimeLight: sql.NullFloat64{0.5728569113612444, true}, TimeDeep: sql.NullFloat64{0.6061209920973862, true}, NumAwakes: sql.NullFloat64{1.3952507407830153, true}, Score: sql.NullFloat64{1.0775738653627638, true}, SleepOnset: sql.NullFloat64{0.8044866024399333, true}, SleepEfficiency: sql.NullFloat64{-0.7390181479721253, true},  Id: 2, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), ProductUuid: "UKPDhycbsCSxcHEr", Source: "eB2ZGsjctzfAZj4R", Days: 1197853366, UserCount: 900296972, Tst: sql.NullFloat64{-0.07116467049688235, true}, TimeAwake: sql.NullFloat64{0.13577417357089033, true}, TimeRem: sql.NullFloat64{-2.1946943836746513, true}, TimeLight: sql.NullFloat64{-0.07509588651616494, true}, TimeDeep: sql.NullFloat64{-0.3344788012685082, true}, NumAwakes: sql.NullFloat64{0.4534466636201801, true}, Score: sql.NullFloat64{-0.9139760818632057, true}, SleepOnset: sql.NullFloat64{-0.21665553245368763, true}, SleepEfficiency: sql.NullFloat64{0.03584014140030245, true} }

var updateCommunity_product_trial_stat = Community_product_trial_stat Id: 1, CreatedAt: time.Now().UTC().Truncate(time.Microsecond), ProductUuid: "NiW4MdJr32Gb3rZk", Source: "WFNaTsGWB48K9d6x", Days: 802428632, UserCount: 1214842608, Tst: sql.NullFloat64{0.6977263811175791, true}, TimeAwake: sql.NullFloat64{0.15151294229353546, true}, TimeRem: sql.NullFloat64{-0.13286503934738436, true}, TimeLight: sql.NullFloat64{-1.0352185204315951, true}, TimeDeep: sql.NullFloat64{0.06522956506116828, true}, NumAwakes: sql.NullFloat64{-0.9570415394482031, true}, Score: sql.NullFloat64{-0.5205128716938747, true}, SleepOnset: sql.NullFloat64{0.041021308273598955, true}, SleepEfficiency: sql.NullFloat64{0.5591887371341402, true}

//compare functions
var compareCommunity_product_trial_stats = map[string]compareType{
	"Id": defaultCompare,
	"CreatedAt": stringCompare,
	"ProductUuid": defaultCompare,
	"Source": defaultCompare,
	"Days": defaultCompare,
	"UserCount": defaultCompare,
	"Tst": defaultCompare,
	"TimeAwake": defaultCompare,
	"TimeRem": defaultCompare,
	"TimeLight": defaultCompare,
	"TimeDeep": defaultCompare,
	"NumAwakes": defaultCompare,
	"Score": defaultCompare,
	"SleepOnset": defaultCompare,
	"SleepEfficiency": defaultCompare,

}

func reverseCommunity_product_trial_stats(community_product_trial_stats []Community_product_trial_stat) (result []Community_product_trial_stat) {

	for i := len(community_product_trial_stats) - 1; i >= 0; i-- {
		result = append(result, community_product_trial_stats[i])
	}
	return
}

// ======= tests: Community_product_trial_stat =======

func TestCreateTableCommunity_product_trial_stats(t *testing.T) {
	fmt.Println("==CreateTableCommunity_product_trial_stats")
	err := CreateTableCommunity_product_trial_stats(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableCommunity_product_trial_stats " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableCommunity_product_trial_stats")
	}
	exists, err := tableExists(testDb, "community_product_trial_stats")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(community_product_trial_stats) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateCommunity_product_trial_stat(t *testing.T) {
	fmt.Println("==CreateCommunity_product_trial_stat")
	result, err := testCommunity_product_trial_stat[0].CreateCommunity_product_trial_stat(testDb)
	if err != nil {
		t.Errorf("cannot CreateCommunity_product_trial_stat " + err.Error())
	} else {
		fmt.Println("  Done: CreateCommunity_product_trial_stat")
	}
	err = equalField(result, testCommunity_product_trial_stat[0], compareCommunity_product_trial_stats)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveCommunity_product_trial_stat(t *testing.T) {
	fmt.Println("==RetrieveCommunity_product_trial_stat")
	result, err := testCommunity_product_trial_stat[0].RetrieveCommunity_product_trial_stat(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveCommunity_product_trial_stat " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveCommunity_product_trial_stat")
	}
	err = equalField(result, testCommunity_product_trial_stat[0], compareCommunity_product_trial_stats)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllCommunity_product_trial_stats(t *testing.T) {
	fmt.Println("==RetrieveAllCommunity_product_trial_stats")
	_, err := testCommunity_product_trial_stat[1].CreateCommunity_product_trial_stat(testDb)
	if err != nil {
		t.Errorf("cannot CreateCommunity_product_trial_stat " + err.Error())
	} else {
		fmt.Println("  Done: CreateCommunity_product_trial_stat")
	}
	result, err := RetrieveAllCommunity_product_trial_stats(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllCommunity_product_trial_stats " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllCommunity_product_trial_stats")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseCommunity_product_trial_stats(testCommunity_product_trial_stat[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareCommunity_product_trial_stats)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateCommunity_product_trial_stat(t *testing.T) {
	fmt.Println("==UpdateCommunity_product_trial_stat")
	result, err := updateCommunity_product_trial_stat.UpdateCommunity_product_trial_stat(testDb)
	if err != nil {
		t.Errorf("cannot UpdateCommunity_product_trial_stat " + err.Error())
	} else {
		fmt.Println("  Done: UpdateCommunity_product_trial_stat")
	}
	err = equalField(result, updateCommunity_product_trial_stat, compareCommunity_product_trial_stats)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const habitstableName = "habits"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testHabit = [2]Habit{  Id: 1, Uuid: "gob12HPAONtEeBpG", UserId: 1, ProductUuid: "3bma6ou5JNMkcSk5", StartDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, EndDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), Archived: true,  Id: 2, Uuid: "hop37AxOn42ACAEr", UserId: 2, ProductUuid: "pZ1AM5XBwy7TIaBJ", StartDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, EndDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), Archived: false }

var updateHabit = Habit Id: 1, Uuid: "WBY4eGlbQ83TnlWr", UserId: 1, ProductUuid: "wgckM3KIof1dDNbj", StartDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, EndDate: sql.NullTime{time.Now().UTC().Truncate(time.Microsecond), true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond), Archived: true

//compare functions
var compareHabits = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"UserId": defaultCompare,
	"ProductUuid": defaultCompare,
	"StartDate": stringCompare,
	"EndDate": stringCompare,
	"CreationDate": stringCompare,
	"Archived": defaultCompare,

}

func reverseHabits(habits []Habit) (result []Habit) {

	for i := len(habits) - 1; i >= 0; i-- {
		result = append(result, habits[i])
	}
	return
}

// ======= tests: Habit =======

func TestCreateTableHabits(t *testing.T) {
	fmt.Println("==CreateTableHabits")
	err := CreateTableHabits(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableHabits " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableHabits")
	}
	exists, err := tableExists(testDb, "habits")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(habits) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateHabit(t *testing.T) {
	fmt.Println("==CreateHabit")
	result, err := testHabit[0].CreateHabit(testDb)
	if err != nil {
		t.Errorf("cannot CreateHabit " + err.Error())
	} else {
		fmt.Println("  Done: CreateHabit")
	}
	err = equalField(result, testHabit[0], compareHabits)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveHabit(t *testing.T) {
	fmt.Println("==RetrieveHabit")
	result, err := testHabit[0].RetrieveHabit(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveHabit " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveHabit")
	}
	err = equalField(result, testHabit[0], compareHabits)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllHabits(t *testing.T) {
	fmt.Println("==RetrieveAllHabits")
	_, err := testHabit[1].CreateHabit(testDb)
	if err != nil {
		t.Errorf("cannot CreateHabit " + err.Error())
	} else {
		fmt.Println("  Done: CreateHabit")
	}
	result, err := RetrieveAllHabits(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllHabits " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllHabits")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseHabits(testHabit[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareHabits)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateHabit(t *testing.T) {
	fmt.Println("==UpdateHabit")
	result, err := updateHabit.UpdateHabit(testDb)
	if err != nil {
		t.Errorf("cannot UpdateHabit " + err.Error())
	} else {
		fmt.Println("  Done: UpdateHabit")
	}
	err = equalField(result, updateHabit, compareHabits)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const community_product_resultstableName = "community_product_results"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testCommunity_product_result = [2]Community_product_result{  Id: 1, UserCount: 125134000, Count: 609561802, ProductUuid: "HznAOnAwiNvbzkW9", Tst: sql.NullFloat64{0.29734566533104356, true}, TimeAwake: sql.NullFloat64{-1.4515869238354309, true}, TimeRem: sql.NullFloat64{-1.5113987059766827, true}, TimeLight: sql.NullFloat64{1.8052735143067178, true}, TimeDeep: sql.NullFloat64{0.8071934164518434, true}, NumAwakes: sql.NullFloat64{0.007765461609626351, true}, Score: sql.NullFloat64{-0.3664494712550386, true}, SleepOnset: sql.NullFloat64{-0.06421326798317917, true}, SleepEfficiency: sql.NullFloat64{-1.9186630861421314, true},  Id: 2, UserCount: 222348744, Count: 164667837, ProductUuid: "Aj8aumZHYXN1iQez", Tst: sql.NullFloat64{-0.0409131146486969, true}, TimeAwake: sql.NullFloat64{-0.24519719966355824, true}, TimeRem: sql.NullFloat64{-1.6116378313247859, true}, TimeLight: sql.NullFloat64{1.9599464046882802, true}, TimeDeep: sql.NullFloat64{-0.08094378853347362, true}, NumAwakes: sql.NullFloat64{0.4719918582335914, true}, Score: sql.NullFloat64{0.35994248882982016, true}, SleepOnset: sql.NullFloat64{0.4491216721248541, true}, SleepEfficiency: sql.NullFloat64{1.3289723376858724, true} }

var updateCommunity_product_result = Community_product_result Id: 1, UserCount: 246003637, Count: 2097010076, ProductUuid: "t4a7K5H5a6EShAU4", Tst: sql.NullFloat64{1.1981170547624433, true}, TimeAwake: sql.NullFloat64{0.6504394269105004, true}, TimeRem: sql.NullFloat64{1.3128597477787938, true}, TimeLight: sql.NullFloat64{1.4059939202651757, true}, TimeDeep: sql.NullFloat64{-0.7495366170829235, true}, NumAwakes: sql.NullFloat64{0.5043072379639186, true}, Score: sql.NullFloat64{0.3869674202670298, true}, SleepOnset: sql.NullFloat64{0.8092946362457321, true}, SleepEfficiency: sql.NullFloat64{1.2410932922630176, true}

//compare functions
var compareCommunity_product_results = map[string]compareType{
	"Id": defaultCompare,
	"UserCount": defaultCompare,
	"Count": defaultCompare,
	"ProductUuid": defaultCompare,
	"Tst": defaultCompare,
	"TimeAwake": defaultCompare,
	"TimeRem": defaultCompare,
	"TimeLight": defaultCompare,
	"TimeDeep": defaultCompare,
	"NumAwakes": defaultCompare,
	"Score": defaultCompare,
	"SleepOnset": defaultCompare,
	"SleepEfficiency": defaultCompare,

}

func reverseCommunity_product_results(community_product_results []Community_product_result) (result []Community_product_result) {

	for i := len(community_product_results) - 1; i >= 0; i-- {
		result = append(result, community_product_results[i])
	}
	return
}

// ======= tests: Community_product_result =======

func TestCreateTableCommunity_product_results(t *testing.T) {
	fmt.Println("==CreateTableCommunity_product_results")
	err := CreateTableCommunity_product_results(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableCommunity_product_results " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableCommunity_product_results")
	}
	exists, err := tableExists(testDb, "community_product_results")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(community_product_results) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateCommunity_product_result(t *testing.T) {
	fmt.Println("==CreateCommunity_product_result")
	result, err := testCommunity_product_result[0].CreateCommunity_product_result(testDb)
	if err != nil {
		t.Errorf("cannot CreateCommunity_product_result " + err.Error())
	} else {
		fmt.Println("  Done: CreateCommunity_product_result")
	}
	err = equalField(result, testCommunity_product_result[0], compareCommunity_product_results)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveCommunity_product_result(t *testing.T) {
	fmt.Println("==RetrieveCommunity_product_result")
	result, err := testCommunity_product_result[0].RetrieveCommunity_product_result(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveCommunity_product_result " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveCommunity_product_result")
	}
	err = equalField(result, testCommunity_product_result[0], compareCommunity_product_results)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllCommunity_product_results(t *testing.T) {
	fmt.Println("==RetrieveAllCommunity_product_results")
	_, err := testCommunity_product_result[1].CreateCommunity_product_result(testDb)
	if err != nil {
		t.Errorf("cannot CreateCommunity_product_result " + err.Error())
	} else {
		fmt.Println("  Done: CreateCommunity_product_result")
	}
	result, err := RetrieveAllCommunity_product_results(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllCommunity_product_results " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllCommunity_product_results")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseCommunity_product_results(testCommunity_product_result[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareCommunity_product_results)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateCommunity_product_result(t *testing.T) {
	fmt.Println("==UpdateCommunity_product_result")
	result, err := updateCommunity_product_result.UpdateCommunity_product_result(testDb)
	if err != nil {
		t.Errorf("cannot UpdateCommunity_product_result " + err.Error())
	} else {
		fmt.Println("  Done: UpdateCommunity_product_result")
	}
	err = equalField(result, updateCommunity_product_result, compareCommunity_product_results)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


//delete all data in reverse order to accommodate foreign keys


func TestDeleteCommunity_product_result(t *testing.T) {
	fmt.Println("==DeleteCommunity_product_result")
	err := testCommunity_product_result[0].DeleteCommunity_product_result(testDb)
	if err != nil {
		t.Errorf("cannot DeleteCommunity_product_result " + err.Error())
	} else {
		fmt.Println("  Done: DeleteCommunity_product_result")
	}
	_, err = testCommunity_product_result[0].RetrieveCommunity_product_result(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveCommunity_product_result with no result")
		} else {
			t.Errorf("cannot RetrieveCommunity_product_result " + err.Error())
		}
	}
}

func TestDeleteAllCommunity_product_results(t *testing.T) {
	fmt.Println("==DeleteAllCommunity_product_results")
	err := DeleteAllCommunity_product_results(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllCommunity_product_results " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllCommunity_product_results")
	}
	result, err := RetrieveAllCommunity_product_results(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllCommunity_product_results " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllCommunity_product_results with no result")
	}
}


func TestDeleteHabit(t *testing.T) {
	fmt.Println("==DeleteHabit")
	err := testHabit[0].DeleteHabit(testDb)
	if err != nil {
		t.Errorf("cannot DeleteHabit " + err.Error())
	} else {
		fmt.Println("  Done: DeleteHabit")
	}
	_, err = testHabit[0].RetrieveHabit(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveHabit with no result")
		} else {
			t.Errorf("cannot RetrieveHabit " + err.Error())
		}
	}
}

func TestDeleteAllHabits(t *testing.T) {
	fmt.Println("==DeleteAllHabits")
	err := DeleteAllHabits(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllHabits " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllHabits")
	}
	result, err := RetrieveAllHabits(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllHabits " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllHabits with no result")
	}
}


func TestDeleteCommunity_product_trial_stat(t *testing.T) {
	fmt.Println("==DeleteCommunity_product_trial_stat")
	err := testCommunity_product_trial_stat[0].DeleteCommunity_product_trial_stat(testDb)
	if err != nil {
		t.Errorf("cannot DeleteCommunity_product_trial_stat " + err.Error())
	} else {
		fmt.Println("  Done: DeleteCommunity_product_trial_stat")
	}
	_, err = testCommunity_product_trial_stat[0].RetrieveCommunity_product_trial_stat(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveCommunity_product_trial_stat with no result")
		} else {
			t.Errorf("cannot RetrieveCommunity_product_trial_stat " + err.Error())
		}
	}
}

func TestDeleteAllCommunity_product_trial_stats(t *testing.T) {
	fmt.Println("==DeleteAllCommunity_product_trial_stats")
	err := DeleteAllCommunity_product_trial_stats(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllCommunity_product_trial_stats " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllCommunity_product_trial_stats")
	}
	result, err := RetrieveAllCommunity_product_trial_stats(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllCommunity_product_trial_stats " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllCommunity_product_trial_stats with no result")
	}
}


func TestDeleteProduct_comment(t *testing.T) {
	fmt.Println("==DeleteProduct_comment")
	err := testProduct_comment[0].DeleteProduct_comment(testDb)
	if err != nil {
		t.Errorf("cannot DeleteProduct_comment " + err.Error())
	} else {
		fmt.Println("  Done: DeleteProduct_comment")
	}
	_, err = testProduct_comment[0].RetrieveProduct_comment(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveProduct_comment with no result")
		} else {
			t.Errorf("cannot RetrieveProduct_comment " + err.Error())
		}
	}
}

func TestDeleteAllProduct_comments(t *testing.T) {
	fmt.Println("==DeleteAllProduct_comments")
	err := DeleteAllProduct_comments(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllProduct_comments " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllProduct_comments")
	}
	result, err := RetrieveAllProduct_comments(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_comments " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllProduct_comments with no result")
	}
}


func TestDeleteHypnostat(t *testing.T) {
	fmt.Println("==DeleteHypnostat")
	err := testHypnostat[0].DeleteHypnostat(testDb)
	if err != nil {
		t.Errorf("cannot DeleteHypnostat " + err.Error())
	} else {
		fmt.Println("  Done: DeleteHypnostat")
	}
	_, err = testHypnostat[0].RetrieveHypnostat(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveHypnostat with no result")
		} else {
			t.Errorf("cannot RetrieveHypnostat " + err.Error())
		}
	}
}

func TestDeleteAllHypnostats(t *testing.T) {
	fmt.Println("==DeleteAllHypnostats")
	err := DeleteAllHypnostats(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllHypnostats " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllHypnostats")
	}
	result, err := RetrieveAllHypnostats(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllHypnostats " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllHypnostats with no result")
	}
}


func TestDeleteProduct_trial(t *testing.T) {
	fmt.Println("==DeleteProduct_trial")
	err := testProduct_trial[0].DeleteProduct_trial(testDb)
	if err != nil {
		t.Errorf("cannot DeleteProduct_trial " + err.Error())
	} else {
		fmt.Println("  Done: DeleteProduct_trial")
	}
	_, err = testProduct_trial[0].RetrieveProduct_trial(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveProduct_trial with no result")
		} else {
			t.Errorf("cannot RetrieveProduct_trial " + err.Error())
		}
	}
}

func TestDeleteAllProduct_trials(t *testing.T) {
	fmt.Println("==DeleteAllProduct_trials")
	err := DeleteAllProduct_trials(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllProduct_trials " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllProduct_trials")
	}
	result, err := RetrieveAllProduct_trials(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_trials " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllProduct_trials with no result")
	}
}


func TestDeleteArticle(t *testing.T) {
	fmt.Println("==DeleteArticle")
	err := testArticle[0].DeleteArticle(testDb)
	if err != nil {
		t.Errorf("cannot DeleteArticle " + err.Error())
	} else {
		fmt.Println("  Done: DeleteArticle")
	}
	_, err = testArticle[0].RetrieveArticle(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveArticle with no result")
		} else {
			t.Errorf("cannot RetrieveArticle " + err.Error())
		}
	}
}

func TestDeleteAllArticles(t *testing.T) {
	fmt.Println("==DeleteAllArticles")
	err := DeleteAllArticles(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllArticles " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllArticles")
	}
	result, err := RetrieveAllArticles(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllArticles " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllArticles with no result")
	}
}


func TestDeleteProduct(t *testing.T) {
	fmt.Println("==DeleteProduct")
	err := testProduct[0].DeleteProduct(testDb)
	if err != nil {
		t.Errorf("cannot DeleteProduct " + err.Error())
	} else {
		fmt.Println("  Done: DeleteProduct")
	}
	_, err = testProduct[0].RetrieveProduct(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveProduct with no result")
		} else {
			t.Errorf("cannot RetrieveProduct " + err.Error())
		}
	}
}

func TestDeleteAllProducts(t *testing.T) {
	fmt.Println("==DeleteAllProducts")
	err := DeleteAllProducts(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllProducts " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllProducts")
	}
	result, err := RetrieveAllProducts(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProducts " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllProducts with no result")
	}
}


func TestDeleteSolution(t *testing.T) {
	fmt.Println("==DeleteSolution")
	err := testSolution[0].DeleteSolution(testDb)
	if err != nil {
		t.Errorf("cannot DeleteSolution " + err.Error())
	} else {
		fmt.Println("  Done: DeleteSolution")
	}
	_, err = testSolution[0].RetrieveSolution(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveSolution with no result")
		} else {
			t.Errorf("cannot RetrieveSolution " + err.Error())
		}
	}
}

func TestDeleteAllSolutions(t *testing.T) {
	fmt.Println("==DeleteAllSolutions")
	err := DeleteAllSolutions(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllSolutions " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllSolutions")
	}
	result, err := RetrieveAllSolutions(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSolutions " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllSolutions with no result")
	}
}


func TestDeleteVendor(t *testing.T) {
	fmt.Println("==DeleteVendor")
	err := testVendor[0].DeleteVendor(testDb)
	if err != nil {
		t.Errorf("cannot DeleteVendor " + err.Error())
	} else {
		fmt.Println("  Done: DeleteVendor")
	}
	_, err = testVendor[0].RetrieveVendor(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveVendor with no result")
		} else {
			t.Errorf("cannot RetrieveVendor " + err.Error())
		}
	}
}

func TestDeleteAllVendors(t *testing.T) {
	fmt.Println("==DeleteAllVendors")
	err := DeleteAllVendors(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllVendors " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllVendors")
	}
	result, err := RetrieveAllVendors(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllVendors " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllVendors with no result")
	}
}


func TestDeletePhoto_url(t *testing.T) {
	fmt.Println("==DeletePhoto_url")
	err := testPhoto_url[0].DeletePhoto_url(testDb)
	if err != nil {
		t.Errorf("cannot DeletePhoto_url " + err.Error())
	} else {
		fmt.Println("  Done: DeletePhoto_url")
	}
	_, err = testPhoto_url[0].RetrievePhoto_url(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrievePhoto_url with no result")
		} else {
			t.Errorf("cannot RetrievePhoto_url " + err.Error())
		}
	}
}

func TestDeleteAllPhoto_urls(t *testing.T) {
	fmt.Println("==DeleteAllPhoto_urls")
	err := DeleteAllPhoto_urls(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllPhoto_urls " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllPhoto_urls")
	}
	result, err := RetrieveAllPhoto_urls(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllPhoto_urls " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllPhoto_urls with no result")
	}
}


func TestDeleteProduct_listing(t *testing.T) {
	fmt.Println("==DeleteProduct_listing")
	err := testProduct_listing[0].DeleteProduct_listing(testDb)
	if err != nil {
		t.Errorf("cannot DeleteProduct_listing " + err.Error())
	} else {
		fmt.Println("  Done: DeleteProduct_listing")
	}
	_, err = testProduct_listing[0].RetrieveProduct_listing(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveProduct_listing with no result")
		} else {
			t.Errorf("cannot RetrieveProduct_listing " + err.Error())
		}
	}
}

func TestDeleteAllProduct_listings(t *testing.T) {
	fmt.Println("==DeleteAllProduct_listings")
	err := DeleteAllProduct_listings(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllProduct_listings " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllProduct_listings")
	}
	result, err := RetrieveAllProduct_listings(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_listings " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllProduct_listings with no result")
	}
}


func TestDeleteProduct_review(t *testing.T) {
	fmt.Println("==DeleteProduct_review")
	err := testProduct_review[0].DeleteProduct_review(testDb)
	if err != nil {
		t.Errorf("cannot DeleteProduct_review " + err.Error())
	} else {
		fmt.Println("  Done: DeleteProduct_review")
	}
	_, err = testProduct_review[0].RetrieveProduct_review(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveProduct_review with no result")
		} else {
			t.Errorf("cannot RetrieveProduct_review " + err.Error())
		}
	}
}

func TestDeleteAllProduct_reviews(t *testing.T) {
	fmt.Println("==DeleteAllProduct_reviews")
	err := DeleteAllProduct_reviews(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllProduct_reviews " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllProduct_reviews")
	}
	result, err := RetrieveAllProduct_reviews(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllProduct_reviews " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllProduct_reviews with no result")
	}
}


func TestDeleteCommunity_result(t *testing.T) {
	fmt.Println("==DeleteCommunity_result")
	err := testCommunity_result[0].DeleteCommunity_result(testDb)
	if err != nil {
		t.Errorf("cannot DeleteCommunity_result " + err.Error())
	} else {
		fmt.Println("  Done: DeleteCommunity_result")
	}
	_, err = testCommunity_result[0].RetrieveCommunity_result(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveCommunity_result with no result")
		} else {
			t.Errorf("cannot RetrieveCommunity_result " + err.Error())
		}
	}
}

func TestDeleteAllCommunity_results(t *testing.T) {
	fmt.Println("==DeleteAllCommunity_results")
	err := DeleteAllCommunity_results(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllCommunity_results " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllCommunity_results")
	}
	result, err := RetrieveAllCommunity_results(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllCommunity_results " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllCommunity_results with no result")
	}
}
