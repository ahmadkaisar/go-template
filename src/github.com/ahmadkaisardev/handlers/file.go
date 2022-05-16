package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func File(r *http.Request) (string, error) {
	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
	    return "", errors.New("failed on upload file")
	}

	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
	    return "", errors.New("failed on get directory")
	}

	filename := fmt.Sprintf("%s%s", strconv.FormatInt(time.Now().UnixNano(), 10), filepath.Ext(handler.Filename))

	fileLocation := filepath.Join(dir, "file", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
	    return "", errors.New("failed on open blank file")
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
	    return "", errors.New("failed on writing file")
	}

	return filename, nil
}