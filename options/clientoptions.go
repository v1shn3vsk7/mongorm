package options

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/youmark/pkcs8"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
)

type ClientOptions struct {
	opts *mongo_options.ClientOptions
}

// Client creates a new ClientOptions instance.
func Client() *ClientOptions {
	return &ClientOptions{
		opts: mongo_options.Client(),
	}
}

// Validate validates the client options. This method will return the first error found.
func (c *ClientOptions) Validate() error {
	return c.opts.Validate()
}

// GetURI returns the original URI used to configure the ClientOptions instance. If ApplyURI was not called during
// construction, this returns "".
func (c *ClientOptions) GetURI() string {
	return c.opts.GetURI()
}

// ApplyURI parses the given URI and sets options accordingly. The URI can contain host names, IPv4/IPv6 literals, or
// an SRV record that will be resolved when the Client is created. When using an SRV record, TLS support is
// implictly enabled. Specify the "tls=false" URI option to override this.
//
// If the connection string contains any options that have previously been set, it will overwrite them. Options that
// correspond to multiple URI parameters, such as WriteConcern, will be completely overwritten if any of the query
// parameters are specified. If an option is set on ClientOptions after this method is called, that option will override
// any option applied via the connection string.
//
// If the URI format is incorrect or there are conflicting options specified in the URI an error will be recorded and
// can be retrieved by calling Validate.
//
// For more information about the URI format, see https://www.mongodb.com/docs/manual/reference/connection-string/. See
// mongo.Connect documentation for examples of using URIs for different Client configurations.
func (c *ClientOptions) ApplyURI(uri string) *ClientOptions {
	c.opts.ApplyURI(uri)

	return c
}

// SetAppName specifies an application name that is sent to the server when creating new connections. It is used by the
// server to log connection and profiling information (e.g. slow query logs). This can also be set through the "appName"
// URI option (e.g "appName=example_application"). The default is empty, meaning no app name will be sent.
func (c *ClientOptions) SetAppName(s string) *ClientOptions {
	c.opts.SetAppName(s)

	return c
}

// SetAuth specifies a Credential containing options for configuring authentication. See the options.Credential
// documentation for more information about Credential fields. The default is an empty Credential, meaning no
// authentication will be configured.
func (c *ClientOptions) SetAuth(auth _Credential) *ClientOptions {
	c.opts.SetAuth(auth)

	return c
}

// SetCompressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
// 3. "zstd" - requires server version >= 4.2, and driver version >= 1.2.0 with cgo support enabled or driver
// version >= 1.3.0 without cgo.
//
// If this option is specified, the driver will perform a negotiation with the server to determine a common list of of
// compressors and will use the first one in that list when performing operations. See
// https://www.mongodb.com/docs/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about configuring compression on the server and the server-side defaults.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func (c *ClientOptions) SetCompressors(comps []string) *ClientOptions {
	c.opts.SetCompressors(comps)

	return c
}

// SetConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
// specified through SetDialer, this option must not be used. This can be set through ApplyURI with the
// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
// seconds.
func (c *ClientOptions) SetConnectTimeout(d time.Duration) *ClientOptions {
	c.opts.SetConnectTimeout(d)

	return c
}

// SetDialer specifies a custom ContextDialer to be used to create new connections to the server. The default is a
// net.Dialer with the Timeout field set to ConnectTimeout. See https://golang.org/pkg/net/#Dialer for more information
// about the net.Dialer type.
func (c *ClientOptions) SetDialer(d _ContextDialer) *ClientOptions {
	c.opts.SetDialer(d)

	return c
}

// SetDirect specifies whether or not a direct connect should be made. If set to true, the driver will only connect to
// the host provided in the URI and will not discover other hosts in the cluster. This can also be set through the
// "directConnection" URI option. This option cannot be set to true if multiple hosts are specified, either through
// ApplyURI or SetHosts, or an SRV URI is used.
//
// As of driver version 1.4, the "connect" URI option has been deprecated and replaced with "directConnection". The
// "connect" URI option has two values:
//
// 1. "connect=direct" for direct connections. This corresponds to "directConnection=true".
//
// 2. "connect=automatic" for automatic discovery. This corresponds to "directConnection=false"
//
// If the "connect" and "directConnection" URI options are both specified in the connection string, their values must
// not conflict. Direct connections are not valid if multiple hosts are specified or an SRV URI is used. The default
// value for this option is false.
func (c *ClientOptions) SetDirect(b bool) *ClientOptions {
	c.opts.SetDirect(b)

	return c
}

// SetHeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
func (c *ClientOptions) SetHeartbeatInterval(d time.Duration) *ClientOptions {
	c.opts.SetHeartbeatInterval(d)

	return c
}

// SetHosts specifies a list of host names or IP addresses for servers in a cluster. Both IPv4 and IPv6 addresses are
// supported. IPv6 literals must be enclosed in '[]' following RFC-2732 syntax.
//
// Hosts can also be specified as a comma-separated list in a URI. For example, to include "localhost:27017" and
// "localhost:27018", a URI could be "mongodb://localhost:27017,localhost:27018". The default is ["localhost:27017"]
func (c *ClientOptions) SetHosts(s []string) *ClientOptions {
	c.opts.SetHosts(s)

	return c
}

// SetLoadBalanced specifies whether or not the MongoDB deployment is hosted behind a load balancer. This can also be
// set through the "loadBalanced" URI option. The driver will error during Client configuration if this option is set
// to true and one of the following conditions are met:
//
// 1. Multiple hosts are specified, either via the ApplyURI or SetHosts methods. This includes the case where an SRV
// URI is used and the SRV record resolves to multiple hostnames.
// 2. A replica set name is specified, either via the URI or the SetReplicaSet method.
// 3. The options specify whether or not a direct connection should be made, either via the URI or the SetDirect method.
//
// The default value is false.
func (c *ClientOptions) SetLoadBalanced(lb bool) *ClientOptions {
	c.opts.SetLoadBalanced(lb)

	return c
}

// SetLocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
// "localThresholdMS=15000"). The default is 15 milliseconds.
func (c *ClientOptions) SetLocalThreshold(d time.Duration) *ClientOptions {
	c.opts.SetLocalThreshold(d)

	return c
}

// SetLoggerOptions specifies a LoggerOptions containing options for
// configuring a logger.
func (c *ClientOptions) SetLoggerOptions(opts *_LoggerOptions) *ClientOptions {
	c.opts.SetLoggerOptions(opts)

	return c
}

// SetMaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
func (c *ClientOptions) SetMaxConnIdleTime(d time.Duration) *ClientOptions {
	c.opts.SetMaxConnIdleTime(d)

	return c
}

// SetMaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
// (e.g. "maxPoolSize=100"). If this is 0, maximum connection pool size is not limited. The default is 100.
func (c *ClientOptions) SetMaxPoolSize(u uint64) *ClientOptions {
	c.opts.SetMaxPoolSize(u)

	return c
}

// SetMinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
func (c *ClientOptions) SetMinPoolSize(u uint64) *ClientOptions {
	c.opts.SetMinPoolSize(u)

	return c
}

// SetMaxConnecting specifies the maximum number of connections a connection pool may establish simultaneously. This can
// also be set through the "maxConnecting" URI option (e.g. "maxConnecting=2"). If this is 0, the default is used. The
// default is 2. Values greater than 100 are not recommended.
func (c *ClientOptions) SetMaxConnecting(u uint64) *ClientOptions {
	c.opts.SetMaxConnecting(u)

	return c
}

// SetPoolMonitor specifies a PoolMonitor to receive connection pool events. See the event.PoolMonitor documentation
// for more information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetPoolMonitor(m *_event.PoolMonitor) *ClientOptions {
	c.opts.SetPoolMonitor(m)

	return c
}

// SetMonitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetMonitor(m *_event.CommandMonitor) *ClientOptions {
	c.opts.SetMonitor(m)

	return c
}

// SetServerMonitor specifies an SDAM monitor used to monitor SDAM events.
func (c *ClientOptions) SetServerMonitor(m *_event.ServerMonitor) *ClientOptions {
	c.opts.SetServerMonitor(m)

	return c
}

