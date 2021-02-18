package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func SendMail(mailData []byte, route string) error{
	token := loginApiMail()

	req, err := http.NewRequest("POST", constants.ApiMail+ "/" + route, bytes.NewBuffer(mailData))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.Post(helper.ApiMail + "/newpassword", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("error in read body: %s", err))
		return err
	}
	byt := []byte(string(body))
	if !strings.Contains(string(body), "accepted") {
		return errors.New("error de respuesta")
	}
	if strings.Contains(string(body), "quota exceeded") || strings.Contains(string(body), "ECONNECTION"){
		return errors.New("cuota de envios diarios excedida")
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		log.Println(fmt.Sprintf("error in unarchall json: %s", err))
		return err
	}
	fmt.Println(dat)
	return nil
}

func SendFileMail(mail string, route string, filename string) error {
	values := map[string]io.Reader{
		"file":  mustOpen(filename), // lets assume its this file
		"mail": strings.NewReader(mail),
	}
	err := Upload(constants.ApiMail + "/" + route, values)
	if err != nil {
		return err
	}
	return nil
	/*data, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer data.Close()
	req, err := http.NewRequest("POST", constants.ApiMail + "/" + route, data)
	if err != nil {
		log.Fatal(err)
	}
	token := loginApiMail()
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	byt := []byte(string(body))
	fmt.Println(byt)
	return nil*/
}

func Upload(url string, values map[string]io.Reader) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	_ = w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	token := loginApiMail()
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	fmt.Println("sending!")
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func loginApiMail() string{
	secret, err := json.Marshal(map[string]string{
		"secret": constants.SecretApiMail,
	})
	if err != nil {
		fmt.Println(err)
	}
	respToken, err := http.Post(constants.ApiMail+ "/login", "application/json", bytes.NewBuffer(secret))
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
