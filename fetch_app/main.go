package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	rds "gitlab.com/erwinnekketsu/fetch_app.git/redis"
)

// Users struct
type Users struct {
	User struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Phone     string    `json:"phone"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"user"`
	Success bool `json:"success"`
}

// Commodity struct
type CommodityList struct {
	CommodityList []Commodity `json:"commodity_list"`
}

// Commodity struct
type Commodity struct {
	UUID         string `json:"uuid"`
	Komoditas    string `json:"komoditas"`
	AreaProvinsi string `json:"area_provinsi"`
	AreaKota     string `json:"area_kota"`
	Size         string `json:"size"`
	Price        string `json:"price"`
	PriceUSD     string `json:"price_usd"`
	TglParsed    string `json:"tgl_parsed"`
	Timestamp    string `json:"timestamp"`
}

// CommodityAggr struct
type CommodityAggr struct {
	AreaProvinsi  string          `json:"area_provinsi"`
	TglParsed     string          `json:"tgl_parsed"`
	CommodityData []CommodityData `json:"commodity_data"`
	MinPrice      string          `json:"min_price"`
	MaxPrice      string          `json:"max_price"`
	MedPrice      string          `json:"med_price"`
	AvgPrice      string          `json:"avg_price"`
}

type CommodityData struct {
	UUID      string `json:"uuid"`
	Komoditas string `json:"komoditas"`
	AreaKota  string `json:"area_kota"`
	Size      string `json:"size"`
	Price     string `json:"price"`
	PriceUSD  string `json:"price_usd"`
	Timestamp string `json:"timestamp"`
}

// USDtoIDR
type USDtoIDR struct {
	USDIDR float64 `json:"USD_IDR"`
}

type sortByPrice []Commodity

func (s sortByPrice) Len() int           { return len(s) }
func (s sortByPrice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByPrice) Less(i, j int) bool { return s[i].AreaProvinsi < s[j].AreaProvinsi }

func tokenCheck(token string, endpoint string, w http.ResponseWriter) (resp []byte) {
	log.Println(os.Getenv("AUTH_APP") + "api/")
	reqURL, _ := url.Parse(os.Getenv("AUTH_APP") + "api/" + endpoint)

	// create a request object
	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {token},
		},
	}

	// send an HTTP request using `req` object
	res, err := http.DefaultClient.Do(req)

	// check for response error
	if err != nil {
		log.Fatal("Error:", err)
	}

	// read response body
	data, _ := ioutil.ReadAll(res.Body)

	// close response body
	res.Body.Close()

	// print response status and body
	fmt.Printf("status: %d\n", res.StatusCode)
	if res.StatusCode == 401 {
		fmt.Fprintf(w, "%s\n", data)
	}
	return data
}

func fetchCommodityList() (commoditylist CommodityList) {
	var usdIdr float64
	resp, err := http.Get(os.Getenv("STEIN_EFISHERY_HOST"))
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// sb := string(body)
	// log.Println(sb)
	data := []Commodity{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
		//json: Unmarshal(non-pointer main.Request)
	}

	usdKey := "usd_to_idr"
	v, _ := rds.RedisClient.Get(context.Background(), usdKey).Bytes()
	err = json.Unmarshal(v, &usdIdr)
	var usdPrice float64
	if int(usdIdr) != 0 {
		usdPrice = usdIdr
	} else {
		usdPrice = getUSDtoIDR()
		rds.RedisClient.Set(context.Background(), usdKey, usdPrice, 5*time.Minute)
	}
	commoditylist.CommodityList = MappingList(data, usdPrice)

	return commoditylist
}

