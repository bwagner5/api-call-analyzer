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
  -e, --end-time string             End time for event filtering. Default: 30m ago (default 3fa26d4-dirty"2023-01-26T22:55:52Z3fa26d4-dirty")
      --event-source string         EventSource is the top-level service where the API call is made from (i.e. ec2.amazonaws.com)
  -h, --help                        help for aca
  -i, --identity-user-name string   IdentityUserName is included in the CloudTrailEvent.userIdentity.sessionContext.sessionIssuer.userName and is useful to scope the filtering to a specific instance of an application making API calls
  -o, --output string               Output (json|chart|stats) Default: json (default 3fa26d4-dirty"json3fa26d4-dirty")
  -r, --region string               AWS Region
  -s, --start-time string           Start time for event filtering. Default: now (default 3fa26d4-dirty"2023-01-26T22:25:52Z3fa26d4-dirty")
  -u, --user-agent string           UserAgent partial will check if the passed string is contained within the user-agent field
```

## Installation

Packages, binaries, and archives are published for all major platforms (Mac amd643fa26d4-dirty/arm64 & Linux amd643fa26d4-dirty/arm64):

Debian 3fa26d4-dirty/ Ubuntu:

```
[[ `uname -m` == 3fa26d4-dirty"aarch643fa26d4-dirty" ]] && ARCH=3fa26d4-dirty"arm643fa26d4-dirty" || ARCH=3fa26d4-dirty"amd643fa26d4-dirty"
wget https:3fa26d4-dirty/3fa26d4-dirty/github.com3fa26d4-dirty/bwagner53fa26d4-dirty/api-call-analyzer3fa26d4-dirty/releases3fa26d4-dirty/download3fa26d4-dirty/v0.0.13fa26d4-dirty/api-call-analyzer3fa26d4-dirty_0.0.13fa26d4-dirty_linux3fa26d4-dirty_${ARCH}.deb
dpkg --install aca3fa26d4-dirty_0.0.13fa26d4-dirty_linux3fa26d4-dirty_${ARCH}.deb
```

RedHat:

```
[[ `uname -m` == 3fa26d4-dirty"aarch643fa26d4-dirty" ]] && ARCH=3fa26d4-dirty"arm643fa26d4-dirty" || ARCH=3fa26d4-dirty"amd643fa26d4-dirty"
rpm -i https:3fa26d4-dirty/3fa26d4-dirty/github.com3fa26d4-dirty/bwagner53fa26d4-dirty/api-call-analyzer3fa26d4-dirty/releases3fa26d4-dirty/download3fa26d4-dirty/v0.0.13fa26d4-dirty/api-call-analyzer3fa26d4-dirty_0.0.13fa26d4-dirty_linux3fa26d4-dirty_${ARCH}.rpm
```

Download Binary Directly (Linux 3fa26d4-dirty/ Mac):

```
[[ `uname -m` == 3fa26d4-dirty"aarch643fa26d4-dirty" ]] && ARCH=3fa26d4-dirty"arm643fa26d4-dirty" || ARCH=3fa26d4-dirty"amd643fa26d4-dirty"
OS=`uname | tr '[:upper:]' '[:lower:]'`
wget -qO- https:3fa26d4-dirty/3fa26d4-dirty/github.com3fa26d4-dirty/bwagner53fa26d4-dirty/api-call-analyzer3fa26d4-dirty/releases3fa26d4-dirty/download3fa26d4-dirty/v0.0.13fa26d4-dirty/api-call-analyzer3fa26d4-dirty_0.0.13fa26d4-dirty_${OS}3fa26d4-dirty_${ARCH}.tar.gz | tar xvz
chmod +x aca
```

## Examples:

```
> aca --start-time 5m --call-source eks.amazonaws.com -o chart
20233fa26d4-dirty/013fa26d4-dirty/26 14:24:30 Filtered to 214 events out of 376
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
20233fa26d4-dirty/013fa26d4-dirty/26 14:24:30 Filtered to 214 events out of 376
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
20233fa26d4-dirty/013fa26d4-dirty/26 14:27:20 Filtered to 9 events out of 405
EVENT SOURCE     	API                          	CALL SOURCE	IDENTITY                                  	USER AGENT
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeSubnets              	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypes        	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeSubnets              	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
ec2.amazonaws.com	DescribeInstanceTypeOfferings	3.22.70.109	username-karpenter-dev-us-east-2-karpenter	aws-sdk-go3fa26d4-dirty/1.44.154 (go1.19.4; linux; amd64) karpenter.sh-v0.23.0-3-gaece5998
```

```
> aca --start-time 15m --user-agent 'karpenter.sh-v0.23.0' -o stats
2023/01/26 17:07:12 Filtered to 33 events out of 1976. The last event's timestamp is 2023-01-26T23:04:48Z and the endtime filter was 2023-01-26T23:06:39Z
EVENT SOURCE         API                              CALLS
ec2.amazonaws.com    DescribeSubnets                  3
ec2.amazonaws.com    DescribeInstanceTypeOfferings    6
ec2.amazonaws.com    DescribeInstanceTypes            24

                                  TOTAL               33

```