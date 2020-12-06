//https://code-database.com/knowledges/87 htmlをgoで立てたサーバで展開する方法
//https://qiita.com/TakahiRoyte/items/949f4e88caecb02119aa#:~:text=REST(REpresentational%20State%20Transfer)%E3%81%AF,%E3%81%AE%E9%80%81%E5%8F%97%E4%BF%A1%E3%82%92%E8%A1%8C%E3%81%84%E3%81%BE%E3%81%99%E3%80%82
//↑RESTについて
package main

import (
	"database/sql"
	"encoding/json"
	"fmt" //標準入力など(デバッグ用なので最終的にはいらない...?)
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"      //サーバを立てるた目に必要
	"text/template" //htmlファイルを展開するために必要
	"time"
)

type Unit struct {
	UnitID     int        `json:"unit_id"`
	UnitType   string     `json:"unit_type"`
	Purpose    string     `json:"purpose"`
	BmsVersion string     `json:"bms_version"`
	UnitState  string     `json:"unit_state"`
	Time       *time.Time `json:"time"`
	IsWorking  string     `json:"is_working"`
	SOC        int        `json:"soc"`
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("resttest.html") //htmlからテンプレートを作成
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil) //テンプレートを実行(ブラウザに表示)
	if err != nil {
		panic(err)
	}
}
func GetUnitsInfo(w http.ResponseWriter, r *http.Request) {//本坊くんstyle
	db, err := sql.Open("mysql", "test_user:test_pass@tcp(10.0.1.229:3306)/test_db?parseTime=True")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	results, err := db.Query("SELECT * FROM units")
	if err != nil {
		panic(err.Error())
	}

	//DBからの返答unitの全てをunits([]Units型)に格納
	var units []Unit
	for results.Next() {
		var unit Unit
		err = results.Scan(&unit.UnitID, &unit.UnitType, &unit.Purpose, &unit.BmsVersion, &unit.UnitState, &unit.Time, &unit.IsWorking, &unit.SOC)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(unit)
		units = append(units, unit)
	}

	fmt.Println(units)
	responseBody, err := json.Marshal(units)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
func fetchAllUnits(w http.ResponseWriter, r *http.Request) {//kitao style
	var units []Unit
	GetAllUnits(&units)
	responseBody, err := json.Marshal(units)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func fetchSingleUnit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["unit_id"]

    var unit Unit
    // modelの呼び出し
    GetSingleUnit(&unit, id)
    responseBody, err := json.Marshal(unit)
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(responseBody)
}

func main() {

	StartWebServer()

	//本坊くん　gorillaを使用しないパタン（これでもできると思われる）
	//http.HandleFunc("/", MainHandler)
	//http.HandleFunc("/api/v1/units/", GetUnitsInfo)
	//fmt.Println("Server Started Port 443")
	//log.Fatal(http.ListenAndServe(":443", nil))
	// http://18.180.144.98:443/
}

func StartWebServer() error {
	fmt.Println("Rest API with Mux Routers")
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc({ エンドポイント }, { レスポンス関数 }).Methods({ リクエストメソッド（複数可能） })
	router.HandleFunc("/", MainHandler)
	//router.HandleFunc("/api/v1/units/", GetUnitsInfo).Methods("GET")
	router.HandleFunc("/api/v1/units/", fetchAllUnits).Methods("GET")
	router.HandleFunc("/api/v1/unit/{id}/", fetchSingleUnit).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", 443), router)
}

func init() {
	var err error
	dbConnectInfo := fmt.Sprintf(
		//`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		`test_user:test_pass@tcp(10.0.1.229:3306)/test_db?parseTime=True`,
		//config.Config.DbUserName,
		//config.Config.DbUserPassword,
		//config.Config.DbHost,
		//config.Config.DbPort,
		//config.Config.DbName,
	)


	// configから読み込んだ情報を元に、データベースに接続
	//Db, err = gorm.Open(config.Config.DbDriverName, dbConnectInfo)
	Db, err = gorm.Open("mysql", dbConnectInfo)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Successfully connect database..")
	}

	/*
		// 接続したデータベースにunitテーブルを作成：多分必要ない
		Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&Unit{})
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Successfully created table..")
		}
	*/
}

var Db *gorm.DB

func GetAllUnits(units *[]Unit) {
	Db.Find(&units)
}

func GetSingleUnit(unit *Unit, key string) {
	Db.First(&unit, key)
}
