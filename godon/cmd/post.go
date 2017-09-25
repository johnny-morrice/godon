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

var statusPostCmd = &cobra.Command{
	Use:   "post",
	Short: "Update your status",
	Run: func(cmd *cobra.Command, args []string) {
		loadViperConfig(statusPostCmdParams)

		godon := makeGodon(authorizeCmdParams)
		err := godon.OauthLogin()

		if err != nil {
			die(err)
		}

		message := *statusPostCmdParams.String(__STATUS_MESSAGE_FLAG)

		err = godon.PostStatus(message)

		if err != nil {
			die(err)
		}
	},
}

var statusPostCmdParams *Parameters = &Parameters{}

func init() {
	statusCmd.AddCommand(statusPostCmd)

	statusPostCmd.Flags().StringVar(statusPostCmdParams.String(__STATUS_MESSAGE_FLAG), __STATUS_MESSAGE_FLAG, __DEFAULT_STATUS_MESSAGE, __STATUS_MESSAGE_USAGE)
}

const __STATUS_MESSAGE_FLAG = "message"
const __DEFAULT_STATUS_MESSAGE = "Test."
const __STATUS_MESSAGE_USAGE = "Status update"
