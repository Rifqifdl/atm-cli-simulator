package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var currentUser string
var currentId int

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/dibimbing_atm"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal membuka koneksi: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal tersambung ke database: ", err)
	}

	fmt.Println("\n===== Selamat Datang di ATM XYZ =====")
	fmt.Println("1. Registrasi")
	fmt.Println("2. Login")

	var menu int
	fmt.Print("\nPilih Menu: ")
	fmt.Scan(&menu)

	switch menu {
	case 1:
		registrasi(db)
	case 2:
		login(db)
	default:
		fmt.Println("Menu tidak tersedia.")
	}

}

func registrasi(db *sql.DB) {
	fmt.Println("\n===== REGISTRASI AKUN =====")

	//Input data user
	var name string
	var pin int
	fmt.Print("Masukan Nama: ")
	fmt.Scan(&name)
	fmt.Print("Masukan PIN (6 digit): ")
	fmt.Scan(&pin)

	//Insert Database
	query := "INSERT INTO accounts (name, pin) VALUES (?, ?)"
	result, err := db.Exec(query, name, pin)
	if err != nil {
		log.Fatal("Registrasi Gagal: ", err)
		main()
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Gagal mengambil ID:", err)
		main()
	}

	fmt.Println("\nRegistrasi berhasil! Nomor Akun: ", lastID, "Nama: ", name)
	main()
}

func login(db *sql.DB) {
	fmt.Println("\n===== LOGIN =====")

	var id, pin int

	fmt.Print("Masukan Nomor Akun: ")
	fmt.Scan(&id)

	fmt.Print("Masukan PIN: ")
	fmt.Scan(&pin)

	// Cek ke database
	var storedPin int
	var name string
	query := "SELECT pin, name FROM accounts WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&storedPin, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Username tidak ditemukan.")
			main()
			return
		}
		fmt.Println("Terjadi kesalahan saat login: ", err)
		main()
		return
	}

	// Cocokkan PIN
	if storedPin != pin {
		fmt.Println("PIN salah.")
		main()
		return
	}

	// Login sukses
	currentUser = name
	currentId = id
	fmt.Println("\nLogin berhasil!")
	dashboard(db)
}

func dashboard(db *sql.DB) {
	fmt.Println("\n===== Selamat datang,", currentUser, " =====")
	fmt.Println("1. Cek Saldo")
	fmt.Println("2. Setor Tunai")
	fmt.Println("3. Tarik Tunai")
	fmt.Println("4. Transfer")
	fmt.Println("5. Riwayat Transaksi")
	fmt.Println("6. Logout")

	var menu int
	fmt.Print("\nPilih Menu: ")
	fmt.Scan(&menu)

	switch menu {
	case 1:
		balance(db)
	case 2:
		deposit(db)
	case 3:
		withdraw(db)
	case 4:
		transfer(db)
	case 5:
		history(db)
	case 6:
		main()
	default:
		fmt.Println("Menu tidak tersedia.")
	}
}

func balance(db *sql.DB) {
	//get data saldo
	var saldo int
	query := "SELECT balance FROM accounts WHERE id = ?"
	err := db.QueryRow(query, currentId).Scan(&saldo)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Akun tidak ditemukan.")
		} else {
			fmt.Println("Terjadi kesalahan saat mengecek saldo: ", err)
		}
		dashboard(db)
		return
	}
	fmt.Print("\nSaldo anda saat ini: Rp", saldo, "\n")
	dashboard(db)
}

func deposit(db *sql.DB) {
	fmt.Println("\n===== Setor Tunai =====")

	var nominal int
	fmt.Print("Masukkan jumlah setor tunai: Rp")
	fmt.Scan(&nominal)

	//validasi input
	if nominal <= 0 {
		fmt.Println("Nominal tidak valid!")
		deposit(db)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Gagal memulai transaksi:", err)
		dashboard(db)
		return
	}

	//Update saldo
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", nominal, currentId)
	if err != nil {
		tx.Rollback()
		fmt.Println("Update saldo gagal:", err)
		dashboard(db)
		return
	}

	//Insert tabel transactions
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount) VALUES (?, 'deposit', ?)", currentId, nominal)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal menyimpan histori transaksi:", err)
		dashboard(db)
		return
	}

	//Commit transaksi
	err = tx.Commit()
	if err != nil {
		fmt.Println("Gagal commit transaksi:", err)
		dashboard(db)
		return
	}

	fmt.Println("Setor tunai berhasil! Saldo Anda bertambah Rp", nominal)
	dashboard(db)
}

