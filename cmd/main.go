/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Options struct {
	// CallSource maps to SourceIP in CloudTrail
	// but AWS services will include a named source IP like eks.amazonaws.com or autoscaling.amazonaws.com
	CallSource string
	// EventSource is the top-level service where the API call is made from (i.e. ec2.amazonaws.com)
	EventSource string
	// API maps to EventName within CloudTrail
	// Examples are DescribeInstances, TerminateInstances, etc
	API string
	// IdentityUserName is included in the CloudTrailEvent.userIdentity.sessionContext.sessionIssuer.userName and is useful
	// to scope the filtering to a specific instance of an application making API calls
	IdentityUserName string
	// StartTime to filter events
	StartTime time.Time
	startTime string
	// EndTime to filter events
	EndTime time.Time
	endTime string
	// UserAgent partial will check if the passed string is contained within the user-agent field
	UserAgent string
	// AWS Region to use to contact the cloudtrail API
	Region string
	// Version
	Version bool
	// Output json or chart
	Output string
}

var (
	OutputJSON  = "json"
	OutputChart = "chart"
	OutputStats = "stats"
)

type Stat struct {
	EventSource string `json:"eventSource"`
	API         string `json:"api"`
	Calls       int    `json:"calls"`
}

type CloudTrailEvent struct {
	EventVersion       string       `json:"eventVersion"`
	UserIdentity       UserIdentity `json:"userIdentity"`
	EventTime          string       `json:"eventTime"`
	EventSource        string       `json:"eventSource"`
	EventName          string       `json:"eventName"`
	AWSRegion          string       `json:"awsRegion"`
	SourceIPAddress    string       `json:"sourceIPAddress"`
	UserAgent          string       `json:"userAgent"`
	RequestParameters  any          `json:"requestParameters"`
	ResponseElements   any          `json:"responseElements"`
	RequestID          string       `json:"requestID"`
	EventID            string       `json:"eventID"`
	ReadOnly           bool         `json:"readOnly"`
	EventType          string       `json:"eventType"`
	ManagementEvent    bool         `json:"managementEvent"`
	RecipientAccountID string       `json:"recipientAccountId"`
	EventCategory      string       `json:"eventCategory"`
}

type UserIdentity struct {
	Type           string         `json:"type"`
	PrincipalID    string         `json:"principalId"`
	ARN            string         `json:"arn"`
	AccountID      string         `json:"accountId"`
	SessionContext SessionContext `json:"sessionContext"`
	InvokedBy      string         `json:"invokedBy"`
}

type SessionContext struct {
	SessionIssuer SessionIssuer `json:"sessionIssuer"`
	Attributes    Attributes    `json:"attributes"`
}

type Attributes struct {
	CreationDate     string `json:"creationDate"`
	MFAAuthenticated string `json:"mfaAuthenticated"`
}

type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	ARN         string `json:"arn"`
	AccountID   string `json:"accountId"`
	UserName    string `json:"userName"`
}

var (
	opts    = &Options{}
	version = ""
	commit  = ""
)

