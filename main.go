package main

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	templates = template.Must(template.ParseFiles(
		"templates/login.html",
		"templates/register.html",
		"templates/student.html",
		"templates/teacher.html",
		"templates/setquestions.html",
		"templates/viewresults.html",
		"templates/quiz.html",
	))

	client *mongo.Client
)

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"` // "student" or "teacher"
}

type Question struct {
	ID       string   `bson:"_id,omitempty"`
	Question string   `bson:"question"`
	Options  []string `bson:"options"`
	Answer   string   `bson:"answer"`
}

type StudentAnswer struct {
	Username string `bson:"username"`
	Question string `bson:"question"`
	Answer   string `bson:"answer"`
}

func main() {
	// MongoDB Atlas connection string
	mongoURI := "mongodb+srv://jidaar718:tRelmEXYu7NEcGFz@cluster0.j9n5kuh.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	// Connect to MongoDB Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Ensure connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")
	fmt.Println("Server is running")

	// Set the random seed once at the start
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/student", studentHandler)
	http.HandleFunc("/teacher", teacherHandler)
	http.HandleFunc("/setQuestions", setQuestionsHandler)
	http.HandleFunc("/viewResults", viewResultsHandler)
	http.HandleFunc("/downloadResults", downloadResultsHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/submitQuiz", submitQuizHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func generateCaptchaCode() string {
	code := rand.Intn(999999)
	return fmt.Sprintf("%06d", code)
}

func sendCaptchaEmail(email, captchaCode string) error {
	from := "userstest323@gmail.com"
	password := "Abbas@0703"
	to := email
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Verification Code\n\n" +
		"Your verification code is: " + captchaCode

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println("Failed to send email:", err)
	}
	return err
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := hashPassword(r.FormValue("password"))
		role := r.FormValue("role")

		user := User{
			Username: username,
			Password: password,
			Role:     role,
		}

		collection := client.Database("quizdb").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := collection.InsertOne(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := hashPassword(r.FormValue("password"))

		collection := client.Database("quizdb").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var user User
		err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil || user.Password != password {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		switch user.Role {
		case "student":
			http.Redirect(w, r, "/quiz", http.StatusSeeOther)
		case "teacher":
			http.Redirect(w, r, "/teacher", http.StatusSeeOther)
		default:
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return
	}
	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		answer := r.FormValue("answer")
		questionID := r.FormValue("questionID")

		collection := client.Database("quizdb").Collection("student_answers")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := collection.InsertOne(ctx, StudentAnswer{
			Username: username,
			Question: questionID,
			Answer:   answer,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/quiz", http.StatusSeeOther)
		return
	}

	collection := client.Database("quizdb").Collection("questions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var questions []Question
	if err = cursor.All(ctx, &questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "student.html", questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func teacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/setQuestions", http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "teacher.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func setQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle CSV upload
		file, _, err := r.FormFile("csvfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collection := client.Database("quizdb").Collection("questions")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		for _, record := range records {
			if len(record) < 3 {
				continue
			}
			_, err := collection.InsertOne(ctx, Question{
				Question: record[0],
				Options:  []string{record[1], record[2]},
				Answer:   record[1], // Assuming the first option is the correct answer
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/teacher", http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "setquestions.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewResultsHandler(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("quizdb").Collection("student_answers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var answers []StudentAnswer
	if err = cursor.All(ctx, &answers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "viewresults.html", answers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func downloadResultsHandler(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("quizdb").Collection("student_answers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	file, err := os.Create("results.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Username", "Question", "Answer"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var answers []StudentAnswer
	if err = cursor.All(ctx, &answers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, answer := range answers {
		if err := writer.Write([]string{answer.Username, answer.Question, answer.Answer}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename=results.csv")
	w.Header().Set("Content-Type", "text/csv")
	http.ServeFile(w, r, "results.csv")
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("quizdb").Collection("questions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var questions []Question
	if err = cursor.All(ctx, &questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "quiz.html", questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func submitQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		question := r.FormValue("question")
		answer := r.FormValue("answer")

		collection := client.Database("quizdb").Collection("student_answers")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := collection.InsertOne(ctx, StudentAnswer{
			Username: username,
			Question: question,
			Answer:   answer,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/quiz", http.StatusSeeOther)
	}
}
