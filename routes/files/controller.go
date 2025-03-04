package files

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paigexx/telegram-go-server/services"
)

type Handler struct {
	service services.FilesService
}

func newHandler(filesService services.FilesService) *Handler {
	return &Handler{
		service: filesService,
	}
}

func (h Handler) Upload(c *gin.Context) {
    err := c.Request.ParseMultipartForm(10 << 20) // 10MB
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form data"})
        return
    }

	// Retrieve the file from the form data
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		http.Error(c.Writer, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	chatID := c.Request.FormValue("chat_id")

	id, err := h.service.Upload(c, file, handler.Filename, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
  	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h Handler) List(c *gin.Context) {
	chatID := c.Param("chat_id")
	pageToken := c.Query("pageToken")
	files, err := h.service.List(c, chatID, pageToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func (h Handler) GetSignedUrl(c *gin.Context) {
	cid := c.Param("cid")
	url, err := h.service.GetSignedUrl(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}