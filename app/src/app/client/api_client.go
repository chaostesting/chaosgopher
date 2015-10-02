package client

import "strconv"
import "bytes"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "net/url"
import "app/data"

type APIClient struct {
	*http.Client
	Host string
}

func (api *APIClient) GetAllTodos() ([]*data.TodoItem, error) {
	result := []*data.TodoItem{}
	if e := api.request(&result, "GET", "/todos", nil); e != nil {
		return nil, e
	}

	return result, nil
}

func (api *APIClient) PatchTodo(id int64, updates *data.TodoItem) (*data.APIResult, error) {
	result := &data.APIResult{}
	if e := api.request(result, "PATCH", "/todos/"+strconv.FormatInt(id, 10), updates); e != nil {
		return nil, e
	}

	return result, nil
}

func (api *APIClient) request(result interface{}, method string, path string, body interface{}) error {
	u, e := url.Parse(api.Host + path)
	if e != nil {
		return e
	}

	buf := &bytes.Buffer{}
	if body != nil {
		if e := json.NewEncoder(buf).Encode(body); e != nil {
			return e
		}
	}

	req := &http.Request{
		Method: method,
		URL:    u,

		Body:          ioutil.NopCloser(buf),
		ContentLength: int64(buf.Len()),
	}

	if api.Client == nil {
		api.Client = &http.Client{}
	}

	resp, e := api.Client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if e != nil {
		return e

	} else if resp.StatusCode != 200 {
		result := &data.APIResult{}
		buf, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return e
		}

		if e := json.Unmarshal(buf, result); e != nil {
			// make sure to still report the response string even if it's malformed JSON
			result.Success = false
			result.ErrorMessage = string(buf)
		}

		return result
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