func getUSDtoIDR() (idr float64) {
	resp, err := http.Get(os.Getenv("CURR_CONV_HOST") + "?q=USD_IDR&compact=ultra&apiKey=" + os.Getenv("CURR_CONV_API_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := USDtoIDR{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	idr = data.USDIDR
	return idr

}

func MappingList(data []Commodity, idr float64) (dataMap []Commodity) {
	dataMap = make([]Commodity, 0)

	for _, val := range data {
		if val.UUID == "" {
			continue
		}
		intPrice, _ := strconv.Atoi(val.Price)
		usdPrice := float64(intPrice) * idr
		strPrice := strconv.Itoa(int(usdPrice))

		commo := Commodity{
			UUID:         val.UUID,
			Komoditas:    val.Komoditas,
			AreaProvinsi: val.AreaProvinsi,
			AreaKota:     val.AreaKota,
			Size:         val.Size,
			Price:        val.Price,
			PriceUSD:     strPrice,
			TglParsed:    val.TglParsed,
			Timestamp:    val.Timestamp,
		}
		dataMap = append(dataMap, commo)
	}

	return dataMap
}

//groupByProv grouping
func groupByProv(list []Commodity) [][]Commodity {
	sort.Sort(sortByPrice(list))
	returnData := make([][]Commodity, 0)
	i := 0
	var j int
	for {
		if i >= len(list) {
			break
		}
		for j = i + 1; j < len(list) && list[i].AreaProvinsi == list[j].AreaProvinsi; j++ {
		}

		returnData = append(returnData, list[i:j])
		i = j
	}
	return returnData
}

func MappingAggregation(data [][]Commodity) (dataMap []CommodityAggr) {
	dataMap = make([]CommodityAggr, 0)

	i := 0
	for _, val := range data {

		dataComm := make([]CommodityData, 0)
		for _, valu := range val {
			dataCom := CommodityData{
				UUID:      valu.UUID,
				Komoditas: valu.Komoditas,
				AreaKota:  valu.AreaKota,
				Price:     valu.Price,
				PriceUSD:  valu.PriceUSD,
				Size:      valu.Size,
			}
			dataComm = append(dataComm, dataCom)
		}
		// log.Println(data[i][0].AreaProvinsi)
		min, max := FindMaxAndMin(data[i])
		avg := GetAverage(data[i])
		med := GetMedian(data[i])
		comm := CommodityAggr{
			AreaProvinsi:  data[i][0].AreaProvinsi,
			TglParsed:     data[i][0].TglParsed,
			CommodityData: dataComm,
			MinPrice:      min.Price,
			MaxPrice:      max.Price,
			MedPrice:      med,
			AvgPrice:      strconv.FormatFloat(avg, 'f', 6, 64),
		}
		dataMap = append(dataMap, comm)
		i++
	}

	return dataMap
}

func FindMaxAndMin(commodities []Commodity) (min Commodity, max Commodity) {
	min = commodities[0]
	max = commodities[0]
	for _, commodity := range commodities {
		comPrice, _ := strconv.ParseFloat(commodity.Price, 64)
		maxPrice, _ := strconv.ParseFloat(max.Price, 64)
		minPrice, _ := strconv.ParseFloat(min.Price, 64)
		if comPrice > maxPrice {
			max = commodity
		}
		if comPrice < minPrice {
			min = commodity
		}
	}
	return min, max
}

func GetAverage(commodities []Commodity) (avg float64) {
	sum := 0
	for i := 0; i < len(commodities); i++ {
		minPrice, _ := strconv.ParseFloat(commodities[i].Price, 64)
		sum += (int(minPrice))
	}

	avg = (float64(sum)) / (float64(len(commodities)))

	return avg
}

func GetMedian(commodities []Commodity) (med string) {
	// var median float64
	median := len(commodities) / 2
	// rounded := math.Round(median)
	med = commodities[median].Price
	return med
}

func comodityList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	fmt.Println("Endpoint Hit: Commodity Page")
	endpoint := "users"
	// get request URL
	resp := tokenCheck(token, endpoint, w)
	data := Users{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println(err.Error())
		//json: Unmarshal(non-pointer main.Request)
	}
	response := fetchCommodityList()
	json.NewEncoder(w).Encode(response)
}

func comodityAggregation(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	fmt.Println("Endpoint Hit: Commodity Aggregation Page")
	endpoint := "users"
	// get request URL
	resp := tokenCheck(token, endpoint, w)
	data := Users{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println(err.Error())
		//json: Unmarshal(non-pointer main.Request)
	}
	if data.User.Role == "admin" {
		response := fetchCommodityList()
		resp := groupByProv(response.CommodityList)
		mapp := MappingAggregation(resp)
		json.NewEncoder(w).Encode(mapp)
	} else {
		response := fetchCommodityList()
		json.NewEncoder(w).Encode(response)
	}
}

func privateClaim(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	fmt.Println("Endpoint Hit: Commodity Aggregation Page")
	endpoint := "users"
	// get request URL
	resp := tokenCheck(token, endpoint, w)
	data := Users{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println(err.Error())
		//json: Unmarshal(non-pointer main.Request)
	}

	json.NewEncoder(w).Encode(data)
}

func handleRequests() {
	http.HandleFunc("/commodity_list", comodityList)
	http.HandleFunc("/commodity_aggregation", comodityAggregation)
	http.HandleFunc("/private_claim", privateClaim)
	log.Println("Server localhost:8080 is running")
	rds.Initiate()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
