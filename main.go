package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"a21hc3NpZ25tZW50/helper"
	"a21hc3NpZ25tZW50/model"
)

type StudentManager interface {
	Login(id string, name string) error
	Register(id string, name string, studyProgram string) error
	GetStudyProgram(code string) (string, error)
	ModifyStudent(name string, fn model.StudentModifier) error
}

type InMemoryStudentManager struct {
	sync.Mutex
	students             []model.Student
	studentStudyPrograms map[string]string
	failedLoginAttempts  map[string]int
	//add map for tracking login attempts here
	// TODO: answer here
}

func NewInMemoryStudentManager() *InMemoryStudentManager {
	return &InMemoryStudentManager{
		students: []model.Student{
			{
				ID:           "A12345",
				Name:         "Aditira",
				StudyProgram: "TI",
			},
			{
				ID:           "B21313",
				Name:         "Dito",
				StudyProgram: "TK",
			},
			{
				ID:           "A34555",
				Name:         "Afis",
				StudyProgram: "MI",
			},
		},
		studentStudyPrograms: map[string]string{
			"TI": "Teknik Informatika",
			"TK": "Teknik Komputer",
			"SI": "Sistem Informasi",
			"MI": "Manajemen Informasi",
		},
		failedLoginAttempts: make(map[string]int),
		//inisialisasi failedLoginAttempts di sini:
		// TODO: answer here
	}
}

func ReadStudentsFromCSV(filename string) ([]model.Student, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3 // ID, Name and StudyProgram

	var students []model.Student
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		student := model.Student{
			ID:           record[0],
			Name:         record[1],
			StudyProgram: record[2],
		}
		students = append(students, student)
	}
	return students, nil
}

func (sm *InMemoryStudentManager) GetStudents() []model.Student {
	return sm.students // TODO: replace this
}

func (sm *InMemoryStudentManager) Login(id string, name string) (string, error) {
	if len(id) == 0 || len(name) == 0 {
		return "", fmt.Errorf("ID or Name is undefined")
	}
	if sm.failedLoginAttempts[id] >= 3 {
		return "", fmt.Errorf("Login gagal: Batas maksimum login terlampaui")
	}
	// sm = NewInMemoryStudentManager()
	// students := sm.students
	// fmt.Println(students)
	// var nameExist string
	var FoundStudent *model.Student

	for i := range sm.students {
		if sm.students[i].ID == id && sm.students[i].Name == name {
			FoundStudent = &sm.students[i]
		}
	}
	// fmt.Println("FoundStudent : ", FoundStudent)
	if FoundStudent != nil {
		sm.failedLoginAttempts[id] = 0
		studyProgram := sm.studentStudyPrograms[FoundStudent.StudyProgram]
		return fmt.Sprintf("Login berhasil: Selamat datang %s! Kamu terdaftar di program studi: %s", name, studyProgram), nil
	}

	// for _, student := range sm.students {
	// 	if student.ID == id || student.Name == name { //JIKA NAMA & ID ADA LOOPING BERHENTI
	// 		FoundUntilEnd = true
	// 		break
	// 	} else {
	// 		FoundUntilEnd = false
	// 	}
	// }

	// if FoundUntilEnd { //jika nama ditemukan
	// 	return "Login berhasil: Selamat datang " + name +"! Kamu terdaftar di program studi: Teknik Informatika", nil
	// } else {
	// 	return "", fmt.Errorf("Login gagal: data mahasiswa tidak ditemukan")
	// }
	sm.failedLoginAttempts[id]++
	return "", fmt.Errorf("Login gagal: data mahasiswa tidak ditemukan") // TODO: replace this
}

func (sm *InMemoryStudentManager) RegisterLongProcess() {
	// 30ms delay to simulate slow processing
	time.Sleep(30 * time.Millisecond)
}

