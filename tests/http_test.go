package tests

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/http"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestHttp(t *testing.T) {

	p := peer.NewPeer("http.Acceptor")
	pset := p.(cellnet.PropertySet)
	pset.SetProperty("Name", "httpserver")
	pset.SetProperty("Address", "127.0.0.1:8081")
	proc.BindProcessor(p, "http", func(raw cellnet.Event) {

		switch raw.Message().(type) {
		case *HttpEchoREQ:

			raw.(interface {
				Send(interface{})
			}).Send(&HttpEchoACK{
				Status: 0,
				Token:  "ok",
			})

		}

	})

	p.Start()

	requestor(t)

	//fmt.Scanln()
	//requestForm(t)
	//postForm(t)
}

func requestor(t *testing.T) {
	p := peer.NewPeer("http.Connector")
	pset := p.(cellnet.PropertySet)
	pset.SetProperty("Name", "httpclient")
	pset.SetProperty("Address", "127.0.0.1:8081")

	raw, err := p.(interface {
		Request(method string, raw interface{}) (interface{}, error)
	}).Request("GET", &HttpEchoREQ{
		UserName: "kitty",
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	msg := raw.(*HttpEchoACK)
	if msg.Token != "ok" {
		t.Log("unexpect token result", err)
		t.FailNow()
	}

}

func requestForm(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8081/hello?UserName=kitty")
	if err != nil {
		t.Log("http req failed", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("http response failed", err)
		t.FailNow()
	}

	var ack HttpEchoACK
	if err := json.Unmarshal(body, &ack); err != nil {
		t.Log("json unmarshal failed", err)
		t.FailNow()
	}

	if ack.Token != "ok" {
		t.Log("unexpect token result", err)
		t.FailNow()
	}
}

func postForm(t *testing.T) {
	resp, err := http.PostForm("http://127.0.0.1:8081/hello",
		url.Values{"UserName": {"kitty"}})

	if err != nil {
		t.Log("http req failed", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("http response failed", err)
		t.FailNow()
	}

	var ack HttpEchoACK
	if err := json.Unmarshal(body, &ack); err != nil {
		t.Log("json unmarshal failed", err)
		t.FailNow()
	}

	if ack.Token != "ok" {
		t.Log("unexpect token result", err)
		t.FailNow()
	}

}

type HttpEchoREQ struct {
	UserName string
}

type HttpEchoACK struct {
	Token  string
	Status int32
}

func (self *HttpEchoREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *HttpEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		URL:          "/hello",
		Method:       "GET",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		ResponseCodec: codec.MustGetCodec("json"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		URL:          "/hello",
		Method:       "POST",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		ResponseCodec: codec.MustGetCodec("json"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

}
