// Copyright © 2017 Johnny Morrice <john@functorama.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/johnny-morrice/godon"
)

// timelineStreamCmd represents the timeline_get command
var timelineStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Display a timeline",
	Run: func(cmd *cobra.Command, args []string) {
		loadViperConfig(timelineStreamCmdParams)

		g := makeClient(timelineStreamCmdParams)
		defer saveConfig(g, timelineStreamCmdParams)

		err := g.Login(userFindsToken)

		if err != nil {
			die(err)
		}

		drawTimeline(cmd, timelineStreamCmdParams, g)
	},
}

func drawTimeline(cmd *cobra.Command, params *Parameters, client *godon.Client) {
	const timelineName = "SingleTimeline"

	params = params.Merge(timelineCmdParams)
	timeline := *params.String(__TIMELINE_FLAG)
	limit := *params.Int(__TIMELINE_GET_LIMIT_FLAG)
	local := *params.Bool(__TIMELINE_GET_LOCAL_FLAG)

	widget := &godon.TimelineWidget{
		Limit:    limit,
		Local:    local,
		Timeline: timeline,
	}

	publicTimelineResource := godon.Resource{
		Name:     timelineName,
		Bufferer: widget,
	}

	term := &godon.Terminal{
		Client: client,
	}

	err := term.AddResource(publicTimelineResource)

	if err != nil {
		die(err)
	}

	showTimeline := godon.Command{
		OpCode:     godon.ShowWidget,
		Parameters: []string{timelineName},
	}

	stopper := term.Render()

	err = term.ExecuteCommand(showTimeline)

	if err != nil {
		die(err)
	}

	err = stopper.Wait()

	if err != nil {
		die(err)
	}
}

var timelineStreamCmdParams *Parameters = &Parameters{}

func init() {
	timelineCmd.AddCommand(timelineStreamCmd)

	addTimelineStreamParams(timelineCmd, timelineStreamCmdParams)
}

func addTimelineStreamParams(cmd *cobra.Command, params *Parameters) {
	cmd.Flags().IntVar(params.Int(__TIMELINE_GET_LIMIT_FLAG), __TIMELINE_GET_LIMIT_FLAG, __DEFAULT_TIMELINE_GET_LIMIT, __TIMELINE_GET_LIMIT_FLAG_USAGE)
	cmd.Flags().BoolVar(params.Bool(__TIMELINE_GET_LOCAL_FLAG), __TIMELINE_GET_LOCAL_FLAG, __DEFAULT_TIMELINE_GET_LOCAL, __TIMELINE_GET_LOCAL_FLAG_USAGE)
}

const __TIMELINE_GET_LOCAL_FLAG = "local"
const __TIMELINE_GET_LOCAL_FLAG_USAGE = "Get public timelines from local instance only"
const __DEFAULT_TIMELINE_GET_LOCAL = false

const __TIMELINE_GET_LIMIT_FLAG = "limit"
const __TIMELINE_GET_LIMIT_FLAG_USAGE = "Number of entries or '-1' for all"
const __DEFAULT_TIMELINE_GET_LIMIT = 20