// SetReadConcern specifies the read concern to use for read operations. A read concern level can also be set through
// the "readConcernLevel" URI option (e.g. "readConcernLevel=majority"). The default is nil, meaning the server will use
// its configured default.
func (c *ClientOptions) SetReadConcern(rc *_readconcern.ReadConcern) *ClientOptions {
	c.opts.SetReadConcern(rc)

	return c
}

// SetReadPreference specifies the read preference to use for read operations. This can also be set through the
// following URI options:
//
// 1. "readPreference" - Specify the read preference mode (e.g. "readPreference=primary").
//
// 2. "readPreferenceTags": Specify one or more read preference tags
// (e.g. "readPreferenceTags=region:south,datacenter:A").
//
// 3. "maxStalenessSeconds" (or "maxStaleness"): Specify a maximum replication lag for reads from secondaries in a
// replica set (e.g. "maxStalenessSeconds=10").
//
// The default is readpref.Primary(). See https://www.mongodb.com/docs/manual/core/read-preference/#read-preference for
// more information about read preferences.
func (c *ClientOptions) SetReadPreference(rp *_readpref.ReadPref) *ClientOptions {
	c.opts.SetReadPreference(rp)

	return c
}

// SetBSONOptions configures optional BSON marshaling and unmarshaling behavior.
func (c *ClientOptions) SetBSONOptions(opts *_BSONOptions) *ClientOptions {
	c.opts.SetBSONOptions(opts)

	return c
}

// SetRegistry specifies the BSON registry to use for BSON marshalling/unmarshalling operations. The default is
// bson.DefaultRegistry.
func (c *ClientOptions) SetRegistry(registry *_bsoncodec.Registry) *ClientOptions {
	c.opts.SetRegistry(registry)

	return c
}

// SetReplicaSet specifies the replica set name for the cluster. If specified, the cluster will be treated as a replica
// set and the driver will automatically discover all servers in the set, starting with the nodes specified through
// ApplyURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be
// considered as part of the set by the Client. This can also be set through the "replicaSet" URI option (e.g.
// "replicaSet=replset"). The default is empty.
func (c *ClientOptions) SetReplicaSet(s string) *ClientOptions {
	c.opts.SetReplicaSet(s)

	return c
}

// SetRetryWrites specifies whether supported write operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are InsertOne, UpdateOne, ReplaceOne, DeleteOne, FindOneAndDelete, FindOneAndReplace,
// FindOneAndDelete, InsertMany, and BulkWrite. Note that BulkWrite requests must not include UpdateManyModel or
// DeleteManyModel instances to be considered retryable. Unacknowledged writes will not be retried, even if this option
// is set to true.
//
// This option requires server version >= 3.6 and a replica set or sharded cluster and will be ignored for any other
// cluster type. This can also be set through the "retryWrites" URI option (e.g. "retryWrites=true"). The default is
// true.
func (c *ClientOptions) SetRetryWrites(b bool) *ClientOptions {
	c.opts.SetRetryWrites(b)

	return c
}

// SetRetryReads specifies whether supported read operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are Find, FindOne, Aggregate without a $out stage, Distinct, CountDocuments,
// EstimatedDocumentCount, Watch (for Client, Database, and Collection), ListCollections, and ListDatabases. Note that
// operations run through RunCommand are not retried.
//
// This option requires server version >= 3.6 and driver version >= 1.1.0. The default is true.
func (c *ClientOptions) SetRetryReads(b bool) *ClientOptions {
	c.opts.SetRetryReads(b)

	return c
}

// SetServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
func (c *ClientOptions) SetServerSelectionTimeout(d time.Duration) *ClientOptions {
	c.opts.SetServerSelectionTimeout(d)

	return c
}

