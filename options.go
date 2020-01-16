package gocheck

import (
	"fmt"
	"strconv"
)

//Option is a Check Option
type Option func(*Check) error

//SetFilePath sets the path to save the config file at.
func SetFilePath(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().FilePath = s
		c.SAMHTTPProxy.Config().FilePath = s
		return nil
	}
}

//SetHostsFile sets the path to look for the hosts file
func SetHostsFile(s string) func(*Check) error {
	return func(c *Check) error {
		c.hostsfile = s
		return nil
	}
}

//SetType sets the type of the forwarder server
func SetType(s string) func(*Check) error {
	return func(c *Check) error {
		if s == "http" {
			c.SAMForwarder.Config().Type = s
			c.SAMHTTPProxy.Config().Type = s
			return nil
		} else {
			c.SAMForwarder.Config().Type = "server"
			c.SAMHTTPProxy.Config().Type = "httpclient"
			return nil
		}
	}
}

//SetSigType sets the type of the forwarder server
func SetSigType(s string) func(*Check) error {
	return func(c *Check) error {
		if s == "" {
			c.SAMForwarder.Config().SigType = ""
		} else if s == "DSA_SHA1" {
			c.SAMForwarder.Config().SigType = "DSA_SHA1"
		} else if s == "ECDSA_SHA256_P256" {
			c.SAMForwarder.Config().SigType = "ECDSA_SHA256_P256"
		} else if s == "ECDSA_SHA384_P384" {
			c.SAMForwarder.Config().SigType = "ECDSA_SHA384_P384"
		} else if s == "ECDSA_SHA512_P521" {
			c.SAMForwarder.Config().SigType = "ECDSA_SHA512_P521"
		} else if s == "EdDSA_SHA512_Ed25519" {
			c.SAMForwarder.Config().SigType = "EdDSA_SHA512_Ed25519"
		} else {
			c.SAMForwarder.Config().SigType = "EdDSA_SHA512_Ed25519"
		}
		return nil
	}
}

//SetSaveFile tells the router to save the tunnel's keys long-term
func SetSaveFile(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().SaveFile = b
		c.SAMHTTPProxy.Config().SaveFile = b
		return nil
	}
}

//SetHost sets the host of the service to forward
func SetHost(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().TargetHost = s
		c.SAMHTTPProxy.Config().TargetHost = s
		return nil
	}
}

//SetPort sets the port of the service to forward
func SetPort(s string) func(*Check) error {
	return func(c *Check) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid TCP Server Target Port %s; non-number ", s)
		}
		clientport := port + 1
		if port < 65536 && port > -1 {
			c.SAMForwarder.Config().TargetPort = s
			c.SAMHTTPProxy.Config().TargetPort = strconv.Itoa(clientport)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetSAMHost sets the host of the Check's SAM bridge
func SetSAMHost(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().SamHost = s
		c.SAMHTTPProxy.Config().SamHost = s
		return nil
	}
}

//SetSAMPort sets the port of the Check's SAM bridge using a string
func SetSAMPort(s string) func(*Check) error {
	return func(c *Check) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid SAM Port %s; non-number", s)
		}
		if port < 65536 && port > -1 {
			c.SAMForwarder.Config().SamPort = s
			c.SAMHTTPProxy.Config().SamPort = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetName sets the host of the Check's SAM bridge
func SetName(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().TunName = s
		c.SAMHTTPProxy.Config().TunName = s + "httpclient"
		return nil
	}
}

