package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/oschwald/geoip2-golang"
	"gopkg.in/yaml.v2"
)

type Nodes map[string][]string

type Config struct {
	Addr    string
	Maxmind string
	Nodes   Nodes
}

type Server struct {
	Config *Config

	geodb *geoip2.Reader
}

func ServerWithConfig(cfg *Config) (*Server, error) {
	db, err := geoip2.Open(cfg.Maxmind)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config: cfg,
		geodb:  db,
	}, nil
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true

	e.GET("api_nodes", s.serveAPINodes)

	return e.Start(s.Config.Addr)
}

type apiNodesReply struct {
	Country string   `json:"country"`
	Nodes   []string `json:"nodes"`
}

func (s *Server) serveAPINodes(c echo.Context) error {
	addr := c.Request().RemoteAddr
	parts := strings.Split(addr, ":")

	host := parts[0]

	ip := net.ParseIP(host)
	if ip == nil {
		return errors.New("Invalid client IP")
	}

	country, err := s.geodb.Country(ip)
	if err != nil {
		return err
	}

	countryCode := strings.ToLower(country.Country.IsoCode)
	if countryCode == "" {
		countryCode = "default"
	}

	log.Println("lookup:", host, countryCode)

	nodes, ok := s.Config.Nodes[countryCode]

	if !ok {
		nodes = s.Config.Nodes["default"]
	}

	return c.JSON(http.StatusOK, apiNodesReply{
		Country: countryCode,
		Nodes:   nodes,
	})
}

func run(cfgFile string) error {
	cfg, err := loadConfig(cfgFile)
	if err != nil {
		return err
	}

	srv, err := ServerWithConfig(cfg)
	if err != nil {
		return err
	}

	return srv.Start()
}

func loadConfig(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg Config

	yaml.Unmarshal(b, &cfg)

	return &cfg, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Specify config file as first argument")
	}

	nodesFile := os.Args[1]

	err := run(nodesFile)
	if err != nil {
		log.Fatalln(err)
	}
}
