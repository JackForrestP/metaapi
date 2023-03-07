//Auto generated with MetaApi https://github.com/exyzzy/metaapi
package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)
import	"time"




// ======= Community_result =======

//Create Table
func CreateTableCommunity_results(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS community_results CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table community_results ( id integer generated always as identity primary key , uuid varchar(36) not null UNIQUE , product_uuid varchar(36) not null UNIQUE REFERENCES products(uuid) , metric_graphs text not null ) ; `)
	return
}

//Drop Table
func DropTableCommunity_results(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS community_results CASCADE")
	return
}

//Struct
type Community_result struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	MetricGraphs string`xml:"MetricGraphs" json:"metricgraphs"`

}

//Create
func (community_result *Community_result) CreateCommunity_result(db *sql.DB) (result Community_result, err error) {
	stmt, err := db.Prepare("INSERT INTO community_results ( uuid, product_uuid, metric_graphs) VALUES ($1,$2,$3) RETURNING id, uuid, product_uuid, metric_graphs")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( community_result.Uuid, community_result.ProductUuid, community_result.MetricGraphs).Scan( &result.Id, &result.Uuid, &result.ProductUuid, &result.MetricGraphs)
	return
}

//Retrieve
func (community_result *Community_result) RetrieveCommunity_result(db *sql.DB) (result Community_result, err error) {
	result = Community_result{}
	err = db.QueryRow("SELECT id, uuid, product_uuid, metric_graphs FROM community_results WHERE (id = $1)", community_result.Id).Scan( &result.Id, &result.Uuid, &result.ProductUuid, &result.MetricGraphs)
	return
}

//RetrieveAll
func RetrieveAllCommunity_results(db *sql.DB) (community_results []Community_result, err error) {
	rows, err := db.Query("SELECT id, uuid, product_uuid, metric_graphs FROM community_results ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Community_result{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.ProductUuid, &result.MetricGraphs); err != nil {
			return
		}
		community_results = append(community_results, result)
	}
	rows.Close()
	return
}

//Update
func (community_result *Community_result) UpdateCommunity_result(db *sql.DB) (result Community_result, err error) {
	stmt, err := db.Prepare("UPDATE community_results SET uuid = $2, product_uuid = $3, metric_graphs = $4 WHERE (id = $1) RETURNING id, uuid, product_uuid, metric_graphs")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( community_result.Id, community_result.Uuid, community_result.ProductUuid, community_result.MetricGraphs).Scan( &result.Id, &result.Uuid, &result.ProductUuid, &result.MetricGraphs)
	return
}

//Delete
func (community_result *Community_result) DeleteCommunity_result(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM community_results WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(community_result.Id)
	return
}

//DeleteAll
func DeleteAllCommunity_results(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM community_results")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Product_review =======

//Create Table
func CreateTableProduct_reviews(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_reviews CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table product_reviews ( id integer generated always as identity primary key , uuid varchar(36) not null UNIQUE , title text not null UNIQUE , author text not null , creation_date timestamptz not null , last_updated timestamptz not null , hero_image text not null , summary text not null , rating text not null , score decimal not null , body_html text not null , user_id integer not null REFERENCES users(id) , product_uuid varchar(36) not null REFERENCES products(uuid) , expected_tst integer , expected_deep integer , expected_rem integer , expected_onset integer , expected_wakefulness integer , expected_efficiency integer , expected_accuracy integer , reviewed_results text , validated boolean not null , start_date timestamptz not null , end_date timestamptz not null , is_preview boolean not null ) ; `)
	return
}

//Drop Table
func DropTableProduct_reviews(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_reviews CASCADE")
	return
}

//Struct
type Product_review struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	Title string`xml:"Title" json:"title"`
	Author string`xml:"Author" json:"author"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`
	LastUpdated time.Time`xml:"LastUpdated" json:"lastupdated"`
	HeroImage string`xml:"HeroImage" json:"heroimage"`
	Summary string`xml:"Summary" json:"summary"`
	Rating string`xml:"Rating" json:"rating"`
	Score float64`xml:"Score" json:"score"`
	BodyHtml string`xml:"BodyHtml" json:"bodyhtml"`
	UserId int32`xml:"UserId" json:"userid"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	ExpectedTst sql.NullInt32`xml:"ExpectedTst" json:"expectedtst"`
	ExpectedDeep sql.NullInt32`xml:"ExpectedDeep" json:"expecteddeep"`
	ExpectedRem sql.NullInt32`xml:"ExpectedRem" json:"expectedrem"`
	ExpectedOnset sql.NullInt32`xml:"ExpectedOnset" json:"expectedonset"`
	ExpectedWakefulness sql.NullInt32`xml:"ExpectedWakefulness" json:"expectedwakefulness"`
	ExpectedEfficiency sql.NullInt32`xml:"ExpectedEfficiency" json:"expectedefficiency"`
	ExpectedAccuracy sql.NullInt32`xml:"ExpectedAccuracy" json:"expectedaccuracy"`
	ReviewedResults sql.NullString`xml:"ReviewedResults" json:"reviewedresults"`
	Validated bool`xml:"Validated" json:"validated"`
	StartDate time.Time`xml:"StartDate" json:"startdate"`
	EndDate time.Time`xml:"EndDate" json:"enddate"`
	IsPreview bool`xml:"IsPreview" json:"ispreview"`

}

//Create
func (product_review *Product_review) CreateProduct_review(db *sql.DB) (result Product_review, err error) {
	stmt, err := db.Prepare("INSERT INTO product_reviews ( uuid, title, author, creation_date, last_updated, hero_image, summary, rating, score, body_html, user_id, product_uuid, expected_tst, expected_deep, expected_rem, expected_onset, expected_wakefulness, expected_efficiency, expected_accuracy, reviewed_results, validated, start_date, end_date, is_preview) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) RETURNING id, uuid, title, author, creation_date, last_updated, hero_image, summary, rating, score, body_html, user_id, product_uuid, expected_tst, expected_deep, expected_rem, expected_onset, expected_wakefulness, expected_efficiency, expected_accuracy, reviewed_results, validated, start_date, end_date, is_preview")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( product_review.Uuid, product_review.Title, product_review.Author, product_review.CreationDate, product_review.LastUpdated, product_review.HeroImage, product_review.Summary, product_review.Rating, product_review.Score, product_review.BodyHtml, product_review.UserId, product_review.ProductUuid, product_review.ExpectedTst, product_review.ExpectedDeep, product_review.ExpectedRem, product_review.ExpectedOnset, product_review.ExpectedWakefulness, product_review.ExpectedEfficiency, product_review.ExpectedAccuracy, product_review.ReviewedResults, product_review.Validated, product_review.StartDate, product_review.EndDate, product_review.IsPreview).Scan( &result.Id, &result.Uuid, &result.Title, &result.Author, &result.CreationDate, &result.LastUpdated, &result.HeroImage, &result.Summary, &result.Rating, &result.Score, &result.BodyHtml, &result.UserId, &result.ProductUuid, &result.ExpectedTst, &result.ExpectedDeep, &result.ExpectedRem, &result.ExpectedOnset, &result.ExpectedWakefulness, &result.ExpectedEfficiency, &result.ExpectedAccuracy, &result.ReviewedResults, &result.Validated, &result.StartDate, &result.EndDate, &result.IsPreview)
	return
}

//Retrieve
func (product_review *Product_review) RetrieveProduct_review(db *sql.DB) (result Product_review, err error) {
	result = Product_review{}
	err = db.QueryRow("SELECT id, uuid, title, author, creation_date, last_updated, hero_image, summary, rating, score, body_html, user_id, product_uuid, expected_tst, expected_deep, expected_rem, expected_onset, expected_wakefulness, expected_efficiency, expected_accuracy, reviewed_results, validated, start_date, end_date, is_preview FROM product_reviews WHERE (id = $1)", product_review.Id).Scan( &result.Id, &result.Uuid, &result.Title, &result.Author, &result.CreationDate, &result.LastUpdated, &result.HeroImage, &result.Summary, &result.Rating, &result.Score, &result.BodyHtml, &result.UserId, &result.ProductUuid, &result.ExpectedTst, &result.ExpectedDeep, &result.ExpectedRem, &result.ExpectedOnset, &result.ExpectedWakefulness, &result.ExpectedEfficiency, &result.ExpectedAccuracy, &result.ReviewedResults, &result.Validated, &result.StartDate, &result.EndDate, &result.IsPreview)
	return
}

//RetrieveAll
func RetrieveAllProduct_reviews(db *sql.DB) (product_reviews []Product_review, err error) {
	rows, err := db.Query("SELECT id, uuid, title, author, creation_date, last_updated, hero_image, summary, rating, score, body_html, user_id, product_uuid, expected_tst, expected_deep, expected_rem, expected_onset, expected_wakefulness, expected_efficiency, expected_accuracy, reviewed_results, validated, start_date, end_date, is_preview FROM product_reviews ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Product_review{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.Title, &result.Author, &result.CreationDate, &result.LastUpdated, &result.HeroImage, &result.Summary, &result.Rating, &result.Score, &result.BodyHtml, &result.UserId, &result.ProductUuid, &result.ExpectedTst, &result.ExpectedDeep, &result.ExpectedRem, &result.ExpectedOnset, &result.ExpectedWakefulness, &result.ExpectedEfficiency, &result.ExpectedAccuracy, &result.ReviewedResults, &result.Validated, &result.StartDate, &result.EndDate, &result.IsPreview); err != nil {
			return
		}
		product_reviews = append(product_reviews, result)
	}
	rows.Close()
	return
}

//Update
func (product_review *Product_review) UpdateProduct_review(db *sql.DB) (result Product_review, err error) {
	stmt, err := db.Prepare("UPDATE product_reviews SET uuid = $2, title = $3, author = $4, creation_date = $5, last_updated = $6, hero_image = $7, summary = $8, rating = $9, score = $10, body_html = $11, user_id = $12, product_uuid = $13, expected_tst = $14, expected_deep = $15, expected_rem = $16, expected_onset = $17, expected_wakefulness = $18, expected_efficiency = $19, expected_accuracy = $20, reviewed_results = $21, validated = $22, start_date = $23, end_date = $24, is_preview = $25 WHERE (id = $1) RETURNING id, uuid, title, author, creation_date, last_updated, hero_image, summary, rating, score, body_html, user_id, product_uuid, expected_tst, expected_deep, expected_rem, expected_onset, expected_wakefulness, expected_efficiency, expected_accuracy, reviewed_results, validated, start_date, end_date, is_preview")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( product_review.Id, product_review.Uuid, product_review.Title, product_review.Author, product_review.CreationDate, product_review.LastUpdated, product_review.HeroImage, product_review.Summary, product_review.Rating, product_review.Score, product_review.BodyHtml, product_review.UserId, product_review.ProductUuid, product_review.ExpectedTst, product_review.ExpectedDeep, product_review.ExpectedRem, product_review.ExpectedOnset, product_review.ExpectedWakefulness, product_review.ExpectedEfficiency, product_review.ExpectedAccuracy, product_review.ReviewedResults, product_review.Validated, product_review.StartDate, product_review.EndDate, product_review.IsPreview).Scan( &result.Id, &result.Uuid, &result.Title, &result.Author, &result.CreationDate, &result.LastUpdated, &result.HeroImage, &result.Summary, &result.Rating, &result.Score, &result.BodyHtml, &result.UserId, &result.ProductUuid, &result.ExpectedTst, &result.ExpectedDeep, &result.ExpectedRem, &result.ExpectedOnset, &result.ExpectedWakefulness, &result.ExpectedEfficiency, &result.ExpectedAccuracy, &result.ReviewedResults, &result.Validated, &result.StartDate, &result.EndDate, &result.IsPreview)
	return
}

//Delete
func (product_review *Product_review) DeleteProduct_review(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_reviews WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product_review.Id)
	return
}

//DeleteAll
func DeleteAllProduct_reviews(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_reviews")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Product_listing =======

//Create Table
func CreateTableProduct_listings(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_listings CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table product_listings ( id integer generated always as identity primary key , uuid varchar(36) not null UNIQUE , vendor_uuid varchar(36) not null REFERENCES vendors(uuid) , creation_date timestamptz not null , last_updated timestamptz not null , purchase_link text not null , price decimal not null , product_uuid varchar(36) not null REFERENCES products(uuid) , is_manufacterer boolean not null ) ; `)
	return
}

