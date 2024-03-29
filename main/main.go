package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os/user"
)

import (
	"github.com/eyedeekay/gocheck"
	"github.com/eyedeekay/sam-forwarder/config"
)

/*
func main() {
	check, err := gocheck.NewSAMChecker(*hostsfile)
	if err != nil {
		panic(err)
	}
	check.CheckAll()
	check.Serve()
}
*/

var cfg = &tls.Config{
	MinVersion:               tls.VersionTLS12,
	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	PreferServerCipherSuites: true,
	CipherSuites: []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	},
}

var (
	host               = flag.String("a", "127.0.0.1", "hostname to serve on")
	port               = flag.String("p", "9778", "port to serve locally on")
	otherproxy         = flag.String("h", "127.0.0.1:4444", "Use an external HTTP Proxy(Set this option to 'no' to use SAM")
	samhost            = flag.String("sh", "127.0.0.1", "sam host to connect to")
	samport            = flag.String("sp", "7656", "sam port to connect to")
	directory          = flag.String("d", "./www", "the directory of static files to host(default ./www)")
	usei2p             = flag.Bool("i", true, "save i2p keys(and thus destinations) across reboots")
	servicename        = flag.String("n", "gocheck", "name to give the tunnel(default gocheck)")
	useCompression     = flag.Bool("g", true, "Uze gzip(true or false)")
	accessListType     = flag.String("l", "none", "Type of access list to use, can be \"whitelist\" \"blacklist\" or \"none\".")
	encryptLeaseSet    = flag.Bool("c", false, "Use an encrypted leaseset(true or false)")
	allowZeroHop       = flag.Bool("z", false, "Allow zero-hop, non-anonymous tunnels(true or false)")
	reduceIdle         = flag.Bool("r", false, "Reduce tunnel quantity when idle(true or false)")
	reduceIdleTime     = flag.Int("rt", 600000, "Reduce tunnel quantity after X (milliseconds)")
	reduceIdleQuantity = flag.Int("rc", 3, "Reduce idle tunnel quantity to X (0 to 5)")
	inLength           = flag.Int("il", 3, "Set inbound tunnel length(0 to 7)")
	outLength          = flag.Int("ol", 3, "Set outbound tunnel length(0 to 7)")
	inQuantity         = flag.Int("iq", 2, "Set inbound tunnel quantity(0 to 15)")
	outQuantity        = flag.Int("oq", 2, "Set outbound tunnel quantity(0 to 15)")
	inVariance         = flag.Int("iv", 0, "Set inbound tunnel length variance(-7 to 7)")
	outVariance        = flag.Int("ov", 0, "Set outbound tunnel length variance(-7 to 7)")
	inBackupQuantity   = flag.Int("ib", 1, "Set inbound tunnel backup quantity(0 to 5)")
	outBackupQuantity  = flag.Int("ob", 1, "Set outbound tunnel backup quantity(0 to 5)")
	iniFile            = flag.String("f", "none", "Use an ini file for configuration")
	useTLS             = flag.Bool("t", false, "Generate or use an existing TLS certificate")
	certFile           = flag.String("m", "cert", "Certificate name to use")
	importFile         = flag.String("j", "", "import an existing json Sites history")
	importPeers        = flag.String("J", "", "import an existing json Peers history")
	peersfile          = flag.String("peers", "", "load a list of peers to monitor and query for unknown sites")
	scriptjs           = flag.String("scriptjs", "../script.js", "Serve some javascript on the application page")
	stylecss           = flag.String("stylecss", "../style.css", "Serve some CSS on the application page")
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	hostspath := user.HomeDir + "/.i2p/hosts.txt"
	hostsfile := flag.String("hosts", hostspath, "Hosts file to use.")
	flag.Parse()
	eepsite := &gocheck.Check{
		RegularProxy: *otherproxy,
	}
	config := i2ptunconf.NewI2PBlankTunConf()
	if *iniFile != "none" {
		var err error
		config, err = i2ptunconf.NewI2PTunConf(*iniFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	config.TargetHost = config.GetHost(*host, "127.0.0.1")
	config.TargetPort = config.GetPort(*port, "9778")
	config.SaveFile = config.GetSaveFile(*usei2p, true)
	config.SamHost = config.GetSAMHost(*samhost, "127.0.0.1")
	config.SamPort = config.GetSAMPort(*samport, "7656")
	config.TunName = config.GetKeys(*servicename, "echosam")
	config.InLength = config.GetInLength(*inLength, 3)
	config.OutLength = config.GetOutLength(*outLength, 3)
	config.InVariance = config.GetInVariance(*inVariance, 0)
	config.OutVariance = config.GetOutVariance(*outVariance, 0)
	config.InQuantity = config.GetInQuantity(*inQuantity, 2)
	config.OutQuantity = config.GetOutQuantity(*outQuantity, 2)
	config.InBackupQuantity = config.GetInBackups(*inBackupQuantity, 1)
	config.OutBackupQuantity = config.GetOutBackups(*outBackupQuantity, 1)
	config.EncryptLeaseSet = config.GetEncryptLeaseset(*encryptLeaseSet, false)
	config.InAllowZeroHop = config.GetInAllowZeroHop(*allowZeroHop, false)
	config.OutAllowZeroHop = config.GetOutAllowZeroHop(*allowZeroHop, false)
	config.UseCompression = config.GetUseCompression(*useCompression, true)
	config.ReduceIdle = config.GetReduceOnIdle(*reduceIdle, true)
	config.ReduceIdleTime = config.GetReduceIdleTime(*reduceIdleTime, 600000)
	config.ReduceIdleQuantity = config.GetReduceIdleQuantity(*reduceIdleQuantity, 2)
	config.AccessListType = config.GetAccessListType(*accessListType, "none")
	config.Type = config.GetTypes(false, false, false, "server")

	eepsite, err = gocheck.NewSAMCheckerFromOptions(
		gocheck.SetType(config.Type),
		gocheck.SetSAMHost(config.SamHost),
		gocheck.SetSAMPort(config.SamPort),
		gocheck.SetHost(config.TargetHost),
		gocheck.SetPort(config.TargetPort),
		gocheck.SetSaveFile(config.SaveFile),
		gocheck.SetName(config.TunName),
		gocheck.SetInLength(config.InLength),
		gocheck.SetOutLength(config.OutLength),
		gocheck.SetInVariance(config.InVariance),
		gocheck.SetOutVariance(config.OutVariance),
		gocheck.SetInQuantity(config.InQuantity),
		gocheck.SetOutQuantity(config.OutQuantity),
		gocheck.SetInBackups(config.InBackupQuantity),
		gocheck.SetOutBackups(config.OutBackupQuantity),
		gocheck.SetEncrypt(config.EncryptLeaseSet),
		gocheck.SetAllowZeroIn(config.InAllowZeroHop),
		gocheck.SetAllowZeroOut(config.OutAllowZeroHop),
		gocheck.SetCompress(config.UseCompression),
		gocheck.SetReduceIdle(config.ReduceIdle),
		gocheck.SetReduceIdleTimeMs(config.ReduceIdleTime),
		gocheck.SetReduceIdleQuantity(config.ReduceIdleQuantity),
		gocheck.SetAccessListType(config.AccessListType),
		gocheck.SetAccessList(config.AccessList),
		gocheck.SetHostsFile(*hostsfile),
		gocheck.SetPeersFile(*peersfile),
		gocheck.SetProxy(*otherproxy),
		gocheck.SetJsonImportSites(*importFile),
		gocheck.SetJsonImportPeers(*importPeers),
	)
	if err != nil {
		log.Fatal(err)
	}
	eepsite.ScriptJS = *scriptjs
	eepsite.StyleCSS = *stylecss

	if eepsite != nil {
		log.Println("Starting server")
		if err = eepsite.Serve(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Unable to start, eepsite was", eepsite)
	}
}