//SetInLength sets the number of hops inbound
func SetInLength(u int) func(*Check) error {
	return func(c *Check) error {
		if u < 7 && u >= 0 {
			c.SAMForwarder.Config().InLength = u
			c.SAMHTTPProxy.Config().InLength = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutLength sets the number of hops outbound
func SetOutLength(u int) func(*Check) error {
	return func(c *Check) error {
		if u < 7 && u >= 0 {
			c.SAMForwarder.Config().OutLength = u
			c.SAMHTTPProxy.Config().OutLength = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel length")
	}
}

//SetInVariance sets the variance of a number of hops inbound
func SetInVariance(i int) func(*Check) error {
	return func(c *Check) error {
		if i < 7 && i > -7 {
			c.SAMForwarder.Config().InVariance = i
			c.SAMHTTPProxy.Config().InVariance = i
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutVariance sets the variance of a number of hops outbound
func SetOutVariance(i int) func(*Check) error {
	return func(c *Check) error {
		if i < 7 && i > -7 {
			c.SAMForwarder.Config().OutVariance = i
			c.SAMHTTPProxy.Config().OutVariance = i
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel variance")
	}
}

//SetInQuantity sets the inbound tunnel quantity
func SetInQuantity(u int) func(*Check) error {
	return func(c *Check) error {
		if u <= 16 && u > 0 {
			c.SAMForwarder.Config().InQuantity = u
			c.SAMHTTPProxy.Config().InQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel quantity")
	}
}

//SetOutQuantity sets the outbound tunnel quantity
func SetOutQuantity(u int) func(*Check) error {
	return func(c *Check) error {
		if u <= 16 && u > 0 {
			c.SAMForwarder.Config().OutQuantity = u
			c.SAMHTTPProxy.Config().OutQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel quantity")
	}
}

//SetInBackups sets the inbound tunnel backups
func SetInBackups(u int) func(*Check) error {
	return func(c *Check) error {
		if u < 6 && u >= 0 {
			c.SAMForwarder.Config().InBackupQuantity = u
			c.SAMHTTPProxy.Config().InBackupQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel backup quantity")
	}
}

//SetOutBackups sets the inbound tunnel backups
func SetOutBackups(u int) func(*Check) error {
	return func(c *Check) error {
		if u < 6 && u >= 0 {
			c.SAMForwarder.Config().OutBackupQuantity = u
			c.SAMHTTPProxy.Config().OutBackupQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel backup quantity")
	}
}

//SetEncrypt tells the router to use an encrypted leaseset
func SetEncrypt(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().EncryptLeaseSet = b
		c.SAMHTTPProxy.Config().EncryptLeaseSet = b
		return nil
	}
}

//SetLeaseSetKey sets the host of the Check's SAM bridge
func SetLeaseSetKey(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().LeaseSetKey = s
		c.SAMHTTPProxy.Config().LeaseSetKey = s
		return nil
	}
}

//SetLeaseSetPrivateKey sets the host of the Check's SAM bridge
func SetLeaseSetPrivateKey(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().LeaseSetPrivateKey = s
		c.SAMHTTPProxy.Config().LeaseSetPrivateKey = s
		return nil
	}
}

//SetLeaseSetPrivateSigningKey sets the host of the Check's SAM bridge
func SetLeaseSetPrivateSigningKey(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().LeaseSetPrivateSigningKey = s
		c.SAMHTTPProxy.Config().LeaseSetPrivateSigningKey = s
		return nil
	}
}

//SetMessageReliability sets the host of the Check's SAM bridge
func SetMessageReliability(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().MessageReliability = s
		c.SAMHTTPProxy.Config().MessageReliability = s
		return nil
	}
}

//SetAllowZeroIn tells the tunnel to accept zero-hop peers
func SetAllowZeroIn(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().InAllowZeroHop = b
		c.SAMHTTPProxy.Config().InAllowZeroHop = b
		return nil
	}
}

//SetAllowZeroOut tells the tunnel to accept zero-hop peers
func SetAllowZeroOut(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().OutAllowZeroHop = b
		c.SAMHTTPProxy.Config().OutAllowZeroHop = b
		return nil
	}
}

//SetCompress tells clients to use compression
func SetCompress(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().UseCompression = b
		c.SAMHTTPProxy.Config().UseCompression = b
		return nil
	}
}

//SetFastRecieve tells clients to use compression
func SetFastRecieve(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().FastRecieve = b
		c.SAMHTTPProxy.Config().FastRecieve = b
		return nil
	}
}

//SetReduceIdle tells the connection to reduce it's tunnels during extended idle time.
func SetReduceIdle(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().ReduceIdle = b
		c.SAMHTTPProxy.Config().ReduceIdle = false
		return nil
	}
}

//SetReduceIdleTime sets the time to wait before reducing tunnels to idle levels
func SetReduceIdleTime(u int) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().ReduceIdleTime = 300000
		if u >= 6 {
			c.SAMForwarder.Config().ReduceIdleTime = (u * 60) * 1000
			c.SAMHTTPProxy.Config().ReduceIdleTime = 3000000
			return nil
		}
		return fmt.Errorf("Invalid reduce idle timeout(Measured in minutes) %v", u)
	}
}

