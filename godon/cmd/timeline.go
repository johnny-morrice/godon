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
)

// timelineCmd represents the timeline command
var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Mastodon timeline control",
}

var timelineCmdParams *Parameters = &Parameters{}

func init() {
	RootCmd.AddCommand(timelineCmd)

	addTimelineParams(timelineCmd, timelineCmdParams)
}

func addTimelineParams(cmd *cobra.Command, params *Parameters) {
	cmd.PersistentFlags().StringVar(params.String(__TIMELINE_FLAG), __TIMELINE_FLAG, __DEFAULT_TIMELINE, __TIMELINE_FLAG_USAGE)
}

const __TIMELINE_FLAG = "timeline"
const __TIMELINE_FLAG_USAGE = "(public|home|:hashtag)"
const __DEFAULT_TIMELINE = "public"
