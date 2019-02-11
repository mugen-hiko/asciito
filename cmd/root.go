package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/MakeNowJust/heredoc"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const appName = "asciito"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [source file]", appName),
	Short: "ASCII tables in CLI",
	Long: heredoc.Doc(`
		Asciito generates an ASCII table from delimited data (CSV, TSV
		and other delimiter-separated text files) on a terminal.
	`),
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)

		records, err := sourceData()

		if err != nil {
			log.Fatal(err)
		}

		header, body := splitHeaderAndBody(records)

		if len(header) != 0 {
			table.SetHeader(header)

			table.SetAutoFormatHeaders(false)
		}

		table.AppendBulk(body)

		table.Render()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetFlags(0)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	rootCmd.Flags().StringP("delimiter", "d", ",", "field delimiter character")
	rootCmd.Flags().BoolP("skip-header", "s", false, "ignore the first line")

	viper.BindPFlag("delimiter", rootCmd.Flags().Lookup("delimiter"))
	viper.BindPFlag("skip-header", rootCmd.Flags().Lookup("skip-header"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()

		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		viper.SetConfigName("." + appName)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix(appName)

	viper.AutomaticEnv()
}

func sourceData() ([][]string, error) {
	csv := csv.NewReader(os.Stdin)

	csv.Comma = []rune(viper.GetString("delimiter"))[0]

	return csv.ReadAll()
}

func splitHeaderAndBody(records [][]string) ([]string, [][]string) {
	switch {
	case viper.GetBool("skip-header"):
		return []string{}, records
	case len(records) == 0:
		return []string{}, [][]string{}
	case len(records) == 1:
		return records[0], [][]string{}
	default:
		return records[0], records[1:]
	}
}
