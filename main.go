package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	Email  string `json:"email"`
	Alamat string `json:"alamat"`
}

func main() {
	// Koneksi ke database MySQL (sesuaikan dengan konfigurasi Anda)
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/latihan5_go_web") // Sesuaikan jika diperlukan
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Memastikan koneksi berhasil
	if err := db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("Database connected!")
	}

	router := gin.Default()

	// Endpoint untuk menyimpan data
	router.POST("/tambahpengguna", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Menyimpan data ke database (tanpa memasukkan nilai id)
		result, err := db.Exec("INSERT INTO pengguna (nama, email, alamat) VALUES (?, ?, ?)", person.Nama, person.Email, person.Alamat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Mendapatkan last insert ID (opsional)
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Menambahkan ID ke dalam struct Person
		person.Id = int(lastInsertID)

		c.JSON(http.StatusCreated, person)
	})

	// Endpoint untuk menampilkan daftar pengguna
	router.GET("/daftarpengguna", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, nama, email, alamat FROM pengguna")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var people []Person
		for rows.Next() {
			var person Person
			if err := rows.Scan(&person.Id, &person.Nama, &person.Email, &person.Alamat); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			people = append(people, person)
		}

		c.JSON(http.StatusOK, people)
	})

	router.Run(":8080")
}
