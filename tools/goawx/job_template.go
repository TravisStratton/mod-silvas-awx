package awx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// JobTemplateService implements awx job template apis.
type JobTemplateService struct {
	client *Client
}

// ListJobTemplatesResponse represents `ListJobTemplates` endpoint response.
type ListJobTemplatesResponse struct {
	Pagination
	Results []*JobTemplate `json:"results"`
}

const jobTemplateAPIEndpoint = "/api/v2/job_templates/"

// GetJobTemplateByID shows the details of a job template.
func (jt *JobTemplateService) GetJobTemplateByID(id int, params map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)
	endpoint := fmt.Sprintf("%s%d/", jobTemplateAPIEndpoint, id)
	resp, err := jt.client.Requester.GetJSON(endpoint, result, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	// Pull in credentials related to the job template
	associatedCredentials, err := jt.ListJobTemplateCredentials(id, params)
	if err != nil {
		return nil, err
	}
	result.Credentials = append(result.Credentials, associatedCredentials...)

	sort.Ints(result.Credentials)

	return result, nil
}

// ReadJobTemplateSurveySpec gets the job template's associated survey_spec
func (jt *JobTemplateService) ReadJobTemplateSurveySpec(id int, params map[string]string) (*JobTemplateSurveySpec, error) {
	specResult := new(JobTemplateSurveySpec)
	endpoint := fmt.Sprintf("%s%d/survey_spec/", jobTemplateAPIEndpoint, id)
	resp, err := jt.client.Requester.GetJSON(endpoint, specResult, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return specResult, err
	}

	if err := CheckResponse(resp); err != nil {
		return specResult, err
	}
	/// TRAVIS START HERE - make a []map[string]string and thenloop
	//     and convert each thing to string  , perhaps fmt.Sprintf("%v",var) for each item in each list values' map

	// THIS WOKRS, SORTA --- maybe find way to have it export as jsonencode()
	// or instead of using Springf(%v) just build a string with value with quotes on either side?
	ListmapStringString := make([]map[string]any, 0)

	for i, v := range specResult.Spec {
		fmt.Print(i)
		fmt.Print(v)
		mapStringer := make(map[string]any)
		for k, val := range specResult.Spec[i] {
			mapStringer[k] = fmt.Sprintf("%v", val)
		}
		ListmapStringString = append(ListmapStringString, mapStringer)
	}
	specResult.Spec = ListmapStringString
	return specResult, nil
}

