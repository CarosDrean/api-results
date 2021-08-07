package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func SendMail(mailData []byte, route string, token string) error {
	req, err := http.NewRequest("POST", constants.ApiMail+"/"+route, bytes.NewBuffer(mailData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 500 {
		fmt.Println(resp)
		return nil
	}
	fmt.Println(resp)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("error in read body: %s", err))
		return err
	}

	byt := []byte(string(body))
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		log.Println(fmt.Sprintf("error in unarchall json: %s", err))
		return err
	}
	fmt.Println(dat)
	return nil
}

func SendFileMail(route string, filename string, token string) (models.MailFileRes, error) {
	values := map[string]io.Reader{
		"file": mustOpen(filename), // lets assume its this file
	}

	resApiMail, err := Upload(constants.ApiMail+"/"+route, values, token)
	if err != nil {
		return models.MailFileRes{}, err
	}

	return resApiMail, nil
}

func Upload(url string, values map[string]io.Reader, token string) (models.MailFileRes, error) {
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
				return models.MailFileRes{}, err
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return models.MailFileRes{}, err
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			return models.MailFileRes{}, err
		}
	}

	_ = w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return models.MailFileRes{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.MailFileRes{}, err
	}

	if res.StatusCode >= 500 {
		//fmt.Errorf()	//TODO: aqui nos quedamos :D
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return models.MailFileRes{}, err
	}

	byt := []byte(string(body))
	resApiMail := models.MailFileRes{}

	if err := json.Unmarshal(byt, &resApiMail); err != nil {
		return models.MailFileRes{}, err
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

func loginApiMail() string {
	secret, err := json.Marshal(map[string]string{
		"secret": constants.SecretApiMail,
	})
	if err != nil {
		fmt.Println(err)
	}
	respToken, err := http.Post(constants.ApiMail+"/login", "application/json", bytes.NewBuffer(secret))
	if err != nil {
		log.Panic(err)
	}
	defer respToken.Body.Close()
	body, err := ioutil.ReadAll(respToken.Body)
	if err != nil {
		log.Panic(err)
	}
	byt := []byte(string(body))
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	return dat["token"].(string)
}
