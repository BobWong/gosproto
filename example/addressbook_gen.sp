# Generated by github.com/bobwong8975789757/gosproto/sprotogen
# DO NOT EDIT!


enum MyCar {
		
	Monkey
	Monk
	Pig
}




message PhoneNumber {

	number string

	type int32  
}


message Person {

	name string  
	
	id  int32  
	
	email string
	
	phone PhoneNumber
}

message AddressBook {

	person []Person
}


#  [agent] client -> battle # comment
message MyData {
	
	
	name string
	
	type MyCar
	
	int32 int32 //  extend standard
	#  extend standard
	uint32 uint32
	
	int64 int64  
	
	uint64 uint64  
}


message MyProfile {
	
	
	nameField MyData
	
	nameArray sMyData
}


