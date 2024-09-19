package user

import (
	"errors"

	"log"
	"task-api/internal/auth"
	"task-api/internal/model"

	

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	Repository Repository
	secret     string
}

func NewService(db *gorm.DB, secret string) Service {
	return Service{
		Repository: NewRepository(db),
		secret:     secret,
	}
}

// internal/user/service.go
func (service Service) Login(req model.RequestLogin) (string, error) {
	// TODO: Check username and password here
	user, err := service.Repository.FindOneByUsername(req.Username)
	if err != nil {
		return "", errors.New("invalid user or password")
	}
	// req.Password // req password
	// user.Password // hashed password
	//ตรงนี้จะเป็นส่วนของ logic ที่ทำงานที่ไปร้องขอ password มาจาก repo แล้วมาเทียบโดยฟังชั่น checkPasswordHash
	if ok := checkPasswordHash(req.Password, user.Password); !ok {
		return "", errors.New("invalid user or password")
	}
	
	// TODO: Create token here
	token, err := auth.CreateToken(user.Username, service.secret)
	if err != nil {
		log.Println("Fail to create token")
		return "", errors.New("something went wrong")
	}
	return token, nil
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
//แล้วทำไม อา.เอาฟังก์ชั่นสร้างtoken ไปไว้ที่ aunt.go ละ?