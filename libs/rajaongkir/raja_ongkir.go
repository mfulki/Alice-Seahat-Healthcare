package rajaongkir

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"
)

type RajaOngkir interface {
	Get(ctx context.Context, url string) ([]byte, error)
	Post(ctx context.Context, url string, payload any) ([]byte, error)
}

type rajaOngkirImpl struct {
	URL string
	Key string
}

func New(url, key string) (*rajaOngkirImpl, error) {
	ro := &rajaOngkirImpl{
		URL: url,
		Key: key,
	}

	if config.App.Env == constant.Production {
		if err := ro.ping(); err != nil {
			return nil, err
		}
	}

	return ro, nil
}

func (ro *rajaOngkirImpl) fetch(ctx context.Context, method, url string, payload any) ([]byte, error) {
	var reqBody = new(bytes.Reader)

	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, ro.URL+url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("key", ro.Key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func (ro *rajaOngkirImpl) ping() error {
	if _, err := ro.Get(context.Background(), "/province?id=1"); err != nil {
		return err
	}

	return nil
}

func (ro *rajaOngkirImpl) Get(ctx context.Context, url string) ([]byte, error) {
	return ro.fetch(ctx, http.MethodGet, url, nil)
}

func (ro *rajaOngkirImpl) Post(ctx context.Context, url string, payload any) ([]byte, error) {
	return ro.fetch(ctx, http.MethodPost, url, payload)
}
