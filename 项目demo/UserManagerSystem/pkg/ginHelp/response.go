package ginHelp

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response should not reflect the specific error information, here should explain
// the simple description of the error, the specific error information should be
// reflected in the log file

type ResponseBody struct {
	Data interface{}
}

type Download struct {
	FileName string
	Data     *bytes.Buffer
}

// DownloadFile used for download file api
func (d *Download) DownloadFile(c *gin.Context) error {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+d.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", d.Data.Len()))
	_, err := c.Writer.Write(d.Data.Bytes())
	return err
}
