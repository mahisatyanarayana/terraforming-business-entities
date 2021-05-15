resource "zixar_customer" "customer0" {
     name = "iwasterraformedatsatyahome${var.id}"
     customerid = "c2905a6a-523a-4a0a-81e5-ac27d30348${var.id}"
}

resource "zixar_user" "user0" {

     email = "melissa@satya.home"
     firstname = "m"
     middlename = "k" //optional in the schema
     lastname = "satyanarayana"
     nauticalcustomerid = zixar_customer.customer0.customerid

}


resource "zixar_user" "user2" {

     email = "mahi@satya.home"
     firstname = "mahi"
     lastname = "satyanarayana"
     nauticalcustomerid = zixar_customer.customer0.customerid

}

output "user0_guid" {
     value = "${zixar_user.user0.useraccountid}"
}

output "customer_name" {
     value = zixar_customer.customer0.name
}

variable id {
  type    = string
  default = "06"
}