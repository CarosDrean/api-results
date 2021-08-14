package utils

import (
	"bytes"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func SendMail(mailData []byte, route string, token string) ([]byte, error) {
	req, err := http.NewRequest("POST", constants.ApiMail+"/"+route, bytes.NewBuffer(mailData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		e := &models.Error{}

		if err := e.Decode(resp.Body); err != nil {
			return nil, fmt.Errorf("errorResponse.decode(): %w", err)
		}

		return nil, e
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UploadFile(route string, filename string, token string) (models.MailFileResponse, error) {
	values := map[string]io.Reader{
		"file": mustOpen(filename), // lets assume its this file
	}

	resApiMail, err := upload(constants.ApiMail+"/"+route, values, token)
	if err != nil {
		return models.MailFileResponse{}, err
	}

	return resApiMail, nil
}

func upload(url string, values map[string]io.Reader, token string) (models.MailFileResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, r := range values {
		var err error
		var fw io.Writer

		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return models.MailFileResponse{}, err
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return models.MailFileResponse{}, err
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			return models.MailFileResponse{}, err
		}
	}

	_ = w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return models.MailFileResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.MailFileResponse{}, err
	}

	if res.StatusCode >= 400 {
		e := &models.Error{}

		if err := e.Decode(res.Body); err != nil {
			return models.MailFileResponse{}, fmt.Errorf("errorResponse.decode(): %w", err)
		}

		return models.MailFileResponse{}, e
	}

	resApiMail := models.MailFileResponse{}

	if err := resApiMail.Decode(res.Body); err != nil {
		return models.MailFileResponse{}, err
	}

	return resApiMail, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
