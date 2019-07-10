package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"encoding/json"
	"database/sql"
)

// Model
type Pegawai struct {
	
	Kode string `json:"kdpeg"`
	Nama string `json:"nama"`
	Jabatan string `json:"jabatan"`
	Divisi string `json:"divisi"`
}

func getPegawai() []Pegawai {
	rows, err := db.Query("SELECT kdpeg, nama, jabatan, divisi FROM tugas_cockroach.pegawai")
	if err != nil { fmt.Println(err) }
	defer rows.Close()

	x := make([]Pegawai, 0)
	for rows.Next() {
		var kdpeg, nama, jabatan, divisi string
		if err := rows.Scan(&kdpeg, &nama, &jabatan, &divisi); err != nil {
			fmt.Println(err)
		}
		x = append(x, Pegawai{kdpeg, nama, jabatan, divisi})
	}
	return x
}

// Controller
func all_pegawai(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/pegawai.html"))
	data := struct {
		Title string
	}{
		Title: "Data Pegawai",
	}
	tmpl.Execute(w, data)
}

func get_pegawai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	datajson, _ := json.Marshal(getUser())
	fmt.Fprintf(w, string(datajson))
}

func get_pegawai_detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	vars := mux.Vars(r)
	if len(vars["kdpeg"]) == 0 {
		// not valid
		fmt.Fprintf(w, "{\"success\": false, \"message\": \"ID not defined\"}")
	} else {
		var id, name, email string

		rows := db.QueryRow("SELECT id, name, email FROM tugas_cockroach.user where id = $1", string(vars["id"]))
		err := rows.Scan(&id, &name, &email)

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Fprintf(w, "{\"success\": false, \"message\": \"Data not found\"}")
			} else {
				fmt.Fprintf(w, "{\"success\": false, \"message\": \"Query Error\"}")
			}
		} else {
			datajson, _ := json.Marshal(struct {
								Success bool `json:"success"`
								Data User `json:"data"`
							}{Success: true, Data: User{id, name, email} })
			fmt.Fprintf(w, string(datajson))
		}

	}
}

func add_user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	is_success := true 
	if len(r.FormValue("name")) == 0 {
		is_success = false
	} else if len(r.FormValue("email")) == 0 {
		is_success = false
	}

	if is_success {
		if _, err := db.Exec("INSERT INTO tugas_cockroach.user(name, email) VALUES($1, $2)", 
			r.FormValue("name"), r.FormValue("email"));
		err != nil {
			fmt.Println(err)
			is_success = false
		}
	}

	datajson, _ := json.Marshal(struct {
					Success bool `json:"success"`
					Data []User `json:"data"`
				}{is_success, getUser()})
	fmt.Fprintf(w, string(datajson))
}

func update_user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	is_success := true 
	vars := mux.Vars(r)
	if len(vars["id"]) == 0 {
		is_success = false
	} else if len(r.FormValue("name")) == 0 {
		is_success = false
	} else if len(r.FormValue("email")) == 0 {
		is_success = false
	}

	if is_success {
		if _, err := db.Exec("UPDATE tugas_cockroach.user SET name = $1, email = $2 WHERE id = $3", 
			r.FormValue("name"), r.FormValue("email"), string(vars["id"]));
		err != nil {
			is_success = false
		}
	}

	datajson, _ := json.Marshal(struct {
					Success bool `json:"success"`
					Data []User `json:"data"`
				}{is_success, getUser()})
	fmt.Fprintf(w, string(datajson))
}

func delete_user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	is_success := true 
	vars := mux.Vars(r)

	if len(vars["id"]) == 0 {
		is_success = false
	}

	if is_success {
		if _, err := db.Exec("DELETE FROM tugas_cockroach.user WHERE id = $1", string(vars["id"]));
		err != nil {
			is_success = false
		}
	}

	datajson, _ := json.Marshal(struct {
					Success bool `json:"success"`
				}{is_success})
	fmt.Fprintf(w, string(datajson))
}
