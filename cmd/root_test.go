// +build unit

package cmd

import (
	"bytes"
	"github.com/apache/incubator-openwhisk-wskdeploy/cmdImp"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var rootcalled bool

type RootCommand struct {
}

func (rootCommand *RootCommand) Deploy(params cmdImp.DeployParams) error {
	rootcalled = true
	return nil
}

func functionMock() {
	rootCommand := &RootCommand{}
	Deploy = rootCommand.Deploy
}

type resulter struct {
	Error   error
	Output  string
	Command *cobra.Command
}

func initializeWithRootCmd() *cobra.Command {
	rootcalled = false
	functionMock()
	return RootCmd
}

func fullSetupTest(input string) resulter {
	c := initializeWithRootCmd()
	return fullTester(c, input)
}

func fullTester(c *cobra.Command, input string) resulter {
	buf := new(bytes.Buffer)
	c.SetArgs(strings.Split(input, " "))
	err := c.Execute()
	output := buf.String()
	return resulter{err, output, c}
}

type Auth_flags struct {
	ApiHost    string
	Auth       string
	ApiVersion string
}

var expected_auth_flags Auth_flags

type Input struct {
	CfgFile        string
	Verbose        bool
	ProjectPath    string
	DeploymentPath string
	ManifestPath   string
	UseDefaults    bool
	UseInteractive bool
}

var expected_input Input

func initializeParameters() {
	rootcalled = false

	expected_auth_flags.ApiHost = "fake_api_host"
	expected_auth_flags.Auth = "fake_auth"
	expected_auth_flags.ApiVersion = "fake_api_version"

	expected_input.CfgFile = "fake_config_file"
	expected_input.Verbose = true
	expected_input.UseDefaults = true
	expected_input.UseInteractive = true
	expected_input.ProjectPath = "fake_project_path"
	expected_input.DeploymentPath = "fake_deployment_path"
	expected_input.ManifestPath = "fake_manifest_path"
}

func checkValidAuthInfo(t *testing.T, expected_auth_flags Auth_flags) {
	assert.Equal(t, expected_auth_flags.ApiHost, utils.Flags.ApiHost, "ApiHost does not match.")
	assert.Equal(t, expected_auth_flags.Auth, utils.Flags.Auth, "Auth does not match.")
	assert.Equal(t, expected_auth_flags.ApiVersion, utils.Flags.ApiVersion, "ApiVersion does not match.")
}

func checkValidInputInfo(t *testing.T, expected_input Input) {
	assert.Equal(t, expected_input.CfgFile, cmdImp.CfgFile, "CfgFile does not match.")
	assert.Equal(t, expected_input.Verbose, cmdImp.Verbose, "Verbose does not match.")
	assert.Equal(t, expected_input.UseDefaults, cmdImp.UseDefaults, "UseDefaults does not match.")
	assert.Equal(t, expected_input.UseInteractive, cmdImp.UseInteractive, "ApiHoUseInteractivest does not match.")
	assert.Equal(t, expected_input.ProjectPath, cmdImp.ProjectPath, "ProjectPath does not match.")
	assert.Equal(t, expected_input.DeploymentPath, cmdImp.DeploymentPath, "DeploymentPath does not match.")
	assert.Equal(t, expected_input.ManifestPath, cmdImp.ManifestPath, "ManifestPath does not match.")
}

func composeCommand(auth Auth_flags, input Input) string {
	cmd := ""
	if len(auth.ApiHost) > 0 {
		cmd = cmd + "--apihost " + auth.ApiHost + " "
	}
	if len(auth.Auth) > 0 {
		cmd = cmd + "--auth " + auth.Auth + " "
	}
	if len(auth.ApiVersion) > 0 {
		cmd = cmd + "--apiversion " + auth.ApiVersion + " "
	}
	if len(input.CfgFile) > 0 {
		cmd = cmd + "--config " + input.CfgFile + " "
	}
	if len(input.ProjectPath) > 0 {
		cmd = cmd + "-p " + input.ProjectPath + " "
	}
	if len(input.ManifestPath) > 0 {
		cmd = cmd + "-m " + input.ManifestPath + " "
	}
	if len(input.DeploymentPath) > 0 {
		cmd = cmd + "-d " + input.DeploymentPath + " "
	}
	if input.Verbose {
		cmd = cmd + "-v "
	}
	if input.UseDefaults {
		cmd = cmd + "-a "
	}
	if input.UseInteractive {
		cmd = cmd + "-i "
	}
	return cmd
}

func TestRootCommand(t *testing.T) {
	initializeParameters()
	command := composeCommand(expected_auth_flags, expected_input)
	output := fullSetupTest(command)

	if !rootcalled {
		t.Errorf("Root Function was not called\n out:%v", output.Error)
	}
	checkValidAuthInfo(t, expected_auth_flags)
	checkValidInputInfo(t, expected_input)
}
