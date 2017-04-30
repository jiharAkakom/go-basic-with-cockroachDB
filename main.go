package main

import "database/sql"
import "net/http"
import "fmt"
import _ "github.com/lib/pq"
import "html/template"
import "os"
import "encoding/json"

type mhs struct {
	Number int
	Nim    string
	Nama   string
}

type jurusan struct {
	Ka []mhs
}

func connect() *sql.DB {
	var db, err = sql.Open("postgres", "postgresql://root@localhost:26257/KA?sslmode=disable") //default cockroach setting
	err = db.Ping()
	if err != nil {
		fmt.Println("database tidak bisa dihubungi")
		os.Exit(0)

	}
	return db

}

func tampil_data(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()

	type maba []mhs
	var mahasiswa maba
	var num = 1

	rows, _ := db.Query("select nim,nama from mhs")

	for rows.Next() {
		var nim, nama string
		rows.Scan(&nim, &nama)
		data := mhs{num, nim, nama}
		mahasiswa = append(mahasiswa, data)
		num++
	}
	halaman, _ := template.ParseFiles("index.html")
	halaman.Execute(res, mahasiswa)

}

func tampil_data_json(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()

	type maba []mhs
	var mahasiswa maba
	var num = 1

	rows, _ := db.Query("select nim,nama from mhs")

	for rows.Next() {
		var nim, nama string
		rows.Scan(&nim, &nama)
		data := mhs{num, nim, nama}
		mahasiswa = append(mahasiswa, data)
		num++
	}
	json_mhs := jurusan{mahasiswa}
	json.NewEncoder(res).Encode(json_mhs)

}

func input_data(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()

	nim := req.FormValue("nim")
	nama := req.FormValue("nama")

	db.Exec("insert into mhs values ($1,$2)", nim, nama)

	http.Redirect(res, req, "/", 301)
}

func edit_data(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()
	nim_sebelumnya := req.FormValue("nim_sebelumnya")

	nim := req.FormValue("nim")
	nama := req.FormValue("nama")

	if nim != "" {
		db.Exec("update mhs set nim = $1 where nim = $2", nim, nim_sebelumnya)
	}

	if nama != "" {
		db.Exec("update mhs set nama = $1 where nim = $2", nama, nim_sebelumnya)
	}

	http.Redirect(res, req, "/", 301)
}

func hapus_data(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()

	nim := req.FormValue("hapus")

	db.Exec("delete from mhs where nim = $1", nim)

	http.Redirect(res, req, "/", 301)
}

func main() {
	http.HandleFunc("/", tampil_data)
	http.HandleFunc("/json", tampil_data_json)
	http.HandleFunc("/input_data", input_data)
	http.HandleFunc("/ubah_data", edit_data)
	http.HandleFunc("/hapus_data", hapus_data)
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))
	fmt.Println("running on port :80....")
	http.ListenAndServe(":80", nil)
}