// SetSocketTimeout specifies how long the driver will wait for a socket read or write to return before returning a
// network error. This can also be set through the "socketTimeoutMS" URI option (e.g. "socketTimeoutMS=1000"). The
// default value is 0, meaning no timeout is used and socket operations can block indefinitely.
//
// NOTE(benjirewis): SocketTimeout will be deprecated in a future release. The more general Timeout option may be used
// in its place to control the amount of time that a single operation can run before returning an error. Setting
// SocketTimeout and Timeout on a single client will result in undefined behavior.
func (c *ClientOptions) SetSocketTimeout(d time.Duration) *ClientOptions {
	c.opts.SetSocketTimeout(d)

	return c
}

// SetTimeout specifies the amount of time that a single operation run on this Client can execute before returning an error.
// The deadline of any operation run through the Client will be honored above any Timeout set on the Client; Timeout will only
// be honored if there is no deadline on the operation Context. Timeout can also be set through the "timeoutMS" URI option
// (e.g. "timeoutMS=1000"). The default value is nil, meaning operations do not inherit a timeout from the Client.
//
// If any Timeout is set (even 0) on the Client, the values of MaxTime on operation options, TransactionOptions.MaxCommitTime and
// SessionOptions.DefaultMaxCommitTime will be ignored. Setting Timeout and SocketTimeout or WriteConcern.wTimeout will result
// in undefined behavior.
//
// NOTE(benjirewis): SetTimeout represents unstable, provisional API. The behavior of the driver when a Timeout is specified is
// subject to change.
func (c *ClientOptions) SetTimeout(d time.Duration) *ClientOptions {
	c.opts.SetTimeout(d)

	return c
}

// SetTLSConfig specifies a tls.Config instance to use use to configure TLS on all connections created to the cluster.
// This can also be set through the following URI options:
//
// 1. "tls" (or "ssl"): Specify if TLS should be used (e.g. "tls=true").
//
// 2. Either "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile") or a combination of "tlsCertificateFile" and
// "tlsPrivateKeyFile". The "tlsCertificateKeyFile" option specifies a path to the client certificate and private key,
// which must be concatenated into one file. The "tlsCertificateFile" and "tlsPrivateKey" combination specifies separate
// paths to the client certificate and private key, respectively. Note that if "tlsCertificateKeyFile" is used, the
// other two options must not be specified. Only the subject name of the first certificate is honored as the username
// for X509 auth in a file with multiple certs.
//
// 3. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 4. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// 5. "tlsInsecure" (or "sslInsecure"): Specifies whether or not certificates and hostnames received from the server
// should be validated. If true (e.g. "tlsInsecure=true"), the TLS library will accept any certificate presented by the
// server and any host name in that certificate. Note that setting this to true makes TLS susceptible to
// man-in-the-middle attacks and should only be done for testing.
//
// The default is nil, meaning no TLS will be enabled.
func (c *ClientOptions) SetTLSConfig(cfg *tls.Config) *ClientOptions {
	c.opts.SetTLSConfig(cfg)

	return c
}

// SetHTTPClient specifies the http.Client to be used for any HTTP requests.
//
// This should only be used to set custom HTTP client configurations. By default, the connection will use an internal.DefaultHTTPClient.
func (c *ClientOptions) SetHTTPClient(client *http.Client) *ClientOptions {
	c.opts.SetHTTPClient(client)

	return c
}

// SetWriteConcern specifies the write concern to use to for write operations. This can also be set through the following
// URI options:
//
// 1. "w": Specify the number of nodes in the cluster that must acknowledge write operations before the operation
// returns or "majority" to specify that a majority of the nodes must acknowledge writes. This can either be an integer
// (e.g. "w=10") or the string "majority" (e.g. "w=majority").
//
// 2. "wTimeoutMS": Specify how long write operations should wait for the correct number of nodes to acknowledge the
// operation (e.g. "wTimeoutMS=1000").
//
// 3. "journal": Specifies whether or not write operations should be written to an on-disk journal on the server before
// returning (e.g. "journal=true").
//
// The default is nil, meaning the server will use its configured default.
func (c *ClientOptions) SetWriteConcern(wc *_writeconcern.WriteConcern) *ClientOptions {
	c.opts.SetWriteConcern(wc)

	return c
}

