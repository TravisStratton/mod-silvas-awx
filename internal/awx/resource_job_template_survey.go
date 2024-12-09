package awx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagJobTemplateSurveyTitle = "Job Template Survey"

//nolint:funlen
func resourceJobTemplateSurvey() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_job_template_survey` manages job template surveys within AWX.",
		CreateContext: resourceJobTemplateSurveyCreate,
		ReadContext:   resourceJobTemplateSurveyRead,
		UpdateContext: resourceJobTemplateSurveyUpdate,
		DeleteContext: resourceJobTemplateSurveyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type:      schema.TypeString,
						StateFunc: utils.Normalize,
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceJobTemplateSurveyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagJobTemplateSurveyTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := client.JobTemplateService.ReadJobTemplateSurveySpec(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagJobTemplateSurveyTitle, id, err)
	}

	if err := d.Set("name", res.Name); err != nil {
		return utils.DiagSet("name", id, err)
	}
	if err := d.Set("description", res.Description); err != nil {
		return utils.DiagSet("description", id, err)
	}

	//err = d.Set("spec", utils.Normalize(res.Spec))
	//if err != nil {
	//	fmt.Println("Error setting extra_vars", err)
	//}
	//temp := res.Spec
	//temp2 := d.Get("spec").([]any)
	//fmt.Print(temp2)
	//err = d.Set("spec", temp)
	//if err != nil {
	//	return utils.DiagSet("spec", id, err)
	//}

	//tempSlice := make([]map[string]any, 0)

	//tempSlice := make([]interface{}, 0)

	// tempGetter := d.Get("spec").([]interface{}) // interface {}([]interface {}) []
	// fmt.Print(tempGetter)

	//for i, v := range temp {
	//	fmt.Print(i)
	//	fmt.Print(v)
	//		temp2 = append(temp2, string(v.(string)))
	// temp is a list of maps
	// v is a map
	// for j, w := range v {
	// 	fmt.Print(j)
	// 	fmt.Print(w)
	// }
	//	}

	// 	tempGetter = append(tempGetter, temp[i])
	// 	// tempMap := make(map[string]any)
	// 	// tempMap["Max"] = v.Max
	// 	// tempMap["Min"] = v.Min
	// 	// tempMap["Type"] = v.Type
	// 	// tempMap["Choices"] = v.Choices
	// 	// tempMap["Default"] = v.Default
	// 	// tempMap["Required"] = v.Required
	// 	// tempMap["Variable"] = v.Variable
	// 	// tempMap["QuestionName"] = v.QuestionName
	// 	// tempMap["QuestionDescription"] = v.QuestionDescription

	// 	// tempGetter = append(tempSlice, tempMap)

	// }

	// //if err := d.Set("spec", temp.(*schema.Set)); err != nil {
	if err := d.Set("spec", res.Spec); err != nil {
		return utils.DiagSet("spec", id, err)
	}

	return nil
}

// type JobTemplateSurveySpec struct {
// 	Name        string       `json:"name"`
// 	Description string       `json:"description"`
// 	Spec        []SurveySpec `json:"spec"`
// }

// type SurveySpec struct {
// 	Max      int    `json:"max"`
// 	Min      int    `json:"min"`
// 	Type     string `json:"type"`
// 	Choices  any    `json:"choices"`
// 	Default  any    `json:"default"`
// 	Required bool   `json:"required"`
// 	Variable string `json:"variable"`
// 	//`json:"new_question"`: true,
// 	QuestionName        string `json:"question_name"`
// 	QuestionDescription string `json:"question_description"`
// }

func resourceJobTemplateSurveyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceJobTemplateSurveyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceJobTemplateSurveyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
