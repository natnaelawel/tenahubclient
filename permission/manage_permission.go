package permission

import (
	"strings"
	"fmt"
)

type permission struct {
	roles   []string
	methods []string
}

type authority map[string]permission

var authorities = authority{
	"/": permission{
		roles:   []string{"USER", ""},
		methods: []string{"GET"},
	},
	"/favicon.ico":permission{
		roles: []string{"USER", "ADMIN", "AGENT", "HEALTH_CENTER", ""},
		methods: []string{"GET", "POST"},
	},
	"/signup": permission{
		roles:   []string{"USER", ""},
		methods: []string{"GET", "POST"},
	},
	"/login": permission{
		roles:   []string{"USER", "","ADMIN","AGENT"},
		methods: []string{"GET", "POST"},
	},
	"/about": permission{
		roles:   []string{"USER", ""},
		methods: []string{"GET"},
	},
	"/logout": permission{
		roles:   []string{"USER", "","ADMIN","AGENT"},
		methods: []string{"GET"},
	},
	"/search": permission{
		roles:   []string{"USER", ""},
		methods: []string{"POST", "GET"},
	},
	"/home": permission{
		roles:   []string{"USER", ""},
		methods: []string{"GET"},
	},
	"/healthcenters": permission{
		roles:   []string{"USER", ""},
		methods: []string{"GET"},
	},
	"/feedback": permission{
		roles:   []string{"USER", ""},
		methods: []string{"GET", "POST"},
	},
	"/admin": permission{
		roles:   []string{"ADMIN"},
		methods: []string{"GET", "POST"},
	},
	"/agent": permission{
		roles:   []string{"AGENT"},
		methods: []string{"GET", "POST"},
	},
	"/healthcenter": permission{
		roles:   []string{"HEALTH_CENTER"},
		methods: []string{"GET", "POST"},
	},
}

// HasPermission checks if a given role has permission to access a given route for a given method
func HasPermission(path string, role string, method string) bool {
	fmt.Printf("Path: %s, Role: %s, Method: %s\n", path, role, method)
	if strings.HasPrefix(path, "/admin") {
		path = "/admin"
	}else if strings.HasPrefix(path, "/agent"){
		path = "/agent"
	}else if strings.HasPrefix(path, "/healthcenters"){
		path = "/healthcenters"
	}else if strings.HasPrefix(path, "/healthcenter"){
		path = "/healthcenter"
	}
	fmt.Println("path: " + path)
	perm := authorities[path]
	checkedRole := checkRole(role, perm.roles)
	checkedMethod := checkMethod(method, perm.methods)
	if !checkedRole || !checkedMethod {
		return false
	}
	return true
}

func checkRole(role string, roles []string) bool {
	for _, r := range roles {
		if strings.ToUpper(r) == strings.ToUpper(role) {
			return true
		}
	}
	return false
}

func checkMethod(method string, methods []string) bool {
	for _, m := range methods {
		if strings.ToUpper(m) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}
