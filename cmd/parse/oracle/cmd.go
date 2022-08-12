package oracle

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewOracleCmd returns the Cobra command that allows to refresh the things related to the x/oracle module
func NewOracleCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle",
		Short: "Refresh things related to the x/oracle module",
	}

	//TODO: REQ REPORTS
	cmd.AddCommand(
		//requestCmd(parseConfig),
		//requestsCmd(parseConfig),
		//dataSourceCmd(parseConfig),
		//dataSourcesCmd(parseConfig),
		//oracleScriptCmd(parseConfig),
		//oracleScriptsCmd(parseConfig),
		//-_-_-_--_-_--__---__--_--_--_--
		alternativeRequestCmd(parseConfig),
		alternativeRequestsCmd(parseConfig),
	)

	return cmd
}
