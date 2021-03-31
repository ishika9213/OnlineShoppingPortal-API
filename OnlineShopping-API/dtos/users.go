package dtos

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
)

type RegisterRequestDto struct {
	Username             string `form:"username" json:"username" xml:"username"  binding:"required"`
	FirstName            string `form:"first_name" json:"first_name" xml:"first_name" binding:"required"`
	LastName             string `form:"last_name" json:"last_name" xml:"last_name" binding:"required"`
	Email                string `form:"email" json:"email" xml:"email" binding:"required"`
	Password             string `form:"password" json:"password" xml:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" xml:"password-confirmation" binding:"required"`
}

type LoginRequestDto struct {
	// Username string `form:"username" json:"username" xml:"username" binding:"exists,username"`
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password"json:"password" binding:"exists,min=8,max=255"`

	userModel models.User `json:"-"`
}

func CreateLoginSuccessful(user *models.User) map[string]interface{} {
	var roles = make([]string, len(user.Roles))

	for i := 0; i < len(user.Roles); i++ {
		roles[i] = user.Roles[i].Name
	}

	return map[string]interface{}{
		"success": true,
		"token":   user.GenerateJwtToken(),
		"user": map[string]interface{}{
			"username": user.Username,
			"id":       user.ID,
			"roles":    roles,
		},
	}
}

func GetUserBasicInfo(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	}
}

func CreatedUserPagedResponse(request *http.Request, products []models.User, page, page_size, count int, commentsCount []int) interface{} {
	var resources = make([]interface{}, len(products))
	for index, user := range user {
		resources[index] = CreateUserDto(&product, commentsCount[index])
	}
	return CreatePagedResponse(request, resources, "products", page, page_size, count)
}

func CreateProductDto(product *models.User, commentCount int) map[string]interface{} {

	var tags = make([]map[string]interface{}, len(product.Tags))
	var categories = make([]map[string]interface{}, len(product.Categories))
	var images = make([]string, len(product.Images))

	for index, tag := range product.Tags {
		tags[index] = map[string]interface{}{
			"id":   tag.ID,
			"name": tag.Name,
			"slug": tag.Slug,
		}
	}

	for index, category := range product.Categories {
		categories[index] = map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
			"slug": category.Slug,
		}
	}
	replaceAllFlag := -1
	for index, image := range product.Images {
		images[index] = strings.Replace(image.FilePath, "\\", "/", replaceAllFlag)
	}

	for index, tag := range product.Tags {
		tags[index] = map[string]interface{}{
			"id":   tag.ID,
			"name": tag.Name,
			"slug": tag.Slug,
		}
	}

	result := map[string]interface{}{
		"id":         product.ID,
		"name":       product.Name,
		"slug":       product.Slug,
		"price":      product.Price,
		"stock":      product.Stock,
		"tags":       tags,
		"categories": categories,
		"image_urls": images,
		"created_at": product.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		"updated_at": product.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}

	if commentCount >= 0 {
		// "comments_count": product.CommentsCount,
		result["comments_count"] = commentCount
	}
	return result
}