//Drop Table
func DropTableProduct_listings(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_listings CASCADE")
	return
}

//Struct
type Product_listing struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	VendorUuid string`xml:"VendorUuid" json:"vendoruuid"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`
	LastUpdated time.Time`xml:"LastUpdated" json:"lastupdated"`
	PurchaseLink string`xml:"PurchaseLink" json:"purchaselink"`
	Price float64`xml:"Price" json:"price"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	IsManufacterer bool`xml:"IsManufacterer" json:"ismanufacterer"`

}

//Create
func (product_listing *Product_listing) CreateProduct_listing(db *sql.DB) (result Product_listing, err error) {
	stmt, err := db.Prepare("INSERT INTO product_listings ( uuid, vendor_uuid, creation_date, last_updated, purchase_link, price, product_uuid, is_manufacterer) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, uuid, vendor_uuid, creation_date, last_updated, purchase_link, price, product_uuid, is_manufacterer")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( product_listing.Uuid, product_listing.VendorUuid, product_listing.CreationDate, product_listing.LastUpdated, product_listing.PurchaseLink, product_listing.Price, product_listing.ProductUuid, product_listing.IsManufacterer).Scan( &result.Id, &result.Uuid, &result.VendorUuid, &result.CreationDate, &result.LastUpdated, &result.PurchaseLink, &result.Price, &result.ProductUuid, &result.IsManufacterer)
	return
}

//Retrieve
func (product_listing *Product_listing) RetrieveProduct_listing(db *sql.DB) (result Product_listing, err error) {
	result = Product_listing{}
	err = db.QueryRow("SELECT id, uuid, vendor_uuid, creation_date, last_updated, purchase_link, price, product_uuid, is_manufacterer FROM product_listings WHERE (id = $1)", product_listing.Id).Scan( &result.Id, &result.Uuid, &result.VendorUuid, &result.CreationDate, &result.LastUpdated, &result.PurchaseLink, &result.Price, &result.ProductUuid, &result.IsManufacterer)
	return
}

//RetrieveAll
func RetrieveAllProduct_listings(db *sql.DB) (product_listings []Product_listing, err error) {
	rows, err := db.Query("SELECT id, uuid, vendor_uuid, creation_date, last_updated, purchase_link, price, product_uuid, is_manufacterer FROM product_listings ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Product_listing{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.VendorUuid, &result.CreationDate, &result.LastUpdated, &result.PurchaseLink, &result.Price, &result.ProductUuid, &result.IsManufacterer); err != nil {
			return
		}
		product_listings = append(product_listings, result)
	}
	rows.Close()
	return
}

//Update
func (product_listing *Product_listing) UpdateProduct_listing(db *sql.DB) (result Product_listing, err error) {
	stmt, err := db.Prepare("UPDATE product_listings SET uuid = $2, vendor_uuid = $3, creation_date = $4, last_updated = $5, purchase_link = $6, price = $7, product_uuid = $8, is_manufacterer = $9 WHERE (id = $1) RETURNING id, uuid, vendor_uuid, creation_date, last_updated, purchase_link, price, product_uuid, is_manufacterer")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( product_listing.Id, product_listing.Uuid, product_listing.VendorUuid, product_listing.CreationDate, product_listing.LastUpdated, product_listing.PurchaseLink, product_listing.Price, product_listing.ProductUuid, product_listing.IsManufacterer).Scan( &result.Id, &result.Uuid, &result.VendorUuid, &result.CreationDate, &result.LastUpdated, &result.PurchaseLink, &result.Price, &result.ProductUuid, &result.IsManufacterer)
	return
}

//Delete
func (product_listing *Product_listing) DeleteProduct_listing(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_listings WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product_listing.Id)
	return
}

//DeleteAll
func DeleteAllProduct_listings(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_listings")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Photo_url =======

//Create Table
func CreateTablePhoto_urls(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS photo_urls CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table photo_urls ( id integer generated always as identity primary key , uuid varchar(36) not null UNIQUE , url text not null , title text not null , alt_text text not null , keywords text not null , user_id integer not null REFERENCES users(id) , creation_date timestamptz not null ) ; `)
	return
}

//Drop Table
func DropTablePhoto_urls(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS photo_urls CASCADE")
	return
}

//Struct
type Photo_url struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	Url string`xml:"Url" json:"url"`
	Title string`xml:"Title" json:"title"`
	AltText string`xml:"AltText" json:"alttext"`
	Keywords string`xml:"Keywords" json:"keywords"`
	UserId int32`xml:"UserId" json:"userid"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`

}

