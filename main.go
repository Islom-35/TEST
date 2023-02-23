package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Hospital struct {
	Name      string
	Staffs    []*Staff
	Patients  []*Patient
	Addresses []*Address
}
type HospitalResponse struct {
	Id   int64
	Name string
}
type Address struct {
	HospitalId int64
	Region     string
	Street     string
}
type AddressRes struct{
	Id int64
	Region string
}

type Staff struct {
	HospitalID  int64
	FullName    string
	PhoneNumber string
}

type StaffRes struct {
	ID       int64
	FullName string
}
type Patient struct {
	HospitalId  int64
	FullName    string
	PatientInfo string
	PhoneNumber string
}
type PatientRes struct {
	ID       int64
	FullName string
}

const (
	PostgresHost     = "localhost"
	PostgresPort     = 5432
	PostgresUser     = "postgres"
	PostgresPassword = "1234"
	PostgresDatabase = "hospital"
)

func check(err error) {
	panic(err)
}
func main() {
	connDB := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDatabase,
	)
	db, err := sql.Open("postgres", connDB)
	if err != nil {
		check(err)
	}
	defer db.Close()
	// CreateInfo(db)
	// SelectInfo(db)
	//Update(db)
	// delete(2, db)
	UpdateInfo(db)

}

func CreateInfo(db *sql.DB) {
	info := &Hospital{
		Name: "Medical",
		Staffs: []*Staff{
			{
				FullName:    "Bobir Davlatov",
				PhoneNumber: "+998919711724",
			},
			{
				FullName:    "Komil Sattorov",
				PhoneNumber: "+998912934345",
			},
		},
		Patients: []*Patient{
			{
				FullName:    "Lobarxon",
				PhoneNumber: "+99899434565",
				PatientInfo: "Bu bemor hozirda o'rta",
			},
			{
				FullName:    "Gavhar Ismoilova",
				PhoneNumber: "+99899456565",
				PatientInfo: "Bu bemor hozirda yaxshi",
			},
		},
		Addresses: []*Address{
			{
				Region: "Chilonzor",
				Street: "Qatortol 68",
			},
			{
				Region: "Yunusobod",
				Street: "Amir Temur 78",
			},
		},
	}
	tr, err := db.Begin()
	if err != nil {
		tr.Rollback()
		log.Println("Error begin trasaction", err)
	}
	query1 := `
	INSERT INTO
	Hospital(name)
	Values
	($1)
	RETURNING
	id,name`
	var response HospitalResponse
	err = tr.QueryRow(query1, info.Name).Scan(&response.Id, &response.Name)
	if err != nil {
		check(err)
	}
	fmt.Println(response)

	res1 := StaffRes{}
	for _, staff := range info.Staffs {
		query2 := `
		insert into
		staff(hospital_id, full_name,phone_number)
		values
		($1,$2,$3)
		Returning
		hospital_id, full_name`
		err = tr.QueryRow(query2, response.Id, staff.FullName, staff.PhoneNumber).Scan(&res1.ID, &res1.FullName)
		if err != nil {
			tr.Rollback()
			fmt.Println("Error while inserting staff info", err)
		}
	}
	fmt.Println("Staff info", res1)
	res2 := PatientRes{}
	for _, patient := range info.Patients {
		query3 := `
		insert into
		patients(hospital_id,full_name,patient_info, phone_number)
		VALUES
		($1,$2,$3,$4)
		returning
		id,
		full_name
		`
		err = tr.QueryRow(query3, response.Id, patient.FullName, patient.PatientInfo, patient.PhoneNumber).
			Scan(&res2.ID, &res2.FullName)
		if err != nil {
			tr.Rollback()
			fmt.Println("Error while inserting patients info:", err)
		}
	}
	fmt.Println("Patient info", res2)
	res3 := AddressRes{}
	for _, address := range info.Addresses {
		query4 := `
		insert into
		addresses(hospital_id,regional, street)
		values
		($1,$2,$3)
		RETURNING hospital_id, regional`
		err = tr.QueryRow(query4, response.Id, address.Region, address.Street).
			Scan(&res3.Id, &res3.Region)
		if err != nil {
			tr.Rollback()
			fmt.Println("Error while inserting address info:", err)
		}
		fmt.Println("Address info : ", res3)
	}

}

func SelectInfo(db *sql.DB) {
	var id int
	print("select id: ")
	fmt.Scan(&id)
	var respon HospitalResponse
	err := db.QueryRow(`select id, name from hospital where id=$1`, id).Scan(&respon.Id, &respon.Name)
	if err != nil {
		check(err)
	}
	fmt.Println(respon)

}
func Update(db *sql.DB) {
	var id int
	print("enter id: ")
	fmt.Scan(&id)
	addres := Address{}
	addres.Region = "olmazor"
	result, err := db.Exec("update addresses set regional=$1 where id=$2", addres.Region, id)
	if err != nil {
		check(err)
	}
	fmt.Println(addres)
	fmt.Println(result)

}

func delete(ID int, db *sql.DB) {
	var del int
	print("which from table do you want delete?\n ")
	print("1 addresses, 2 patients, 3staff\n")
	fmt.Scan(&del)
	if del == 1 {
		_, err := db.Exec("delete from addresses where id=$1", ID)
		if err != nil {
			check(err)
		}
	} else if del == 2 {
		_, err := db.Exec("delete from patients where id=$1", ID)
		if err != nil {
			check(err)
		}
	} else if del == 3 {
		_, err := db.Exec("delete from addresses", ID)
		if err != nil {
			check(err)
		}
	} else {
		print("this table does not exist")

	}
	print("ok")

}

func UpdateInfo(db *sql.DB){
	var id int64
	print("enter id\n")
	fmt.Scan((&id))
	hospital:=&HospitalResponse{Id: id, Name: "Update hospital name"}
	tx,err :=db.Begin()
	if err !=nil{
		fmt.Println("Error begin transaction: ",err)
		return
	}
	res:=HospitalResponse{}
	err =tx.QueryRow(`Update hospital set name = $1 where id = $2 returning id, name`,hospital.Name,hospital.Id).Scan(&res.Id,&res.Name)
	if err !=nil{
		tx.Rollback()
		log.Println("Error update hospital info",err)
	}
	err = tx.Commit()
	if err !=nil{
		fmt.Println("Error commin tr:", err)
	}
}
