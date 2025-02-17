package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdManager struct {
	client *clientv3.Client
}

var (
	etcdEndpoints  *string
	pathToCACert *string
	pathToCert *string
	pathToKey *string
	listenPort *string
	debug *bool
)

func loadTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:     caCertPool,
	}, nil
}

func NewEtcdManager(certFile, keyFile, caFile string, endpoints []string) (*EtcdManager, error) {
	tlsConfig, err := loadTLSConfig(certFile, keyFile, caFile)
	if err != nil {
		return nil, err
	}

	cfg := clientv3.Config{
		Endpoints:   endpoints,
		TLS:        tlsConfig,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdManager{client: client}, nil
}

func (em *EtcdManager) getEndpointStatus(c *gin.Context) {
	ctx := context.Background()
	var statusInfo []map[string]interface{}

	for _, ep := range em.client.Endpoints() {
		status, err := em.client.Status(ctx, ep)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		info := map[string]interface{}{
			//"endpoint":     status.Endpoint,
			"leader":       status.Header.MemberId == status.Leader,
			"version":      status.Version,
			"dbSize":      status.DbSize,
			"dbSizeInUse": status.DbSizeInUse,
		}
		statusInfo = append(statusInfo, info)
	}

	c.JSON(http.StatusOK, statusInfo)
}

func (em *EtcdManager) compactEtcd(c *gin.Context) {
	ctx := context.Background()

	// Get current revision
	resp, err := em.client.Status(ctx, em.client.Endpoints()[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compact up to current revision
	_, err = em.client.Compact(ctx, resp.Header.Revision)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compaction successful"})
}

func (em *EtcdManager) defragEndpoints(c *gin.Context) {
	ctx := context.Background()

	for _, ep := range em.client.Endpoints() {
		_, err := em.client.Defragment(ctx, ep)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"endpoint": ep,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Defragmentation successful"})
}

func set_parameters() {
	// Set defaults from environment variables
	/*
	defaultPort := 8080
	if portEnv := os.Getenv("SERVER_PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil {
			defaultPort = p
		}
	}
	*/
	defaultListenPort := os.Getenv("LISTEN_PORT")
	if defaultListenPort == "" {
		defaultListenPort = "8080"
	}
	defaultEtcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if defaultEtcdEndpoints == "" {
		defaultEtcdEndpoints = "localhost:2379"
	}
	defaultPathToCACert := os.Getenv("CA_CERT")
	if defaultPathToCACert == "" {
		defaultPathToCACert = "./cacert.pem"
	}
	defaultPathToCert := os.Getenv("CERT")
	if defaultPathToCert == "" {
		defaultPathToCert = "./cert.pem"
	}
	defaultPathToKey := os.Getenv("KEY")
	if defaultPathToKey == "" {
		defaultPathToKey = "./key.pem"
	}

	// Define command-line flags
	listenPort = flag.String("listen-port", defaultListenPort, "Listening port")
	etcdEndpoints = flag.String("etcd-endpoints", defaultEtcdEndpoints, "Etcd endpoints")
	pathToCACert = flag.String("cacert", defaultPathToCACert, "Etcd CA certificate")
	pathToCert = flag.String("cert", defaultPathToCert, "Certificate")
	pathToKey = flag.String("key", defaultPathToKey, "Private key")
	debug = flag.Bool("debug", false, "Enable debug mode")

	// Parse the flags
	flag.Parse()
}

func main() {
	set_parameters()

	// Initialize Gin
	r := gin.Default()

	// Configure CORS
	r.Use(cors.Default())
	/*
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	*/

	// Initialize EtcdManager with your certificates
	em, err := NewEtcdManager(
		*pathToCert,
		*pathToKey,
		*pathToCACert,
		[]string{*etcdEndpoints}, // Your etcd endpoints
	)
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	r.GET("/api/status", em.getEndpointStatus)
	r.POST("/api/compact", em.compactEtcd)
	r.POST("/api/defrag", em.defragEndpoints)

	// Start server
	r.Run(*listenPort)
}

