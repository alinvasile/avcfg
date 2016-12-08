# avcfg
HTTP simple configuration server that exposes properties from a json file

Considering we have this sample json configuration file called aws-integration-feature.json:

```json
{
	"enable": "true",
	"credentials":{
		"key":"ABffds7yggf",
		"secret":"hb65555"		
	},
	"location":{
		"region":"eu-west-1",
		"bucket":"my-log-bucket"
	},
	"caching":{
		"key.ttl.ms":"120000",
		"key.idle.ms":"60000"
	}
}
```

We can fetch the properties in the following way:

1. `enable: http://localhost:8080/json/aws-integration-feature/enable`
1. `key in credentials: http://localhost:8080/json/aws-integration-feature/credentials/key`
1. `key.ttl.ms in caching: http://localhost:8080/json/aws-integration-feature/caching/key.ttl.ms`

The properties are cached for 5 minutes. If any json files are added or updated, there's no need to restart the server.

General request pattern:


`http://<IP address>:<PORT>/json/<JSON filename without extension>/<property path in the json file>`


