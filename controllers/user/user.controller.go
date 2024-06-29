package user

import (
	"final-task-pbi-fullstackdev/database/dbconfig"
	"final-task-pbi-fullstackdev/helpers"
	"final-task-pbi-fullstackdev/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary Login user
// @Description Login user dengan payload input dan akan membuat jwt token baru yang akan di store di cookie
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} swagger.LoginValue
// @Failure 400 {object} swagger.ReturnValue
// @Router /users/login [get]
// =====LOGIN=====
func Login(ctx *gin.Context) {
	var userInput models.User //initialize data user yang diinput

	if err := ctx.ShouldBindBodyWithJSON(&userInput); err != nil { //validasi untuk memastikan bahwa struct sudah sesuai di binding pada json
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	var user models.User //initialize spesifik data user

	if err := dbconfig.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound: //cek email, apakah email valid/sudah terdaftar, jika belum tampilkan error di response
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email belum terdaftar",
			})
			return
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ //cek apakah password match dengan yg ada di database atau tidak, jika tidak tampilkan error di response
			"message": "Password tidak sesuai",
		})
		return
	}

	dbconfig.DB.Preload("Photo").Find(&user)                               //preload photo(untuk populate constraints yang terhubung pada table photo)
	token := helpers.GenerateJWT(user.Email)                               // generate token jwt
	ctx.SetCookie("jwt-token", token, 3600, "/", "localhost", false, true) // set secure cookie
	ctx.JSON(http.StatusOK, gin.H{                                         // menampilkan response sukses beserta data usernya
		"message": "Success login user",
		"data": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"photo":    user.Photo,
		},
	})
}

// @Summary Register user
// @Description Register user dengan input payload, lalu payload akan divalidasi. Kemudian user baru akan dibuat dan disimpan ke database.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} swagger.ReturnValue
// @Failure 400 {object} swagger.ReturnValue
// @Failure 409 {object} swagger.ReturnValue
// @Failure 500 {object} swagger.ReturnValue
// @Router /users/register [post]
// =====REGISTER=====
func Register(ctx *gin.Context) {
	var user models.User        // initialize data user yang diinput
	user.ID = uuid.New()        // generate userId by UUID
	validate := validator.New() // initialize validator

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil { //validasi untuk memastikan bahwa struct sudah sesuai di binding pada json
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	if err := validate.Struct(&user); err != nil { // validasi setiap field, jika ada error tampilkan di response
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorsMap := make(map[string]string)
			for _, fieldError := range validationErrors {
				field := fieldError.Field()
				tag := fieldError.Tag()
				errorsMap[field] = helpers.DisplayValidationErrors(field, tag)
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"message": errorsMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validasi error"})
		return
	}

	if err := dbconfig.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") { // cek apabila email sudah terdaftar, maka akan menampilkan pesan error
			ctx.JSON(http.StatusConflict, gin.H{"message": "Email sudah terdaftar"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal daftar akun baru"})
		}
		return
	}

	hashPassword, _ := helpers.HashPassword(user.Password) //generate hash password
	user.Password = hashPassword                           //simpan hash password di field Password
	dbconfig.DB.Save(&user)                                //store user baru ke dalam database
	ctx.JSON(http.StatusOK, gin.H{                         // tampilkan status dan response sukses
		"message": "Berhasil membuat akun baru",
	})
}

// @Summary Logout user
// @Description Logout user dan cookie user jwt akan dihapus
// @Tags users
// @Produce json
// @Success 200 {object} swagger.ReturnValue
// @Router /users/logout [get]
// =====LOGOUT=====
func Logout(ctx *gin.Context) {
	//handler ini sebagai utilitas saja untuk menghapus cookie user jwt
	ctx.SetCookie("jwt-token", "", -1, "/", "localhost", false, true) //reset cookie
	ctx.JSON(http.StatusOK, gin.H{                                    // tampilkan status dan response sukses
		"message": "Berhasil logout",
	})
}

// @Summary Update user
// @Description Update user dengan input payload dan parameter, lalu payload akan divalidasi. Kemudian data user akan diupdate.
// @Tags users
// @Accept json
// @Produce json
// @Param userId query string true "User ID"
// @Success 200 {object} swagger.ReturnValue
// @Failure 400 {object} swagger.ReturnValue
// @Router /users/:userId [put]
// =====UPDATE=====
func UpdateUser(ctx *gin.Context) {
	var user models.User                  // initialize data user yang diinput
	userId := ctx.Param("userId")         // ambil parameter userId
	userIdParse := uuid.MustParse(userId) //parsing userId string menjadi uuid
	validate := validator.New()           //initialize validator

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil { //validasi untuk memastikan bahwa struct sudah sesuai di binding pada json
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	if err := validate.Struct(&user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok { // validasi setiap field, jika ada error tampilkan di response
			errorsMap := make(map[string]string)
			for _, fieldError := range validationErrors {
				field := fieldError.Field()
				tag := fieldError.Tag()
				errorsMap[field] = helpers.DisplayValidationErrors(field, tag)
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"message": errorsMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validasi error"})
		return
	}

	hashPassword, _ := helpers.HashPassword(user.Password)                     //generate hash password
	dbconfig.DB.Model(&user).Where("id = ?", userIdParse).Updates(models.User{ //update user
		Username: user.Username,
		Email:    user.Email,
		Password: hashPassword,
	})
	ctx.JSON(http.StatusOK, gin.H{ //tampilkan response sukses
		"message": "Berhasil update user",
	})
}

// @Summary Delete user
// @Description Delete user dengan input payload, lalu payload akan divalidasi. Kemudian data user akan diupdate.
// @Tags users
// @Accept json
// @Produce json
// @Param userId query string true "User ID"
// @Success 200 {object} swagger.ReturnValue
// @Failure 400 {object} swagger.ReturnValue
// @Router /users/:userId [delete]
// =====DELETE=====
func DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")          //ambil parameter userId
	userIdParse, err := uuid.Parse(userId) //parsing parameter userId menjadi uuid

	if err != nil { //validasi apakah parameter valid atau tidak
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter tidak valid",
		})
		return
	}

	if err := dbconfig.DB.Delete(&models.User{}, userIdParse).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ //cek apakah ada error atau tidak, apakah userId valid atau tidak
			"message": "Parameter tidak valid",
		})
		return
	}

	dbconfig.DB.Delete(&models.User{}, userIdParse) //hapus user yang ada di database berdasarkan parameter userId
	ctx.JSON(http.StatusOK, gin.H{                  // tampilkan response sukses
		"message": "Berhasil hapus user",
	})
}
