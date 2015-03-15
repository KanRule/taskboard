package impl

import (
	"bytes"
	"io"

	"github.com/dgrijalva/jwt-go"
	"github.com/quintans/goSQL/db"
	"github.com/quintans/goSQL/dbx"
	trx "github.com/quintans/goSQL/translators"
	tk "github.com/quintans/toolkit"
	"github.com/quintans/toolkit/log"
	"github.com/quintans/toolkit/web"
	"github.com/quintans/toolkit/web/app"
	"github.com/quintans/toolkit/web/poller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/msbranco/goconfig"

	T "github.com/quintans/taskboard/biz/tables"
	"github.com/quintans/taskboard/common/entity"
	"github.com/quintans/taskboard/common/lov"

	"compress/gzip"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	UNKNOWN = "TB00"
	DBFAIL  = "TB01"
	DBLOCK  = "TB02"

	PRINCIPAL_KEY = "principal"
	JWT_TIMEOUT   = 15
)

var (
	logger *log.Logger
	TM     db.ITransactionManager

	ContentDir     string
	IpPort         string
	postLimit      int64
	sessionTimeout time.Duration

	Poll *poller.Poller

	SmtpUser string
	SmtpPass string
	SmtpHost string
	SmtpPort string
	SmtpFrom string
)

var varguard = regexp.MustCompile("\\${[^{}]+}")

func replaceOsEnv(val string) string {
	// finds all keys in the OS environment
	values := make(map[string]string)
	keys := varguard.FindAllString(val, -1)
	for _, k := range keys {
		v := os.Getenv(k[2 : len(k)-1])
		values[k] = v
	}
	// replace the keys found
	for k, v := range values {
		val = strings.Replace(val, k, v, -1)
	}
	return val
}

