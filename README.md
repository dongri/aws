# eb-deploy
AWS Elastic Beanstalk Deploy Tool

# Requirement

EB CLI 3.x

[Getting Set Up with EB Command Line Interface (CLI) 3.x](http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/eb-cli3-getting-set-up.html)

# Install (Settings ENVs)

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