// SetZlibLevel specifies the level for the zlib compressor. This option is ignored if zlib is not specified as a
// compressor through ApplyURI or SetCompressors. Supported values are -1 through 9, inclusive. -1 tells the zlib
// library to use its default, 0 means no compression, 1 means best speed, and 9 means best compression.
// This can also be set through the "zlibCompressionLevel" URI option (e.g. "zlibCompressionLevel=-1"). Defaults to -1.
func (c *ClientOptions) SetZlibLevel(level int) *ClientOptions {
	c.opts.SetZlibLevel(level)

	return c
}

// SetZstdLevel sets the level for the zstd compressor. This option is ignored if zstd is not specified as a compressor
// through ApplyURI or SetCompressors. Supported values are 1 through 20, inclusive. 1 means best speed and 20 means
// best compression. This can also be set through the "zstdCompressionLevel" URI option. Defaults to 6.
func (c *ClientOptions) SetZstdLevel(level int) *ClientOptions {
	c.opts.SetZstdLevel(level)

	return c
}

// SetAutoEncryptionOptions specifies an AutoEncryptionOptions instance to automatically encrypt and decrypt commands
// and their results. See the options.AutoEncryptionOptions documentation for more information about the supported
// options.
func (c *ClientOptions) SetAutoEncryptionOptions(opts *_AutoEncryptionOptions) *ClientOptions {
	c.opts.SetAutoEncryptionOptions(opts)

	return c
}

// SetDisableOCSPEndpointCheck specifies whether or not the driver should reach out to OCSP responders to verify the
// certificate status for certificates presented by the server that contain a list of OCSP responders.
//
// If set to true, the driver will verify the status of the certificate using a response stapled by the server, if there
// is one, but will not send an HTTP request to any responders if there is no staple. In this case, the driver will
// continue the connection even though the certificate status is not known.
//
// This can also be set through the tlsDisableOCSPEndpointCheck URI option. Both this URI option and tlsInsecure must
// not be set at the same time and will error if they are. The default value is false.
func (c *ClientOptions) SetDisableOCSPEndpointCheck(disableCheck bool) *ClientOptions {
	c.opts.SetDisableOCSPEndpointCheck(disableCheck)

	return c
}

// SetServerAPIOptions specifies a ServerAPIOptions instance used to configure the API version sent to the server
// when running commands. See the options.ServerAPIOptions documentation for more information about the supported
// options.
func (c *ClientOptions) SetServerAPIOptions(opts *_ServerAPIOptions) *ClientOptions {
	c.opts.SetServerAPIOptions(opts)

	return c
}

// SetSRVMaxHosts specifies the maximum number of SRV results to randomly select during polling. To limit the number
// of hosts selected in SRV discovery, this function must be called before ApplyURI. This can also be set through
// the "srvMaxHosts" URI option.
func (c *ClientOptions) SetSRVMaxHosts(srvMaxHosts int) *ClientOptions {
	c.opts.SetSRVMaxHosts(srvMaxHosts)

	return c
}

// SetSRVServiceName specifies a custom SRV service name to use in SRV polling. To use a custom SRV service name
// in SRV discovery, this function must be called before ApplyURI. This can also be set through the "srvServiceName"
// URI option.
func (c *ClientOptions) SetSRVServiceName(srvName string) *ClientOptions {
	c.opts.SetSRVServiceName(srvName)

	return c
}

// addCACertFromFile adds a root CA certificate to the configuration given a path
// to the containing file.
func addCACertFromFile(cfg *tls.Config, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if cfg.RootCAs == nil {
		cfg.RootCAs = x509.NewCertPool()
	}
	if !cfg.RootCAs.AppendCertsFromPEM(data) {
		return errors.New("the specified CA file does not contain any valid certificates")
	}

	return nil
}

func addClientCertFromSeparateFiles(cfg *tls.Config, keyFile, certFile, keyPassword string) (string, error) {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", err
	}
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return "", err
	}

	data := make([]byte, 0, len(keyData)+len(certData)+1)
	data = append(data, keyData...)
	data = append(data, '\n')
	data = append(data, certData...)
	return addClientCertFromBytes(cfg, data, keyPassword)
}

