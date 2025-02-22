package iterator

import (
	"fmt"
	"strings"

	"giantswarm.io/projectctl/generated"
)

type ItemSummary struct {
	ID          string   `json:"id" yaml:"id"`
	Title       string   `json:"title" yaml:"title"`
	FieldValues []string `json:"field_values" yaml:"field_values"`
}

type IterationConfig struct {
	ID        string      `json:"id" yaml:"id"`
	Title     string      `json:"title" yaml:"title"`
	StartDate string      `json:"start_date" yaml:"start_date"`
	Duration  interface{} `json:"duration" yaml:"duration"`
}

type Option struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Color       string `json:"color" yaml:"color"`
	Description string `json:"description" yaml:"description"`
}

type FieldSummary struct {
	Kind             string            `json:"kind" yaml:"kind"`
	ID               string            `json:"id" yaml:"id"`
	Name             string            `json:"name" yaml:"name"`
	DataType         string            `json:"data_type" yaml:"data_type"`
	IterationConfigs []IterationConfig `json:"iteration_configs,omitempty" yaml:"iteration_configs,omitempty"`
	Options          []Option          `json:"options,omitempty" yaml:"options,omitempty"`
}

// IterateItems returns a summary slice of project items.
func IterateItems(p *generated.GetProjectItemsNodeProjectV2) []ItemSummary {
	var items []ItemSummary
	for _, item := range p.GetItems().Nodes {
		if content, ok := item.GetContent().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemContentIssue); ok {
			var details []string
			for _, fieldValue := range item.FieldValues.Nodes {
				switch v := fieldValue.(type) {
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldDateValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldDateValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %s", field.GetName(), v.GetDate()))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldIterationValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldIterationValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %+v", field.GetName(), v))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldLabelValue:
					var labels []string
					for _, label := range v.GetLabels().Nodes {
						labels = append(labels, fmt.Sprintf("%+v", label.GetName()))
					}
					details = append(details, "Labels: "+strings.Join(labels, ", "))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldMilestoneValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldMilestoneValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %+v", field.GetName(), v))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldNumberValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldNumberValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %+v", field.GetName(), v))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldPullRequestValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldPullRequestValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %+v", field.GetName(), v))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldRepositoryValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldRepositoryValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %s", field.GetName(), v.GetRepository().Name))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldReviewerValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldReviewerValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %+v", field.GetName(), v))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldSingleSelectValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldSingleSelectValueFieldProjectV2SingleSelectField)
					details = append(details, fmt.Sprintf("%s: %s", field.GetName(), v.GetName()))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldTextValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldTextValueFieldProjectV2Field)
					details = append(details, fmt.Sprintf("%s: %s", field.GetName(), v.GetText()))
				case *generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldUserValue:
					field := v.GetField().(*generated.GetProjectItemsNodeProjectV2ItemsProjectV2ItemConnectionNodesProjectV2ItemFieldValuesProjectV2ItemFieldValueConnectionNodesProjectV2ItemFieldUserValueFieldProjectV2Field)
					var users []string
					for _, user := range v.GetUsers().Nodes {
						users = append(users, user.Login)
					}
					details = append(details, fmt.Sprintf("%s: %s", field.GetName(), strings.Join(users, ", ")))
				default:
					details = append(details, fmt.Sprintf("Unexpected type: %T", v))
				}
			}
			items = append(items, ItemSummary{
				ID:          item.Id,
				Title:       content.Title,
				FieldValues: details,
			})
		}
	}
	return items
}

// IterateFields returns a summary slice of project fields.
func IterateFields(p *generated.GetProjectFieldsNodeProjectV2) []FieldSummary {
	var fields []FieldSummary
	for _, f := range p.Fields.Nodes {
		switch field := f.(type) {
		case *generated.GetProjectFieldsNodeProjectV2FieldsProjectV2FieldConfigurationConnectionNodesProjectV2Field:
			fields = append(fields, FieldSummary{
				Kind:             "Basic Field",
				ID:               field.Id,
				Name:             field.Name,
				DataType:         string(field.DataType),
				IterationConfigs: nil,
				Options:          nil,
			})
		case *generated.GetProjectFieldsNodeProjectV2FieldsProjectV2FieldConfigurationConnectionNodesProjectV2IterationField:
			var configs []IterationConfig
			if len(field.Configuration.Iterations) > 0 {
				for _, iter := range field.Configuration.Iterations {
					configs = append(configs, IterationConfig{
						ID:        iter.Id,
						Title:     iter.Title,
						StartDate: iter.StartDate,
						Duration:  iter.Duration,
					})
				}
			}
			fields = append(fields, FieldSummary{
				Kind:             "Iteration Field",
				ID:               field.Id,
				Name:             field.Name,
				DataType:         string(field.DataType),
				IterationConfigs: configs,
				Options:          nil,
			})
		case *generated.GetProjectFieldsNodeProjectV2FieldsProjectV2FieldConfigurationConnectionNodesProjectV2SingleSelectField:
			var opts []Option
			if len(field.Options) > 0 {
				for _, opt := range field.Options {
					opts = append(opts, Option{
						ID:          opt.Id,
						Name:        opt.Name,
						Color:       string(opt.Color),
						Description: opt.Description,
					})
				}
			}
			fields = append(fields, FieldSummary{
				Kind:             "Single Select Field",
				ID:               field.Id,
				Name:             field.Name,
				DataType:         string(field.DataType),
				IterationConfigs: nil,
				Options:          opts,
			})
		default:
			fields = append(fields, FieldSummary{
				Kind:             fmt.Sprintf("Unknown field type: %T", field),
				IterationConfigs: nil,
				Options:          nil,
			})
		}
	}
	return fields
}