//Create
func (photo_url *Photo_url) CreatePhoto_url(db *sql.DB) (result Photo_url, err error) {
	stmt, err := db.Prepare("INSERT INTO photo_urls ( uuid, url, title, alt_text, keywords, user_id, creation_date) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, uuid, url, title, alt_text, keywords, user_id, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( photo_url.Uuid, photo_url.Url, photo_url.Title, photo_url.AltText, photo_url.Keywords, photo_url.UserId, photo_url.CreationDate).Scan( &result.Id, &result.Uuid, &result.Url, &result.Title, &result.AltText, &result.Keywords, &result.UserId, &result.CreationDate)
	return
}

//Retrieve
func (photo_url *Photo_url) RetrievePhoto_url(db *sql.DB) (result Photo_url, err error) {
	result = Photo_url{}
	err = db.QueryRow("SELECT id, uuid, url, title, alt_text, keywords, user_id, creation_date FROM photo_urls WHERE (id = $1)", photo_url.Id).Scan( &result.Id, &result.Uuid, &result.Url, &result.Title, &result.AltText, &result.Keywords, &result.UserId, &result.CreationDate)
	return
}

//RetrieveAll
func RetrieveAllPhoto_urls(db *sql.DB) (photo_urls []Photo_url, err error) {
	rows, err := db.Query("SELECT id, uuid, url, title, alt_text, keywords, user_id, creation_date FROM photo_urls ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Photo_url{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.Url, &result.Title, &result.AltText, &result.Keywords, &result.UserId, &result.CreationDate); err != nil {
			return
		}
		photo_urls = append(photo_urls, result)
	}
	rows.Close()
	return
}

//Update
func (photo_url *Photo_url) UpdatePhoto_url(db *sql.DB) (result Photo_url, err error) {
	stmt, err := db.Prepare("UPDATE photo_urls SET uuid = $2, url = $3, title = $4, alt_text = $5, keywords = $6, user_id = $7, creation_date = $8 WHERE (id = $1) RETURNING id, uuid, url, title, alt_text, keywords, user_id, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( photo_url.Id, photo_url.Uuid, photo_url.Url, photo_url.Title, photo_url.AltText, photo_url.Keywords, photo_url.UserId, photo_url.CreationDate).Scan( &result.Id, &result.Uuid, &result.Url, &result.Title, &result.AltText, &result.Keywords, &result.UserId, &result.CreationDate)
	return
}

//Delete
func (photo_url *Photo_url) DeletePhoto_url(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM photo_urls WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(photo_url.Id)
	return
}

//DeleteAll
func DeleteAllPhoto_urls(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM photo_urls")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Vendor =======

//Create Table
func CreateTableVendors(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS vendors CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table vendors ( id integer generated always as identity primary key , uuid varchar(36) not null UNIQUE , name text not null UNIQUE , link text , creation_date timestamptz not null , last_updated timestamptz not null , affiliate_id text ) ; `)
	return
}

//Drop Table
func DropTableVendors(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS vendors CASCADE")
	return
}

//Struct
type Vendor struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	Name string`xml:"Name" json:"name"`
	Link sql.NullString`xml:"Link" json:"link"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`
	LastUpdated time.Time`xml:"LastUpdated" json:"lastupdated"`
	AffiliateId sql.NullString`xml:"AffiliateId" json:"affiliateid"`

}

//Create
func (vendor *Vendor) CreateVendor(db *sql.DB) (result Vendor, err error) {
	stmt, err := db.Prepare("INSERT INTO vendors ( uuid, name, link, creation_date, last_updated, affiliate_id) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, uuid, name, link, creation_date, last_updated, affiliate_id")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( vendor.Uuid, vendor.Name, vendor.Link, vendor.CreationDate, vendor.LastUpdated, vendor.AffiliateId).Scan( &result.Id, &result.Uuid, &result.Name, &result.Link, &result.CreationDate, &result.LastUpdated, &result.AffiliateId)
	return
}

//Retrieve
func (vendor *Vendor) RetrieveVendor(db *sql.DB) (result Vendor, err error) {
	result = Vendor{}
	err = db.QueryRow("SELECT id, uuid, name, link, creation_date, last_updated, affiliate_id FROM vendors WHERE (id = $1)", vendor.Id).Scan( &result.Id, &result.Uuid, &result.Name, &result.Link, &result.CreationDate, &result.LastUpdated, &result.AffiliateId)
	return
}

//RetrieveAll
func RetrieveAllVendors(db *sql.DB) (vendors []Vendor, err error) {
	rows, err := db.Query("SELECT id, uuid, name, link, creation_date, last_updated, affiliate_id FROM vendors ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Vendor{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.Name, &result.Link, &result.CreationDate, &result.LastUpdated, &result.AffiliateId); err != nil {
			return
		}
		vendors = append(vendors, result)
	}
	rows.Close()
	return
}

//Update
func (vendor *Vendor) UpdateVendor(db *sql.DB) (result Vendor, err error) {
	stmt, err := db.Prepare("UPDATE vendors SET uuid = $2, name = $3, link = $4, creation_date = $5, last_updated = $6, affiliate_id = $7 WHERE (id = $1) RETURNING id, uuid, name, link, creation_date, last_updated, affiliate_id")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( vendor.Id, vendor.Uuid, vendor.Name, vendor.Link, vendor.CreationDate, vendor.LastUpdated, vendor.AffiliateId).Scan( &result.Id, &result.Uuid, &result.Name, &result.Link, &result.CreationDate, &result.LastUpdated, &result.AffiliateId)
	return
}

//Delete
func (vendor *Vendor) DeleteVendor(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM vendors WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(vendor.Id)
	return
}

//DeleteAll
func DeleteAllVendors(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM vendors")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Solution =======

//Create Table
func CreateTableSolutions(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS solutions CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table solutions ( id integer generated always as identity primary key , uuid varchar(36) not null unique , name text not null , product_uuid varchar(36) , technique_uuid varchar(36) not null , user_id integer not null references users(id) , created_at timestamptz not null , is_public boolean not null , validated boolean not null , last_updated timestamptz not null , archived boolean not null , description text ) ; `)
	return
}

//Drop Table
func DropTableSolutions(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS solutions CASCADE")
	return
}

//Struct
type Solution struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	Name string`xml:"Name" json:"name"`
	ProductUuid sql.NullString`xml:"ProductUuid" json:"productuuid"`
	TechniqueUuid string`xml:"TechniqueUuid" json:"techniqueuuid"`
	UserId int32`xml:"UserId" json:"userid"`
	CreatedAt time.Time`xml:"CreatedAt" json:"createdat"`
	IsPublic bool`xml:"IsPublic" json:"ispublic"`
	Validated bool`xml:"Validated" json:"validated"`
	LastUpdated time.Time`xml:"LastUpdated" json:"lastupdated"`
	Archived bool`xml:"Archived" json:"archived"`
	Description sql.NullString`xml:"Description" json:"description"`

}

//Create
func (solution *Solution) CreateSolution(db *sql.DB) (result Solution, err error) {
	stmt, err := db.Prepare("INSERT INTO solutions ( uuid, name, product_uuid, technique_uuid, user_id, created_at, is_public, validated, last_updated, archived, description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, uuid, name, product_uuid, technique_uuid, user_id, created_at, is_public, validated, last_updated, archived, description")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( solution.Uuid, solution.Name, solution.ProductUuid, solution.TechniqueUuid, solution.UserId, solution.CreatedAt, solution.IsPublic, solution.Validated, solution.LastUpdated, solution.Archived, solution.Description).Scan( &result.Id, &result.Uuid, &result.Name, &result.ProductUuid, &result.TechniqueUuid, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.LastUpdated, &result.Archived, &result.Description)
	return
}

//Retrieve
func (solution *Solution) RetrieveSolution(db *sql.DB) (result Solution, err error) {
	result = Solution{}
	err = db.QueryRow("SELECT id, uuid, name, product_uuid, technique_uuid, user_id, created_at, is_public, validated, last_updated, archived, description FROM solutions WHERE (id = $1)", solution.Id).Scan( &result.Id, &result.Uuid, &result.Name, &result.ProductUuid, &result.TechniqueUuid, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.LastUpdated, &result.Archived, &result.Description)
	return
}

//RetrieveAll
func RetrieveAllSolutions(db *sql.DB) (solutions []Solution, err error) {
	rows, err := db.Query("SELECT id, uuid, name, product_uuid, technique_uuid, user_id, created_at, is_public, validated, last_updated, archived, description FROM solutions ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Solution{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.Name, &result.ProductUuid, &result.TechniqueUuid, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.LastUpdated, &result.Archived, &result.Description); err != nil {
			return
		}
		solutions = append(solutions, result)
	}
	rows.Close()
	return
}

