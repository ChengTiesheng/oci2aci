all:
	godep go build -o oci2aci .
clean:
	rm oci2aci
