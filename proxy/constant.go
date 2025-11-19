/*
 * @Description: constant
 */
package proxy

var allowModifyHosts map[string]bool = map[string]bool{
	"www.github.com": true,
}

// cert download url
const CertDownloadUrl = "/download-proxy-cert"

// Event
const (
	EventNameMessage = "sendMessage" // send toast message to front

	EventNameProxyRemoteAddr = "newRemoteAddr" // new remote-addr
	EventNameProxyHost       = "newHost"       // new host

	EventNameReloadFlow = "reloadFlowList" // refresh table data
)

const (
	cacheKeyFlowChangeMark = "flowChangeMark" // bool, mark has new flow
	cacheKeyFrontPause     = "frontPause"     // bool, mark front is inactive
)

// HTTPS use
const CERTIFICATE = `-----BEGIN CERTIFICATE-----
MIIDNTCCAh2gAwIBAgIUcWuIU767yZ8g3ffg1xBlRmCH5QEwDQYJKoZIhvcNAQEL
BQAwKDESMBAGA1UEAwwJbWl0bXByb3h5MRIwEAYDVQQKDAltaXRtcHJveHkwHhcN
MjUwMTA1MTEzMzAwWhcNMzUwMTA1MTEzMzAwWjAoMRIwEAYDVQQDDAltaXRtcHJv
eHkxEjAQBgNVBAoMCW1pdG1wcm94eTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBALEmh+tJPHyLiYAq38G/h774ZF2YqxCCI/SLFm2EKkOesaVDwWTdxifS
yzGCp+vR10jP5TJngoyxmOmzUq8utXtwnMbuRAMiH7k/Rx8D+TXrp+Wef/a1i8lH
CRPiRDhM0J3WbgDVKzI9oB6qE8UvtWnE4+xASqcnJnJivc39pSmI/tbCelN22Z9I
b8AxMHv+Yzv/Z/ximGojxaSOi38dA7kurmuQNjoz38uERRpG4lBEjgVqeRIbS0L6
pn+384AcuHQNzXYXv0w4w1KgaSdpzBM180maFVnYeBXeVnp01K/e+rWtZ0ZrWPq8
wVVAr4cdeJuXpszZ0kLSsgcAJf6ZsZsCAwEAAaNXMFUwDwYDVR0TAQH/BAUwAwEB
/zATBgNVHSUEDDAKBggrBgEFBQcDATAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYE
FKBt3U6gS3WOgXINNfJfIHLCYKJAMA0GCSqGSIb3DQEBCwUAA4IBAQBn6sNP8eeG
KNNoidsTRe8PFgo4EuOe921RrTBG4I/E84umZbQHtaoYoNRq5cmha+lYdK9z7PgH
bVpHg2Ijkt97cVOqGyDqQ6exfZt2APXtNP8hcb9sPmktZwOzai8dJrN+OCcQeta1
RsXX92SmSa38rVaUP2X89CAFx2n342njMIV04zyRiBZyp3dKmUVwgjciC6QVbBlT
+KSBMHOtE6Sq+6U6N3zov7tCSE0LPEZXTZK/Jqtw6vGyef+reO1vWf2TD+Zyzp23
dq/MtdNLxcHSs6kQPyPlKq29bIRp+vnQklAykeeAMc8C9oDrF6JccdSPUu0sCVIa
GHWcI918sS35
-----END CERTIFICATE-----`

const RSA_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAsSaH60k8fIuJgCrfwb+HvvhkXZirEIIj9IsWbYQqQ56xpUPB
ZN3GJ9LLMYKn69HXSM/lMmeCjLGY6bNSry61e3Ccxu5EAyIfuT9HHwP5Neun5Z5/
9rWLyUcJE+JEOEzQndZuANUrMj2gHqoTxS+1acTj7EBKpycmcmK9zf2lKYj+1sJ6
U3bZn0hvwDEwe/5jO/9n/GKYaiPFpI6Lfx0DuS6ua5A2OjPfy4RFGkbiUESOBWp5
EhtLQvqmf7fzgBy4dA3Ndhe/TDjDUqBpJ2nMEzXzSZoVWdh4Fd5WenTUr976ta1n
RmtY+rzBVUCvhx14m5emzNnSQtKyBwAl/pmxmwIDAQABAoIBAEEBKxH/E0d+aovR
776oWhcNDcM8oNUTdD8phKbWPy4F/xjuRnqNWMfHNZGq6JSDsVSCGRAIZKLdDof/
KgGniafMuHLU/rcJoAt5eU8bxahwG3GKWll839kRjHaz6iWgplKs2zvJ/SpX67gk
39FVQuQXgC1yoBI7LwP9ULnUb/VPtqoo7EeJ30nVKCMRk1D54IWOJnDXkQUwxoAx
raTayyhIasjdkGX/8tWNWyrjmRdXXdz03tZynpbA0kjauGn5jugJX54iCidtKOVf
8g40EAJVZU2Os3IA8Vk0J6kTKui545Euot8dnqiPFn1cMEiCoxouaTSnw1GoOh9k
+2eV+d0CgYEA9LqmUyxqeSAHBzaQHzKRxanDQ8/XQ1IUbKSPtpdbA5hB2SP7P6dh
K3i0bFoQx1ivYONIEjEi2WQKYb8qA0dxNU35GcfnSXCABJlVLVpTA6y74JQpmeGz
KN/f6inMfWMzz9l6kch7mw2dA+qyAeSBEyRHuJBL3CNnH0Iaacj/DCUCgYEAuU8h
mpmCPe93dvsObSoi2DivoXDheXXeAroF3q1dTs+71ox34hS3hLflcYA4co5cc5BF
oP4cr5OsRbFAggq8hkcE8gYRlh0MMbvS5A/rA1PheolpNi2GPuq85fcnqvoTNGh5
EfLLNG552GWdt+UfaQ+MZtBBosxzjeBa/8COer8CgYBLK3GfJotlvQKrooTOMFg6
IJTLpXF3SiTf4gQKj57zEbdYHQc1XBIqqh6xEVEKCiGEXwtqGxGbOHE/6Nncu/r3
5oVW9+IA2r9XqklRbucEv+/NYXUNwReovGson8Ih2XC2rRRK5wugaqQJ6fLXRfqn
iINHG3tRmMY6AbHfQE+myQKBgB+rJlMFaPysr/3ladG+IbZR+bHdWKspOYe9wsML
XabaWX+RAevhHnvP8aY0A3GQYrxdK4kX5E8dLZO7dwrCSvWnrXCvGNJOpIbbUm8n
SUKIU0Bdu11G6mpWM3IZiu3tTtrl/8rRc2jCeJBCzXFf9r49qgDoFk98IUNd1C9g
+PiBAoGAUaCMV2ZclidZhathQs266GI8ncf53kuEMi7ESkIJPQd8Is8z6+vF4Feh
cJiI4oFmVg29tveZ7ZhZf4/fR/He2zOmXhxYDSUjNu16EmXvuQqCAbR1Yv/7H2vU
Uwi4bxsNbZpmvBzkZ8+5+kpA6G8DGiWduofg2654HnrjIFLHPDQ=
-----END RSA PRIVATE KEY-----`