//Update
func (solution *Solution) UpdateSolution(db *sql.DB) (result Solution, err error) {
	stmt, err := db.Prepare("UPDATE solutions SET uuid = $2, name = $3, product_uuid = $4, technique_uuid = $5, user_id = $6, created_at = $7, is_public = $8, validated = $9, last_updated = $10, archived = $11, description = $12 WHERE (id = $1) RETURNING id, uuid, name, product_uuid, technique_uuid, user_id, created_at, is_public, validated, last_updated, archived, description")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( solution.Id, solution.Uuid, solution.Name, solution.ProductUuid, solution.TechniqueUuid, solution.UserId, solution.CreatedAt, solution.IsPublic, solution.Validated, solution.LastUpdated, solution.Archived, solution.Description).Scan( &result.Id, &result.Uuid, &result.Name, &result.ProductUuid, &result.TechniqueUuid, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.LastUpdated, &result.Archived, &result.Description)
	return
}

//Delete
func (solution *Solution) DeleteSolution(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM solutions WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(solution.Id)
	return
}

//DeleteAll
func DeleteAllSolutions(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM solutions")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Product =======

//Create Table
func CreateTableProducts(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS products CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table products ( id integer generated always as identity primary key , uuid varchar(36) not null unique , name text not null , model text not null , category integer not null references product_categories(id) , manufacturer text not null , image_link text not null , purchase_link text not null , affiliated_id text not null , description text not null , user_id integer not null references users(id) , created_at timestamptz not null , is_public boolean not null , validated boolean not null , archived boolean not null , manufacturer_uuid varchar(36) refrences vendor(uuid) , manufacturer_product_id text , manufacturer_product_link text , msrp decimal ) ; `)
	return
}

//Drop Table
func DropTableProducts(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS products CASCADE")
	return
}

//Struct
type Product struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	Name string`xml:"Name" json:"name"`
	Model string`xml:"Model" json:"model"`
	Category int32`xml:"Category" json:"category"`
	Manufacturer string`xml:"Manufacturer" json:"manufacturer"`
	ImageLink string`xml:"ImageLink" json:"imagelink"`
	PurchaseLink string`xml:"PurchaseLink" json:"purchaselink"`
	AffiliatedId string`xml:"AffiliatedId" json:"affiliatedid"`
	Description string`xml:"Description" json:"description"`
	UserId int32`xml:"UserId" json:"userid"`
	CreatedAt time.Time`xml:"CreatedAt" json:"createdat"`
	IsPublic bool`xml:"IsPublic" json:"ispublic"`
	Validated bool`xml:"Validated" json:"validated"`
	Archived bool`xml:"Archived" json:"archived"`
	ManufacturerUuid sql.NullString`xml:"ManufacturerUuid" json:"manufactureruuid"`
	ManufacturerProductId sql.NullString`xml:"ManufacturerProductId" json:"manufacturerproductid"`
	ManufacturerProductLink sql.NullString`xml:"ManufacturerProductLink" json:"manufacturerproductlink"`
	Msrp sql.NullFloat64`xml:"Msrp" json:"msrp"`

}

//Create
func (product *Product) CreateProduct(db *sql.DB) (result Product, err error) {
	stmt, err := db.Prepare("INSERT INTO products ( uuid, name, model, category, manufacturer, image_link, purchase_link, affiliated_id, description, user_id, created_at, is_public, validated, archived, manufacturer_uuid, manufacturer_product_id, manufacturer_product_link, msrp) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING id, uuid, name, model, category, manufacturer, image_link, purchase_link, affiliated_id, description, user_id, created_at, is_public, validated, archived, manufacturer_uuid, manufacturer_product_id, manufacturer_product_link, msrp")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( product.Uuid, product.Name, product.Model, product.Category, product.Manufacturer, product.ImageLink, product.PurchaseLink, product.AffiliatedId, product.Description, product.UserId, product.CreatedAt, product.IsPublic, product.Validated, product.Archived, product.ManufacturerUuid, product.ManufacturerProductId, product.ManufacturerProductLink, product.Msrp).Scan( &result.Id, &result.Uuid, &result.Name, &result.Model, &result.Category, &result.Manufacturer, &result.ImageLink, &result.PurchaseLink, &result.AffiliatedId, &result.Description, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.Archived, &result.ManufacturerUuid, &result.ManufacturerProductId, &result.ManufacturerProductLink, &result.Msrp)
	return
}

//Retrieve
func (product *Product) RetrieveProduct(db *sql.DB) (result Product, err error) {
	result = Product{}
	err = db.QueryRow("SELECT id, uuid, name, model, category, manufacturer, image_link, purchase_link, affiliated_id, description, user_id, created_at, is_public, validated, archived, manufacturer_uuid, manufacturer_product_id, manufacturer_product_link, msrp FROM products WHERE (id = $1)", product.Id).Scan( &result.Id, &result.Uuid, &result.Name, &result.Model, &result.Category, &result.Manufacturer, &result.ImageLink, &result.PurchaseLink, &result.AffiliatedId, &result.Description, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.Archived, &result.ManufacturerUuid, &result.ManufacturerProductId, &result.ManufacturerProductLink, &result.Msrp)
	return
}

//RetrieveAll
func RetrieveAllProducts(db *sql.DB) (products []Product, err error) {
	rows, err := db.Query("SELECT id, uuid, name, model, category, manufacturer, image_link, purchase_link, affiliated_id, description, user_id, created_at, is_public, validated, archived, manufacturer_uuid, manufacturer_product_id, manufacturer_product_link, msrp FROM products ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Product{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.Name, &result.Model, &result.Category, &result.Manufacturer, &result.ImageLink, &result.PurchaseLink, &result.AffiliatedId, &result.Description, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.Archived, &result.ManufacturerUuid, &result.ManufacturerProductId, &result.ManufacturerProductLink, &result.Msrp); err != nil {
			return
		}
		products = append(products, result)
	}
	rows.Close()
	return
}

//Update
func (product *Product) UpdateProduct(db *sql.DB) (result Product, err error) {
	stmt, err := db.Prepare("UPDATE products SET uuid = $2, name = $3, model = $4, category = $5, manufacturer = $6, image_link = $7, purchase_link = $8, affiliated_id = $9, description = $10, user_id = $11, created_at = $12, is_public = $13, validated = $14, archived = $15, manufacturer_uuid = $16, manufacturer_product_id = $17, manufacturer_product_link = $18, msrp = $19 WHERE (id = $1) RETURNING id, uuid, name, model, category, manufacturer, image_link, purchase_link, affiliated_id, description, user_id, created_at, is_public, validated, archived, manufacturer_uuid, manufacturer_product_id, manufacturer_product_link, msrp")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( product.Id, product.Uuid, product.Name, product.Model, product.Category, product.Manufacturer, product.ImageLink, product.PurchaseLink, product.AffiliatedId, product.Description, product.UserId, product.CreatedAt, product.IsPublic, product.Validated, product.Archived, product.ManufacturerUuid, product.ManufacturerProductId, product.ManufacturerProductLink, product.Msrp).Scan( &result.Id, &result.Uuid, &result.Name, &result.Model, &result.Category, &result.Manufacturer, &result.ImageLink, &result.PurchaseLink, &result.AffiliatedId, &result.Description, &result.UserId, &result.CreatedAt, &result.IsPublic, &result.Validated, &result.Archived, &result.ManufacturerUuid, &result.ManufacturerProductId, &result.ManufacturerProductLink, &result.Msrp)
	return
}

//Delete
func (product *Product) DeleteProduct(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM products WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Id)
	return
}

//DeleteAll
func DeleteAllProducts(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM products")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Article =======

//Create Table
func CreateTableArticles(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS articles CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table articles ( id integer generated always as identity primary key , uuid varchar(36) not null unique , title text not null , summary text not null , urls text not null , author text not null , category text not null , image_link text not null , creation_date timestamptz not null , user_id integer not null references users(id) , article_link text not null , body_html text not null , validated boolean not null ) ; `)
	return
}

//Drop Table
func DropTableArticles(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS articles CASCADE")
	return
}