func init() {
	c, err := goconfig.ReadConfigFile("taskboard.ini")
	if err != nil {
		panic(err)
	}
	// web
	ContentDir, _ = c.GetString("web", "dir")
	ContentDir = replaceOsEnv(ContentDir)
	fmt.Println("[web]dir=", ContentDir)
	IpPort, _ = c.GetString("web", "ip_port")
	IpPort = replaceOsEnv(IpPort)
	fmt.Println("[web]ip_port=", IpPort)
	timeout, _ := c.GetInt64("web", "session_timeout")
	fmt.Println("[web]session_timeout=", timeout)
	sessionTimeout = time.Duration(timeout) * time.Minute
	postLimit, _ = c.GetInt64("web", "post_limit")
	fmt.Println("[web]post_limit=", postLimit)

	//// log configuration - its the same the default
	level, _ := c.GetString("log", "level")
	level = replaceOsEnv(level)
	fmt.Println("[log]level=", level)
	file, _ := c.GetString("log", "file")
	file = replaceOsEnv(file)
	fmt.Println("[log]file=", file)
	count, _ := c.GetInt64("log", "count")
	fmt.Println("[log]count=", count)
	size, _ := c.GetInt64("log", "size")
	fmt.Println("[log]size=", size)
	logToConsole, _ := c.GetBool("log", "console")
	fmt.Println("[log]console=", logToConsole)

	logLevel := log.ParseLevel(level, log.ERROR)
	if logLevel <= log.INFO {
		log.ShowCaller(true)
	}

	// setting root writers
	writers := make([]log.LogWriter, 0)
	writers = append(writers, log.NewRollingFileAppender(file, size, int(count), true))
	if logToConsole {
		writers = append(writers, log.NewConsoleAppender(false))
	}
	log.Register("/", logLevel, writers...)

	//log.Register("/", logLevel, log.NewRollingFileAppender(file, size, int(count), true))

	//master.SetLevel("pqp", log.LogLevel(level))
	logger = log.LoggerFor("taskboad/biz/impl")

	// smtp
	SmtpHost, _ = c.GetString("smtp", "host")
	SmtpHost = replaceOsEnv(SmtpHost)
	fmt.Println("[smtp]host=", SmtpHost)
	SmtpPort, _ = c.GetString("smtp", "port")
	SmtpPort = replaceOsEnv(SmtpPort)
	fmt.Println("[smtp]port=", SmtpPort)
	SmtpUser, _ = c.GetString("smtp", "user")
	SmtpUser = replaceOsEnv(SmtpUser)
	fmt.Println("[smtp]user=", SmtpUser)
	SmtpPass, _ = c.GetString("smtp", "pass")
	SmtpPass = replaceOsEnv(SmtpPass)
	SmtpFrom, _ = c.GetString("smtp", "from")
	SmtpFrom = replaceOsEnv(SmtpFrom)
	fmt.Println("[smtp]from=", SmtpFrom)

	/*
	 * =======================
	 * BEGIN DATABASE CONFIG
	 * =======================
	 */
	// database configuration
	driverName, _ := c.GetString("database", "driver_name")
	driverName = replaceOsEnv(driverName)
	fmt.Println("[database]driver_name=", driverName)
	dataSourceName, _ := c.GetString("database", "data_source_name")
	dataSourceName = replaceOsEnv(dataSourceName)
	fmt.Println("[database]data_source_name=", dataSourceName)
	statementCache, _ := c.GetInt64("database", "statement_cache")
	fmt.Println("[database]statement_cache=", statementCache)
	idle, _ := c.GetInt64("database", "idle_connections")
	fmt.Println("[database]idle_connections=", idle)

	appDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	if idle > 0 {
		appDB.SetMaxIdleConns(int(idle))
	}

	// wake up the database pool
	err = appDB.Ping()
	if err != nil {
		panic(err)
	}

	TM = db.NewTransactionManager(
		// database
		appDB,
		// databse context factory
		func(inTx *bool, c dbx.IConnection) db.IDb {
			return db.NewDb(inTx, c, trx.NewMySQL5Translator())
		},
		// statement cache
		int(statementCache),
	)
	/*
	 * =======================
	 * END DATABASE CONFIG
	 * =======================
	 */

	Poll = poller.NewPoller(30 * time.Second)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

type Principal struct {
	UserId   int64
	Username string
	Roles    []lov.ERole
	Version  int64
}

func (this Principal) HasRole(role lov.ERole) bool {
	for _, r := range this.Roles {
		if r == role {
			return true
		}
	}
	return false
}

var secret, _ = GenerateRandomBytes(64)

func KeyFunction(token *jwt.Token) (interface{}, error) {
	return secret, nil
}

func serializePrincipal(p Principal) (string, error) {
	principal, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims[PRINCIPAL_KEY] = string(principal)
	// Expiration time in minutes
	token.Claims["exp"] = time.Now().Add(time.Minute * JWT_TIMEOUT).Unix()
	return token.SignedString(secret)
}

func deserializePrincipal(r *http.Request) *Principal {
	token, err := jwt.ParseFromRequest(r, KeyFunction)
	if err == nil && token.Valid {
		p := token.Claims[PRINCIPAL_KEY].(string)
		principal := &Principal{}
		json.Unmarshal([]byte(p), principal)
		return principal
	} else {
		return nil
	}
}

func PingFilter(ctx web.IContext) error {
	ctx.GetResponse().Header().Set("Content-Type", "text/html; charset=utf-8")

	// TODO what happens if this is called right after the user submited a password change
	// and the reply was not yet delivered/processed

	p := ctx.GetPrincipal().(Principal)
	// check if version is valid - if the user changed the password the version changes
	var uid int64
	store := ctx.(*AppCtx).Store
	ok, err := store.Query(T.USER).Column(T.USER_C_ID).
		Where(T.USER_C_ID.Matches(p.UserId).And(T.USER_C_VERSION.Matches(p.Version))).
		SelectInto(&uid)
	if !ok {
		// version is different
		logger.Debugf("Unable to revalidate the token, because the user no longer exists or it was changed (different version)")
		http.Error(ctx.GetResponse(), "Unauthorized", http.StatusUnauthorized)
		return nil
	} else if err != nil {
		return err
	}
	tokenString, err := serializePrincipal(p)
	if err != nil {
		return err
	}
	ctx.GetResponse().Write([]byte(tokenString))

	return nil
}

func LoginFilter(ctx web.IContext) error {
	ctx.GetResponse().Header().Set("Content-Type", "text/html; charset=utf-8")

	//logger.Debugf("serving static(): " + ctx.GetRequest().URL.Path)
	username := ctx.GetRequest().FormValue("username")
	pass := ctx.GetRequest().FormValue("password")
	var err error
	if username != "" && pass != "" {
		store := ctx.(*AppCtx).Store
		// usernames are stored in lowercase
		username = strings.ToLower(username)
		var user entity.User
		var ok bool
		if ok, err = store.Query(T.USER).
			All().
			Inner(T.USER_A_ROLES).Fetch().
			Where(
			db.And(T.USER_C_USERNAME.Matches(username),
				T.USER_C_DEAD.Matches(app.NOT_DELETED),
				T.USER_C_PASSWORD.Matches(pass))).
			SelectTree(&user); ok && err == nil {
			// role array
			roles := make([]lov.ERole, len(user.Roles))
			for k, v := range user.Roles {
				roles[k] = v.Kind
			}
			tokenString, err := serializePrincipal(Principal{
				*user.Id,
				*user.Username,
				roles,
				*user.Version,
			})
			if err != nil {
				return err
			}
			ctx.GetResponse().Write([]byte(tokenString))
		}
	}

	return err
}

func AuthenticationFilter(ctx web.IContext) error {
	p := deserializePrincipal(ctx.GetRequest())
	if p != nil {
		// for authorizations and business logic
		ctx.(*AppCtx).SetPrincipal(*p)
		return ctx.Proceed()
	} else {
		logger.Debugf("Unable to proceed: invalid token!")
		http.Error(ctx.GetResponse(), "Unauthorized", http.StatusUnauthorized)
	}

	return nil
}

func TransactionFilter(ctx web.IContext) error {
	return TM.Transaction(func(DB db.IDb) error {
		appCtx := ctx.(*AppCtx)
		appCtx.Store = DB
		p := ctx.GetPrincipal()
		if p != nil {
			appCtx.Store.SetAttribute(entity.ATTR_USERID, p.(Principal).UserId)
		}

		return ctx.Proceed()
	})
}

func NoTransactionFilter(ctx web.IContext) error {
	return TM.NoTransaction(func(DB db.IDb) error {
		appCtx := ctx.(*AppCtx)
		appCtx.Store = DB

		return ctx.Proceed()
	})
}

func ContextFactory(w http.ResponseWriter, r *http.Request) web.IContext {
	return NewAppCtx(w, r)
}

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// limits the body of a post, compress response and format eventual errors
func Limit(ctx web.IContext) (err error) {
	/*
		Very Important: Before compressing the response, the "Content-Type" header must be properly set!
	*/
	if strings.Contains(fmt.Sprint(ctx.GetRequest().Header["Accept-Encoding"]), "gzip") {
		appCtx := ctx.(*AppCtx)
		w := appCtx.Response
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		appCtx.Response = gzipResponseWriter{Writer: gz, ResponseWriter: w}
	}

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(runtime.Error); ok {
				logger.Errorf("%s\n========== Begin Stack Trace ==========\n%s\n========== End Stack Trace ==========\n", e, debug.Stack())
			}
			err = formatError(ctx.GetResponse(), r.(error))
		}
	}()

	logger.Debugf("requesting " + ctx.GetRequest().URL.Path)
	ctx.GetRequest().Body = http.MaxBytesReader(ctx.GetResponse(), ctx.GetRequest().Body, postLimit)
	err = ctx.Proceed()
	if err != nil {
		err = formatError(ctx.GetResponse(), err)
	}
	return err
}

