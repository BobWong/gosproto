-- Generated by github.com/davyxu/gosproto/sprotogen
-- DO NOT EDIT!

-- Enum:
--[[

-- MyCar 	
local MyCar_Monkey = 1 	
local MyCar_Monk = 2 	
local MyCar_Pig = 3 

]]

local sproto = {
	Schema = [[

.PhoneNumber {		
	number 0 : string 	
	type 1 : integer 
}

.Person {		
	name 0 : string 	
	id 1 : integer 	
	email 2 : string 	
	phone 3 : *PhoneNumber 
}

.AddressBook {		
	person 0 : *Person 
}

.MyData {		
	name 1 : string 	
	type 2 : integer 	
	int32 3 : integer 	
	uint32 4 : integer 	
	int64 5 : integer 	
	uint64 6 : integer 	
	bool 7 : boolean 
}

.MyProfile {		
	nameField 1 : MyData 	
	nameArray 2 : *MyData 	
	nameMap 3 : *MyData(type) 
}

	]],

	NameByID = { 
		[4271979557] = "PhoneNumber",
		[1498745430] = "Person",
		[2618161298] = "AddressBook",
		[2244887298] = "MyData",
		[438153711] = "MyProfile",
	},
	
	IDByName = {},
}

local t = sproto.IDByName
for k, v in pairs(sproto.NameByID) do
	t[v] = k
end

return sproto

