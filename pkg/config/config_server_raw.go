package config

// !!!!!!! DON'T CHANGE THIS FILE !!!!!!
// This file is copied from istio-broker-proxy
// Do changes in istio-broker-proxy, copy it to here
// and replace DestinationWeight by  RouteDestination

import (
	"fmt"
	"istio.io/api/networking/v1alpha3"
)

const ingressCertName = "cf-service"

func createRawServiceEntryForExternalService(endpointAddress string, portNumber uint32, serviceName string) *v1alpha3.ServiceEntry {
	hosts := []string{createServiceHost(serviceName)}
	portName := fmt.Sprintf("%s-%d", serviceName, portNumber)

	ports := v1alpha3.Port{Number: portNumber, Name: portName, Protocol: "TCP"}
	resolution := v1alpha3.ServiceEntry_STATIC
	endpoint := v1alpha3.ServiceEntry_Endpoint{Address: endpointAddress}
	return &v1alpha3.ServiceEntry{Hosts: hosts, Ports: []*v1alpha3.Port{&ports}, Resolution: resolution,
		Endpoints: []*v1alpha3.ServiceEntry_Endpoint{&endpoint}}
}

func createServiceHost(serviceName string) string {
	serviceHost := serviceName + ".service-fabrik"
	return serviceHost
}

func createRawIngressVirtualServiceForExternalService(hostName string, port uint32, serviceName string) *v1alpha3.VirtualService {
	destination := v1alpha3.Destination{Host: createServiceHost(serviceName),
		Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: port}}}
	route := v1alpha3.TCPRoute{Route: []*v1alpha3.RouteDestination{{Destination: &destination}}}
	tcpRoutes := []*v1alpha3.TCPRoute{&route}
	hosts := []string{hostName}
	gateways := []string{serviceName + "-gateway"}
	return &v1alpha3.VirtualService{Tcp: tcpRoutes, Hosts: hosts, Gateways: gateways}
}

func createRawIngressGatewayForExternalService(hostName string, portNumber uint32, clientName string) *v1alpha3.Gateway {
	port := v1alpha3.Port{Number: portNumber, Name: "tls", Protocol: "TLS"}
	hosts := []string{hostName}
	certPath := "/var/vcap/jobs/envoy/config/certs/"
	tls := v1alpha3.Server_TLSOptions{Mode: v1alpha3.Server_TLSOptions_MUTUAL,
		ServerCertificate: certPath + ingressCertName + ".crt",
		PrivateKey:        certPath + ingressCertName + ".key",
		CaCertificates:    certPath + "ca.crt",
		SubjectAltNames:   []string{clientName}}
	return &v1alpha3.Gateway{Servers: []*v1alpha3.Server{{Port: &port, Hosts: hosts, Tls: &tls}}}
}