func (sm *InMemoryStudentManager) Register(id string, name string, studyProgram string) (string, error) {
	if len(id) == 0 || len(name) == 0 || len(studyProgram) == 0 {
		return "", fmt.Errorf("ID, Name or StudyProgram is undefined!")
	}
	var validProgram bool
	// fmt.Println("valid :", validProgram)
	// studentProg := sm.studentStudyPrograms [studyProgram]
	// // fmt.Println("studentprog", studentProg)
	if _, codeProg := sm.studentStudyPrograms[studyProgram]; !codeProg { //MENGECEK APABILA STUDYPROGRAM SAMA DALAM MAP, !CODEPROG ARTINYA TIDAK DITEMUKAN MAKA BERNILAI FALSE DAN PROGRAM DIdalmnya DIJALANKAN
		// fmt.Println("codeprog : ", codeProg)
		validProgram = false
		return "", fmt.Errorf("Study program " + studyProgram + " is not found")
	} else {
		validProgram = true
	}

	var IDFound bool
	// fmt.Println("found", IDFound)
	for _, student := range sm.students {
		if student.ID == id {
			IDFound = true
			break
		} else {
			IDFound = false
		}
	}

	if IDFound {
		return "", fmt.Errorf("Registrasi gagal: id sudah digunakan")
	}
	// fmt.Println("IDfound", IDFound)
	// fmt.Println("validProgram", validProgram)
	if validProgram && !IDFound { //JIKA PROGRAM VALID DAN ID TIDAK DITEMUKAN MAKA DIJJALANKAN PROGRAM , !IDFOUND KARENA IDFOUND SENDIRI FALSE JIKA ID BELUM DIGUNAKAN SEHINGGA AGAR DAPAT SAMA MENGGUNAKAN OPERATOR && DIBUAT SAMA
		newStudent := model.Student{ //AGAR TIPE DATA YANG DITAMBAHKAN SESUAI
			ID:           id,
			Name:         name,
			StudyProgram: studyProgram,
		}
		sm.students = append(sm.students, newStudent)
		return "Registrasi berhasil: " + name + " (" + studyProgram + ")", nil
	}

	// 30ms delay to simulate slow processing. DO NOT REMOVE THIS LINE
	sm.RegisterLongProcess()
	// Below lock is needed to prevent data race error. DO NOT REMOVE BELOW 2 LINES
	sm.Lock()
	defer sm.Unlock()

	return "", nil // TODO: replace this
}

func (sm *InMemoryStudentManager) GetStudyProgram(code string) (string, error) {
	if len(code) == 0 {
		return "", fmt.Errorf("ID, Name or StudyProgram is undefined!")
	}
	if _, codeProg := sm.studentStudyPrograms[code]; !codeProg { //akan mencari input code dalam mapp, jika tidak ditemukan (!codeprogram), program didalamnya dijalankna
		return "", fmt.Errorf("Kode program studi tidak ditemukan")

	} else {
		return sm.studentStudyPrograms[code], nil //KETIKA LAINNYA ATAU KODE SESUAI DENGAN DI MAP MAKA AKAN MENGEMBALIKAN VALUE DARI KEY(CODE) YANG DIBERIKAN
	}
}

func (sm *InMemoryStudentManager) ModifyStudent(name string, fn model.StudentModifier) (string, error) {
	var nameFound bool
	for i, student := range sm.students {
		if student.Name == name {
			nameFound = true
			if err := fn(&sm.students[i]); err != nil { //MEMANGGIL FUNGSI MODIF KETIKA NAMA DITEMUKAN
				return "", err
			}
			break
		}
	}
	if !nameFound {
		return "", fmt.Errorf("Mahasiswa tidak ditemukan.")
	}
	return "Program studi mahasiswa berhasil diubah.", nil // TODO: replace this
}

func (sm *InMemoryStudentManager) ChangeStudyProgram(programStudi string) model.StudentModifier {
	return func(s *model.Student) error {
		if _, correct := sm.studentStudyPrograms[programStudi]; !correct { //MEMERIKSA APAKAH INPUT PROGRAM STUDI SAMA
			return fmt.Errorf("Kode program studi tidak ditemukan")
		}
		s.StudyProgram = programStudi
		return nil // TODO: replace this
	}
}

func (sm *InMemoryStudentManager) ImportStudents(filenames []string) error {
	var wg sync.WaitGroup
	wg.Add(len(filenames))
	// fmt.Println(len(filenames))
	for _, filename := range filenames {
		go func(filename string) {
			defer wg.Done()
			students, err := ReadStudentsFromCSV(filename)
			fmt.Println("student", students)
			if err != nil {
				fmt.Println("Error reading CSV:", err)
				return
			}
			for _, student := range students {
				_, err := sm.Register(student.ID, student.Name, student.StudyProgram)
				if err != nil {
					fmt.Printf("Error registering student %s: %s\n", student.Name, err)

				}
			}
		}(filename)
	}
	wg.Wait()
	return nil // TODO: replace this
}

func (sm *InMemoryStudentManager) SubmitAssignmentLongProcess() {
	// 3000ms delay to simulate slow processing
	time.Sleep(30 * time.Millisecond)
}

