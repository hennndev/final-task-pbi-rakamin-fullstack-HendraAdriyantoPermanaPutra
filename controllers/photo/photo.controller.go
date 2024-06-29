package photo

import (
	"final-task-pbi-fullstackdev/database/dbconfig"
	"final-task-pbi-fullstackdev/helpers"
	"final-task-pbi-fullstackdev/models"
	"net/http"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// @Summary Get Photos
// @Description Get photos untuk mengembalikan response semua foto dalam bentuk array.
// @Tags photos
// @Produce json
// @Success 200 {object} swagger.GetPhotosValue
// @Failure 400 {object} swagger.ReturnValue
// @Router /photos [get]
// =====GET PHOTOS=====
func GetPhotos(ctx *gin.Context) {
	var photos []models.Photo      // initialize data photos yang akan ditampilkan dalam bentuk slice
	dbconfig.DB.Find(&photos)      //query data photos
	ctx.JSON(http.StatusOK, gin.H{ //kemudian tampilkan response dan data photos
		"message": "Berhasil request data photos",
		"photos":  photos,
	})
}

// @Summary Post Photo
// @Description Post photo untuk mengupload photo.
// @Tags photos
// @Accept json
// @Produce json
// @Success 200 {object} swagger.ReturnValue
// @Failure 400 {object} swagger.ReturnValue
// @Failure 409 {object} swagger.ReturnValue
// @Router /photos [post]
// =====POST PHOTO=====
func PostPhoto(ctx *gin.Context) {
	var user models.User        // initialize data user
	var photo models.Photo      // initialize data photo
	validate := validator.New() // initialize validator

	if err := ctx.ShouldBindBodyWithJSON(&photo); err != nil { //validasi untuk memastikan bahwa struct sudah sesuai di binding pada json
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	if err := validate.Struct(&photo); err != nil { // validasi setiap field, jika ada error tampilkan di response
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

	if err := dbconfig.DB.First(&user, &photo.UserID).Error; err != nil { // memastikan apakah userID invalid atau tidak
		ctx.JSON(http.StatusConflict, gin.H{"message": "UserID invalid"})
		return
	}

	photo.ID = uuid.New()          //kemduain assign id dengan value UUID
	dbconfig.DB.Save(&photo)       //lalu store data photo ke dalam database
	ctx.JSON(http.StatusOK, gin.H{ // tampilkan response suksesnya
		"message": "Berhasil post foto",
	})
}

// @Summary Update Photo
// @Description Update photo untuk mengedit data photo yang ada di database.
// @Tags photos
// @Produce json
// @Param photoId query string true "Photo ID"
// @Success 200 {object} swagger.ReturnValue
// @Failure 400 {object} swagger.ReturnValue
// @Failure 409 {object} swagger.ReturnValue
// @Router /photos/:photoId [put]
// =====UPDATE PHOTO=====
func UpdatePhoto(ctx *gin.Context) {
	var user models.User                     // initialize data user
	var photo models.Photo                   // initialize data photo
	validate := validator.New()              // initialize validator
	photoId := ctx.Param("photoId")          //ambil parameter photoId
	photoIdParse, err := uuid.Parse(photoId) //parsing photoId string menjadi uuid

	if err != nil { // cek apakah parameter valid atau tidak
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter tidak valid",
		})
		return
	}

	if err := ctx.ShouldBindBodyWithJSON(&photo); err != nil { //validasi untuk memastikan bahwa struct sudah sesuai di binding pada json
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	if err := validate.Struct(&photo); err != nil { // validasi setiap field, jika ada error tampilkan di response
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

	if err := dbconfig.DB.First(&user, &photo.UserID).Error; err != nil { // cek apakah userID invalid atau tidak
		ctx.JSON(http.StatusConflict, gin.H{"message": "UserID invalid"})
		return
	}

	dbconfig.DB.Model(&photo).Where("id = ?", photoIdParse).Updates(models.Photo{ //update photo
		Title:    photo.Title,
		Caption:  photo.Caption,
		PhotoUrl: photo.PhotoUrl,
		UserID:   photo.UserID,
	})

	if err := dbconfig.DB.First(&photo, photoIdParse).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ //cek apakah ada photo dengan userId tersebut atau tidak, jika tidak ada tampilkan response
			"message": "Parameter tidak valid",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{ //tampilkan response suksesnya
		"message": "Berhasil update foto",
	})
}

// @Summary Delete Photo
// @Description Delete photo untuk menghapus data photo yang ada di database.
// @Tags photos
// @Produce json
// @Param photoId query string true "Photo ID"
// @Success 200 {object} swagger.ReturnValue
// @Failure 400 {object} swagger.ReturnValue
// @Router /photos/:photoId [delete]
// =====DELETE PHOTO=====
func DeletePhoto(ctx *gin.Context) {
	var photo models.Photo                   //initialize photo data
	photoId := ctx.Param("photoId")          // mengambil parameter photoId
	photoIdParse, err := uuid.Parse(photoId) //parsing parameter photoId

	if err != nil { // cek apakah parameter valid atau tidak
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter tidak valid",
		})
		return
	}

	if err := dbconfig.DB.First(&photo, photoIdParse).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ //cek apakah ada photo dengan userId tersebut atau tidak, jika tidak ada tampilkan response
			"message": "Parameter tidak valid",
		})
		return
	}

	if err := dbconfig.DB.Delete(&models.Photo{}, photoIdParse).Error; err != nil { //cek apakah ada error atau tidak, apakah photoId valid atau tidak
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter tidak valid",
		})
		return
	}
	dbconfig.DB.Delete(&models.Photo{}, photoIdParse) //lalu delete photo

	ctx.JSON(http.StatusOK, gin.H{ // kemudian tampilkan response sukses
		"message": "Berhasil hapus foto",
	})
}
