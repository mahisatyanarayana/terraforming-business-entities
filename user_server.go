package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func userServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"middlename": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nauticalcustomerid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"useraccountid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	email := d.Get("email").(string)
	firstname := d.Get("firstname").(string)
	middlename := d.Get("middlename").(string)
	lastname := d.Get("lastname").(string)
	nauticalcustomerid := d.Get("nauticalcustomerid").(string)

	createUser(email, firstname, middlename, lastname, nauticalcustomerid)

	d.SetId(email)
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {

	url := fmt.Sprintf(`https://users-dev.dev.appriver.corp/api/v1/users/%s/relationships/profiledata`, d.Id())
	method := "GET"
	payload := strings.NewReader("")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", `application/x-www-form-urlencoded`)
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var data Data
	json.Unmarshal(body, &data)
	fmt.Println(data.Data.Attributes.UserId)

	d.Set("email", data.Data.Attributes.UserId)
	d.Set("firstname", data.Data.Attributes.FirstName)
	d.Set("lastname", data.Data.Attributes.LastName)
	d.Set("nauticalcustomerid", data.Data.Attributes.NauticalCustomerId)
	d.Set("useraccountid", data.Data.Attributes.UserAccountId)

	return err
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {

	email := d.Get("email").(string)
	firstname := d.Get("firstname").(string)
	middlename := d.Get("middlename").(string)
	lastname := d.Get("lastname").(string)
	nauticalcustomerid := d.Get("nauticalcustomerid").(string)

	updateUser(email, firstname, middlename, lastname, nauticalcustomerid)
	d.SetId(email)
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {

	email := d.Id()
	deleteUser(email)

	d.SetId("")
	return nil
}

func createUser(email, firstname, middlename, lastname, nauticalcustomerid string) {
	url := "https://users-dev.dev.appriver.corp/AdministrationWS.svc"
	method := "POST"
	formattedbody := fmt.Sprintf(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
	<s:Body>
	   <AsyncCreateUser xmlns="http://appriver.com/users/administration">
		   <request xmlns:a="http://appriver.com/users/administration" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
				<CorrelationId i:nil="true" xmlns="http://appriver.com/users"/>
				<IpAddress i:nil="true" xmlns="http://appriver.com/users"/>
				<Requestor xmlns="http://appriver.com/users">terraformplugin@appriver.com</Requestor>
  
				<Email>%s</Email>
				<FirstName>%s</FirstName>
			  
			   <LastName>%s</LastName>
				<Middle>%s</Middle>
			   <NauticalCustomerId>%s</NauticalCustomerId>
				<Password>pa55word</Password>
				<Roles xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
					<b:string>CUSTOMER_SUPERADMIN</b:string>
				</Roles>

			   <Username>%s</Username>
			</request>
		</AsyncCreateUser>
	</s:Body>
</s:Envelope>`, email, firstname, lastname, middlename, nauticalcustomerid, email)

	payload := strings.NewReader(formattedbody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", `http://appriver.com/users/administration/IAdministrationService/AsyncCreateUser`)

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	return

}

func updateUser(email, firstname, middlename, lastname, nauticalcustomerid string) {
	url := "https://users-dev.dev.appriver.corp/AdministrationWS.svc"
	method := "POST"
	formattedbody := fmt.Sprintf(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
    <s:Body>
        <UpdatePersonalInformationV2 xmlns="http://appriver.com/users/administration">
            <request xmlns:a="http://appriver.com/users/administration" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
                <CorrelationId i:nil="true" xmlns="http://appriver.com/users"/>
                <IpAddress i:nil="true" xmlns="http://appriver.com/users"/>
                <Requestor xmlns="http://appriver.com/users">terraformplugin@appriver.com</Requestor>
                <UpdatedFirstName>%s</UpdatedFirstName>
				<UpdatedLastName>%s</UpdatedLastName>
				<UpdatedMiddle>%s</UpdatedMiddle>
                <UserId>%s</UserId>
            </request>
        </UpdatePersonalInformationV2>
    </s:Body>
</s:Envelope>`, firstname, lastname, middlename, email)

	payload := strings.NewReader(formattedbody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", `http://appriver.com/users/administration/IAdministrationService/UpdatePersonalInformationV2`)

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	return

}

func deleteUser(email string) {
	url := "https://users-dev.dev.appriver.corp/AdministrationWS.svc"
	method := "POST"
	formattedbody := fmt.Sprintf(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
    <s:Body>
        <DeleteUser xmlns="http://appriver.com/users/administration">
            <request xmlns:a="http://appriver.com/users/administration" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
                <CorrelationId i:nil="true" xmlns="http://appriver.com/users"/>
                <IpAddress i:nil="true" xmlns="http://appriver.com/users"/>
                <Requestor xmlns="http://appriver.com/users">terraformplugin@appriver.com</Requestor>
                <UserId>%s</UserId>
            </request>
        </DeleteUser>
    </s:Body>
</s:Envelope>`, email)

	payload := strings.NewReader(formattedbody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", `http://appriver.com/users/administration/IAdministrationService/DeleteUser`)

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	return

}

//Data ...
type Data struct {
	Data  User   `json:"data"`
	Links string `json:"links"`
}

//User ...
type User struct { //data
	Attributes Attributes `json:"attributes"`
}

//Attributes ...
type Attributes struct {
	UserId             string `json:"UserId"`
	UserAccountId      string `json:UserAccountId`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	NauticalCustomerId string `json:"nauticalCustomerId"`
}
