package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func customerServer() *schema.Resource {
	return &schema.Resource{
		Create: customerServerCreate,
		Read:   customerServerRead,
		Update: customerServerUpdate,
		Delete: customerServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"customerid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func customerServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	nauticalcustomerid := d.Get("customerid").(string)

	createCustomer(name, nauticalcustomerid)

	d.SetId(nauticalcustomerid)
	return resourceServerRead(d, m)
}

func customerServerRead(d *schema.ResourceData, m interface{}) error {

	name := d.Get("name").(string)
	nauticalcustomerid := d.Get("customerid").(string)

	d.Set("name", name)
	d.Set("customerid", nauticalcustomerid)

	return nil
}

func customerServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func customerServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func createCustomer(name, nauticalcustomerid string) {
	url := "https://nautical-dev.dev.appriver.corp/nauticalws.svc"
	method := "POST"
	formattedbody := fmt.Sprintf(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing">
    <s:Header>
        <a:Action s:mustUnderstand="1">msg://appriver/nautical/core/INauticalService/ProcessCommand</a:Action>
       
        <a:To s:mustUnderstand="1">https://nautical-dev.dev.appriver.corp/NauticalWS.svc</a:To>
    </s:Header>
    <s:Body>
        <ProcessCommand xmlns="msg://appriver/nautical/core">
            <request xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
                <Command i:type="b:AddNewDirectCustomerCommand" xmlns:b="msg://appriver/nautical/customers">
                    
                    <IsSynchronized>false</IsSynchronized>
                    <OnBehalfOf i:nil="true"/>
                    <Requestor>terraformplugin@appriver.com</Requestor>
                    <b:AccountManager>
                        <b:AccountManager>msatya@appriver.com</b:AccountManager>
                    </b:AccountManager>
                    <b:AccountNumber i:nil="true"/>
                    <b:Address>
                        <Country>US</Country>
                        <Geocode i:nil="true"/>
                        <Line1>4215 Rosebud Ct</Line1>
                        <Line2/>
                        <Municipality>Pensacola</Municipality>
                        <PostalCode>32504</PostalCode>
                        <StateOrProvince>FL</StateOrProvince>
                    </b:Address>
                    <b:AssignCatalog i:nil="true"/>
                    <b:AssignPriceLists/>
                    <b:BillingAddress>
                        <Country>US</Country>
                        <Geocode i:nil="true"/>
                        <Line1>4215 Rosebud Ct</Line1>
                        <Line2/>
                        <Municipality>Pensacola</Municipality>
                        <PostalCode>32504</PostalCode>
                        <StateOrProvince>FL</StateOrProvince>
                    </b:BillingAddress>
                    <b:Culture>en-US</b:Culture>
                    <b:Currency>USD</b:Currency>
                    <b:CustomerCareSpecialist>
                        <b:CustomerCareSpecialist>msatya@appriver.com</b:CustomerCareSpecialist>
                    </b:CustomerCareSpecialist>
                    <b:CustomerId>%s</b:CustomerId>
                    <b:CustomerNumber i:nil="true"/>
                    <b:CustomerSince>2020-07-23T19:41:34.156Z</b:CustomerSince>
                    <b:CustomerSource i:nil="true"/>
                    <b:Domains/>
                    <b:Industry i:nil="true"/>
                    <b:InitialBalance>0</b:InitialBalance>
                    <b:InvoiceCycle>6</b:InvoiceCycle>
                    <b:IsRenewalDateAligned>false</b:IsRenewalDateAligned>
                    <b:Name>[Test][Terraform]%s</b:Name>
                    <b:PaymentTerms>Net15</b:PaymentTerms>
                    <b:SendWelcomeEmail>false</b:SendWelcomeEmail>
                    <b:ServiceAddress i:nil="true"/>
                    <b:SkipServiceAddressValidation>true</b:SkipServiceAddressValidation>
                    <b:TaxId i:nil="true"/>
                    <b:TimeZone>Central Standard Time</b:TimeZone>
                    <b:Users/>
                    <b:VatRate>0</b:VatRate>
                    <b:VerbalPasswordToken>00000000-0000-0000-0000-000000000000</b:VerbalPasswordToken>
                    <b:Website i:nil="true"/>
                    <b:AccountActivationDate i:nil="true"/>
                    <b:PrimaryContact i:nil="true"/>
                    <b:ServiceProviderId>9D74F0B7-0225-464B-9DB5-E057E03B965B</b:ServiceProviderId>
                </Command>
                <Timeout>PT0S</Timeout>
                <WaitForResult>true</WaitForResult>
            </request>
        </ProcessCommand>
    </s:Body>
</s:Envelope>`, nauticalcustomerid, name)

	payload := strings.NewReader(formattedbody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/soap+xml")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//thank you eventsourcing for eventual consistency.
	//letting aqua denormalize for BOSUN
	fmt.Println("🤦‍♂️")
	time.Sleep(10 * time.Second)

	fmt.Println(string(body))
	return

}
