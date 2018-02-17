package helpers

import (
	"github.com/solo-io/gloo-plugins/service"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
)

func NewTestConfig() *v1.Config {
	upstreams := []*v1.Upstream{
		NewTestUpstream1(),
		NewTestUpstream2(),
	}
	virtualhosts := []*v1.VirtualHost{
		NewTestVirtualHost("my-vhost", NewTestRoute1(), NewTestRoute2()),
		NewTestVirtualHost("my-vhost-2", NewTestRoute1(), NewTestRoute2()),
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func NewTestUpstream1() *v1.Upstream {
	usSpec, _ := protoutil.MarshalStruct(map[string]interface{}{
		"region":         "us-east-1",
		"secret_key_ref": "my-secret-key",
		"access_key_ref": "my-access-key",
	})
	fnSpec, _ := protoutil.MarshalStruct(map[string]interface{}{
		"key": "value",
	})
	return &v1.Upstream{
		Name: "aws",
		Type: "lambda",
		Spec: usSpec,
		Functions: []*v1.Function{
			{
				Name: "my_lambda_function",
				Spec: fnSpec,
			},
		},
	}
}
func NewTestUpstream2() *v1.Upstream {
	return &v1.Upstream{
		Name: "localhost-python",
		Type: service.UpstreamTypeService,
		Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
			Hosts: []service.Host{
				{Addr: "localhost", Port: 8000},
			},
		}),
	}
}

func NewTestVirtualHost(name string, routes ...*v1.Route) *v1.VirtualHost {
	return &v1.VirtualHost{
		Name:   name,
		Routes: routes,
	}
}

func NewTestRoute1() *v1.Route {
	extensions, _ := protoutil.MarshalStruct(map[string]interface{}{
		"auth": map[string]interface{}{
			"credentials": struct {
				Username, Password string
			}{
				Username: "alice",
				Password: "bob",
			},
			"token": "my-12345",
		}})
	return &v1.Route{
		Matcher: &v1.Matcher{
			Path: &v1.Matcher_PathPrefix{
				PathPrefix: "/foo",
			},
			Headers: map[string]string{"x-foo-bar": ""},
			Verbs:   []string{"GET", "POST"},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Extensions: extensions,
	}
}

func NewTestRoute2() *v1.Route {
	extensions, _ := protoutil.MarshalStruct(map[string]interface{}{
		"auth": map[string]interface{}{
			"credentials": struct {
				Username, Password string
			}{
				Username: "alice",
				Password: "bob",
			},
			"token": "my-12345",
		}})
	return &v1.Route{
		Matcher: &v1.Matcher{
			Path: &v1.Matcher_PathExact{
				PathExact: "/bar",
			},
			Verbs: []string{"GET", "POST"},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Upstream{
				Upstream: &v1.UpstreamDestination{
					Name: "my-upstream",
				},
			},
		},
		Extensions: extensions,
	}
}
