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
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/johnny-morrice/godon"
	"github.com/spf13/viper"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register godon with your server",
	Run: func(cmd *cobra.Command, args []string) {
		loadViperConfig(registerCmdParams)

		godon := makeGodon(registerCmdParams)
		defer saveConfig(godon, registerCmdParams)

		err := godon.Register()

		if err != nil {
			die(err)
		}
	},
}

func loadViperConfig(params *Parameters) {
	stringFlags := []string{
		__SERVER_FLAG,
		__CLIENT_NAME_FLAG,
		__CLIENT_ID_FLAG,
		__CLIENT_SECRET_FLAG,
	}

	for _, flag := range stringFlags {
		loadViperString(flag, params)
	}
}

func loadViperString(flag string, params *Parameters) {
	flagP := params.String(flag)
	if *flagP == "" {
		*flagP = viper.GetString(flag)
	}
}

func makeGodon(params *Parameters) *godon.Godon {
	appConfig := &mastodon.AppConfig{
		ClientName: *params.String(__CLIENT_NAME_FLAG),
		Server:     *params.String(__SERVER_FLAG),
		Scopes:     __SCOPES,
	}

	app := &mastodon.Application{
		ClientID:     *params.String(__CLIENT_ID_FLAG),
		ClientSecret: *params.String(__CLIENT_SECRET_FLAG),
	}

	options := godon.Options{
		AppConfig: appConfig,
		App:       app,
	}

	godon, err := godon.New(options)

	if err != nil {
		die(err)
	}

	return godon
}

func saveConfig(godon *godon.Godon, params *Parameters) {
	configFile := createConfigFile(params)
	defer configFile.Close()
	err := writeJsonConfig(godon, configFile)

	if err != nil {
		die(err)
	}
}

func createConfigFile(params *Parameters) io.WriteCloser {
	file, err := os.Create(cfgFile)

	if err != nil {
		die(err)
	}

	return file
}

func die(err error) {
	log.Fatalf("Uncaught error: %s", err.Error())
}

func writeJsonConfig(godon *godon.Godon, w io.Writer) error {
	contents := map[string]interface{}{}
	contents[__SERVER_FLAG] = godon.AppConfig.Server
	contents[__CLIENT_NAME_FLAG] = godon.AppConfig.ClientName
	contents[__CLIENT_ID_FLAG] = godon.App.ClientID
	contents[__CLIENT_SECRET_FLAG] = godon.App.ClientSecret

	bs, err := json.MarshalIndent(contents, "", "  ")

	if err != nil {
		return errors.Wrap(err, "JSON encoding failed")
	}

	buffer := bufio.NewWriter(w)
	_, err = buffer.Write(bs)

	if err != nil {
		return errors.Wrap(err, "JSON buffer write failed")
	}

	err = buffer.Flush()

	if err != nil {
		return errors.Wrap(err, "JSON buffer flush failed")
	}

	return nil
}

var registerCmdParams *Parameters = &Parameters{}

func init() {
	appCmd.AddCommand(registerCmd)

	addRegisterParams(registerCmd, registerCmdParams)
}

func addRegisterParams(cmd *cobra.Command, params *Parameters) {
	cmd.Flags().StringVar(params.String(__SERVER_FLAG), __SERVER_FLAG, __DEFAULT_SERVER, __SERVER_FLAG_USAGE)
	cmd.Flags().StringVar(params.String(__CLIENT_NAME_FLAG), __CLIENT_NAME_FLAG, __DEFAULT_CLIENT_NAME, __CLIENT_FLAG_USAGE)
}

const __SCOPES = "read write follow"

const __SERVER_FLAG = "server"
const __SERVER_FLAG_USAGE = "Mastodon server"
const __DEFAULT_SERVER = ""
const __CLIENT_NAME_FLAG = "client-name"
const __CLIENT_FLAG_USAGE = "The name of this app"
const __DEFAULT_CLIENT_NAME = "godon"
