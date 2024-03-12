package lib

import (
	"io"
	"net/http"

	breaker "github.com/sony/gobreaker"
)

var cb *breaker.CircuitBreaker

type Marshaller[T any] func(data []byte) (T, error)

func init() {
	cb = breaker.NewCircuitBreaker(breaker.Settings{
		Name: "client",
		ReadyToTrip: func(counts breaker.Counts) bool {
			return counts.Requests >= 3
		},
	})
}

func GetUrl[T any](url string, m Marshaller[T]) (T, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	if err != nil {
		return *new(T), err
	}

	return m(body.([]byte))
}

func DownloadFile(url string, followRedir bool) ([]byte, error) {
	res, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if followRedir {
			resp, err = http.DefaultClient.Do(resp.Request)
		}

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		bucket := make([]byte, resp.ContentLength)

		_, err = io.ReadFull(resp.Body, bucket)
		if err != nil {
			return nil, err
		}

		return bucket, nil
		// _, err = io.Copy(os.Stdout, resp.Body)
		// if err != nil {
		// 	return nil, err
		// }

		// return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}