//Struct
type Article struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	Title string`xml:"Title" json:"title"`
	Summary string`xml:"Summary" json:"summary"`
	Urls string`xml:"Urls" json:"urls"`
	Author string`xml:"Author" json:"author"`
	Category string`xml:"Category" json:"category"`
	ImageLink string`xml:"ImageLink" json:"imagelink"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`
	UserId int32`xml:"UserId" json:"userid"`
	ArticleLink string`xml:"ArticleLink" json:"articlelink"`
	BodyHtml string`xml:"BodyHtml" json:"bodyhtml"`
	Validated bool`xml:"Validated" json:"validated"`

}

//Create
func (article *Article) CreateArticle(db *sql.DB) (result Article, err error) {
	stmt, err := db.Prepare("INSERT INTO articles ( uuid, title, summary, urls, author, category, image_link, creation_date, user_id, article_link, body_html, validated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id, uuid, title, summary, urls, author, category, image_link, creation_date, user_id, article_link, body_html, validated")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( article.Uuid, article.Title, article.Summary, article.Urls, article.Author, article.Category, article.ImageLink, article.CreationDate, article.UserId, article.ArticleLink, article.BodyHtml, article.Validated).Scan( &result.Id, &result.Uuid, &result.Title, &result.Summary, &result.Urls, &result.Author, &result.Category, &result.ImageLink, &result.CreationDate, &result.UserId, &result.ArticleLink, &result.BodyHtml, &result.Validated)
	return
}

//Retrieve
func (article *Article) RetrieveArticle(db *sql.DB) (result Article, err error) {
	result = Article{}
	err = db.QueryRow("SELECT id, uuid, title, summary, urls, author, category, image_link, creation_date, user_id, article_link, body_html, validated FROM articles WHERE (id = $1)", article.Id).Scan( &result.Id, &result.Uuid, &result.Title, &result.Summary, &result.Urls, &result.Author, &result.Category, &result.ImageLink, &result.CreationDate, &result.UserId, &result.ArticleLink, &result.BodyHtml, &result.Validated)
	return
}

//RetrieveAll
func RetrieveAllArticles(db *sql.DB) (articles []Article, err error) {
	rows, err := db.Query("SELECT id, uuid, title, summary, urls, author, category, image_link, creation_date, user_id, article_link, body_html, validated FROM articles ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Article{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.Title, &result.Summary, &result.Urls, &result.Author, &result.Category, &result.ImageLink, &result.CreationDate, &result.UserId, &result.ArticleLink, &result.BodyHtml, &result.Validated); err != nil {
			return
		}
		articles = append(articles, result)
	}
	rows.Close()
	return
}

//Update
func (article *Article) UpdateArticle(db *sql.DB) (result Article, err error) {
	stmt, err := db.Prepare("UPDATE articles SET uuid = $2, title = $3, summary = $4, urls = $5, author = $6, category = $7, image_link = $8, creation_date = $9, user_id = $10, article_link = $11, body_html = $12, validated = $13 WHERE (id = $1) RETURNING id, uuid, title, summary, urls, author, category, image_link, creation_date, user_id, article_link, body_html, validated")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( article.Id, article.Uuid, article.Title, article.Summary, article.Urls, article.Author, article.Category, article.ImageLink, article.CreationDate, article.UserId, article.ArticleLink, article.BodyHtml, article.Validated).Scan( &result.Id, &result.Uuid, &result.Title, &result.Summary, &result.Urls, &result.Author, &result.Category, &result.ImageLink, &result.CreationDate, &result.UserId, &result.ArticleLink, &result.BodyHtml, &result.Validated)
	return
}

//Delete
func (article *Article) DeleteArticle(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM articles WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.Id)
	return
}

//DeleteAll
func DeleteAllArticles(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM articles")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Product_trial =======

//Create Table
func CreateTableProduct_trials(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_trials CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table product_trials ( id integer generated always as identity primary key , uuid varchar(36) not null unique , user_id integer not null references users(id) , product_uuid varchar(36) not null references products(uuid) , start_date timestamptz , end_date timestamptz , creation_date timestamptz not null , archived boolean not null ) ; `)
	return
}

//Drop Table
func DropTableProduct_trials(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_trials CASCADE")
	return
}

//Struct
type Product_trial struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	UserId int32`xml:"UserId" json:"userid"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	StartDate sql.NullTime`xml:"StartDate" json:"startdate"`
	EndDate sql.NullTime`xml:"EndDate" json:"enddate"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`
	Archived bool`xml:"Archived" json:"archived"`

}

//Create
func (product_trial *Product_trial) CreateProduct_trial(db *sql.DB) (result Product_trial, err error) {
	stmt, err := db.Prepare("INSERT INTO product_trials ( uuid, user_id, product_uuid, start_date, end_date, creation_date, archived) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( product_trial.Uuid, product_trial.UserId, product_trial.ProductUuid, product_trial.StartDate, product_trial.EndDate, product_trial.CreationDate, product_trial.Archived).Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived)
	return
}

//Retrieve
func (product_trial *Product_trial) RetrieveProduct_trial(db *sql.DB) (result Product_trial, err error) {
	result = Product_trial{}
	err = db.QueryRow("SELECT id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived FROM product_trials WHERE (id = $1)", product_trial.Id).Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived)
	return
}

//RetrieveAll
func RetrieveAllProduct_trials(db *sql.DB) (product_trials []Product_trial, err error) {
	rows, err := db.Query("SELECT id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived FROM product_trials ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Product_trial{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived); err != nil {
			return
		}
		product_trials = append(product_trials, result)
	}
	rows.Close()
	return
}

//Update
func (product_trial *Product_trial) UpdateProduct_trial(db *sql.DB) (result Product_trial, err error) {
	stmt, err := db.Prepare("UPDATE product_trials SET uuid = $2, user_id = $3, product_uuid = $4, start_date = $5, end_date = $6, creation_date = $7, archived = $8 WHERE (id = $1) RETURNING id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( product_trial.Id, product_trial.Uuid, product_trial.UserId, product_trial.ProductUuid, product_trial.StartDate, product_trial.EndDate, product_trial.CreationDate, product_trial.Archived).Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived)
	return
}

//Delete
func (product_trial *Product_trial) DeleteProduct_trial(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_trials WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product_trial.Id)
	return
}

