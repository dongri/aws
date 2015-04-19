# eb-deploy
AWS Elastic Beanstalk Deploy Tool

# Settings ENVs

```
var envs = map[string]interface{}{
	"dev": map[string]string{
		"CNAME":  "******.elasticbeanstalk.com",
		"REGION": "ap-northeast-1",
	},
	"prd": map[string]string{
		"CNAME":  "******.elasticbeanstalk.com",
		"REGION": "us-east-1",
	},
}
```

# Usage

```
$ go run deploy.go dev

or

$ go run deploy.go prd
```
