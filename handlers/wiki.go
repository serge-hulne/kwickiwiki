package handlers

import (
	"encoding/json"
	"fmt"
	"kwickiwiki/models"
	"log"
	"net/http"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func renderTemplate(c *gin.Context, templateName string, data pongo2.Context) {
	tpl, err := pongo2.FromFile("templates/" + templateName)
	if err != nil {
		log.Printf("Template error: %v", err)
		c.String(http.StatusInternalServerError, "Template error: %s", err.Error())
		return
	}

	rendered, err := tpl.Execute(data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		c.String(http.StatusInternalServerError, "Template execution error: %s", err.Error())
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(rendered))
}

// Edit a wiki page
func EditPage(c *gin.Context) {
	title := c.Param("title")
	var page models.Page
	result := models.DB.Where("LOWER(title) = LOWER(?)", title).First(&page)

	if result.Error != nil {
		page = models.Page{Title: title, Content: ""}
	}

	// Extract metadata from JSON
	author := extractMetadata(page.Metadata, "Author")
	//category := extractMetadata(page.Metadata, "Category")
	published := extractMetadata(page.Metadata, "Published")

	renderTemplate(c, "edit.html", pongo2.Context{
		"Title":     page.Title,
		"Content":   page.Content,
		"Author":    author,
		"Published": published,
	})
}

// Save a wiki page
func SavePage(c *gin.Context) {
	title := c.Param("title")
	content := c.PostForm("content")

	// Collect metadata from form
	metadata := map[string]interface{}{
		"Author":    c.PostForm("author"),
		"Category":  c.PostForm("category"),
		"Published": c.PostForm("published") == "on",
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Fatalf("Failed to convert metadata to JSON: %v", err)
	}

	log.Println("Received save request for:", title)

	var page models.Page
	result := models.DB.Where("LOWER(title) = LOWER(?)", title).First(&page)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		// If the page exists without metadata, add it
		if len(page.Metadata) == 0 {
			page.Metadata = datatypes.JSON(metadataJSON)
		}
		page = models.Page{Title: title, Content: content, Metadata: datatypes.JSON(metadataJSON)}
		if err := models.DB.Create(&page).Error; err != nil {
			log.Fatalf("Failed to create new page: %v", err)
		}
	} else {
		page.Content = content
		page.Metadata = datatypes.JSON(metadataJSON)
		models.DB.Save(&page)
	}

	log.Println("Page saved successfully:", page.Title)
	c.Redirect(http.StatusFound, "/"+page.Title)
}

// Extract metadata field from JSON
func extractMetadata(metadata datatypes.JSON, key string) string {
	var meta map[string]interface{}
	if len(metadata) == 0 { // Handle empty metadata case
		return ""
	}
	if err := json.Unmarshal(metadata, &meta); err != nil {
		log.Printf("Error parsing metadata: %v", err)
		return ""
	}

	if value, exists := meta[key]; exists {
		return fmt.Sprintf("%v", value)
	}
	return ""
}

func ShowPage(c *gin.Context) {
	title := c.Param("title")

	// Special case for home page: Show content first, then list articles
	if title == "home" {
		var homePage models.Page
		result := models.DB.Where("LOWER(title) = LOWER(?)", "home").First(&homePage)

		if result.Error != nil {
			log.Printf("Error fetching home page: %v", result.Error)
			homePage.Content = "Welcome to the Wiki!"
		}

		// Fetch all pages for listing
		var pages []models.Page
		pageResult := models.DB.Order("updated_at DESC").Find(&pages)

		if pageResult.Error != nil {
			log.Printf("Error fetching pages: %v", pageResult.Error)
		}

		// Organize pages by category
		pageCategories := make(map[string][]map[string]string) // Store pages as maps
		for _, page := range pages {
			category := extractMetadata(page.Metadata, "Category")
			if category == "" {
				category = "Uncategorized"
			}

			// Convert page into a simple map with formatted date
			pageData := map[string]string{
				"Title":     page.Title,
				"UpdatedAt": page.UpdatedAt.Format("2006-01-02 15:04:05"), // Short date format
			}

			pageCategories[category] = append(pageCategories[category], pageData)
		}

		// Render the home page
		renderTemplate(c, "home.html", pongo2.Context{
			"Title":          "Home",
			"HomeContent":    homePage.Content,
			"PageCategories": pageCategories,
		})
		return
	}

	// Normal page lookup
	var page models.Page
	result := models.DB.Where("LOWER(title) = LOWER(?)", title).First(&page)

	if result.Error != nil {
		renderTemplate(c, "page.html", pongo2.Context{
			"Title":   title,
			"Content": "Page not found",
		})
		return
	}

	// Extract metadata
	author := extractMetadata(page.Metadata, "Author")
	category := extractMetadata(page.Metadata, "Category")
	published := extractMetadata(page.Metadata, "Published")

	renderTemplate(c, "page.html", pongo2.Context{
		"Title":     page.Title,
		"Content":   page.Content,
		"Author":    author,
		"Category":  category,
		"Published": published,
		"UpdatedAt": page.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func SearchPage(c *gin.Context) {
	query := c.Query("q") // Get search query from URL
	log.Println("Search query:", query)

	if query == "" {
		c.Redirect(http.StatusFound, "/home")
		return
	}

	query = strings.ToLower(query) // Convert query to lowercase

	// Perform case-insensitive search in title and content
	var pages []models.Page
	result := models.DB.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ?", "%"+query+"%", "%"+query+"%").Find(&pages)

	if result.Error != nil {
		log.Printf("Error searching pages: %v", result.Error)
	}

	// Organize search results by category
	pageCategories := make(map[string][]map[string]string)
	for _, page := range pages {
		category := extractMetadata(page.Metadata, "Category")
		if category == "" {
			category = "Uncategorized"
		}

		// Convert page into a simple map with formatted date
		pageData := map[string]string{
			"Title":     page.Title,
			"UpdatedAt": page.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		pageCategories[category] = append(pageCategories[category], pageData)
	}

	renderTemplate(c, "search.html", pongo2.Context{
		"Title":          "Search Results",
		"Query":          query,
		"PageCategories": pageCategories,
	})
}

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash password before storing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		models.DB.Preload("Roles").First(&user, userID)

		for _, role := range user.Roles {
			if role.Name == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}