//DeleteAll
func DeleteAllProduct_trials(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_trials")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Hypnostat =======

//Create Table
func CreateTableHypnostats(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS hypnostats CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table hypnostats ( id integer generated always as identity primary key , user_id integer not null references users(id) , created_at timestamptz not null , source text not null , use boolean not null , hypno jsonb not null , hypno_model int not null , motion_len int not null , energy_len int not null , heartrate_len int not null , stand_len int not null , sleepsleep_len int not null , sleepinbed_len int not null , begin_bed_rel decimal not null , begin_bed decimal not null , end_bed decimal not null , hr_min decimal not null , hr_max decimal not null , hr_avg decimal not null , tst decimal not null , time_awake decimal not null , time_rem decimal not null , time_light decimal not null , time_deep decimal not null , num_awakes decimal not null , score decimal not null , score_duration decimal not null , score_efficiency decimal not null , score_continuity decimal not null , utc_offset integer not null , sleep_onset decimal not null , sleep_efficiency decimal not null ) ; `)
	return
}

//Drop Table
func DropTableHypnostats(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS hypnostats CASCADE")
	return
}

//Struct
type Hypnostat struct {
	Id int32`xml:"Id" json:"id"`
	UserId int32`xml:"UserId" json:"userid"`
	CreatedAt time.Time`xml:"CreatedAt" json:"createdat"`
	Source string`xml:"Source" json:"source"`
	Use bool`xml:"Use" json:"use"`
	Hypno string`xml:"Hypno" json:"hypno"`
	HypnoModel int32`xml:"HypnoModel" json:"hypnomodel"`
	MotionLen int32`xml:"MotionLen" json:"motionlen"`
	EnergyLen int32`xml:"EnergyLen" json:"energylen"`
	HeartrateLen int32`xml:"HeartrateLen" json:"heartratelen"`
	StandLen int32`xml:"StandLen" json:"standlen"`
	SleepsleepLen int32`xml:"SleepsleepLen" json:"sleepsleeplen"`
	SleepinbedLen int32`xml:"SleepinbedLen" json:"sleepinbedlen"`
	BeginBedRel float64`xml:"BeginBedRel" json:"beginbedrel"`
	BeginBed float64`xml:"BeginBed" json:"beginbed"`
	EndBed float64`xml:"EndBed" json:"endbed"`
	HrMin float64`xml:"HrMin" json:"hrmin"`
	HrMax float64`xml:"HrMax" json:"hrmax"`
	HrAvg float64`xml:"HrAvg" json:"hravg"`
	Tst float64`xml:"Tst" json:"tst"`
	TimeAwake float64`xml:"TimeAwake" json:"timeawake"`
	TimeRem float64`xml:"TimeRem" json:"timerem"`
	TimeLight float64`xml:"TimeLight" json:"timelight"`
	TimeDeep float64`xml:"TimeDeep" json:"timedeep"`
	NumAwakes float64`xml:"NumAwakes" json:"numawakes"`
	Score float64`xml:"Score" json:"score"`
	ScoreDuration float64`xml:"ScoreDuration" json:"scoreduration"`
	ScoreEfficiency float64`xml:"ScoreEfficiency" json:"scoreefficiency"`
	ScoreContinuity float64`xml:"ScoreContinuity" json:"scorecontinuity"`
	UtcOffset int32`xml:"UtcOffset" json:"utcoffset"`
	SleepOnset float64`xml:"SleepOnset" json:"sleeponset"`
	SleepEfficiency float64`xml:"SleepEfficiency" json:"sleepefficiency"`

}

//Create
func (hypnostat *Hypnostat) CreateHypnostat(db *sql.DB) (result Hypnostat, err error) {
	stmt, err := db.Prepare("INSERT INTO hypnostats ( user_id, created_at, source, use, hypno, hypno_model, motion_len, energy_len, heartrate_len, stand_len, sleepsleep_len, sleepinbed_len, begin_bed_rel, begin_bed, end_bed, hr_min, hr_max, hr_avg, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, score_duration, score_efficiency, score_continuity, utc_offset, sleep_onset, sleep_efficiency) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31) RETURNING id, user_id, created_at, source, use, hypno, hypno_model, motion_len, energy_len, heartrate_len, stand_len, sleepsleep_len, sleepinbed_len, begin_bed_rel, begin_bed, end_bed, hr_min, hr_max, hr_avg, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, score_duration, score_efficiency, score_continuity, utc_offset, sleep_onset, sleep_efficiency")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( hypnostat.UserId, hypnostat.CreatedAt, hypnostat.Source, hypnostat.Use, hypnostat.Hypno, hypnostat.HypnoModel, hypnostat.MotionLen, hypnostat.EnergyLen, hypnostat.HeartrateLen, hypnostat.StandLen, hypnostat.SleepsleepLen, hypnostat.SleepinbedLen, hypnostat.BeginBedRel, hypnostat.BeginBed, hypnostat.EndBed, hypnostat.HrMin, hypnostat.HrMax, hypnostat.HrAvg, hypnostat.Tst, hypnostat.TimeAwake, hypnostat.TimeRem, hypnostat.TimeLight, hypnostat.TimeDeep, hypnostat.NumAwakes, hypnostat.Score, hypnostat.ScoreDuration, hypnostat.ScoreEfficiency, hypnostat.ScoreContinuity, hypnostat.UtcOffset, hypnostat.SleepOnset, hypnostat.SleepEfficiency).Scan( &result.Id, &result.UserId, &result.CreatedAt, &result.Source, &result.Use, &result.Hypno, &result.HypnoModel, &result.MotionLen, &result.EnergyLen, &result.HeartrateLen, &result.StandLen, &result.SleepsleepLen, &result.SleepinbedLen, &result.BeginBedRel, &result.BeginBed, &result.EndBed, &result.HrMin, &result.HrMax, &result.HrAvg, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.ScoreDuration, &result.ScoreEfficiency, &result.ScoreContinuity, &result.UtcOffset, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//Retrieve
func (hypnostat *Hypnostat) RetrieveHypnostat(db *sql.DB) (result Hypnostat, err error) {
	result = Hypnostat{}
	err = db.QueryRow("SELECT id, user_id, created_at, source, use, hypno, hypno_model, motion_len, energy_len, heartrate_len, stand_len, sleepsleep_len, sleepinbed_len, begin_bed_rel, begin_bed, end_bed, hr_min, hr_max, hr_avg, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, score_duration, score_efficiency, score_continuity, utc_offset, sleep_onset, sleep_efficiency FROM hypnostats WHERE (id = $1)", hypnostat.Id).Scan( &result.Id, &result.UserId, &result.CreatedAt, &result.Source, &result.Use, &result.Hypno, &result.HypnoModel, &result.MotionLen, &result.EnergyLen, &result.HeartrateLen, &result.StandLen, &result.SleepsleepLen, &result.SleepinbedLen, &result.BeginBedRel, &result.BeginBed, &result.EndBed, &result.HrMin, &result.HrMax, &result.HrAvg, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.ScoreDuration, &result.ScoreEfficiency, &result.ScoreContinuity, &result.UtcOffset, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//RetrieveAll
func RetrieveAllHypnostats(db *sql.DB) (hypnostats []Hypnostat, err error) {
	rows, err := db.Query("SELECT id, user_id, created_at, source, use, hypno, hypno_model, motion_len, energy_len, heartrate_len, stand_len, sleepsleep_len, sleepinbed_len, begin_bed_rel, begin_bed, end_bed, hr_min, hr_max, hr_avg, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, score_duration, score_efficiency, score_continuity, utc_offset, sleep_onset, sleep_efficiency FROM hypnostats ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Hypnostat{}
		if err = rows.Scan( &result.Id, &result.UserId, &result.CreatedAt, &result.Source, &result.Use, &result.Hypno, &result.HypnoModel, &result.MotionLen, &result.EnergyLen, &result.HeartrateLen, &result.StandLen, &result.SleepsleepLen, &result.SleepinbedLen, &result.BeginBedRel, &result.BeginBed, &result.EndBed, &result.HrMin, &result.HrMax, &result.HrAvg, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.ScoreDuration, &result.ScoreEfficiency, &result.ScoreContinuity, &result.UtcOffset, &result.SleepOnset, &result.SleepEfficiency); err != nil {
			return
		}
		hypnostats = append(hypnostats, result)
	}
	rows.Close()
	return
}

//Update
func (hypnostat *Hypnostat) UpdateHypnostat(db *sql.DB) (result Hypnostat, err error) {
	stmt, err := db.Prepare("UPDATE hypnostats SET user_id = $2, created_at = $3, source = $4, use = $5, hypno = $6, hypno_model = $7, motion_len = $8, energy_len = $9, heartrate_len = $10, stand_len = $11, sleepsleep_len = $12, sleepinbed_len = $13, begin_bed_rel = $14, begin_bed = $15, end_bed = $16, hr_min = $17, hr_max = $18, hr_avg = $19, tst = $20, time_awake = $21, time_rem = $22, time_light = $23, time_deep = $24, num_awakes = $25, score = $26, score_duration = $27, score_efficiency = $28, score_continuity = $29, utc_offset = $30, sleep_onset = $31, sleep_efficiency = $32 WHERE (id = $1) RETURNING id, user_id, created_at, source, use, hypno, hypno_model, motion_len, energy_len, heartrate_len, stand_len, sleepsleep_len, sleepinbed_len, begin_bed_rel, begin_bed, end_bed, hr_min, hr_max, hr_avg, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, score_duration, score_efficiency, score_continuity, utc_offset, sleep_onset, sleep_efficiency")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( hypnostat.Id, hypnostat.UserId, hypnostat.CreatedAt, hypnostat.Source, hypnostat.Use, hypnostat.Hypno, hypnostat.HypnoModel, hypnostat.MotionLen, hypnostat.EnergyLen, hypnostat.HeartrateLen, hypnostat.StandLen, hypnostat.SleepsleepLen, hypnostat.SleepinbedLen, hypnostat.BeginBedRel, hypnostat.BeginBed, hypnostat.EndBed, hypnostat.HrMin, hypnostat.HrMax, hypnostat.HrAvg, hypnostat.Tst, hypnostat.TimeAwake, hypnostat.TimeRem, hypnostat.TimeLight, hypnostat.TimeDeep, hypnostat.NumAwakes, hypnostat.Score, hypnostat.ScoreDuration, hypnostat.ScoreEfficiency, hypnostat.ScoreContinuity, hypnostat.UtcOffset, hypnostat.SleepOnset, hypnostat.SleepEfficiency).Scan( &result.Id, &result.UserId, &result.CreatedAt, &result.Source, &result.Use, &result.Hypno, &result.HypnoModel, &result.MotionLen, &result.EnergyLen, &result.HeartrateLen, &result.StandLen, &result.SleepsleepLen, &result.SleepinbedLen, &result.BeginBedRel, &result.BeginBed, &result.EndBed, &result.HrMin, &result.HrMax, &result.HrAvg, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.ScoreDuration, &result.ScoreEfficiency, &result.ScoreContinuity, &result.UtcOffset, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//Delete
func (hypnostat *Hypnostat) DeleteHypnostat(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM hypnostats WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(hypnostat.Id)
	return
}

//DeleteAll
func DeleteAllHypnostats(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM hypnostats")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Product_comment =======

//Create Table
func CreateTableProduct_comments(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_comments CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table product_comments ( id integer generated always as identity primary key , uuid varchar(36) not null unique , user_id integer not null references users(id) , created_at timestamptz not null , last_updated timestamptz not null , comments text , rating int not null , product_uuid text not null refrences products(uuid) ) ; `)
	return
}

