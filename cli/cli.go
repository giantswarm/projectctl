package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"giantswarm.io/projectctl/generated"
	"giantswarm.io/projectctl/iterator"
	"github.com/Khan/genqlient/graphql"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v2"
)

// githubTokenTransport injects the GITHUB_TOKEN header.
type githubTokenTransport struct {
	token string
	base  http.RoundTripper
}

func (t *githubTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}

// New helper: getClient fetches the token and creates the GraphQL client.
func getClient() (graphql.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	client := graphql.NewClient("https://api.github.com/graphql", &http.Client{
		Transport: &githubTokenTransport{
			token: token,
			base:  http.DefaultTransport,
		},
	})
	return client, nil
}

// Execute sets up the CLI using urfave/cli/v3 and runs it.
func Execute() error {
	// Use helper function to get the client.
	client, err := getClient()
	if err != nil {
		return err
	}

	cmd := &cli.Command{
		Name:  "projectctl",
		Usage: "A CLI tool to interact with GitHub ProjectV2 boards",
		Commands: []*cli.Command{
			{
				Name:  "projects",
				Usage: "List projects for an owner (user or organization)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "owner",
						Usage:    "Owner login (user or organization)",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "owner-type",
						Usage: "Owner type: user or organization",
						Value: "user",
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "Output format: default, table, json, yaml",
						Value: "default",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					owner := cmd.String("owner")
					outputFormat := cmd.String("output")
					var allProjects []generated.GetProjectsOrganizationProjectsV2ProjectV2ConnectionNodesProjectV2
					after := ""
					for {
						projectsResponse, err := generated.GetProjects(ctx, client, owner, 10, after)
						if err != nil {
							errMsg := fmt.Sprintf("Error: failed to list projects: %v", err)
							if strings.Contains(err.Error(), "read:project") {
								errMsg = "Error: Your GitHub token does not have the required 'read:project' scope. " +
									"Please update your token's scopes at: https://github.com/settings/tokens"
							} else if strings.Contains(err.Error(), "403") {
								errMsg = "Error: Authentication failed (403). Please check your GitHub token and its permissions."
							}
							return cli.Exit(errMsg, 1)
						}
						allProjects = append(allProjects, projectsResponse.Organization.ProjectsV2.Nodes...)
						if !projectsResponse.Organization.ProjectsV2.PageInfo.HasNextPage {
							break
						}
						after = projectsResponse.Organization.ProjectsV2.PageInfo.EndCursor
					}
					switch outputFormat {
					case "json":
						data, err := json.MarshalIndent(allProjects, "", "  ")
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to marshal json: %v", err), 1)
						}
						fmt.Println(string(data))
					case "yaml":
						data, err := yaml.Marshal(allProjects)
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to marshal yaml: %v", err), 1)
						}
						fmt.Println(string(data))
					case "table":
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
						fmt.Fprintln(w, "ID\tNumber\tTitle")
						for _, project := range allProjects {
							fmt.Fprintf(w, "%s\t%d\t%s\n", project.Id, project.Number, project.Title)
						}
						w.Flush()
					default:
						// Default output
						fmt.Println("Projects:")
						for _, project := range allProjects {
							fmt.Printf("%s %d %s\n", project.Id, project.Number, project.Title)
						}
					}
					return nil
				},
			},
			{
				Name:  "items",
				Usage: "List items for a given project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-id",
						Usage:    "ProjectV2 ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "Output format: default, table, json, yaml",
						Value: "default",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					projectID := cmd.String("project-id")
					outputFormat := cmd.String("output")
					var allItems []iterator.ItemSummary
					after := ""
					for {
						// Note: added 'after' parameter for pagination
						itemsResponse, err := generated.GetProjectItems(ctx, client, projectID, 10, after)
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to list items: %v", err), 1)
						}
						if project, ok := itemsResponse.Node.(*generated.GetProjectItemsNodeProjectV2); ok {
							pageItems := iterator.IterateItems(project)
							allItems = append(allItems, pageItems...)
							if !project.GetItems().PageInfo.HasNextPage {
								break
							}
							after = project.GetItems().PageInfo.EndCursor
						} else {
							break
						}
					}
					switch outputFormat {
					case "json":
						data, err := json.MarshalIndent(allItems, "", "  ")
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to marshal json: %v", err), 1)
						}
						fmt.Println(string(data))
					case "yaml":
						data, err := yaml.Marshal(allItems)
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to marshal yaml: %v", err), 1)
						}
						fmt.Println(string(data))
					case "table":
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
						fmt.Fprintln(w, "ID\tTitle\tDetails")
						for _, item := range allItems {
							fmt.Fprintf(w, "%s\t%s\t%s\n", item.ID, item.Title, strings.Join(item.FieldValues, " | "))
						}
						w.Flush()
					default:
						for _, item := range allItems {
							fmt.Printf("ID: %s, Title: %s\n", item.ID, item.Title)
							for _, detail := range item.FieldValues {
								fmt.Println(detail)
							}
							fmt.Println()
						}
					}
					return nil
				},
			},
			{
				Name:  "fields",
				Usage: "List all fields for a given project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-id",
						Usage:    "ProjectV2 ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "Output format: default, table, json, yaml",
						Value: "default",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					projectID := cmd.String("project-id")
					outputFormat := cmd.String("output")
					var allFields []iterator.FieldSummary
					after := ""
					for {
						// Fetch a page of fields: assuming GetProjectFields accepts first & after parameters.
						fieldsResponse, err := generated.GetProjectFields(ctx, client, projectID, 10, after)
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to get fields: %v", err), 1)
						}
						if project, ok := fieldsResponse.GetNode().(*generated.GetProjectFieldsNodeProjectV2); ok {
							pageFields := iterator.IterateFields(project)
							allFields = append(allFields, pageFields...)
							if !project.Fields.PageInfo.HasNextPage {
								break
							}
							after = project.Fields.PageInfo.EndCursor
						} else {
							return cli.Exit("unexpected node type, expected ProjectV2", 1)
						}
					}
					switch outputFormat {
					case "json":
						data, err := json.MarshalIndent(allFields, "", "  ")
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to marshal json: %v", err), 1)
						}
						fmt.Println(string(data))
					case "yaml":
						data, err := yaml.Marshal(allFields)
						if err != nil {
							return cli.Exit(fmt.Sprintf("failed to marshal yaml: %v", err), 1)
						}
						fmt.Println(string(data))
					case "table":
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
						fmt.Fprintln(w, "Kind\tID\tName\tDataType\tIterationConfigs\tOptions")
						for _, field := range allFields {
							var iterConfigs []string
							for _, cfg := range field.IterationConfigs {
								iterConfigs = append(iterConfigs, fmt.Sprintf("%s (%s)", cfg.Title, cfg.StartDate))
							}
							var opts []string
							for _, opt := range field.Options {
								opts = append(opts, opt.Name)
							}
							fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", field.Kind, field.ID, field.Name, field.DataType, strings.Join(iterConfigs, " | "), strings.Join(opts, " | "))
						}
						w.Flush()
					default:
						fmt.Println("Fields:")
						for _, field := range allFields {
							fmt.Printf("• %s\n   ID: %s\n   Name: %s\n   DataType: %s\n", field.Kind, field.ID, field.Name, field.DataType)
							if len(field.IterationConfigs) > 0 {
								fmt.Println("   Iteration Configs:")
								for _, cfg := range field.IterationConfigs {
									fmt.Printf("     - %s (ID: %s, Start: %s, Duration: %v)\n", cfg.Title, cfg.ID, cfg.StartDate, cfg.Duration)
								}
							}
							if len(field.Options) > 0 {
								fmt.Println("   Options:")
								for _, opt := range field.Options {
									fmt.Printf("     - %s (ID: %s, Color: %s, Description: %s)\n", opt.Name, opt.ID, opt.Color, opt.Description)
								}
							}
							fmt.Println()
						}
					}
					return nil
				},
			},
		},
	}
	return cmd.Run(context.Background(), os.Args)
}
