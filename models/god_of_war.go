package models

type GodOfWar struct {
	TableName struct{} `json:"-" db:"god_of_war" pk:"Id"`
	Analysis  struct {
		TableName    struct{} `json:"-" db:"analysis" pk:"Id"`
		Active       string   `json:"active" db:"active" type:"string"`
		DpvCmra      string   `json:"dpv_cmra" db:"dpv_cmra" type:"string"`
		DpvFootnotes string   `json:"dpv_footnotes" db:"dpv_footnotes" type:"string"`
		DpvMatchCode string   `json:"dpv_match_code" db:"dpv_match_code" type:"string"`
		DpvVacant    string   `json:"dpv_vacant" db:"dpv_vacant" type:"string"`
	} `json:"analysis" db:"-" type:"-"`
	CandidateIndex int `json:"candidate_index" db:"candidate_index" type:"int"`
	Components     struct {
		TableName               struct{} `json:"-" db:"components" pk:"Id"`
		CityName                string   `json:"city_name" db:"city_name" type:"string"`
		DeliveryPoint           string   `json:"delivery_point" db:"delivery_point" type:"string"`
		DeliveryPointCheckDigit string   `json:"delivery_point_check_digit" db:"delivery_point_check_digit" type:"string"`
		Plus4Code               string   `json:"plus4_code" db:"plus4_code" type:"string"`
		PrimaryNumber           string   `json:"primary_number" db:"primary_number" type:"string"`
		StateAbbreviation       string   `json:"state_abbreviation" db:"state_abbreviation" type:"string"`
		StreetName              string   `json:"street_name" db:"street_name" type:"string"`
		StreetPredirection      string   `json:"street_predirection" db:"street_predirection" type:"string"`
		StreetSuffix            string   `json:"street_suffix" db:"street_suffix" type:"string"`
		Zipcode                 string   `json:"zipcode" db:"zipcode" type:"string"`
	} `json:"components" db:"-" type:"-"`
	DeliveryLine1        string `json:"delivery_line_1" db:"delivery_line_1" type:"string"`
	DeliveryPointBarcode string `json:"delivery_point_barcode" db:"delivery_point_barcode" type:"string"`
	InputIndex           int    `json:"input_index" db:"input_index" type:"int"`
	LastLine             string `json:"last_line" db:"last_line" type:"string"`
	Metadata             struct {
		TableName             struct{} `json:"-" db:"metadata" pk:"Id"`
		CarrierRoute          string   `json:"carrier_route" db:"carrier_route" type:"string"`
		CongressionalDistrict string   `json:"congressional_district" db:"congressional_district" type:"string"`
		CountyFips            string   `json:"county_fips" db:"county_fips" type:"string"`
		CountyName            string   `json:"county_name" db:"county_name" type:"string"`
		Dst                   bool     `json:"dst" db:"dst" type:"bool"`
		ElotSequence          string   `json:"elot_sequence" db:"elot_sequence" type:"string"`
		ElotSort              string   `json:"elot_sort" db:"elot_sort" type:"string"`
		Latitude              float64  `json:"latitude" db:"latitude" type:"float64"`
		Longitude             float64  `json:"longitude" db:"longitude" type:"float64"`
		Precision             string   `json:"precision" db:"precision" type:"string"`
		Rdi                   string   `json:"rdi" db:"rdi" type:"string"`
		RecordType            string   `json:"record_type" db:"record_type" type:"string"`
		TimeZone              string   `json:"time_zone" db:"time_zone" type:"string"`
		UtcOffset             int      `json:"utc_offset" db:"utc_offset" type:"int"`
		ZipType               string   `json:"zip_type" db:"zip_type" type:"string"`
	} `json:"metadata" db:"-" type:"-"`
}
