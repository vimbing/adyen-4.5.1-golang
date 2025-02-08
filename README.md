# Adyen encryption in golang

Based on [THIS  REPO](https://gist.github.com/levi-nz/da0690e204c6150207bd0b52a0897b61)

Tested on adyen secured  fields 4.5.1, 4.5.0

# Example

```go
package  main

import  (
	"fmt"
	"github.com/vimbing/adyen-4.5.1-golang"
)

func  main()  {
	encryptor,  err  :=  adyen.NewEncryptor(
		adyen.NewUrl(
			"live_DY4VMYQL5ZHXXE5NLG4RA5PYKYWDYAU2",
			"4.5.0",
			"aHR0cHM6Ly9jaGVsc2VhZmMuM2RkaWdpdGFsdmVudWUuY29t",
		),
		"10001|E9299A45B34AE878855F3E66136B461664F519E85F36E59B505CD6590311FE96BAF50830BED460FE6EB8AD39B3E4BFCF5028A33A64C518E3BC13F23E49CE9C68B13A3ED9BB9233C166A7572755E62CB67AAF7A6AFC1070CAD7FF3F6FD8C070168FC6ED31E81F3DE10A93D6A9494F9D24900F1499D95264E66E3DC357B4628E02A6DF0ED37196539309AB0B1EA7EEB2BD67452B16289452D617C687867981C3570E0C43C51EB273154011D53F09B2B2E1AAD41B13B686A861D2C095DFEA258AD589AE482CAF9B05EFFF1C16EF182D67CA459B6EBD00E63F170307B56237A6C8AE593EFAD9E58AEC7D560B41B3412DD7D5E64B76BFEF75354DC52BD2138B77F279",
		adyen.Card{
			Number: "5344 3124 6454 4325",
			ExpiryMonth: "01",
			ExpiryYear: "2030",
			CVV: "543", 
	})

	if  err  !=  nil  {
		panic(err)
	}
	
	card,  err  :=  encryptor.Encrypt()

	if  err  !=  nil  {
		panic(err)
	}
	
	fmt.Printf("card: %+v\n",  card)
}
```

# URL Building

Library needs a adyen url for encrypting data. You can probably find url somewhere in website's source code, it will look like this: 
- https://checkoutshopper-live-us.adyen.com/checkoutshopper/securedfields/live_DY4VMYQL5ZHXXE5NLG4RA5PYKYWDYAU2/4.5.0/securedFields.html?type=card&d=aHR0cHM6Ly9jaGVsc2VhZmMuM2RkaWdpdGFsdmVudWUuY29t


You can then use this url directly while initializing new encryptor, or use adyen.NewUrl and provide all needed parameters in order to create url. 
### Disclaimer
I'm not sure if  d parameter is needed, if you will not manage to find it, then try to encrypt without it.