func formatError(w http.ResponseWriter, err error) error {
	switch t := err.(type) {
	case *web.HttpFail:
		if t.Status == http.StatusInternalServerError {
			jsonError(w, t.GetCode(), t.GetMessage())
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.Error(w, t.Message, t.Status)
		}
	case *dbx.PersistenceFail:
		logger.Errorf("%s", err)
		jsonError(w, t.GetCode(), t.GetMessage())
	case *dbx.OptimisticLockFail:
		logger.Errorf("%s", err)
		jsonError(w, t.GetCode(), t.GetMessage())
	case tk.Fault:
		jsonError(w, t.GetCode(), t.GetMessage())
	default:
		jsonError(w, UNKNOWN, err.Error())
	}
	return nil
}

func jsonError(w http.ResponseWriter, code string, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	jmsg, _ := json.Marshal(msg)
	http.Error(w, fmt.Sprintf(`{"code":"%s", "message":%s}`, code, jmsg), http.StatusInternalServerError)
}

type MyResponse struct {
	Body   *bytes.Buffer
	header http.Header
	Code   int
}

func NewMyResponse() *MyResponse {
	return &MyResponse{
		Body:   new(bytes.Buffer),
		header: make(http.Header),
	}
}

func (this *MyResponse) Header() http.Header {
	return this.header
}

func (this *MyResponse) Write(data []byte) (int, error) {
	if this.Code == 0 {
		this.Code = http.StatusOK
	}
	return this.Body.Write(data)
}

func (this *MyResponse) WriteHeader(code int) {
	this.Code = code
}

var _ http.ResponseWriter = &MyResponse{}

// buffer the response, permiting setting headers after starting writing the response.
// It also gzips the response if the client browser supports it.
func ResponseBuffer(ctx web.IContext) error {
	appCtx := ctx.(*AppCtx)
	rec := NewMyResponse()
	w := appCtx.Response
	// passing a buffer instead of the original RW
	appCtx.Response = rec
	// restores the original response, even in the case of a panic
	defer func() {
		appCtx.Response = w
	}()
	err := ctx.Proceed()
	if err == nil {
		// we copy the original headers first
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		// status code
		w.WriteHeader(rec.Code)

		// then write out the original body
		w.Write(rec.Body.Bytes())
	}

	return err
}