//Drop Table
func DropTableProduct_comments(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS product_comments CASCADE")
	return
}

//Struct
type Product_comment struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	UserId int32`xml:"UserId" json:"userid"`
	CreatedAt time.Time`xml:"CreatedAt" json:"createdat"`
	LastUpdated time.Time`xml:"LastUpdated" json:"lastupdated"`
	Comments sql.NullString`xml:"Comments" json:"comments"`
	Rating int32`xml:"Rating" json:"rating"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`

}

//Create
func (product_comment *Product_comment) CreateProduct_comment(db *sql.DB) (result Product_comment, err error) {
	stmt, err := db.Prepare("INSERT INTO product_comments ( uuid, user_id, created_at, last_updated, comments, rating, product_uuid) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, uuid, user_id, created_at, last_updated, comments, rating, product_uuid")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( product_comment.Uuid, product_comment.UserId, product_comment.CreatedAt, product_comment.LastUpdated, product_comment.Comments, product_comment.Rating, product_comment.ProductUuid).Scan( &result.Id, &result.Uuid, &result.UserId, &result.CreatedAt, &result.LastUpdated, &result.Comments, &result.Rating, &result.ProductUuid)
	return
}

//Retrieve
func (product_comment *Product_comment) RetrieveProduct_comment(db *sql.DB) (result Product_comment, err error) {
	result = Product_comment{}
	err = db.QueryRow("SELECT id, uuid, user_id, created_at, last_updated, comments, rating, product_uuid FROM product_comments WHERE (id = $1)", product_comment.Id).Scan( &result.Id, &result.Uuid, &result.UserId, &result.CreatedAt, &result.LastUpdated, &result.Comments, &result.Rating, &result.ProductUuid)
	return
}

//RetrieveAll
func RetrieveAllProduct_comments(db *sql.DB) (product_comments []Product_comment, err error) {
	rows, err := db.Query("SELECT id, uuid, user_id, created_at, last_updated, comments, rating, product_uuid FROM product_comments ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Product_comment{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.UserId, &result.CreatedAt, &result.LastUpdated, &result.Comments, &result.Rating, &result.ProductUuid); err != nil {
			return
		}
		product_comments = append(product_comments, result)
	}
	rows.Close()
	return
}

//Update
func (product_comment *Product_comment) UpdateProduct_comment(db *sql.DB) (result Product_comment, err error) {
	stmt, err := db.Prepare("UPDATE product_comments SET uuid = $2, user_id = $3, created_at = $4, last_updated = $5, comments = $6, rating = $7, product_uuid = $8 WHERE (id = $1) RETURNING id, uuid, user_id, created_at, last_updated, comments, rating, product_uuid")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( product_comment.Id, product_comment.Uuid, product_comment.UserId, product_comment.CreatedAt, product_comment.LastUpdated, product_comment.Comments, product_comment.Rating, product_comment.ProductUuid).Scan( &result.Id, &result.Uuid, &result.UserId, &result.CreatedAt, &result.LastUpdated, &result.Comments, &result.Rating, &result.ProductUuid)
	return
}

//Delete
func (product_comment *Product_comment) DeleteProduct_comment(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_comments WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product_comment.Id)
	return
}