// ListJobTemplateCredentials returns a list of int ids for credentials associated
func (jt *JobTemplateService) ListJobTemplateCredentials(id int, params map[string]string) ([]int, error) {
	credentialResult := new(JobTemplateCredentials)
	endpoint := fmt.Sprintf("%s%d/credentials/", jobTemplateAPIEndpoint, id)
	credResp, err := jt.client.Requester.GetJSON(endpoint, credentialResult, params)
	if credResp != nil {
		func() {
			if err := credResp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(credResp); err != nil {
		return nil, err
	}

	credInts := make([]int, credentialResult.Count)

	for i, v := range credentialResult.Results {
		credInts[i] = v.ID
	}
	sort.Ints(credInts)
	return credInts, nil
}

// ListJobTemplates shows a list of job templates.
func (jt *JobTemplateService) ListJobTemplates(params map[string]string) ([]*JobTemplate, *ListJobTemplatesResponse, error) {
	result := new(ListJobTemplatesResponse)
	resp, err := jt.client.Requester.GetJSON(jobTemplateAPIEndpoint, result, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, result, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, result, err
	}

	return result.Results, result, nil
}

// Launch lauchs a job with the job template.
func (jt *JobTemplateService) Launch(id int, data map[string]interface{}, params map[string]string) (*JobLaunch, error) {
	result := new(JobLaunch)
	endpoint := fmt.Sprintf("%s%d/launch/", jobTemplateAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	// in case invalid job id return
	if result.Job == 0 {
		return nil, errors.New("invalid job id 0")
	}

	return result, nil
}

// CreateJobTemplate creates a job template.
func (jt *JobTemplateService) CreateJobTemplate(data map[string]interface{}, params map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)
	mandatoryFields = []string{"name", "job_type", "inventory", "project"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	//start travis code
	//var credentials schema.Set
	credentials := data["credential_ids"].(*schema.Set)

	credentialInts := make([]int, credentials.Len())
	for i, v := range credentials.List() {
		credentialInts[i] = v.(int)
	}

	// delete credentials as it's no part of the AWX api for job templates directly
	delete(data, "credential_ids")
	/// back to existing code

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	//travis := string(payload)
	//fmt.Println(travis)
	resp, err := jt.client.Requester.PostJSON(jobTemplateAPIEndpoint, bytes.NewReader(payload), result, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}
	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	// Associate credentials
	for _, v := range credentialInts {
		err := jt.AssocCredentialToTemplate(result.ID, v)
		if err != nil {
			return nil, err
		}
	}
	// add back the credentials
	result.Credentials = credentialInts
	// end my new code

	return result, nil
}

// UpdateJobTemplate updates a job template.
func (jt *JobTemplateService) UpdateJobTemplate(id int, data map[string]interface{}, params map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)

	// compare value of credential_ids list from data["credential_ids"] to what is in AWX
	// remove credential_ids from data so it doesn't end up in payload
	// then after patchJSON call below - make calls to new functions to associate/dissaciate as neecessary
	// Then put the new list as found in data abobe back into the resp after the parent PatchJSON call

	credentials := data["credential_ids"].(*schema.Set)

	credentialInts := make([]int, credentials.Len())
	for i, v := range credentials.List() {
		credentialInts[i] = v.(int)
	}

	// delete credentials as it's no part of the AWX api for job templates directly
	delete(data, "credential_ids")
	/// back to existing code
	endpoint := fmt.Sprintf("%s%d", jobTemplateAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := jt.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}
	if err := CheckResponse(resp); err != nil {
		return nil, err
	}
	// Get state from AWX and compare to desired TF state
	existingIds, err := jt.ListJobTemplateCredentials(id, params)
	if err != nil {
		return nil, err
	}
	// compare list of desired tf state of associated credentials (credentialInts) to
	//   existing AWX state (existingIds) and then assoc/dissoc credentials from the job template as needed
	if !slices.Equal(existingIds, credentialInts) {
		for _, v := range credentialInts {
			if !slices.Contains(existingIds, v) {
				err := jt.AssocCredentialToTemplate(id, v)
				if err != nil {
					return nil, err
				}
			}
		}
		for _, v := range existingIds {
			if !slices.Contains(credentialInts, v) {
				err := jt.DisassocCredentialToTemplate(id, v)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	// add back the credentials
	result.Credentials = credentialInts
	// end my new code
	return result, nil
}

// DeleteJobTemplate deletes a job template.
func (jt *JobTemplateService) DeleteJobTemplate(id int) (*JobTemplate, error) {
	result := new(JobTemplate)
	endpoint := fmt.Sprintf("%s%d", jobTemplateAPIEndpoint, id)

	resp, err := jt.client.Requester.Delete(endpoint, result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// DisAssociateCredentials remove Credentials form an awx job template.
func (jt *JobTemplateService) DisAssociateCredentials(id int, data map[string]interface{}, _ map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)
	endpoint := fmt.Sprintf("%s%d/credentials/", jobTemplateAPIEndpoint, id)
	data["disassociate"] = true
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// AssocCredentialToTemplate will post an API request to associate one credential by ID to a job template by ID
func (jt *JobTemplateService) AssocCredentialToTemplate(jtId int, credentialId int) error {

	endpoint := fmt.Sprintf("%s%d/credentials/", jobTemplateAPIEndpoint, jtId)

	data := map[string]int{
		"id": credentialId,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), nil, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("responsed with %d, resp: %v", resp.StatusCode, resp)
	}

	return nil
}

// DisassocCredentialToTemplate will post an API request to associate one credential by ID to a job template by ID
func (jt *JobTemplateService) DisassocCredentialToTemplate(jtId int, credentialId int) error {

	endpoint := fmt.Sprintf("%s%d/credentials/", jobTemplateAPIEndpoint, jtId)

	data := map[string]interface{}{
		"id":           credentialId,
		"disassociate": true,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), nil, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("responsed with %d, resp: %v", resp.StatusCode, resp)
	}

	return nil
}

// AssociateCredentials  adding credentials to JobTemplate.
func (jt *JobTemplateService) AssociateCredentials(id int, data map[string]interface{}, _ map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)

	endpoint := fmt.Sprintf("%s%d/credentials/", jobTemplateAPIEndpoint, id)
	data["associate"] = true
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// DisAssociateInstanceGroups remove instance group from an awx job template.
func (jt *JobTemplateService) DisAssociateInstanceGroups(id int, data map[string]interface{}, _ map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)
	endpoint := fmt.Sprintf("%s%d/instance_groups/", jobTemplateAPIEndpoint, id)
	data["disassociate"] = true
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// AssociateInstanceGroups  adding instance group to JobTemplate.
func (jt *JobTemplateService) AssociateInstanceGroups(id int, data map[string]interface{}, _ map[string]string) (*JobTemplate, error) {
	result := new(JobTemplate)

	endpoint := fmt.Sprintf("%s%d/instance_groups/", jobTemplateAPIEndpoint, id)
	data["associate"] = true
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}
