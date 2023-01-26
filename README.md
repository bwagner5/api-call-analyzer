# API Call Analyzer (ACA)

ACA is a tool to help you identify sources of API calls in your AWS account by an application or service. The tool uses cloudtrail to search for events based on predefined filters and aggregate the information for analysis. 

## Usage:

```
aca is a CLI to analyze API calls from applications utilizing CloudTrail as an API audit log

Usage:
  aca [filters] [options] [flags]

Examples:
  aca --region us-east-2 --call-source eks.amazonaws.com --event-source ec2.amazonaws.com --api DescribeInstances --user-agent='karpenter.sh'

Flags:
  -a, --api string                  API maps to EventName within CloudTrail Examples are DescribeInstances, TerminateInstances, etc
  -c, --call-source string          CallSource maps to SourceIP in CloudTrail but AWS services will include a named source IP like eks.amazonaws.com or autoscaling.amazonaws.com
  -e, --end-time string             End time for event filtering. Default: 30m ago (default "2023-01-26T14:28:12-06:00")
      --event-source string         EventSource is the top-level service where the API call is made from (i.e. ec2.amazonaws.com)
  -h, --help                        help for aca
  -i, --identity-user-name string   IdentityUserName is included in the CloudTrailEvent.userIdentity.sessionContext.sessionIssuer.userName and is useful to scope the filtering to a specific instance of an application making API calls
  -o, --output string               Output (json|chart) Default: json (default "json")
  -r, --region string               AWS Region
  -s, --start-time string           Start time for event filtering. Default: now (default "2023-01-26T13:58:12-06:00")
  -u, --user-agent string           UserAgent partial will check if the passed string is contained within the user-agent field
```


## Examples:

```
> aca --start-time 5m --call-source eks.amazonaws.com -o chart
2023/01/26 14:24:30 Filtered to 214 events out of 376
EVENT SOURCE     	API       	CALL SOURCE      	IDENTITY                                                       	USER AGENT
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
...
```

```
> aca --start-time 5m --api DescribeInstances  -o chart
2023/01/26 14:24:30 Filtered to 214 events out of 376
EVENT SOURCE     	API       	CALL SOURCE      	IDENTITY                                                       	USER AGENT
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
ec2.amazonaws.com	CreateTags	eks.amazonaws.com	eksctl-my-demo-us-east-2-clus-ServiceRole-012345678901234567	eks.amazonaws.com
...
```

```
> aca --start-time 5m --user-agent 'karpenter.sh-v0.23.0' -o chart
2023/01/26 14:27:20 Filtered to 9 events out of 405
EVENT SOURCE     	API                          	CALL SOURCE	IDENTITY                                  	USER AGENT
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeSubnets              	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeSubnets              	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
```