func (sm *InMemoryStudentManager) SubmitAssignments(numAssignments int) {

	start := time.Now()
	const numWorkers = 3
	jobs := make(chan int, numAssignments)
	results := make(chan struct{}, numAssignments)

	// Goroutine untuk menambahkan tugas ke dalam job queue
	go func() {
		for i := 0; i < numAssignments; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Goroutine untuk menangani pengiriman tugas
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d: Processing assignment %d\n", workerID, job)
				sm.SubmitAssignmentLongProcess()
				fmt.Printf("Worker %d: Finished assignment %d\n", workerID, job)
				results <- struct{}{}
			}
		}(i + 1)
	}

	// Menunggu semua goroutine selesai
	go func() {
		wg.Wait()
		close(results)
	}()

	// Menghitung jumlah tugas yang selesai
	var completedTasks int
	for range results {
		completedTasks++
	}
	// TODO: answer here

	elapsed := time.Since(start)
	fmt.Printf("Submitting %d assignments took %s\n", numAssignments, elapsed)
}

func main() {
	manager := NewInMemoryStudentManager()

	for {
		helper.ClearScreen()
		students := manager.GetStudents()
		for _, student := range students {
			fmt.Printf("ID: %s\n", student.ID)
			fmt.Printf("Name: %s\n", student.Name)
			fmt.Printf("Study Program: %s\n", student.StudyProgram)
			fmt.Println()
		}

		fmt.Println("Selamat datang di Student Portal!")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Get Study Program")
		fmt.Println("4. Modify Student")
		fmt.Println("5. Bulk Import Student")
		fmt.Println("6. Submit assignment")
		fmt.Println("7. Exit")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Pilih menu: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			helper.ClearScreen()
			fmt.Println("=== Login ===")
			fmt.Print("ID: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)

			fmt.Print("Name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			msg, err := manager.Login(id, name)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
			fmt.Println(msg)
			// Wait until the user presses any key
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')
		case "2":
			helper.ClearScreen()
			fmt.Println("=== Register ===")
			fmt.Print("ID: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)

			fmt.Print("Name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Study Program Code (TI/TK/SI/MI): ")
			code, _ := reader.ReadString('\n')
			code = strings.TrimSpace(code)

			msg, err := manager.Register(id, name, code)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
			fmt.Println(msg)
			// Wait until the user presses any key
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')
		case "3":
			helper.ClearScreen()
			fmt.Println("=== Get Study Program ===")
			fmt.Print("Program Code (TI/TK/SI/MI): ")
			code, _ := reader.ReadString('\n')
			code = strings.TrimSpace(code)

			if studyProgram, err := manager.GetStudyProgram(code); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Program Studi: %s\n", studyProgram)
			}
			// Wait until the user presses any key
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')
		case "4":
			helper.ClearScreen()
			fmt.Println("=== Modify Student ===")
			fmt.Print("Name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Program Studi Baru (TI/TK/SI/MI): ")
			code, _ := reader.ReadString('\n')
			code = strings.TrimSpace(code)

			msg, err := manager.ModifyStudent(name, manager.ChangeStudyProgram(code))
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
			fmt.Println(msg)

			// Wait until the user presses any key
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')
		case "5":
			helper.ClearScreen()
			fmt.Println("=== Bulk Import Student ===")

			// Define the list of CSV file names
			csvFiles := []string{"students1.csv", "students2.csv", "students3.csv"}

			err := manager.ImportStudents(csvFiles)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Println("Import successful!")
			}

			// Wait until the user presses any key
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')

		case "6":
			helper.ClearScreen()
			fmt.Println("=== Submit Assignment ===")

			// Enter how many assignments you want to submit
			fmt.Print("Enter the number of assignments you want to submit: ")
			numAssignments, _ := reader.ReadString('\n')

			// Convert the input to an integer
			numAssignments = strings.TrimSpace(numAssignments)
			numAssignmentsInt, err := strconv.Atoi(numAssignments)

			if err != nil {
				fmt.Println("Error: Please enter a valid number")
			}

			manager.SubmitAssignments(numAssignmentsInt)

			// Wait until the user presses any key
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')
		case "7":
			helper.ClearScreen()
			fmt.Println("Goodbye!")
			return
		default:
			helper.ClearScreen()
			fmt.Println("Pilihan tidak valid!")
			helper.Delay(5)
		}

		fmt.Println()
	}
}
