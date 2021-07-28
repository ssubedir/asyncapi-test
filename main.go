package main

import (
	"io/ioutil"
	"time"
	"github.com/swaggest/go-asyncapi/reflector/asyncapi-2.0.0"
	"github.com/swaggest/go-asyncapi/spec-2.0.0"
)

func main() {

	type SubItem struct {
		Key    string  `json:"key" description:"Item key"`
		Values []int64 `json:"values" uniqueItems:"true" description:"List of item values"`
	}

	type MyMessage struct {
		Name      string    `path:"name" description:"Name"`
		CreatedAt time.Time `json:"createdAt" description:"Creation time"`
		Items     []SubItem `json:"items" description:"List of items"`
	}

	type MyAnotherMessage struct {
		TraceID string  `header:"X-Trace-ID" description:"Tracing header" required:"true"`
		Item    SubItem `json:"item" description:"Some item"`
	}

	reflector := asyncapi.Reflector{
		Schema: &spec.AsyncAPI{
			Servers: map[string]spec.Server{
				"live": {
					URL:             "api.{country}.lovely.com:5672",
					Description:     "Production instance.",
					ProtocolVersion: "0.9.1",
					Protocol:        "amqp",
					Variables: map[string]spec.ServerVariable{
						"country": {
							Enum:        []string{"RU", "US", "DE", "FR"},
							Default:     "US",
							Description: "Country code.",
						},
					},
				},
			},
			Info: &spec.Info{
				Version: "1.2.3", // required
				Title:   "My Lovely Messaging API",
			},
		},
	}
	mustNotFail := func(err error) {
		if err != nil {
			panic(err.Error())
		}
	}

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: "one.{name}.two",
		BaseChannelItem: &spec.ChannelItem{
			Bindings: &spec.ChannelBindingsObject{
				Amqp: &spec.AMQP091ChannelBindingObject{
					Is: spec.AMQP091ChannelBindingObjectIsRoutingKey,
					Exchange: &spec.Exchange{
						Name: "some-exchange",
					},
				},
			},
		},
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "This is a sample schema.",
				Summary:     "Sample publisher",
			},
			MessageSample: new(MyMessage),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: "another.one",
		Subscribe: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "This is another sample schema.",
				Summary:     "Sample consumer",
			},
			MessageSample: new(MyAnotherMessage),
		},
	}))

	yaml, err := reflector.Schema.MarshalYAML()
	mustNotFail(err)

	mustNotFail(ioutil.WriteFile("asyncapi.yaml", yaml, 0644))
}