var rootCmd = &cobra.Command{
	Use:     "aca [filters] [options]",
	Example: `  aca --region us-east-2 --call-source eks.amazonaws.com --event-source ec2.amazonaws.com --api DescribeInstances --user-agent='karpenter.sh'`,
	Short:   "aca is a CLI to analyze API calls from applications utilizing CloudTrail as an API audit log",
	Run: func(cmd *cobra.Command, args []string) {
		if opts.Version {
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			return
		}
		total, events, err := filterEvents(cmd.Context(), processOpts(opts))
		if err != nil {
			log.Fatalln(err.Error())
		}
		if len(events) == 0 {
			log.Printf("All %d events did not match your filters\n", total)
			os.Exit(1)
		}

		log.Printf("Filtered to %d events out of %d. The last event's timestamp is %s and the endtime filter was %s\n",
			len(events), total, events[len(events)-1].EventTime, opts.EndTime.UTC().Format(time.RFC3339))
		if opts.Output == OutputChart {
			outputChart(events)
		} else if opts.Output == OutputStats {
			stats := computeStats(events)
			outputStatsChart(stats)
		} else {
			outputJSON(events)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&opts.Region, "region", "r", "", "AWS Region")
	rootCmd.PersistentFlags().StringVarP(&opts.Output, "output", "o", "json", "Output (json|chart|stats) Default: json")

	rootCmd.PersistentFlags().StringVarP(&opts.CallSource, "call-source", "c", "", "CallSource maps to SourceIP in CloudTrail but AWS services will include a named source IP like eks.amazonaws.com or autoscaling.amazonaws.com")
	rootCmd.PersistentFlags().StringVar(&opts.EventSource, "event-source", "", "EventSource is the top-level service where the API call is made from (i.e. ec2.amazonaws.com)")
	rootCmd.PersistentFlags().StringVarP(&opts.API, "api", "a", "", "API maps to EventName within CloudTrail Examples are DescribeInstances, TerminateInstances, etc")
	rootCmd.PersistentFlags().StringVarP(&opts.IdentityUserName, "identity-user-name", "i", "", "IdentityUserName is included in the CloudTrailEvent.userIdentity.sessionContext.sessionIssuer.userName and is useful to scope the filtering to a specific instance of an application making API calls")
	rootCmd.PersistentFlags().StringVarP(&opts.startTime, "start-time", "s", time.Now().Add(-30*time.Minute).UTC().Format(time.RFC3339), "Start time for event filtering. Default: now")
	rootCmd.PersistentFlags().StringVarP(&opts.endTime, "end-time", "e", time.Now().UTC().Format(time.RFC3339), "End time for event filtering. Default: 30m ago")
	rootCmd.PersistentFlags().StringVarP(&opts.UserAgent, "user-agent", "u", "", "UserAgent partial will check if the passed string is contained within the user-agent field")

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processOpts(opts *Options) *Options {
	start, err := parseTimeOrDuration(opts.startTime)
	if err != nil {
		log.Fatalln(err)
	}
	opts.StartTime = start
	end, err := parseTimeOrDuration(opts.endTime)
	if err != nil {
		log.Fatalln(err)
	}
	opts.EndTime = end
	return opts
}

func parseTimeOrDuration(input string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		dur, derr := time.ParseDuration(input)
		if derr != nil {
			return time.Time{}, fmt.Errorf("unable to parse as RFC3339 time or duration. Time Err: %w Duration Err: %w", err, derr)
		}
		return time.Now().Add(-1 * dur), nil
	}
	return t, nil
}

func filterEvents(ctx context.Context, opts *Options) (int, []*CloudTrailEvent, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return 0, nil, err
	}

	// Create an Amazon CloudTrail service client
	ct := cloudtrail.NewFromConfig(cfg)

	lookupInput := &cloudtrail.LookupEventsInput{
		StartTime: &opts.StartTime,
		EndTime:   &opts.EndTime,
	}
	if opts.API != "" {
		lookupInput.LookupAttributes = append(lookupInput.LookupAttributes, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyEventName,
			AttributeValue: &opts.API,
		})
	} else if opts.EventSource != "" {
		lookupInput.LookupAttributes = append(lookupInput.LookupAttributes, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyEventSource,
			AttributeValue: &opts.EventSource,
		})
	}
	var filteredEvents []*CloudTrailEvent
	rawEvents := 0

	eventPaginator := cloudtrail.NewLookupEventsPaginator(ct, lookupInput)
	for eventPaginator.HasMorePages() {
		output, err := eventPaginator.NextPage(ctx)
		if err != nil {
			return rawEvents, filteredEvents, err
		}
		rawEvents += len(output.Events)
		for _, event := range output.Events {
			var ctEvent *CloudTrailEvent
			if err := json.Unmarshal([]byte(*event.CloudTrailEvent), &ctEvent); err != nil {
				return rawEvents, filteredEvents, fmt.Errorf("unable to unmarshal cloudtrail event: %w", err)
			}
			if opts.EventSource != "" && *event.EventSource != opts.EventSource {
				continue
			}
			if opts.CallSource != "" && ctEvent.SourceIPAddress != opts.CallSource {
				continue
			}
			if opts.API != "" && ctEvent.EventName != opts.API {
				continue
			}
			if opts.IdentityUserName != "" && ctEvent.UserIdentity.SessionContext.SessionIssuer.UserName != opts.IdentityUserName {
				continue
			}
			if opts.UserAgent != "" && !strings.Contains(ctEvent.UserAgent, opts.UserAgent) {
				continue
			}
			filteredEvents = append(filteredEvents, ctEvent)
		}
	}
	sort.Slice(filteredEvents, func(i, j int) bool {
		ti, err := time.Parse(time.RFC3339, filteredEvents[i].EventTime)
		if err != nil {
			return true
		}
		tj, err := time.Parse(time.RFC3339, filteredEvents[j].EventTime)
		if err != nil {
			return true
		}
		return ti.Before(tj)
	})
	return rawEvents, filteredEvents, nil
}

func outputJSON(events []*CloudTrailEvent) {
	log.Printf("Found %d events\n", len(events))
	eventsJSON, err := json.MarshalIndent(events, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(eventsJSON))
}

func outputChart(events []*CloudTrailEvent) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Start Time", "Event Source", "API", "Call Source", "Identity", "User Agent"})
	data := [][]string{}
	for _, event := range events {
		data = append(data, []string{
			event.EventTime,
			event.EventSource,
			event.EventName,
			event.SourceIPAddress,
			event.UserIdentity.SessionContext.SessionIssuer.UserName,
			event.UserAgent,
		})
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("    ") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

func computeStats(events []*CloudTrailEvent) []*Stat {
	stats := map[string]*Stat{}
	var statsList []*Stat
	for _, event := range events {
		key := fmt.Sprintf("%s:%s", event.EventSource, event.EventName)
		stat, ok := stats[key]
		if !ok {
			stat = &Stat{EventSource: event.EventSource, API: event.EventName}
			stats[key] = stat
			statsList = append(statsList, stat)
		}
		stat.Calls++
	}
	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].Calls < statsList[j].Calls
	})
	return statsList
}

func outputStatsChart(stats []*Stat) {
	table := tablewriter.NewWriter(os.Stdout)
	totalCalls := 0
	table.SetHeader([]string{"Event Source", "API", "Calls"})
	data := [][]string{}
	for _, stat := range stats {
		data = append(data, []string{
			stat.EventSource,
			stat.API,
			strconv.Itoa(stat.Calls),
		})
		totalCalls += stat.Calls
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("    ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.SetFooter([]string{"", "TOTAL", strconv.Itoa(totalCalls)})
	table.Render()
}