func withdraw(db *sql.DB) {
	fmt.Println("\n===== Tarik Tunai =====")

	var nominal, saldo int
	fmt.Print("Masukkan jumlah penarikan: Rp")
	fmt.Scan(&nominal)

	//validasi input
	if nominal <= 0 {
		fmt.Println("Nominal tidak valid!")
		withdraw(db)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Gagal memulai transaksi:", err)
		dashboard(db)
		return
	}

	//cek saldo
	query := "SELECT balance FROM accounts WHERE id = ?"
	err = db.QueryRow(query, currentId).Scan(&saldo)
	if err != nil {
		fmt.Println("Terjadi kesalahan saat mengecek saldo: ", err)
		dashboard(db)
		return
	}

	//validasi saldo
	if nominal > saldo {
		fmt.Println("Saldo tidak mencukupi untuk penarikan.")
		dashboard(db)
		return
	}

	//Update saldo
	newSaldo := saldo - nominal
	_, err = tx.Exec("UPDATE accounts SET balance =  ? WHERE id = ?", newSaldo, currentId)
	if err != nil {
		tx.Rollback()
		fmt.Println("Update saldo gagal:", err)
		dashboard(db)
		return
	}

	//Insert tabel transactions
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount) VALUES (?, 'withdraw', ?)", currentId, nominal)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal menyimpan histori transaksi:", err)
		dashboard(db)
		return
	}

	//Commit transaksi
	err = tx.Commit()
	if err != nil {
		fmt.Println("Gagal commit transaksi:", err)
		dashboard(db)
		return
	}

	fmt.Println("Penarikan berhasil!")
	fmt.Println("Sisa Saldo: ", newSaldo)
	dashboard(db)
}

func transfer(db *sql.DB) {
	fmt.Println("\n===== Transfer Tunai =====")

	var targetId, nominal, saldo int

	fmt.Print("Masukan Nomor Akun Tujuan: ")
	fmt.Scan(&targetId)

	if currentId == targetId {
		fmt.Print("Tidak dapat transfer ke akun anda")
		transfer(db)
		return
	}

	//cek targetId
	var cek int
	var name string
	err := db.QueryRow("SELECT id, name FROM accounts WHERE id = ?", targetId).Scan(&cek, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Akun tujuan tidak ditemukan.")
		} else {
			fmt.Println("Terjadi kesalahan:", err)
		}
		dashboard(db)
		return
	}

	fmt.Println("\nTujuan Transfer: ", name)
	fmt.Print("Masukkan jumlah transfer: Rp")
	fmt.Scan(&nominal)

	//validasi input
	if nominal <= 0 {
		fmt.Println("Nominal tidak valid!")
		transfer(db)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Gagal memulai transaksi:", err)
		dashboard(db)
		return
	}

	//cek saldo
	query := "SELECT balance FROM accounts WHERE id = ?"
	err = db.QueryRow(query, currentId).Scan(&saldo)
	if err != nil {
		fmt.Println("Terjadi kesalahan saat mengecek saldo: ", err)
		dashboard(db)
		return
	}

	//validasi saldo
	if nominal > saldo {
		fmt.Println("Saldo tidak mencukupi untuk transfer tunai.")
		dashboard(db)
		return
	}

	//Update saldo pengirim
	newSaldo := saldo - nominal
	_, err = tx.Exec("UPDATE accounts SET balance =  ? WHERE id = ?", newSaldo, currentId)
	if err != nil {
		tx.Rollback()
		fmt.Println("Update saldo gagal:", err)
		dashboard(db)
		return
	}

	//update saldo penerima
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", nominal, targetId)
	if err != nil {
		tx.Rollback()
		fmt.Println("Update saldo gagal:", err)
		dashboard(db)
		return
	}

	// Simpan histori transfer out (pengirim)
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount) VALUES (?, ?, ?)", currentId, "transfer_out", nominal)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal mencatat histori transfer out:", err)
		dashboard(db)
		return
	}

	// Simpan histori transfer in (penerima)
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount) VALUES (?, ?, ?)", targetId, "transfer_in", nominal)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal mencatat histori transfer out:", err)
		dashboard(db)
		return
	}

	//Commit transaksi
	err = tx.Commit()
	if err != nil {
		fmt.Println("Gagal commit transaksi:", err)
		dashboard(db)
		return
	}

	fmt.Println("Transfer berhasil!")
	dashboard(db)
}

func history(db *sql.DB) {
	//get data transaksi
	query := "SELECT type, amount, created_at FROM transactions WHERE account_id = ? ORDER BY created_at DESC"
	rows, err := db.Query(query, currentId)
	if err != nil {
		fmt.Println("Gagal mengambil histori transaksi:", err)
		dashboard(db)
		return
	}
	defer rows.Close()

	//show data transaksi
	fmt.Println("\n===== Histori Transaksi =====")
	hasTransaction := false
	for rows.Next() {
		var tipe string
		var amount int
		var createdAt string

		err := rows.Scan(&tipe, &amount, &createdAt)
		if err != nil {
			fmt.Println("Error membaca data:", err)
			return
		}

		fmt.Printf("[%s] %s Rp%d\n", createdAt, tipe, amount)
		hasTransaction = true
	}

	if !hasTransaction {
		fmt.Println("Belum ada transaksi.")
	}

	dashboard(db)
}