//SetReduceIdleTimeMs sets the time to wait before reducing tunnels to idle levels in milliseconds
func SetReduceIdleTimeMs(u int) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().ReduceIdleTime = 300000
		if u >= 300000 {
			c.SAMForwarder.Config().ReduceIdleTime = u
			c.SAMHTTPProxy.Config().ReduceIdleTime = 3000000
			return nil
		}
		return fmt.Errorf("Invalid reduce idle timeout(Measured in milliseconds) %v", u)
	}
}

//SetReduceIdleQuantity sets minimum number of tunnels to reduce to during idle time
func SetReduceIdleQuantity(u int) func(*Check) error {
	return func(c *Check) error {
		if u < 5 {
			c.SAMForwarder.Config().ReduceIdleQuantity = u
			c.SAMHTTPProxy.Config().ReduceIdleQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid reduce tunnel quantity")
	}
}

//SetCloseIdle tells the connection to close it's tunnels during extended idle time.
func SetCloseIdle(b bool) func(*Check) error {
	return func(c *Check) error {
		c.SAMHTTPProxy.Config().CloseIdle = b
		c.SAMForwarder.Config().CloseIdle = false
		return nil
	}
}

//SetCloseIdleTime sets the time to wait before closing tunnels to idle levels
func SetCloseIdleTime(u int) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().CloseIdleTime = 300000
		if u >= 6 {
			c.SAMForwarder.Config().CloseIdleTime = (u * 60) * 1000
			c.SAMHTTPProxy.Config().CloseIdleTime = 3000000
			return nil
		}
		return fmt.Errorf("Invalid close idle timeout(Measured in minutes) %v", u)
	}
}

//SetCloseIdleTimeMs sets the time to wait before closing tunnels to idle levels in milliseconds
func SetCloseIdleTimeMs(u int) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().CloseIdleTime = 300000
		if u >= 300000 {
			c.SAMForwarder.Config().CloseIdleTime = u
			c.SAMHTTPProxy.Config().CloseIdleTime = 3000000
			return nil
		}
		return fmt.Errorf("Invalid close idle timeout(Measured in milliseconds) %v", u)
	}
}

//SetAccessListType tells the system to treat the accessList as a whitelist
func SetAccessListType(s string) func(*Check) error {
	return func(c *Check) error {
		if s == "whitelist" {
			c.SAMForwarder.Config().AccessListType = "whitelist"
			return nil
		} else if s == "blacklist" {
			c.SAMForwarder.Config().AccessListType = "blacklist"
			return nil
		} else if s == "none" {
			c.SAMForwarder.Config().AccessListType = ""
			return nil
		} else if s == "" {
			c.SAMForwarder.Config().AccessListType = ""
			return nil
		}
		return fmt.Errorf("Invalid Access list type(whitelist, blacklist, none)")
	}
}

//SetAccessList tells the system to treat the accessList as a whitelist
func SetAccessList(s []string) func(*Check) error {
	return func(c *Check) error {
		if len(s) > 0 {
			for _, a := range s {
				c.SAMForwarder.Config().AccessList = append(c.SAMForwarder.Config().AccessList, a)
			}
			return nil
		}
		return nil
	}
}

//SetTargetForPort sets the port of the Check's SAM bridge using a string
/*func SetTargetForPort443(s string) func(*Check) error {
	return func(c *Check) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid Target Port %s; non-number ", s)
		}
		if port < 65536 && port > -1 {
			c.SAMForwarder.Config().TargetForPort443 = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}
*/

//SetKeyFile sets
func SetKeyFile(s string) func(*Check) error {
	return func(c *Check) error {
		c.SAMForwarder.Config().KeyFilePath = s
		c.SAMHTTPProxy.Config().KeyFilePath = s
		return nil
	}
}