func addClientCertFromConcatenatedFile(cfg *tls.Config, certKeyFile, keyPassword string) (string, error) {
	data, err := ioutil.ReadFile(certKeyFile)
	if err != nil {
		return "", err
	}

	return addClientCertFromBytes(cfg, data, keyPassword)
}

// addClientCertFromBytes adds client certificates to the configuration given a path to the
// containing file and returns the subject name in the first certificate.
func addClientCertFromBytes(cfg *tls.Config, data []byte, keyPasswd string) (string, error) {
	var currentBlock *pem.Block
	var certDecodedBlock []byte
	var certBlocks, keyBlocks [][]byte

	remaining := data
	start := 0
	for {
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil {
			break
		}

		if currentBlock.Type == "CERTIFICATE" {
			certBlock := data[start : len(data)-len(remaining)]
			certBlocks = append(certBlocks, certBlock)
			// Assign the certDecodedBlock when it is never set,
			// so only the first certificate is honored in a file with multiple certs.
			if certDecodedBlock == nil {
				certDecodedBlock = currentBlock.Bytes
			}
			start += len(certBlock)
		} else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") {
			isEncrypted := x509.IsEncryptedPEMBlock(currentBlock) || strings.Contains(currentBlock.Type, "ENCRYPTED PRIVATE KEY")
			if isEncrypted {
				if keyPasswd == "" {
					return "", fmt.Errorf("no password provided to decrypt private key")
				}

				var keyBytes []byte
				var err error
				// Process the X.509-encrypted or PKCS-encrypted PEM block.
				if x509.IsEncryptedPEMBlock(currentBlock) {
					// Only covers encrypted PEM data with a DEK-Info header.
					keyBytes, err = x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
					if err != nil {
						return "", err
					}
				} else if strings.Contains(currentBlock.Type, "ENCRYPTED") {
					// The pkcs8 package only handles the PKCS #5 v2.0 scheme.
					decrypted, err := pkcs8.ParsePKCS8PrivateKey(currentBlock.Bytes, []byte(keyPasswd))
					if err != nil {
						return "", err
					}
					keyBytes, err = x509.MarshalPKCS8PrivateKey(decrypted)
					if err != nil {
						return "", err
					}
				}
				var encoded bytes.Buffer
				pem.Encode(&encoded, &pem.Block{Type: currentBlock.Type, Bytes: keyBytes})
				keyBlock := encoded.Bytes()
				keyBlocks = append(keyBlocks, keyBlock)
				start = len(data) - len(remaining)
			} else {
				keyBlock := data[start : len(data)-len(remaining)]
				keyBlocks = append(keyBlocks, keyBlock)
				start += len(keyBlock)
			}
		}
	}
	if len(certBlocks) == 0 {
		return "", fmt.Errorf("failed to find CERTIFICATE")
	}
	if len(keyBlocks) == 0 {
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	}

	cert, err := tls.X509KeyPair(bytes.Join(certBlocks, []byte("\n")), bytes.Join(keyBlocks, []byte("\n")))
	if err != nil {
		return "", err
	}

	cfg.Certificates = append(cfg.Certificates, cert)

	// The documentation for the tls.X509KeyPair indicates that the Leaf certificate is not
	// retained.
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil {
		return "", err
	}

	return crt.Subject.String(), nil
}

func stringSliceContains(source []string, target string) bool {
	for _, str := range source {
		if str == target {
			return true
		}
	}
	return false
}

// create a username for x509 authentication from an x509 certificate subject.
func extractX509UsernameFromSubject(subject string) string {
	// the Go x509 package gives the subject with the pairs in the reverse order from what we want.
	pairs := strings.Split(subject, ",")
	for left, right := 0, len(pairs)-1; left < right; left, right = left+1, right-1 {
		pairs[left], pairs[right] = pairs[right], pairs[left]
	}

	return strings.Join(pairs, ",")
}

func (c *ClientOptions) MongoOptions() *mongo_options.ClientOptions {
	return c.opts
}
