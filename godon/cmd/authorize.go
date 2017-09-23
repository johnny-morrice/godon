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
	"fmt"

	"github.com/spf13/cobra"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Fetch OAuth refresh token",

	Run: func(cmd *cobra.Command, args []string) {
		loadViperConfig(authorizeCmdParams)

		godon := makeGodon(authorizeCmdParams)

		err := godon.Authorize(userFindsToken)

		if err != nil {
			die(err)
		}

		saveConfig(godon, authorizeCmdParams)
	},
}

func userFindsToken(url string) (string, error) {
	fmt.Printf("Visit %s then enter the token:\n", url)
	var token string
	_, err := fmt.Scan(&token)
	return token, err
}

var authorizeCmdParams *Parameters = &Parameters{}

func init() {
	RootCmd.AddCommand(authorizeCmd)

	addRegisterParams(authorizeCmd, authorizeCmdParams)
	addAuthorizeParams(authorizeCmd, authorizeCmdParams)
}

func addAuthorizeParams(cmd *cobra.Command, params *Parameters) {
	cmd.Flags().StringVar(params.String(__CLIENT_ID_FLAG), __CLIENT_ID_FLAG, __DEFAULT_CLIENT_ID, __CLIENT_ID_FLAG_USAGE)
	cmd.Flags().StringVar(params.String(__CLIENT_SECRET_FLAG), __CLIENT_SECRET_FLAG, __DEFAULT_CLIENT_SECRET, __CLIENT_SECRET_FLAG_USAGE)
}

const __CLIENT_ID_FLAG = "client-id"
const __CLIENT_ID_FLAG_USAGE = "Mastadon app ID"
const __DEFAULT_CLIENT_ID = ""
const __CLIENT_SECRET_FLAG = "client-secret"
const __CLIENT_SECRET_FLAG_USAGE = "Mastadon app secret"
const __DEFAULT_CLIENT_SECRET = ""