//DeleteAll
func DeleteAllProduct_comments(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM product_comments")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Community_product_trial_stat =======

//Create Table
func CreateTableCommunity_product_trial_stats(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS community_product_trial_stats CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table community_product_trial_stats ( id integer generated always as identity primary key , created_at timestamptz not null , product_uuid varchar(36) not null references products(uuid) , source text not null , days int not null , user_count int not null , tst decimal , time_awake decimal , time_rem decimal , time_light decimal , time_deep decimal , num_awakes decimal , score decimal , sleep_onset decimal , sleep_efficiency decimal ) ; `)
	return
}

//Drop Table
func DropTableCommunity_product_trial_stats(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS community_product_trial_stats CASCADE")
	return
}

//Struct
type Community_product_trial_stat struct {
	Id int32`xml:"Id" json:"id"`
	CreatedAt time.Time`xml:"CreatedAt" json:"createdat"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	Source string`xml:"Source" json:"source"`
	Days int32`xml:"Days" json:"days"`
	UserCount int32`xml:"UserCount" json:"usercount"`
	Tst sql.NullFloat64`xml:"Tst" json:"tst"`
	TimeAwake sql.NullFloat64`xml:"TimeAwake" json:"timeawake"`
	TimeRem sql.NullFloat64`xml:"TimeRem" json:"timerem"`
	TimeLight sql.NullFloat64`xml:"TimeLight" json:"timelight"`
	TimeDeep sql.NullFloat64`xml:"TimeDeep" json:"timedeep"`
	NumAwakes sql.NullFloat64`xml:"NumAwakes" json:"numawakes"`
	Score sql.NullFloat64`xml:"Score" json:"score"`
	SleepOnset sql.NullFloat64`xml:"SleepOnset" json:"sleeponset"`
	SleepEfficiency sql.NullFloat64`xml:"SleepEfficiency" json:"sleepefficiency"`

}

//Create
func (community_product_trial_stat *Community_product_trial_stat) CreateCommunity_product_trial_stat(db *sql.DB) (result Community_product_trial_stat, err error) {
	stmt, err := db.Prepare("INSERT INTO community_product_trial_stats ( created_at, product_uuid, source, days, user_count, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id, created_at, product_uuid, source, days, user_count, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( community_product_trial_stat.CreatedAt, community_product_trial_stat.ProductUuid, community_product_trial_stat.Source, community_product_trial_stat.Days, community_product_trial_stat.UserCount, community_product_trial_stat.Tst, community_product_trial_stat.TimeAwake, community_product_trial_stat.TimeRem, community_product_trial_stat.TimeLight, community_product_trial_stat.TimeDeep, community_product_trial_stat.NumAwakes, community_product_trial_stat.Score, community_product_trial_stat.SleepOnset, community_product_trial_stat.SleepEfficiency).Scan( &result.Id, &result.CreatedAt, &result.ProductUuid, &result.Source, &result.Days, &result.UserCount, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//Retrieve
func (community_product_trial_stat *Community_product_trial_stat) RetrieveCommunity_product_trial_stat(db *sql.DB) (result Community_product_trial_stat, err error) {
	result = Community_product_trial_stat{}
	err = db.QueryRow("SELECT id, created_at, product_uuid, source, days, user_count, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency FROM community_product_trial_stats WHERE (id = $1)", community_product_trial_stat.Id).Scan( &result.Id, &result.CreatedAt, &result.ProductUuid, &result.Source, &result.Days, &result.UserCount, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//RetrieveAll
func RetrieveAllCommunity_product_trial_stats(db *sql.DB) (community_product_trial_stats []Community_product_trial_stat, err error) {
	rows, err := db.Query("SELECT id, created_at, product_uuid, source, days, user_count, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency FROM community_product_trial_stats ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Community_product_trial_stat{}
		if err = rows.Scan( &result.Id, &result.CreatedAt, &result.ProductUuid, &result.Source, &result.Days, &result.UserCount, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency); err != nil {
			return
		}
		community_product_trial_stats = append(community_product_trial_stats, result)
	}
	rows.Close()
	return
}

//Update
func (community_product_trial_stat *Community_product_trial_stat) UpdateCommunity_product_trial_stat(db *sql.DB) (result Community_product_trial_stat, err error) {
	stmt, err := db.Prepare("UPDATE community_product_trial_stats SET created_at = $2, product_uuid = $3, source = $4, days = $5, user_count = $6, tst = $7, time_awake = $8, time_rem = $9, time_light = $10, time_deep = $11, num_awakes = $12, score = $13, sleep_onset = $14, sleep_efficiency = $15 WHERE (id = $1) RETURNING id, created_at, product_uuid, source, days, user_count, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( community_product_trial_stat.Id, community_product_trial_stat.CreatedAt, community_product_trial_stat.ProductUuid, community_product_trial_stat.Source, community_product_trial_stat.Days, community_product_trial_stat.UserCount, community_product_trial_stat.Tst, community_product_trial_stat.TimeAwake, community_product_trial_stat.TimeRem, community_product_trial_stat.TimeLight, community_product_trial_stat.TimeDeep, community_product_trial_stat.NumAwakes, community_product_trial_stat.Score, community_product_trial_stat.SleepOnset, community_product_trial_stat.SleepEfficiency).Scan( &result.Id, &result.CreatedAt, &result.ProductUuid, &result.Source, &result.Days, &result.UserCount, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//Delete
func (community_product_trial_stat *Community_product_trial_stat) DeleteCommunity_product_trial_stat(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM community_product_trial_stats WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(community_product_trial_stat.Id)
	return
}

//DeleteAll
func DeleteAllCommunity_product_trial_stats(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM community_product_trial_stats")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Habit =======

//Create Table
func CreateTableHabits(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS habits CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table habits ( id integer generated always as identity primary key , uuid varchar(36) not null unique , user_id integer not null references users(id) , product_uuid varchar(36) not null references products(uuid) , start_date timestamptz , end_date timestamptz , creation_date timestamptz not null , archived boolean not null ) ; `)
	return
}

//Drop Table
func DropTableHabits(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS habits CASCADE")
	return
}

//Struct
type Habit struct {
	Id int32`xml:"Id" json:"id"`
	Uuid string`xml:"Uuid" json:"uuid"`
	UserId int32`xml:"UserId" json:"userid"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	StartDate sql.NullTime`xml:"StartDate" json:"startdate"`
	EndDate sql.NullTime`xml:"EndDate" json:"enddate"`
	CreationDate time.Time`xml:"CreationDate" json:"creationdate"`
	Archived bool`xml:"Archived" json:"archived"`

}

//Create
func (habit *Habit) CreateHabit(db *sql.DB) (result Habit, err error) {
	stmt, err := db.Prepare("INSERT INTO habits ( uuid, user_id, product_uuid, start_date, end_date, creation_date, archived) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( habit.Uuid, habit.UserId, habit.ProductUuid, habit.StartDate, habit.EndDate, habit.CreationDate, habit.Archived).Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived)
	return
}

//Retrieve
func (habit *Habit) RetrieveHabit(db *sql.DB) (result Habit, err error) {
	result = Habit{}
	err = db.QueryRow("SELECT id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived FROM habits WHERE (id = $1)", habit.Id).Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived)
	return
}

//RetrieveAll
func RetrieveAllHabits(db *sql.DB) (habits []Habit, err error) {
	rows, err := db.Query("SELECT id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived FROM habits ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Habit{}
		if err = rows.Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived); err != nil {
			return
		}
		habits = append(habits, result)
	}
	rows.Close()
	return
}

//Update
func (habit *Habit) UpdateHabit(db *sql.DB) (result Habit, err error) {
	stmt, err := db.Prepare("UPDATE habits SET uuid = $2, user_id = $3, product_uuid = $4, start_date = $5, end_date = $6, creation_date = $7, archived = $8 WHERE (id = $1) RETURNING id, uuid, user_id, product_uuid, start_date, end_date, creation_date, archived")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( habit.Id, habit.Uuid, habit.UserId, habit.ProductUuid, habit.StartDate, habit.EndDate, habit.CreationDate, habit.Archived).Scan( &result.Id, &result.Uuid, &result.UserId, &result.ProductUuid, &result.StartDate, &result.EndDate, &result.CreationDate, &result.Archived)
	return
}

//Delete
func (habit *Habit) DeleteHabit(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM habits WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(habit.Id)
	return
}

//DeleteAll
func DeleteAllHabits(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM habits")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}


// ======= Community_product_result =======

//Create Table
func CreateTableCommunity_product_results(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS community_product_results CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table community_product_results ( id integer generated always as identity primary key , user_count integer not null , count integer not null , product_uuid varchar(36) not null references products(uuid) , tst decimal , time_awake decimal , time_rem decimal , time_light decimal , time_deep decimal , num_awakes decimal , score decimal , sleep_onset decimal , sleep_efficiency decimal ) ; `)
	return
}

//Drop Table
func DropTableCommunity_product_results(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS community_product_results CASCADE")
	return
}

//Struct
type Community_product_result struct {
	Id int32`xml:"Id" json:"id"`
	UserCount int32`xml:"UserCount" json:"usercount"`
	Count int32`xml:"Count" json:"count"`
	ProductUuid string`xml:"ProductUuid" json:"productuuid"`
	Tst sql.NullFloat64`xml:"Tst" json:"tst"`
	TimeAwake sql.NullFloat64`xml:"TimeAwake" json:"timeawake"`
	TimeRem sql.NullFloat64`xml:"TimeRem" json:"timerem"`
	TimeLight sql.NullFloat64`xml:"TimeLight" json:"timelight"`
	TimeDeep sql.NullFloat64`xml:"TimeDeep" json:"timedeep"`
	NumAwakes sql.NullFloat64`xml:"NumAwakes" json:"numawakes"`
	Score sql.NullFloat64`xml:"Score" json:"score"`
	SleepOnset sql.NullFloat64`xml:"SleepOnset" json:"sleeponset"`
	SleepEfficiency sql.NullFloat64`xml:"SleepEfficiency" json:"sleepefficiency"`

}

//Create
func (community_product_result *Community_product_result) CreateCommunity_product_result(db *sql.DB) (result Community_product_result, err error) {
	stmt, err := db.Prepare("INSERT INTO community_product_results ( user_count, count, product_uuid, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id, user_count, count, product_uuid, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow( community_product_result.UserCount, community_product_result.Count, community_product_result.ProductUuid, community_product_result.Tst, community_product_result.TimeAwake, community_product_result.TimeRem, community_product_result.TimeLight, community_product_result.TimeDeep, community_product_result.NumAwakes, community_product_result.Score, community_product_result.SleepOnset, community_product_result.SleepEfficiency).Scan( &result.Id, &result.UserCount, &result.Count, &result.ProductUuid, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//Retrieve
func (community_product_result *Community_product_result) RetrieveCommunity_product_result(db *sql.DB) (result Community_product_result, err error) {
	result = Community_product_result{}
	err = db.QueryRow("SELECT id, user_count, count, product_uuid, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency FROM community_product_results WHERE (id = $1)", community_product_result.Id).Scan( &result.Id, &result.UserCount, &result.Count, &result.ProductUuid, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//RetrieveAll
func RetrieveAllCommunity_product_results(db *sql.DB) (community_product_results []Community_product_result, err error) {
	rows, err := db.Query("SELECT id, user_count, count, product_uuid, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency FROM community_product_results ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Community_product_result{}
		if err = rows.Scan( &result.Id, &result.UserCount, &result.Count, &result.ProductUuid, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency); err != nil {
			return
		}
		community_product_results = append(community_product_results, result)
	}
	rows.Close()
	return
}

//Update
func (community_product_result *Community_product_result) UpdateCommunity_product_result(db *sql.DB) (result Community_product_result, err error) {
	stmt, err := db.Prepare("UPDATE community_product_results SET user_count = $2, count = $3, product_uuid = $4, tst = $5, time_awake = $6, time_rem = $7, time_light = $8, time_deep = $9, num_awakes = $10, score = $11, sleep_onset = $12, sleep_efficiency = $13 WHERE (id = $1) RETURNING id, user_count, count, product_uuid, tst, time_awake, time_rem, time_light, time_deep, num_awakes, score, sleep_onset, sleep_efficiency")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow( community_product_result.Id, community_product_result.UserCount, community_product_result.Count, community_product_result.ProductUuid, community_product_result.Tst, community_product_result.TimeAwake, community_product_result.TimeRem, community_product_result.TimeLight, community_product_result.TimeDeep, community_product_result.NumAwakes, community_product_result.Score, community_product_result.SleepOnset, community_product_result.SleepEfficiency).Scan( &result.Id, &result.UserCount, &result.Count, &result.ProductUuid, &result.Tst, &result.TimeAwake, &result.TimeRem, &result.TimeLight, &result.TimeDeep, &result.NumAwakes, &result.Score, &result.SleepOnset, &result.SleepEfficiency)
	return
}

//Delete
func (community_product_result *Community_product_result) DeleteCommunity_product_result(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM community_product_results WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(community_product_result.Id)
	return
}

//DeleteAll
func DeleteAllCommunity_product_results(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM community_product_results")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}

