package main

type Ingress struct {
	Items []struct {
		Spec struct {
			Rules []struct {
				Host string `json:"host"`
				HTTP struct {
					Paths []struct {
						Path     string `json:"path"`
						PathType string `json:"pathType"`
						Backend  struct {
							Service struct {
								Name string `json:"name"`
								Port struct {
									Number int `json:"number"`
								} `json:"port"`
							} `json:"service"`
						} `json:"backend"`
					} `json:"paths"`
				} `json:"http"`
			} `json:"rules"`
		} `json:"spec"`
		Status struct {
			LoadBalancer struct {
				Ingress []struct {
					Hostname string `json:"hostname"`
				} `json:"ingress"`
			} `json:"loadBalancer"`
		} `json:"status"`
	} `json:"items"`
}
